package utils

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"   ", true},
		{"\t\n", true},
		{"hello", false},
		{" hello ", false},
	}

	for _, test := range tests {
		result := IsEmpty(test.input)
		if result != test.expected {
			t.Errorf("IsEmpty(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"   ", false},
		{"\t\n", false},
		{"hello", true},
		{" hello ", true},
	}

	for _, test := range tests {
		result := IsNotEmpty(test.input)
		if result != test.expected {
			t.Errorf("IsNotEmpty(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{"hello", 3, "hel"},
		{"hello", 5, "hello"},
		{"hello", 10, "hello"},
		{"", 5, ""},
	}

	for _, test := range tests {
		result := Truncate(test.input, test.length)
		if result != test.expected {
			t.Errorf("Truncate(%q, %d) = %q, expected %q", test.input, test.length, result, test.expected)
		}
	}
}

func TestTruncateWithEllipsis(t *testing.T) {
	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{"hello", 3, "hel"},
		{"hello", 4, "h..."},
		{"hello", 5, "hello"},
		{"hello", 10, "hello"},
		{"", 5, ""},
	}

	for _, test := range tests {
		result := TruncateWithEllipsis(test.input, test.length)
		if result != test.expected {
			t.Errorf("TruncateWithEllipsis(%q, %d) = %q, expected %q", test.input, test.length, result, test.expected)
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"hello world", "world", true},
		{"hello world", "hello", true},
		{"hello world", "xyz", false},
		{"", "test", false},
		{"test", "", true},
	}

	for _, test := range tests {
		result := Contains(test.s, test.substr)
		if result != test.expected {
			t.Errorf("Contains(%q, %q) = %v, expected %v", test.s, test.substr, result, test.expected)
		}
	}
}

func TestContainsAny(t *testing.T) {
	tests := []struct {
		s          string
		substrings []string
		expected   bool
	}{
		{"hello world", []string{"world", "xyz"}, true},
		{"hello world", []string{"abc", "xyz"}, false},
		{"hello world", []string{}, false},
		{"", []string{"test"}, false},
	}

	for _, test := range tests {
		result := ContainsAny(test.s, test.substrings...)
		if result != test.expected {
			t.Errorf("ContainsAny(%q, %v) = %v, expected %v", test.s, test.substrings, result, test.expected)
		}
	}
}

func TestStartsWith(t *testing.T) {
	tests := []struct {
		s        string
		prefix   string
		expected bool
	}{
		{"hello world", "hello", true},
		{"hello world", "world", false},
		{"hello", "hello", true},
		{"", "test", false},
		{"test", "", true},
	}

	for _, test := range tests {
		result := StartsWith(test.s, test.prefix)
		if result != test.expected {
			t.Errorf("StartsWith(%q, %q) = %v, expected %v", test.s, test.prefix, result, test.expected)
		}
	}
}

func TestEndsWith(t *testing.T) {
	tests := []struct {
		s        string
		suffix   string
		expected bool
	}{
		{"hello world", "world", true},
		{"hello world", "hello", false},
		{"hello", "hello", true},
		{"", "test", false},
		{"test", "", true},
	}

	for _, test := range tests {
		result := EndsWith(test.s, test.suffix)
		if result != test.expected {
			t.Errorf("EndsWith(%q, %q) = %v, expected %v", test.s, test.suffix, result, test.expected)
		}
	}
}

func TestToTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "Hello World"},
		{"HELLO WORLD", "Hello World"},
		{"hELLo WoRLd", "Hello World"},
		{"", ""},
		{"a", "A"},
	}

	for _, test := range tests {
		result := ToTitle(test.input)
		if result != test.expected {
			t.Errorf("ToTitle(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "helloWorld"},
		{"hello_world", "helloWorld"},
		{"hello-world", "hello-world"},
		{"Hello World", "helloWorld"},
		{"", ""},
		{"a", "a"},
	}

	for _, test := range tests {
		result := ToCamelCase(test.input)
		if result != test.expected {
			t.Errorf("ToCamelCase(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"helloWorld", "hello_world"},
		{"HelloWorld", "hello_world"},
		{"hello world", "hello world"},
		{"hello-world", "hello-world"},
		{"", ""},
		{"a", "a"},
	}

	for _, test := range tests {
		result := ToSnakeCase(test.input)
		if result != test.expected {
			t.Errorf("ToSnakeCase(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}
