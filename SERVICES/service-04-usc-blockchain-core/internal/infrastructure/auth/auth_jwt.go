package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"service-04/internal/infrastructure/database"

	"github.com/usc-platform/shared/auth"
	"github.com/usc-platform/shared/config"
	"github.com/usc-platform/shared/errors"
	"github.com/usc-platform/shared/logging"
	"google.golang.org/grpc/codes"
)

// JWTService handles JWT token operations using shared library
type JWTService struct {
	jwtService   *auth.JWTService
	redisManager *database.RedisManager
	logger       logging.Logger
}

// NewJWTService creates a new JWT service instance
func NewJWTService(cfg *config.Config, logger logging.Logger, redisManager *database.RedisManager) (*JWTService, error) {
	// Parse durations from config with safe defaults
	accessDur, err := time.ParseDuration(cfg.Auth.JWTExpiry)
	if err != nil || accessDur <= 0 {
		accessDur = 24 * time.Hour
	}
	refreshDur, err := time.ParseDuration(cfg.Auth.RefreshExpiry)
	if err != nil || refreshDur <= 0 {
		refreshDur = 7 * 24 * time.Hour
	}

	// Get JWT secret from environment variable first, fallback to config
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = cfg.Auth.JWTSecret
	}
	if jwtSecret == "" {
		return nil, errors.NewGRPCError(codes.Internal, "JWT secret not configured", "JWT_SECRET environment variable or config.auth.jwt_secret must be set")
	}

	// Create JWT service using shared library
	jwtSvc, err := auth.NewJWTService(auth.Config{
		SecretKey:     jwtSecret,
		Issuer:        cfg.Auth.Issuer,
		AccessExpiry:  accessDur,
		RefreshExpiry: refreshDur,
	})
	if err != nil {
		return nil, errors.NewGRPCError(codes.Internal, "failed to initialize JWT service", err.Error())
	}

	return &JWTService{
		jwtService:   jwtSvc,
		redisManager: redisManager,
		logger:       logger,
	}, nil
}

// GenerateTokenPair generates access and refresh tokens
func (s *JWTService) GenerateTokenPair(ctx context.Context, userID, email, role string, permissions []string) (*auth.TokenPair, error) {
	s.logger.Info("Generating token pair",
		logging.String("user_id", userID),
		logging.String("email", email),
		logging.String("role", role))

	tokenPair, err := s.jwtService.GenerateTokenPair(userID, email, role, permissions, make(map[string]string))
	if err != nil {
		s.logger.Error("Failed to generate token pair", logging.Error(err))
		return nil, errors.NewGRPCError(codes.Internal, "failed to generate tokens", err.Error())
	}

	s.logger.Info("Token pair generated successfully",
		logging.String("user_id", userID))

	return tokenPair, nil
}

// ValidateToken validates a JWT token and returns claims
func (s *JWTService) ValidateToken(ctx context.Context, tokenString string) (*auth.Claims, error) {
	s.logger.Debug("Validating token")

	// Check blacklist first (before validating token)
	if s.redisManager != nil {
		tokenHash := sha256.Sum256([]byte(tokenString))
		tokenHashStr := hex.EncodeToString(tokenHash[:])
		blacklistKey := fmt.Sprintf("token:blacklist:%s", tokenHashStr)

		exists, err := s.redisManager.Exists(ctx, blacklistKey)
		if err == nil && exists {
			s.logger.Warn("Token found in blacklist",
				logging.String("token_hash_prefix", tokenHashStr[:16]))
			return nil, errors.NewGRPCError(codes.Unauthenticated, "token has been revoked", "token is blacklisted")
		}
	}

	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		s.logger.Error("Token validation failed", logging.Error(err))
		return nil, errors.NewGRPCError(codes.Unauthenticated, "invalid token", err.Error())
	}

	s.logger.Debug("Token validated successfully",
		logging.String("user_id", claims.UserID))

	return claims, nil
}

// RefreshToken generates a new access token using refresh token
func (s *JWTService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	s.logger.Info("Refreshing token")

	newAccessToken, err := s.jwtService.RefreshToken(refreshToken)
	if err != nil {
		s.logger.Error("Token refresh failed", logging.Error(err))
		return "", errors.NewGRPCError(codes.Unauthenticated, "invalid refresh token", err.Error())
	}

	s.logger.Info("Token refreshed successfully")
	return newAccessToken, nil
}

// RevokeToken revokes a token (adds to blacklist)
func (s *JWTService) RevokeToken(ctx context.Context, tokenString string) error {
	s.logger.Info("Revoking token",
		logging.String("token_prefix", tokenString[:min(10, len(tokenString))]))

	// Validate token first to get expiry time
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != nil {
		s.logger.Warn("Failed to validate token before revocation",
			logging.Error(err))
		// Still try to blacklist even if validation fails (defensive)
	}

	// Calculate token hash for blacklist key
	tokenHash := sha256.Sum256([]byte(tokenString))
	tokenHashStr := hex.EncodeToString(tokenHash[:])
	blacklistKey := fmt.Sprintf("token:blacklist:%s", tokenHashStr)

	// Calculate expiry time - use token expiry if available, otherwise use default
	var expiry time.Duration
	if claims != nil && !claims.ExpiresAt.IsZero() {
		expiry = time.Until(claims.ExpiresAt.Time)
		if expiry <= 0 {
			// Token already expired, use short TTL
			expiry = 1 * time.Hour
		}
	} else {
		// Default to 24 hours if expiry not available
		expiry = 24 * time.Hour
	}

	// Add token to Redis blacklist
	if s.redisManager != nil {
		if err := s.redisManager.Set(ctx, blacklistKey, "1", expiry); err != nil {
			s.logger.Warn("Failed to add token to blacklist",
				logging.Error(err),
				logging.String("token_hash_prefix", tokenHashStr[:16]))
			// Don't fail revocation, just log warning
			// Token will still be invalid after expiry
		} else {
			s.logger.Info("Token added to blacklist successfully",
				logging.String("token_hash_prefix", tokenHashStr[:16]),
				logging.Duration("expiry", expiry))
		}
	} else {
		s.logger.Warn("Redis manager not available, token blacklist not implemented",
			logging.String("note", "Token will only be invalid after expiry"))
	}

	s.logger.Info("Token revoked successfully",
		logging.String("token_hash_prefix", tokenHashStr[:16]))
	return nil
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ExtractUserID extracts user ID from token without full validation
func (s *JWTService) ExtractUserID(ctx context.Context, tokenString string) (string, error) {
	s.logger.Debug("Extracting user ID from token")

	userID, err := s.jwtService.ExtractUserID(tokenString)
	if err != nil {
		s.logger.Error("Failed to extract user ID", logging.Error(err))
		return "", errors.NewGRPCError(codes.Unauthenticated, "invalid token", err.Error())
	}

	return userID, nil
}

// HasPermission checks if user has specific permission
func (s *JWTService) HasPermission(ctx context.Context, claims *auth.Claims, permission string) bool {
	s.logger.Debug("Checking permission",
		logging.String("user_id", claims.UserID),
		logging.String("permission", permission))

	hasPermission := s.jwtService.HasPermission(claims, permission)

	s.logger.Debug("Permission check result",
		logging.String("user_id", claims.UserID),
		logging.String("permission", permission),
		logging.Bool("has_permission", hasPermission))

	return hasPermission
}

// GetUserPermissions returns all permissions for a user
func (s *JWTService) GetUserPermissions(ctx context.Context, claims *auth.Claims) []string {
	s.logger.Debug("Getting user permissions",
		logging.String("user_id", claims.UserID))

	// Return permissions from claims
	permissions := claims.Permissions

	s.logger.Debug("User permissions retrieved",
		logging.String("user_id", claims.UserID),
		logging.Int("permission_count", len(permissions)))

	return permissions
}

// Close closes the JWT service
func (s *JWTService) Close() error {
	s.logger.Info("Closing JWT service")
	// JWT service doesn't need explicit closing, but you can add cleanup logic here if needed
	return nil
}
