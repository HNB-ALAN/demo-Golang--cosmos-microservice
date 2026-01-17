// Package auth provides authentication and authorization utilities for USC platform services.
package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService handles JWT token operations
type JWTService struct {
	secretKey     []byte
	privateKey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
	issuer        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// Claims represents JWT claims
type Claims struct {
	UserID      string            `json:"user_id"`
	Email       string            `json:"email"`
	Role        string            `json:"role"`
	Permissions []string          `json:"permissions"`
	Metadata    map[string]string `json:"metadata"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh token pair
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Config represents JWT service configuration
type Config struct {
	SecretKey     string        `mapstructure:"secret_key"`
	PrivateKey    string        `mapstructure:"private_key"`
	PublicKey     string        `mapstructure:"public_key"`
	Issuer        string        `mapstructure:"issuer"`
	AccessExpiry  time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
}

// NewJWTService creates a new JWT service
func NewJWTService(cfg Config) (*JWTService, error) {
	if cfg.SecretKey == "" {
		return nil, errors.New("secret key is required")
	}

	service := &JWTService{
		secretKey:     []byte(cfg.SecretKey),
		issuer:        cfg.Issuer,
		accessExpiry:  cfg.AccessExpiry,
		refreshExpiry: cfg.RefreshExpiry,
	}

	// Load RSA keys if provided
	if cfg.PrivateKey != "" && cfg.PublicKey != "" {
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.PublicKey))
		if err != nil {
			return nil, fmt.Errorf("failed to parse public key: %w", err)
		}

		service.privateKey = privateKey
		service.publicKey = publicKey
	}

	return service, nil
}

// GenerateTokenPair generates both access and refresh tokens
func (j *JWTService) GenerateTokenPair(userID, email, role string, permissions []string, metadata map[string]string) (*TokenPair, error) {
	// Generate access token
	accessToken, err := j.GenerateAccessToken(userID, email, role, permissions, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := j.GenerateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(j.accessExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// GenerateAccessToken generates an access token
func (j *JWTService) GenerateAccessToken(userID, email, role string, permissions []string, metadata map[string]string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:      userID,
		Email:       email,
		Role:        role,
		Permissions: permissions,
		Metadata:    metadata,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID,
			Audience:  []string{"usc-platform"},
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessExpiry)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// GenerateRefreshToken generates a refresh token
func (j *JWTService) GenerateRefreshToken(userID string) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID,
			Audience:  []string{"usc-platform"},
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshExpiry)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ValidateToken validates a JWT token and returns claims
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// RefreshToken validates a refresh token and generates new access token
func (j *JWTService) RefreshToken(refreshToken string) (string, error) {
	claims, err := j.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generate new access token with same user info
	accessToken, err := j.GenerateAccessToken(claims.UserID, claims.Email, claims.Role, claims.Permissions, claims.Metadata)
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	return accessToken, nil
}

// ExtractUserID extracts user ID from token without full validation
func (j *JWTService) ExtractUserID(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	return claims.UserID, nil
}

// HasPermission checks if user has specific permission
func (j *JWTService) HasPermission(claims *Claims, permission string) bool {
	for _, p := range claims.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// HasRole checks if user has specific role
func (j *JWTService) HasRole(claims *Claims, role string) bool {
	return claims.Role == role
}

// IsExpired checks if token is expired
func (j *JWTService) IsExpired(tokenString string) bool {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return true
	}

	return claims.ExpiresAt.Before(time.Now())
}

// GetTokenInfo returns token information without validation
func (j *JWTService) GetTokenInfo(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return map[string]interface{}{
		"user_id":     claims.UserID,
		"email":       claims.Email,
		"role":        claims.Role,
		"permissions": claims.Permissions,
		"metadata":    claims.Metadata,
		"expires_at":  claims.ExpiresAt,
		"issued_at":   claims.IssuedAt,
		"issuer":      claims.Issuer,
	}, nil
}
