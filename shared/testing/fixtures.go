// Package testing provides testing utilities for USC platform services.
package testing

import (
	"time"
)

// TestUser represents a test user
type TestUser struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Role        string    `json:"role"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsActive    bool      `json:"is_active"`
}

// TestContent represents test content
type TestContent struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	AuthorID    string    `json:"author_id"`
	Category    string    `json:"category"`
	Tags        []string  `json:"tags"`
	Status      string    `json:"status"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TestOrder represents a test order
type TestOrder struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Items     []TestOrderItem `json:"items"`
	Total     float64         `json:"total"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// TestOrderItem represents a test order item
type TestOrderItem struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// TestProduct represents a test product
type TestProduct struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	Stock       int       `json:"stock"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TestSession represents a test session
type TestSession struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TestConfig represents test configuration
type TestConfig struct {
	DatabaseURL string `json:"database_url"`
	RedisURL    string `json:"redis_url"`
	JWTSecret   string `json:"jwt_secret"`
	Port        int    `json:"port"`
	Environment string `json:"environment"`
}

// TestFixtures provides test fixtures
type TestFixtures struct {
	Users    []TestUser
	Content  []TestContent
	Orders   []TestOrder
	Products []TestProduct
	Sessions []TestSession
	Config   TestConfig
}

// NewTestFixtures creates new test fixtures
func NewTestFixtures() *TestFixtures {
	now := time.Now()

	return &TestFixtures{
		Users: []TestUser{
			{
				ID:          "user-1",
				Email:       "admin@example.com",
				Username:    "admin",
				Password:    "password123",
				Role:        "admin",
				Permissions: []string{"user:read", "user:write", "user:delete", "content:read", "content:write", "content:delete"},
				CreatedAt:   now,
				UpdatedAt:   now,
				IsActive:    true,
			},
			{
				ID:          "user-2",
				Email:       "user@example.com",
				Username:    "user",
				Password:    "password123",
				Role:        "user",
				Permissions: []string{"user:read", "content:read", "content:write"},
				CreatedAt:   now,
				UpdatedAt:   now,
				IsActive:    true,
			},
			{
				ID:          "user-3",
				Email:       "moderator@example.com",
				Username:    "moderator",
				Password:    "password123",
				Role:        "moderator",
				Permissions: []string{"user:read", "content:read", "content:write", "content:delete"},
				CreatedAt:   now,
				UpdatedAt:   now,
				IsActive:    true,
			},
		},
		Content: []TestContent{
			{
				ID:          "content-1",
				Title:       "Test Article 1",
				Content:     "This is a test article content.",
				AuthorID:    "user-1",
				Category:    "technology",
				Tags:        []string{"test", "article", "technology"},
				Status:      "published",
				PublishedAt: now,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "content-2",
				Title:       "Test Article 2",
				Content:     "This is another test article content.",
				AuthorID:    "user-2",
				Category:    "business",
				Tags:        []string{"test", "article", "business"},
				Status:      "draft",
				PublishedAt: time.Time{},
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		Orders: []TestOrder{
			{
				ID:     "order-1",
				UserID: "user-2",
				Items: []TestOrderItem{
					{
						ID:        "item-1",
						ProductID: "product-1",
						Quantity:  2,
						Price:     10.99,
					},
					{
						ID:        "item-2",
						ProductID: "product-2",
						Quantity:  1,
						Price:     15.99,
					},
				},
				Total:     37.97,
				Status:    "completed",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Products: []TestProduct{
			{
				ID:          "product-1",
				Name:        "Test Product 1",
				Description: "This is a test product",
				Price:       10.99,
				Category:    "electronics",
				Stock:       100,
				IsActive:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "product-2",
				Name:        "Test Product 2",
				Description: "This is another test product",
				Price:       15.99,
				Category:    "books",
				Stock:       50,
				IsActive:    true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		},
		Sessions: []TestSession{
			{
				ID:        "session-1",
				UserID:    "user-1",
				Token:     "test-token-1",
				ExpiresAt: now.Add(24 * time.Hour),
				CreatedAt: now,
			},
			{
				ID:        "session-2",
				UserID:    "user-2",
				Token:     "test-token-2",
				ExpiresAt: now.Add(24 * time.Hour),
				CreatedAt: now,
			},
		},
		Config: TestConfig{
			DatabaseURL: "postgres://test:test@localhost:5432/test_db",
			RedisURL:    "redis://localhost:6379/0",
			JWTSecret:   "test-jwt-secret",
			Port:        8080,
			Environment: "test",
		},
	}
}

// GetUserByID returns a user by ID
func (tf *TestFixtures) GetUserByID(id string) *TestUser {
	for _, user := range tf.Users {
		if user.ID == id {
			return &user
		}
	}
	return nil
}

// GetUserByEmail returns a user by email
func (tf *TestFixtures) GetUserByEmail(email string) *TestUser {
	for _, user := range tf.Users {
		if user.Email == email {
			return &user
		}
	}
	return nil
}

// GetUserByUsername returns a user by username
func (tf *TestFixtures) GetUserByUsername(username string) *TestUser {
	for _, user := range tf.Users {
		if user.Username == username {
			return &user
		}
	}
	return nil
}

// GetContentByID returns content by ID
func (tf *TestFixtures) GetContentByID(id string) *TestContent {
	for _, content := range tf.Content {
		if content.ID == id {
			return &content
		}
	}
	return nil
}

// GetContentByAuthor returns content by author ID
func (tf *TestFixtures) GetContentByAuthor(authorID string) []TestContent {
	var result []TestContent
	for _, content := range tf.Content {
		if content.AuthorID == authorID {
			result = append(result, content)
		}
	}
	return result
}

// GetOrderByID returns an order by ID
func (tf *TestFixtures) GetOrderByID(id string) *TestOrder {
	for _, order := range tf.Orders {
		if order.ID == id {
			return &order
		}
	}
	return nil
}

// GetOrdersByUser returns orders by user ID
func (tf *TestFixtures) GetOrdersByUser(userID string) []TestOrder {
	var result []TestOrder
	for _, order := range tf.Orders {
		if order.UserID == userID {
			result = append(result, order)
		}
	}
	return result
}

// GetProductByID returns a product by ID
func (tf *TestFixtures) GetProductByID(id string) *TestProduct {
	for _, product := range tf.Products {
		if product.ID == id {
			return &product
		}
	}
	return nil
}

// GetProductsByCategory returns products by category
func (tf *TestFixtures) GetProductsByCategory(category string) []TestProduct {
	var result []TestProduct
	for _, product := range tf.Products {
		if product.Category == category {
			result = append(result, product)
		}
	}
	return result
}

// GetSessionByID returns a session by ID
func (tf *TestFixtures) GetSessionByID(id string) *TestSession {
	for _, session := range tf.Sessions {
		if session.ID == id {
			return &session
		}
	}
	return nil
}

// GetSessionByToken returns a session by token
func (tf *TestFixtures) GetSessionByToken(token string) *TestSession {
	for _, session := range tf.Sessions {
		if session.Token == token {
			return &session
		}
	}
	return nil
}

// GetActiveSessions returns all active sessions
func (tf *TestFixtures) GetActiveSessions() []TestSession {
	var result []TestSession
	now := time.Now()
	for _, session := range tf.Sessions {
		if session.ExpiresAt.After(now) {
			result = append(result, session)
		}
	}
	return result
}

// AddUser adds a new user to fixtures
func (tf *TestFixtures) AddUser(user TestUser) {
	tf.Users = append(tf.Users, user)
}

// AddContent adds new content to fixtures
func (tf *TestFixtures) AddContent(content TestContent) {
	tf.Content = append(tf.Content, content)
}

// AddOrder adds a new order to fixtures
func (tf *TestFixtures) AddOrder(order TestOrder) {
	tf.Orders = append(tf.Orders, order)
}

// AddProduct adds a new product to fixtures
func (tf *TestFixtures) AddProduct(product TestProduct) {
	tf.Products = append(tf.Products, product)
}

// AddSession adds a new session to fixtures
func (tf *TestFixtures) AddSession(session TestSession) {
	tf.Sessions = append(tf.Sessions, session)
}

// Clear clears all fixtures
func (tf *TestFixtures) Clear() {
	tf.Users = []TestUser{}
	tf.Content = []TestContent{}
	tf.Orders = []TestOrder{}
	tf.Products = []TestProduct{}
	tf.Sessions = []TestSession{}
}

// Reset resets fixtures to default values
func (tf *TestFixtures) Reset() {
	tf.Clear()
	*tf = *NewTestFixtures()
}

// TestDataGenerator provides test data generation utilities
type TestDataGenerator struct {
	fixtures *TestFixtures
}

// NewTestDataGenerator creates a new test data generator
func NewTestDataGenerator() *TestDataGenerator {
	return &TestDataGenerator{
		fixtures: NewTestFixtures(),
	}
}

// GenerateUser generates a test user
func (tdg *TestDataGenerator) GenerateUser(overrides ...func(*TestUser)) TestUser {
	user := TestUser{
		ID:          "generated-user-" + time.Now().Format("20060102150405"),
		Email:       "test@example.com",
		Username:    "testuser",
		Password:    "password123",
		Role:        "user",
		Permissions: []string{"user:read", "content:read"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		IsActive:    true,
	}

	// Apply overrides
	for _, override := range overrides {
		override(&user)
	}

	return user
}

// GenerateContent generates test content
func (tdg *TestDataGenerator) GenerateContent(overrides ...func(*TestContent)) TestContent {
	content := TestContent{
		ID:          "generated-content-" + time.Now().Format("20060102150405"),
		Title:       "Generated Test Article",
		Content:     "This is generated test content.",
		AuthorID:    "user-1",
		Category:    "technology",
		Tags:        []string{"test", "generated"},
		Status:      "draft",
		PublishedAt: time.Time{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Apply overrides
	for _, override := range overrides {
		override(&content)
	}

	return content
}

// GenerateOrder generates a test order
func (tdg *TestDataGenerator) GenerateOrder(overrides ...func(*TestOrder)) TestOrder {
	order := TestOrder{
		ID:     "generated-order-" + time.Now().Format("20060102150405"),
		UserID: "user-1",
		Items: []TestOrderItem{
			{
				ID:        "generated-item-1",
				ProductID: "product-1",
				Quantity:  1,
				Price:     10.99,
			},
		},
		Total:     10.99,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Apply overrides
	for _, override := range overrides {
		override(&order)
	}

	return order
}

// GenerateProduct generates a test product
func (tdg *TestDataGenerator) GenerateProduct(overrides ...func(*TestProduct)) TestProduct {
	product := TestProduct{
		ID:          "generated-product-" + time.Now().Format("20060102150405"),
		Name:        "Generated Test Product",
		Description: "This is a generated test product",
		Price:       9.99,
		Category:    "electronics",
		Stock:       10,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Apply overrides
	for _, override := range overrides {
		override(&product)
	}

	return product
}

// GenerateSession generates a test session
func (tdg *TestDataGenerator) GenerateSession(overrides ...func(*TestSession)) TestSession {
	session := TestSession{
		ID:        "generated-session-" + time.Now().Format("20060102150405"),
		UserID:    "user-1",
		Token:     "generated-token-" + time.Now().Format("20060102150405"),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	// Apply overrides
	for _, override := range overrides {
		override(&session)
	}

	return session
}
