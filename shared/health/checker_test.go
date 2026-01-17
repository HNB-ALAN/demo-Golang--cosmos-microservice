package health

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewChecker(t *testing.T) {
	name := "test-checker"
	description := "Test health checker"
	checkFunc := func(ctx context.Context) error {
		return nil
	}

	checker := NewChecker(name, description, checkFunc)

	if checker.Name() != name {
		t.Errorf("Expected name %s, got %s", name, checker.Name())
	}

	if checker.Description() != description {
		t.Errorf("Expected description %s, got %s", description, checker.Description())
	}
}

func TestChecker_Check(t *testing.T) {
	t.Run("successful check", func(t *testing.T) {
		checker := NewChecker("test", "test", func(ctx context.Context) error {
			return nil
		})

		err := checker.Check(context.Background())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("failed check", func(t *testing.T) {
		expectedErr := errors.New("check failed")
		checker := NewChecker("test", "test", func(ctx context.Context) error {
			return expectedErr
		})

		err := checker.Check(context.Background())
		if err == nil {
			t.Error("Expected error, got nil")
		}

		if err.Error() != expectedErr.Error() {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})

	t.Run("nil check function", func(t *testing.T) {
		checker := &Checker{
			name:        "test",
			description: "test",
			checkFunc:   nil,
		}

		err := checker.Check(context.Background())
		if err == nil {
			t.Error("Expected error for nil check function, got nil")
		}
	})
}

func TestNewDatabaseChecker(t *testing.T) {
	name := "db-checker"
	description := "Database health checker"
	pingFunc := func(ctx context.Context) error {
		return nil
	}

	checker := NewDatabaseChecker(name, description, pingFunc)

	if checker.Name() != name {
		t.Errorf("Expected name %s, got %s", name, checker.Name())
	}

	if checker.Description() != description {
		t.Errorf("Expected description %s, got %s", description, checker.Description())
	}
}

func TestDatabaseChecker_Check(t *testing.T) {
	t.Run("successful ping", func(t *testing.T) {
		checker := NewDatabaseChecker("test", "test", func(ctx context.Context) error {
			return nil
		})

		err := checker.Check(context.Background())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("failed ping", func(t *testing.T) {
		expectedErr := errors.New("ping failed")
		checker := NewDatabaseChecker("test", "test", func(ctx context.Context) error {
			return expectedErr
		})

		err := checker.Check(context.Background())
		if err == nil {
			t.Error("Expected error, got nil")
		}

		if err.Error() != expectedErr.Error() {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
	})
}

func TestNewHTTPChecker(t *testing.T) {
	name := "http-checker"
	description := "HTTP health checker"
	url := "http://example.com/health"
	timeout := 5 * time.Second

	checker := NewHTTPChecker(name, description, url, timeout)

	if checker.Name() != name {
		t.Errorf("Expected name %s, got %s", name, checker.Name())
	}

	if checker.Description() != description {
		t.Errorf("Expected description %s, got %s", description, checker.Description())
	}
}

func TestNewCompositeChecker(t *testing.T) {
	name := "composite-checker"
	description := "Composite health checker"
	checker1 := NewChecker("check1", "check1", func(ctx context.Context) error { return nil })
	checker2 := NewChecker("check2", "check2", func(ctx context.Context) error { return nil })

	checker := NewCompositeChecker(name, description, checker1, checker2)

	if checker.Name() != name {
		t.Errorf("Expected name %s, got %s", name, checker.Name())
	}

	if checker.Description() != description {
		t.Errorf("Expected description %s, got %s", description, checker.Description())
	}
}

func TestCompositeChecker_Check(t *testing.T) {
	t.Run("all checks pass", func(t *testing.T) {
		checker1 := NewChecker("check1", "check1", func(ctx context.Context) error { return nil })
		checker2 := NewChecker("check2", "check2", func(ctx context.Context) error { return nil })
		composite := NewCompositeChecker("composite", "composite", checker1, checker2)

		err := composite.Check(context.Background())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("one check fails", func(t *testing.T) {
		checker1 := NewChecker("check1", "check1", func(ctx context.Context) error { return nil })
		checker2 := NewChecker("check2", "check2", func(ctx context.Context) error { return errors.New("failed") })
		composite := NewCompositeChecker("composite", "composite", checker1, checker2)

		err := composite.Check(context.Background())
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
