package auth

import (
	"testing"
	"time"
)

func TestNewJWTService(t *testing.T) {
	config := Config{
		SecretKey:     "test-secret-key",
		Issuer:        "test-issuer",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	service, err := NewJWTService(config)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if service == nil {
		t.Fatal("Expected service to be created")
	}

	if service.issuer != config.Issuer {
		t.Errorf("Expected issuer %s, got %s", config.Issuer, service.issuer)
	}

	if service.accessExpiry != config.AccessExpiry {
		t.Errorf("Expected access expiry %v, got %v", config.AccessExpiry, service.accessExpiry)
	}

	if service.refreshExpiry != config.RefreshExpiry {
		t.Errorf("Expected refresh expiry %v, got %v", config.RefreshExpiry, service.refreshExpiry)
	}
}

func TestNewJWTServiceWithEmptySecret(t *testing.T) {
	config := Config{
		SecretKey: "",
		Issuer:    "test-issuer",
	}

	_, err := NewJWTService(config)
	if err == nil {
		t.Error("Expected error for empty secret key, got nil")
	}
}

func TestJWTService_GenerateTokenPair(t *testing.T) {
	config := Config{
		SecretKey:     "test-secret-key",
		Issuer:        "test-issuer",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	service, err := NewJWTService(config)
	if err != nil {
		t.Fatalf("Failed to create JWT service: %v", err)
	}

	userID := "user-123"
	email := "user@example.com"
	role := "user"
	permissions := []string{"read", "write"}
	metadata := map[string]string{"department": "engineering"}

	tokenPair, err := service.GenerateTokenPair(userID, email, role, permissions, metadata)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if tokenPair.AccessToken == "" {
		t.Error("Expected access token to be generated")
	}

	if tokenPair.RefreshToken == "" {
		t.Error("Expected refresh token to be generated")
	}

	if tokenPair.TokenType != "Bearer" {
		t.Errorf("Expected token type 'Bearer', got %s", tokenPair.TokenType)
	}

	if tokenPair.ExpiresIn <= 0 {
		t.Error("Expected expires in to be positive")
	}
}

func TestJWTService_ValidateToken(t *testing.T) {
	config := Config{
		SecretKey:     "test-secret-key",
		Issuer:        "test-issuer",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	service, err := NewJWTService(config)
	if err != nil {
		t.Fatalf("Failed to create JWT service: %v", err)
	}

	userID := "user-123"
	email := "user@example.com"
	role := "user"
	permissions := []string{"read", "write"}
	metadata := map[string]string{"department": "engineering"}

	tokenPair, err := service.GenerateTokenPair(userID, email, role, permissions, metadata)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Test valid token
	validatedClaims, err := service.ValidateToken(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if validatedClaims.UserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, validatedClaims.UserID)
	}

	if validatedClaims.Email != email {
		t.Errorf("Expected email %s, got %s", email, validatedClaims.Email)
	}

	if validatedClaims.Role != role {
		t.Errorf("Expected role %s, got %s", role, validatedClaims.Role)
	}

	// Test invalid token
	_, err = service.ValidateToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestJWTService_RefreshToken(t *testing.T) {
	config := Config{
		SecretKey:     "test-secret-key",
		Issuer:        "test-issuer",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	service, err := NewJWTService(config)
	if err != nil {
		t.Fatalf("Failed to create JWT service: %v", err)
	}

	userID := "user-123"
	email := "user@example.com"
	role := "user"
	permissions := []string{"read", "write"}
	metadata := map[string]string{"department": "engineering"}

	tokenPair, err := service.GenerateTokenPair(userID, email, role, permissions, metadata)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Test valid refresh token
	newAccessToken, err := service.RefreshToken(tokenPair.RefreshToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if newAccessToken == "" {
		t.Error("Expected new access token to be generated")
	}

	// Test invalid refresh token
	_, err = service.RefreshToken("invalid-refresh-token")
	if err == nil {
		t.Error("Expected error for invalid refresh token, got nil")
	}
}

func TestJWTService_ExtractUserID(t *testing.T) {
	config := Config{
		SecretKey:     "test-secret-key",
		Issuer:        "test-issuer",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	service, err := NewJWTService(config)
	if err != nil {
		t.Fatalf("Failed to create JWT service: %v", err)
	}

	userID := "user-123"
	email := "user@example.com"
	role := "user"
	permissions := []string{"read", "write"}
	metadata := map[string]string{"department": "engineering"}

	tokenPair, err := service.GenerateTokenPair(userID, email, role, permissions, metadata)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Test extracting user ID from valid token
	extractedUserID, err := service.ExtractUserID(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if extractedUserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, extractedUserID)
	}

	// Test extracting user ID from invalid token
	_, err = service.ExtractUserID("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestJWTService_HasPermission(t *testing.T) {
	config := Config{
		SecretKey:     "test-secret-key",
		Issuer:        "test-issuer",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	service, err := NewJWTService(config)
	if err != nil {
		t.Fatalf("Failed to create JWT service: %v", err)
	}

	userID := "user-123"
	email := "user@example.com"
	role := "user"
	permissions := []string{"read", "write"}
	metadata := map[string]string{"department": "engineering"}

	tokenPair, err := service.GenerateTokenPair(userID, email, role, permissions, metadata)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	claims, err := service.ValidateToken(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Test has permission
	if !service.HasPermission(claims, "read") {
		t.Error("Expected user to have read permission")
	}

	if !service.HasPermission(claims, "write") {
		t.Error("Expected user to have write permission")
	}

	if service.HasPermission(claims, "admin") {
		t.Error("Expected user to not have admin permission")
	}
}

func TestJWTService_HasRole(t *testing.T) {
	config := Config{
		SecretKey:     "test-secret-key",
		Issuer:        "test-issuer",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	service, err := NewJWTService(config)
	if err != nil {
		t.Fatalf("Failed to create JWT service: %v", err)
	}

	userID := "user-123"
	email := "user@example.com"
	role := "user"
	permissions := []string{"read", "write"}
	metadata := map[string]string{"department": "engineering"}

	tokenPair, err := service.GenerateTokenPair(userID, email, role, permissions, metadata)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	claims, err := service.ValidateToken(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Test has role
	if !service.HasRole(claims, "user") {
		t.Error("Expected user to have user role")
	}

	if service.HasRole(claims, "admin") {
		t.Error("Expected user to not have admin role")
	}
}

func TestJWTService_IsExpired(t *testing.T) {
	config := Config{
		SecretKey:     "test-secret-key",
		Issuer:        "test-issuer",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	service, err := NewJWTService(config)
	if err != nil {
		t.Fatalf("Failed to create JWT service: %v", err)
	}

	userID := "user-123"
	email := "user@example.com"
	role := "user"
	permissions := []string{"read", "write"}
	metadata := map[string]string{"department": "engineering"}

	tokenPair, err := service.GenerateTokenPair(userID, email, role, permissions, metadata)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Test non-expired token
	if service.IsExpired(tokenPair.AccessToken) {
		t.Error("Expected token to not be expired")
	}

	// Test invalid token
	if !service.IsExpired("invalid-token") {
		t.Error("Expected invalid token to be considered expired")
	}
}

func TestJWTService_GetTokenInfo(t *testing.T) {
	config := Config{
		SecretKey:     "test-secret-key",
		Issuer:        "test-issuer",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	service, err := NewJWTService(config)
	if err != nil {
		t.Fatalf("Failed to create JWT service: %v", err)
	}

	userID := "user-123"
	email := "user@example.com"
	role := "user"
	permissions := []string{"read", "write"}
	metadata := map[string]string{"department": "engineering"}

	tokenPair, err := service.GenerateTokenPair(userID, email, role, permissions, metadata)
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	// Test getting token info
	tokenInfo, err := service.GetTokenInfo(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if tokenInfo == nil {
		t.Error("Expected token info to be returned")
	}

	// Test getting info from invalid token
	_, err = service.GetTokenInfo("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

// Benchmark tests
func BenchmarkJWTService_GenerateTokenPair(b *testing.B) {
	config := Config{
		SecretKey:     "test-secret-key",
		Issuer:        "test-issuer",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	service, _ := NewJWTService(config)
	userID := "user-123"
	email := "user@example.com"
	role := "user"
	permissions := []string{"read", "write"}
	metadata := map[string]string{"department": "engineering"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GenerateTokenPair(userID, email, role, permissions, metadata)
	}
}

func BenchmarkJWTService_ValidateToken(b *testing.B) {
	config := Config{
		SecretKey:     "test-secret-key",
		Issuer:        "test-issuer",
		AccessExpiry:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
	}

	service, _ := NewJWTService(config)
	userID := "user-123"
	email := "user@example.com"
	role := "user"
	permissions := []string{"read", "write"}
	metadata := map[string]string{"department": "engineering"}

	tokenPair, _ := service.GenerateTokenPair(userID, email, role, permissions, metadata)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ValidateToken(tokenPair.AccessToken)
	}
}
