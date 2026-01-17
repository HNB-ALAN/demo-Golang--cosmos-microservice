package testing

import (
	"testing"
	"time"
)

func TestNewTestFixtures(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}
	if len(fixtures.Users) == 0 {
		t.Error("Expected users to be populated")
	}
	if len(fixtures.Content) == 0 {
		t.Error("Expected content to be populated")
	}
	if len(fixtures.Orders) == 0 {
		t.Error("Expected orders to be populated")
	}
	if len(fixtures.Products) == 0 {
		t.Error("Expected products to be populated")
	}
	if len(fixtures.Sessions) == 0 {
		t.Error("Expected sessions to be populated")
	}
}

func TestTestFixtures_GetUserByID(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Test existing user
	user := fixtures.GetUserByID("user-1")
	if user == nil {
		t.Error("Expected user, got nil")
		return
	}
	if user.ID != "user-1" {
		t.Errorf("Expected user ID 'user-1', got %s", user.ID)
	}

	// Test non-existent user
	user = fixtures.GetUserByID("nonexistent")
	if user != nil {
		t.Error("Expected nil for non-existent user")
	}
}

func TestTestFixtures_GetUserByEmail(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Test existing user
	user := fixtures.GetUserByEmail("admin@example.com")
	if user == nil {
		t.Error("Expected user, got nil")
		return
	}
	if user.Email != "admin@example.com" {
		t.Errorf("Expected email 'admin@example.com', got %s", user.Email)
	}

	// Test non-existent user
	user = fixtures.GetUserByEmail("nonexistent@example.com")
	if user != nil {
		t.Error("Expected nil for non-existent user")
	}
}

func TestTestFixtures_GetUserByUsername(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Test existing user
	user := fixtures.GetUserByUsername("admin")
	if user == nil {
		t.Error("Expected user, got nil")
		return
	}
	if user.Username != "admin" {
		t.Errorf("Expected username 'admin', got %s", user.Username)
	}

	// Test non-existent user
	user = fixtures.GetUserByUsername("nonexistent")
	if user != nil {
		t.Error("Expected nil for non-existent user")
	}
}

func TestTestFixtures_GetContentByID(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Test existing content
	content := fixtures.GetContentByID("content-1")
	if content == nil {
		t.Error("Expected content, got nil")
		return
	}
	if content.ID != "content-1" {
		t.Errorf("Expected content ID 'content-1', got %s", content.ID)
	}

	// Test non-existent content
	content = fixtures.GetContentByID("nonexistent")
	if content != nil {
		t.Error("Expected nil for non-existent content")
	}
}

func TestTestFixtures_GetContentByAuthor(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Test existing author
	content := fixtures.GetContentByAuthor("user-1")
	if len(content) == 0 {
		t.Error("Expected content for author")
	}
	if content[0].AuthorID != "user-1" {
		t.Errorf("Expected author ID 'user-1', got %s", content[0].AuthorID)
	}

	// Test non-existent author
	content = fixtures.GetContentByAuthor("nonexistent")
	if len(content) != 0 {
		t.Error("Expected no content for non-existent author")
	}
}

func TestTestFixtures_GetOrderByID(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Test existing order
	order := fixtures.GetOrderByID("order-1")
	if order == nil {
		t.Error("Expected order, got nil")
		return
	}
	if order.ID != "order-1" {
		t.Errorf("Expected order ID 'order-1', got %s", order.ID)
	}

	// Test non-existent order
	order = fixtures.GetOrderByID("nonexistent")
	if order != nil {
		t.Error("Expected nil for non-existent order")
	}
}

func TestTestFixtures_GetOrdersByUser(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Test existing user
	orders := fixtures.GetOrdersByUser("user-2")
	if len(orders) == 0 {
		t.Error("Expected orders for user")
	}
	if orders[0].UserID != "user-2" {
		t.Errorf("Expected user ID 'user-2', got %s", orders[0].UserID)
	}

	// Test non-existent user
	orders = fixtures.GetOrdersByUser("nonexistent")
	if len(orders) != 0 {
		t.Error("Expected no orders for non-existent user")
	}
}

func TestTestFixtures_GetProductByID(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Test existing product
	product := fixtures.GetProductByID("product-1")
	if product == nil {
		t.Error("Expected product, got nil")
		return
	}
	if product.ID != "product-1" {
		t.Errorf("Expected product ID 'product-1', got %s", product.ID)
	}

	// Test non-existent product
	product = fixtures.GetProductByID("nonexistent")
	if product != nil {
		t.Error("Expected nil for non-existent product")
	}
}

func TestTestFixtures_GetProductsByCategory(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Test existing category
	products := fixtures.GetProductsByCategory("electronics")
	if len(products) == 0 {
		t.Error("Expected products for category")
	}
	if products[0].Category != "electronics" {
		t.Errorf("Expected category 'electronics', got %s", products[0].Category)
	}

	// Test non-existent category
	products = fixtures.GetProductsByCategory("nonexistent")
	if len(products) != 0 {
		t.Error("Expected no products for non-existent category")
	}
}

func TestTestFixtures_GetSessionByID(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Test existing session
	session := fixtures.GetSessionByID("session-1")
	if session == nil {
		t.Error("Expected session, got nil")
		return
	}
	if session.ID != "session-1" {
		t.Errorf("Expected session ID 'session-1', got %s", session.ID)
	}

	// Test non-existent session
	session = fixtures.GetSessionByID("nonexistent")
	if session != nil {
		t.Error("Expected nil for non-existent session")
	}
}

func TestTestFixtures_GetSessionByToken(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Test existing token
	session := fixtures.GetSessionByToken("test-token-1")
	if session == nil {
		t.Error("Expected session, got nil")
		return
	}
	if session.Token != "test-token-1" {
		t.Errorf("Expected token 'test-token-1', got %s", session.Token)
	}

	// Test non-existent token
	session = fixtures.GetSessionByToken("nonexistent")
	if session != nil {
		t.Error("Expected nil for non-existent token")
	}
}

func TestTestFixtures_GetActiveSessions(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Get active sessions
	sessions := fixtures.GetActiveSessions()
	if len(sessions) == 0 {
		t.Error("Expected active sessions")
	}

	// All sessions should be active (not expired)
	now := time.Now()
	for _, session := range sessions {
		if session.ExpiresAt.Before(now) {
			t.Error("Expected session to be active")
		}
	}
}

func TestTestFixtures_AddUser(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}
	initialCount := len(fixtures.Users)

	// Add new user
	newUser := TestUser{
		ID:       "new-user",
		Email:    "new@example.com",
		Username: "newuser",
		Password: "password",
		Role:     "user",
	}
	fixtures.AddUser(newUser)

	// Verify addition
	if len(fixtures.Users) != initialCount+1 {
		t.Errorf("Expected %d users, got %d", initialCount+1, len(fixtures.Users))
	}

	// Verify user was added
	user := fixtures.GetUserByID("new-user")
	if user == nil {
		t.Error("Expected new user to be found")
		return
	}
}

func TestTestFixtures_AddContent(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}
	initialCount := len(fixtures.Content)

	// Add new content
	newContent := TestContent{
		ID:       "new-content",
		Title:    "New Article",
		Content:  "New content",
		AuthorID: "user-1",
		Category: "test",
		Status:   "draft",
	}
	fixtures.AddContent(newContent)

	// Verify addition
	if len(fixtures.Content) != initialCount+1 {
		t.Errorf("Expected %d content items, got %d", initialCount+1, len(fixtures.Content))
	}

	// Verify content was added
	content := fixtures.GetContentByID("new-content")
	if content == nil {
		t.Error("Expected new content to be found")
		return
	}
}

func TestTestFixtures_AddOrder(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}
	initialCount := len(fixtures.Orders)

	// Add new order
	newOrder := TestOrder{
		ID:     "new-order",
		UserID: "user-1",
		Total:  99.99,
		Status: "pending",
	}
	fixtures.AddOrder(newOrder)

	// Verify addition
	if len(fixtures.Orders) != initialCount+1 {
		t.Errorf("Expected %d orders, got %d", initialCount+1, len(fixtures.Orders))
	}

	// Verify order was added
	order := fixtures.GetOrderByID("new-order")
	if order == nil {
		t.Error("Expected new order to be found")
		return
	}
}

func TestTestFixtures_AddProduct(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}
	initialCount := len(fixtures.Products)

	// Add new product
	newProduct := TestProduct{
		ID:          "new-product",
		Name:        "New Product",
		Description: "New product description",
		Price:       19.99,
		Category:    "test",
		Stock:       10,
		IsActive:    true,
	}
	fixtures.AddProduct(newProduct)

	// Verify addition
	if len(fixtures.Products) != initialCount+1 {
		t.Errorf("Expected %d products, got %d", initialCount+1, len(fixtures.Products))
	}

	// Verify product was added
	product := fixtures.GetProductByID("new-product")
	if product == nil {
		t.Error("Expected new product to be found")
		return
	}
}

func TestTestFixtures_AddSession(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}
	initialCount := len(fixtures.Sessions)

	// Add new session
	newSession := TestSession{
		ID:        "new-session",
		UserID:    "user-1",
		Token:     "new-token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	fixtures.AddSession(newSession)

	// Verify addition
	if len(fixtures.Sessions) != initialCount+1 {
		t.Errorf("Expected %d sessions, got %d", initialCount+1, len(fixtures.Sessions))
	}

	// Verify session was added
	session := fixtures.GetSessionByID("new-session")
	if session == nil {
		t.Error("Expected new session to be found")
		return
	}
}

func TestTestFixtures_Clear(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Verify initial data exists
	if len(fixtures.Users) == 0 {
		t.Error("Expected initial users to exist")
	}

	// Clear fixtures
	fixtures.Clear()

	// Verify clearing
	if len(fixtures.Users) != 0 {
		t.Error("Expected users to be cleared")
	}
	if len(fixtures.Content) != 0 {
		t.Error("Expected content to be cleared")
	}
	if len(fixtures.Orders) != 0 {
		t.Error("Expected orders to be cleared")
	}
	if len(fixtures.Products) != 0 {
		t.Error("Expected products to be cleared")
	}
	if len(fixtures.Sessions) != 0 {
		t.Error("Expected sessions to be cleared")
	}
}

func TestTestFixtures_Reset(t *testing.T) {
	fixtures := NewTestFixtures()
	if fixtures == nil {
		t.Error("Expected fixtures, got nil")
		return
	}

	// Clear fixtures first
	fixtures.Clear()

	// Verify clearing
	if len(fixtures.Users) != 0 {
		t.Error("Expected users to be cleared")
	}

	// Reset fixtures
	fixtures.Reset()

	// Verify reset
	if len(fixtures.Users) == 0 {
		t.Error("Expected users to be reset")
	}
	if len(fixtures.Content) == 0 {
		t.Error("Expected content to be reset")
	}
	if len(fixtures.Orders) == 0 {
		t.Error("Expected orders to be reset")
	}
	if len(fixtures.Products) == 0 {
		t.Error("Expected products to be reset")
	}
	if len(fixtures.Sessions) == 0 {
		t.Error("Expected sessions to be reset")
	}
}

func TestNewTestDataGenerator(t *testing.T) {
	generator := NewTestDataGenerator()
	if generator == nil {
		t.Error("Expected generator, got nil")
		return
	}
	if generator.fixtures == nil {
		t.Error("Expected fixtures to be initialized")
	}
}

func TestTestDataGenerator_GenerateUser(t *testing.T) {
	generator := NewTestDataGenerator()

	// Generate user without overrides
	user := generator.GenerateUser()
	if user.ID == "" {
		t.Error("Expected user ID to be set")
	}
	if user.Email == "" {
		t.Error("Expected user email to be set")
	}
	if user.Username == "" {
		t.Error("Expected user username to be set")
	}
	if user.Role == "" {
		t.Error("Expected user role to be set")
	}

	// Generate user with overrides
	user = generator.GenerateUser(func(u *TestUser) {
		u.Email = "override@example.com"
		u.Role = "admin"
	})
	if user.Email != "override@example.com" {
		t.Errorf("Expected email 'override@example.com', got %s", user.Email)
	}
	if user.Role != "admin" {
		t.Errorf("Expected role 'admin', got %s", user.Role)
	}
}

func TestTestDataGenerator_GenerateContent(t *testing.T) {
	generator := NewTestDataGenerator()

	// Generate content without overrides
	content := generator.GenerateContent()
	if content.ID == "" {
		t.Error("Expected content ID to be set")
	}
	if content.Title == "" {
		t.Error("Expected content title to be set")
	}
	if content.Content == "" {
		t.Error("Expected content content to be set")
	}
	if content.AuthorID == "" {
		t.Error("Expected content author ID to be set")
	}

	// Generate content with overrides
	content = generator.GenerateContent(func(c *TestContent) {
		c.Title = "Override Title"
		c.Status = "published"
	})
	if content.Title != "Override Title" {
		t.Errorf("Expected title 'Override Title', got %s", content.Title)
	}
	if content.Status != "published" {
		t.Errorf("Expected status 'published', got %s", content.Status)
	}
}

func TestTestDataGenerator_GenerateOrder(t *testing.T) {
	generator := NewTestDataGenerator()

	// Generate order without overrides
	order := generator.GenerateOrder()
	if order.ID == "" {
		t.Error("Expected order ID to be set")
	}
	if order.UserID == "" {
		t.Error("Expected order user ID to be set")
	}
	if len(order.Items) == 0 {
		t.Error("Expected order items to be set")
	}
	if order.Total == 0 {
		t.Error("Expected order total to be set")
	}

	// Generate order with overrides
	order = generator.GenerateOrder(func(o *TestOrder) {
		o.Status = "completed"
		o.Total = 99.99
	})
	if order.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", order.Status)
	}
	if order.Total != 99.99 {
		t.Errorf("Expected total 99.99, got %f", order.Total)
	}
}

func TestTestDataGenerator_GenerateProduct(t *testing.T) {
	generator := NewTestDataGenerator()

	// Generate product without overrides
	product := generator.GenerateProduct()
	if product.ID == "" {
		t.Error("Expected product ID to be set")
	}
	if product.Name == "" {
		t.Error("Expected product name to be set")
	}
	if product.Description == "" {
		t.Error("Expected product description to be set")
	}
	if product.Price == 0 {
		t.Error("Expected product price to be set")
	}

	// Generate product with overrides
	product = generator.GenerateProduct(func(p *TestProduct) {
		p.Name = "Override Product"
		p.Price = 29.99
	})
	if product.Name != "Override Product" {
		t.Errorf("Expected name 'Override Product', got %s", product.Name)
	}
	if product.Price != 29.99 {
		t.Errorf("Expected price 29.99, got %f", product.Price)
	}
}

func TestTestDataGenerator_GenerateSession(t *testing.T) {
	generator := NewTestDataGenerator()

	// Generate session without overrides
	session := generator.GenerateSession()
	if session.ID == "" {
		t.Error("Expected session ID to be set")
	}
	if session.UserID == "" {
		t.Error("Expected session user ID to be set")
	}
	if session.Token == "" {
		t.Error("Expected session token to be set")
	}
	if session.ExpiresAt.IsZero() {
		t.Error("Expected session expires at to be set")
	}

	// Generate session with overrides
	session = generator.GenerateSession(func(s *TestSession) {
		s.UserID = "override-user"
		s.Token = "override-token"
	})
	if session.UserID != "override-user" {
		t.Errorf("Expected user ID 'override-user', got %s", session.UserID)
	}
	if session.Token != "override-token" {
		t.Errorf("Expected token 'override-token', got %s", session.Token)
	}
}

func TestTestUser_Fields(t *testing.T) {
	now := time.Now()
	user := TestUser{
		ID:          "test-id",
		Email:       "test@example.com",
		Username:    "testuser",
		Password:    "password123",
		Role:        "user",
		Permissions: []string{"read", "write"},
		CreatedAt:   now,
		UpdatedAt:   now,
		IsActive:    true,
	}

	// Test all fields
	if user.ID != "test-id" {
		t.Errorf("Expected ID 'test-id', got %s", user.ID)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got %s", user.Email)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %s", user.Username)
	}
	if user.Password != "password123" {
		t.Errorf("Expected password 'password123', got %s", user.Password)
	}
	if user.Role != "user" {
		t.Errorf("Expected role 'user', got %s", user.Role)
	}
	if len(user.Permissions) != 2 {
		t.Errorf("Expected 2 permissions, got %d", len(user.Permissions))
	}
	if !user.CreatedAt.Equal(now) {
		t.Error("Expected created at to match")
	}
	if !user.UpdatedAt.Equal(now) {
		t.Error("Expected updated at to match")
	}
	if !user.IsActive {
		t.Error("Expected user to be active")
	}
}

func TestTestContent_Fields(t *testing.T) {
	now := time.Now()
	content := TestContent{
		ID:          "content-id",
		Title:       "Test Title",
		Content:     "Test content",
		AuthorID:    "author-id",
		Category:    "test",
		Tags:        []string{"tag1", "tag2"},
		Status:      "published",
		PublishedAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Test all fields
	if content.ID != "content-id" {
		t.Errorf("Expected ID 'content-id', got %s", content.ID)
	}
	if content.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got %s", content.Title)
	}
	if content.Content != "Test content" {
		t.Errorf("Expected content 'Test content', got %s", content.Content)
	}
	if content.AuthorID != "author-id" {
		t.Errorf("Expected author ID 'author-id', got %s", content.AuthorID)
	}
	if content.Category != "test" {
		t.Errorf("Expected category 'test', got %s", content.Category)
	}
	if len(content.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(content.Tags))
	}
	if content.Status != "published" {
		t.Errorf("Expected status 'published', got %s", content.Status)
	}
	if !content.PublishedAt.Equal(now) {
		t.Error("Expected published at to match")
	}
	if !content.CreatedAt.Equal(now) {
		t.Error("Expected created at to match")
	}
	if !content.UpdatedAt.Equal(now) {
		t.Error("Expected updated at to match")
	}
}

func TestTestOrder_Fields(t *testing.T) {
	now := time.Now()
	order := TestOrder{
		ID:     "order-id",
		UserID: "user-id",
		Items: []TestOrderItem{
			{ID: "item1", ProductID: "product1", Quantity: 1, Price: 10.99},
			{ID: "item2", ProductID: "product2", Quantity: 2, Price: 5.99},
		},
		Total:     22.97,
		Status:    "completed",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test all fields
	if order.ID != "order-id" {
		t.Errorf("Expected ID 'order-id', got %s", order.ID)
	}
	if order.UserID != "user-id" {
		t.Errorf("Expected user ID 'user-id', got %s", order.UserID)
	}
	if len(order.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(order.Items))
	}
	if order.Total != 22.97 {
		t.Errorf("Expected total 22.97, got %f", order.Total)
	}
	if order.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", order.Status)
	}
	if !order.CreatedAt.Equal(now) {
		t.Error("Expected created at to match")
	}
	if !order.UpdatedAt.Equal(now) {
		t.Error("Expected updated at to match")
	}
}

func TestTestProduct_Fields(t *testing.T) {
	now := time.Now()
	product := TestProduct{
		ID:          "product-id",
		Name:        "Test Product",
		Description: "Test description",
		Price:       19.99,
		Category:    "electronics",
		Stock:       100,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Test all fields
	if product.ID != "product-id" {
		t.Errorf("Expected ID 'product-id', got %s", product.ID)
	}
	if product.Name != "Test Product" {
		t.Errorf("Expected name 'Test Product', got %s", product.Name)
	}
	if product.Description != "Test description" {
		t.Errorf("Expected description 'Test description', got %s", product.Description)
	}
	if product.Price != 19.99 {
		t.Errorf("Expected price 19.99, got %f", product.Price)
	}
	if product.Category != "electronics" {
		t.Errorf("Expected category 'electronics', got %s", product.Category)
	}
	if product.Stock != 100 {
		t.Errorf("Expected stock 100, got %d", product.Stock)
	}
	if !product.IsActive {
		t.Error("Expected product to be active")
	}
	if !product.CreatedAt.Equal(now) {
		t.Error("Expected created at to match")
	}
	if !product.UpdatedAt.Equal(now) {
		t.Error("Expected updated at to match")
	}
}

func TestTestSession_Fields(t *testing.T) {
	now := time.Now()
	session := TestSession{
		ID:        "session-id",
		UserID:    "user-id",
		Token:     "session-token",
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}

	// Test all fields
	if session.ID != "session-id" {
		t.Errorf("Expected ID 'session-id', got %s", session.ID)
	}
	if session.UserID != "user-id" {
		t.Errorf("Expected user ID 'user-id', got %s", session.UserID)
	}
	if session.Token != "session-token" {
		t.Errorf("Expected token 'session-token', got %s", session.Token)
	}
	if !session.ExpiresAt.After(now) {
		t.Error("Expected expires at to be in the future")
	}
	if !session.CreatedAt.Equal(now) {
		t.Error("Expected created at to match")
	}
}

func TestTestConfig_Fields(t *testing.T) {
	config := TestConfig{
		DatabaseURL: "postgres://test:test@localhost:5432/test",
		RedisURL:    "redis://localhost:6379/0",
		JWTSecret:   "test-secret",
		Port:        8080,
		Environment: "test",
	}

	// Test all fields
	if config.DatabaseURL != "postgres://test:test@localhost:5432/test" {
		t.Errorf("Expected database URL, got %s", config.DatabaseURL)
	}
	if config.RedisURL != "redis://localhost:6379/0" {
		t.Errorf("Expected Redis URL, got %s", config.RedisURL)
	}
	if config.JWTSecret != "test-secret" {
		t.Errorf("Expected JWT secret, got %s", config.JWTSecret)
	}
	if config.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", config.Port)
	}
	if config.Environment != "test" {
		t.Errorf("Expected environment 'test', got %s", config.Environment)
	}
}
