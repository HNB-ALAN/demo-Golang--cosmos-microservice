// Package auth provides authentication and authorization utilities for USC platform services.
package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// ContextKey represents context key for user information
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// UserEmailKey is the context key for user email
	UserEmailKey ContextKey = "user_email"
	// UserRoleKey is the context key for user role
	UserRoleKey ContextKey = "user_role"
	// UserPermissionsKey is the context key for user permissions
	UserPermissionsKey ContextKey = "user_permissions"
	// UserMetadataKey is the context key for user metadata
	UserMetadataKey ContextKey = "user_metadata"
	// ClaimsKey is the context key for JWT claims
	ClaimsKey ContextKey = "claims"
)

// AuthMiddleware represents authentication middleware
type AuthMiddleware struct {
	jwtService *JWTService
	required   bool
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtService *JWTService, required bool) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		required:   required,
	}
}

// HTTPMiddleware returns HTTP authentication middleware
func (m *AuthMiddleware) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			if m.required {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		// Parse Bearer token
		token := m.extractBearerToken(authHeader)
		if token == "" {
			if m.required {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			if m.required {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		// Add user information to context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		ctx = context.WithValue(ctx, UserPermissionsKey, claims.Permissions)
		ctx = context.WithValue(ctx, UserMetadataKey, claims.Metadata)
		ctx = context.WithValue(ctx, ClaimsKey, claims)

		// Continue with authenticated request
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GRPCUnaryInterceptor returns gRPC unary authentication interceptor
func (m *AuthMiddleware) GRPCUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			if m.required {
				return nil, status.Error(codes.Unauthenticated, "metadata not found")
			}
			return handler(ctx, req)
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			if m.required {
				return nil, status.Error(codes.Unauthenticated, "authorization header not found")
			}
			return handler(ctx, req)
		}

		// Parse Bearer token
		token := m.extractBearerToken(authHeader[0])
		if token == "" {
			if m.required {
				return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
			}
			return handler(ctx, req)
		}

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			if m.required {
				return nil, status.Error(codes.Unauthenticated, "invalid token")
			}
			return handler(ctx, req)
		}

		// Add user information to context
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		ctx = context.WithValue(ctx, UserPermissionsKey, claims.Permissions)
		ctx = context.WithValue(ctx, UserMetadataKey, claims.Metadata)
		ctx = context.WithValue(ctx, ClaimsKey, claims)

		// Continue with authenticated request
		return handler(ctx, req)
	}
}

// GRPCStreamInterceptor returns gRPC stream authentication interceptor
func (m *AuthMiddleware) GRPCStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			if m.required {
				return status.Error(codes.Unauthenticated, "metadata not found")
			}
			return handler(srv, ss)
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			if m.required {
				return status.Error(codes.Unauthenticated, "authorization header not found")
			}
			return handler(srv, ss)
		}

		// Parse Bearer token
		token := m.extractBearerToken(authHeader[0])
		if token == "" {
			if m.required {
				return status.Error(codes.Unauthenticated, "invalid authorization header format")
			}
			return handler(srv, ss)
		}

		// Validate token
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			if m.required {
				return status.Error(codes.Unauthenticated, "invalid token")
			}
			return handler(srv, ss)
		}

		// Add user information to context
		ctx := context.WithValue(ss.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
		ctx = context.WithValue(ctx, UserPermissionsKey, claims.Permissions)
		ctx = context.WithValue(ctx, UserMetadataKey, claims.Metadata)
		ctx = context.WithValue(ctx, ClaimsKey, claims)

		// Create new stream with authenticated context
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		// Continue with authenticated request
		return handler(srv, wrappedStream)
	}
}

// extractBearerToken extracts Bearer token from authorization header
func (m *AuthMiddleware) extractBearerToken(authHeader string) string {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return ""
	}
	return strings.TrimSpace(authHeader[len(bearerPrefix):])
}

// wrappedServerStream wraps gRPC ServerStream with custom context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the custom context
func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}

// GetUserEmailFromContext extracts user email from context
func GetUserEmailFromContext(ctx context.Context) (string, error) {
	email, ok := ctx.Value(UserEmailKey).(string)
	if !ok {
		return "", errors.New("user email not found in context")
	}
	return email, nil
}

// GetUserRoleFromContext extracts user role from context
func GetUserRoleFromContext(ctx context.Context) (string, error) {
	role, ok := ctx.Value(UserRoleKey).(string)
	if !ok {
		return "", errors.New("user role not found in context")
	}
	return role, nil
}

// GetUserPermissionsFromContext extracts user permissions from context
func GetUserPermissionsFromContext(ctx context.Context) ([]string, error) {
	permissions, ok := ctx.Value(UserPermissionsKey).([]string)
	if !ok {
		return nil, errors.New("user permissions not found in context")
	}
	return permissions, nil
}

// GetUserMetadataFromContext extracts user metadata from context
func GetUserMetadataFromContext(ctx context.Context) (map[string]string, error) {
	metadata, ok := ctx.Value(UserMetadataKey).(map[string]string)
	if !ok {
		return nil, errors.New("user metadata not found in context")
	}
	return metadata, nil
}

// GetClaimsFromContext extracts JWT claims from context
func GetClaimsFromContext(ctx context.Context) (*Claims, error) {
	claims, ok := ctx.Value(ClaimsKey).(*Claims)
	if !ok {
		return nil, errors.New("claims not found in context")
	}
	return claims, nil
}

// RequirePermission creates a middleware that requires specific permission
func (m *AuthMiddleware) RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := GetClaimsFromContext(r.Context())
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			if !m.jwtService.HasPermission(claims, permission) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRole creates a middleware that requires specific role
func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := GetClaimsFromContext(r.Context())
			if err != nil {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			if !m.jwtService.HasRole(claims, role) {
				http.Error(w, "Insufficient role", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermissionGRPC creates a gRPC interceptor that requires specific permission
func (m *AuthMiddleware) RequirePermissionGRPC(permission string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		claims, err := GetClaimsFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		if !m.jwtService.HasPermission(claims, permission) {
			return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
		}

		return handler(ctx, req)
	}
}

// RequireRoleGRPC creates a gRPC interceptor that requires specific role
func (m *AuthMiddleware) RequireRoleGRPC(role string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		claims, err := GetClaimsFromContext(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		if !m.jwtService.HasRole(claims, role) {
			return nil, status.Error(codes.PermissionDenied, "insufficient role")
		}

		return handler(ctx, req)
	}
}
