// Package auth provides authentication and authorization utilities for USC platform services.
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

// TokenType represents the type of token
type TokenType string

const (
	// AccessTokenType represents access token
	AccessTokenType TokenType = "access"
	// RefreshTokenType represents refresh token
	RefreshTokenType TokenType = "refresh"
	// APIKeyType represents API key
	APIKeyType TokenType = "api_key"
	// SessionTokenType represents session token
	SessionTokenType TokenType = "session"
)

// TokenInfo represents token information
type TokenInfo struct {
	Token     string            `json:"token"`
	Type      TokenType         `json:"type"`
	UserID    string            `json:"user_id"`
	ExpiresAt time.Time         `json:"expires_at"`
	CreatedAt time.Time         `json:"created_at"`
	Metadata  map[string]string `json:"metadata"`
}

// TokenStore represents a token storage interface
type TokenStore interface {
	Store(token *TokenInfo) error
	Get(token string) (*TokenInfo, error)
	Delete(token string) error
	DeleteByUserID(userID string) error
	CleanupExpired() error
}

// InMemoryTokenStore implements TokenStore using in-memory storage
type InMemoryTokenStore struct {
	tokens map[string]*TokenInfo
}

// NewInMemoryTokenStore creates a new in-memory token store
func NewInMemoryTokenStore() *InMemoryTokenStore {
	return &InMemoryTokenStore{
		tokens: make(map[string]*TokenInfo),
	}
}

// Store stores a token
func (s *InMemoryTokenStore) Store(token *TokenInfo) error {
	s.tokens[token.Token] = token
	return nil
}

// Get retrieves a token
func (s *InMemoryTokenStore) Get(token string) (*TokenInfo, error) {
	tokenInfo, exists := s.tokens[token]
	if !exists {
		return nil, errors.New("token not found")
	}

	// Check if token is expired
	if time.Now().After(tokenInfo.ExpiresAt) {
		delete(s.tokens, token)
		return nil, errors.New("token expired")
	}

	return tokenInfo, nil
}

// Delete deletes a token
func (s *InMemoryTokenStore) Delete(token string) error {
	delete(s.tokens, token)
	return nil
}

// DeleteByUserID deletes all tokens for a user
func (s *InMemoryTokenStore) DeleteByUserID(userID string) error {
	for token, info := range s.tokens {
		if info.UserID == userID {
			delete(s.tokens, token)
		}
	}
	return nil
}

// CleanupExpired removes expired tokens
func (s *InMemoryTokenStore) CleanupExpired() error {
	now := time.Now()
	for token, info := range s.tokens {
		if now.After(info.ExpiresAt) {
			delete(s.tokens, token)
		}
	}
	return nil
}

// TokenManager manages token operations
type TokenManager struct {
	store         TokenStore
	jwtService    *JWTService
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewTokenManager creates a new token manager
func NewTokenManager(store TokenStore, jwtService *JWTService, accessExpiry, refreshExpiry time.Duration) *TokenManager {
	return &TokenManager{
		store:         store,
		jwtService:    jwtService,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateTokenPair generates access and refresh token pair
func (tm *TokenManager) GenerateTokenPair(userID, email, role string, permissions []string, metadata map[string]string) (*TokenPair, error) {
	// Generate JWT tokens
	tokenPair, err := tm.jwtService.GenerateTokenPair(userID, email, role, permissions, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT tokens: %w", err)
	}

	// Store refresh token
	refreshTokenInfo := &TokenInfo{
		Token:     tokenPair.RefreshToken,
		Type:      RefreshTokenType,
		UserID:    userID,
		ExpiresAt: time.Now().Add(tm.refreshExpiry),
		CreatedAt: time.Now(),
		Metadata:  metadata,
	}

	if err := tm.store.Store(refreshTokenInfo); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return tokenPair, nil
}

// ValidateAccessToken validates an access token
func (tm *TokenManager) ValidateAccessToken(token string) (*Claims, error) {
	return tm.jwtService.ValidateToken(token)
}

// RefreshAccessToken refreshes an access token using refresh token
func (tm *TokenManager) RefreshAccessToken(refreshToken string) (string, error) {
	// Validate refresh token
	claims, err := tm.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if refresh token exists in store
	_, err = tm.store.Get(refreshToken)
	if err != nil {
		return "", fmt.Errorf("refresh token not found in store: %w", err)
	}

	// Generate new access token
	newAccessToken, err := tm.jwtService.GenerateAccessToken(
		claims.UserID,
		claims.Email,
		claims.Role,
		claims.Permissions,
		claims.Metadata,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	return newAccessToken, nil
}

// RevokeToken revokes a token
func (tm *TokenManager) RevokeToken(token string) error {
	return tm.store.Delete(token)
}

// RevokeUserTokens revokes all tokens for a user
func (tm *TokenManager) RevokeUserTokens(userID string) error {
	return tm.store.DeleteByUserID(userID)
}

// GenerateAPIKey generates a new API key
func (tm *TokenManager) GenerateAPIKey(userID string, expiry time.Duration, metadata map[string]string) (string, error) {
	// Generate random API key
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	apiKey := base64.URLEncoding.EncodeToString(keyBytes)

	// Store API key
	tokenInfo := &TokenInfo{
		Token:     apiKey,
		Type:      APIKeyType,
		UserID:    userID,
		ExpiresAt: time.Now().Add(expiry),
		CreatedAt: time.Now(),
		Metadata:  metadata,
	}

	if err := tm.store.Store(tokenInfo); err != nil {
		return "", fmt.Errorf("failed to store API key: %w", err)
	}

	return apiKey, nil
}

// ValidateAPIKey validates an API key
func (tm *TokenManager) ValidateAPIKey(apiKey string) (*TokenInfo, error) {
	tokenInfo, err := tm.store.Get(apiKey)
	if err != nil {
		return nil, err
	}

	if tokenInfo.Type != APIKeyType {
		return nil, errors.New("invalid token type")
	}

	return tokenInfo, nil
}

// GenerateSessionToken generates a session token
func (tm *TokenManager) GenerateSessionToken(userID string, expiry time.Duration, metadata map[string]string) (string, error) {
	// Generate random session token
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	sessionToken := base64.URLEncoding.EncodeToString(keyBytes)

	// Store session token
	tokenInfo := &TokenInfo{
		Token:     sessionToken,
		Type:      SessionTokenType,
		UserID:    userID,
		ExpiresAt: time.Now().Add(expiry),
		CreatedAt: time.Now(),
		Metadata:  metadata,
	}

	if err := tm.store.Store(tokenInfo); err != nil {
		return "", fmt.Errorf("failed to store session token: %w", err)
	}

	return sessionToken, nil
}

// ValidateSessionToken validates a session token
func (tm *TokenManager) ValidateSessionToken(sessionToken string) (*TokenInfo, error) {
	tokenInfo, err := tm.store.Get(sessionToken)
	if err != nil {
		return nil, err
	}

	if tokenInfo.Type != SessionTokenType {
		return nil, errors.New("invalid token type")
	}

	return tokenInfo, nil
}

// GetUserTokens returns all tokens for a user
func (tm *TokenManager) GetUserTokens(userID string) ([]*TokenInfo, error) {
	// This is a simplified implementation
	// In a real implementation, you would query the store for all tokens by user ID
	return nil, errors.New("not implemented")
}

// CleanupExpiredTokens removes expired tokens
func (tm *TokenManager) CleanupExpiredTokens() error {
	return tm.store.CleanupExpired()
}

// TokenBlacklist represents a token blacklist
type TokenBlacklist struct {
	blacklistedTokens map[string]time.Time
}

// NewTokenBlacklist creates a new token blacklist
func NewTokenBlacklist() *TokenBlacklist {
	return &TokenBlacklist{
		blacklistedTokens: make(map[string]time.Time),
	}
}

// AddToken adds a token to the blacklist
func (tb *TokenBlacklist) AddToken(token string, expiry time.Time) {
	tb.blacklistedTokens[token] = expiry
}

// IsBlacklisted checks if a token is blacklisted
func (tb *TokenBlacklist) IsBlacklisted(token string) bool {
	expiry, exists := tb.blacklistedTokens[token]
	if !exists {
		return false
	}

	// Remove expired tokens
	if time.Now().After(expiry) {
		delete(tb.blacklistedTokens, token)
		return false
	}

	return true
}

// RemoveToken removes a token from the blacklist
func (tb *TokenBlacklist) RemoveToken(token string) {
	delete(tb.blacklistedTokens, token)
}

// CleanupExpired removes expired tokens from blacklist
func (tb *TokenBlacklist) CleanupExpired() {
	now := time.Now()
	for token, expiry := range tb.blacklistedTokens {
		if now.After(expiry) {
			delete(tb.blacklistedTokens, token)
		}
	}
}

// TokenValidator validates tokens with blacklist support
type TokenValidator struct {
	jwtService *JWTService
	blacklist  *TokenBlacklist
}

// NewTokenValidator creates a new token validator
func NewTokenValidator(jwtService *JWTService, blacklist *TokenBlacklist) *TokenValidator {
	return &TokenValidator{
		jwtService: jwtService,
		blacklist:  blacklist,
	}
}

// ValidateToken validates a token with blacklist check
func (tv *TokenValidator) ValidateToken(token string) (*Claims, error) {
	// Check if token is blacklisted
	if tv.blacklist != nil && tv.blacklist.IsBlacklisted(token) {
		return nil, errors.New("token is blacklisted")
	}

	// Validate JWT token
	claims, err := tv.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// BlacklistToken adds a token to the blacklist
func (tv *TokenValidator) BlacklistToken(token string, expiry time.Time) {
	if tv.blacklist != nil {
		tv.blacklist.AddToken(token, expiry)
	}
}
