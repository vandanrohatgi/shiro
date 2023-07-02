package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"testing"
)

func TestInspectMethod(t *testing.T) {
	// Create a mock request
	request, err := http.NewRequest("POST", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock rule
	rule := Rules{
		Method: "GET",
	}

	// Call the InspectMethod function
	result := InspectMethod(request, rule)

	// Check if the method in the rule has been updated with the generated regex
	if result.Method == rule.Method {
		t.Errorf("Expected method to be updated, got %s", result.Method)
	}
}

func TestInspectHeaders(t *testing.T) {
	// Create a mock request
	request, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set mock headers in the request
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0")

	// Create a mock rule
	rule := Rules{
		Headers: Headers{
			Key:   "User-Agent",
			Value: "Chrome",
		},
	}

	// Call the InspectHeaders function
	result := InspectHeaders(request, rule)

	// Check if the headers in the rule have been updated with the generated regex
	if result.Headers.Key == rule.Headers.Key {
		t.Errorf("Expected headers key to be updated, got %s", result.Headers.Key)
	}
	if result.Headers.Value == rule.Headers.Value {
		t.Errorf("Expected headers value to be updated, got %s", result.Headers.Value)
	}
}

func TestInspectBody(t *testing.T) {
	// Create a sample request and rule
	requestBody := []byte("Sample body")
	request := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(requestBody)),
	}
	rule := Rules{
		Body: "initial regex",
	}

	// Call the InspectBody function
	updatedRule := InspectBody(request, rule)

	// Assert the updated rule's properties
	if updatedRule.Body == "" {
		t.Errorf("Expected updated rule body to be non-empty, got empty")
	}

	if request.Body == nil {
		t.Errorf("Request body should not be empty after inspection")
	}
}

func TestReadWriteRules(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_rules.yaml")
	if err != nil {
		t.Fatal("Failed to create temporary file:", err)
	}
	defer os.Remove(tmpFile.Name())

	// Create a sample RuleConfig
	ruleConfig := &RuleConfig{
		Path: tmpFile.Name(),
		Rules: map[string]Rules{
			"/": {
				Method: "GET",
				Body:   "",
				Headers: Headers{
					Key:   "Content-Type",
					Value: "application/json",
				},
			},
		},
	}

	// Invoke the WriteRules function
	err = ruleConfig.WriteRules()

	// Check for any errors
	if err != nil {
		t.Errorf("WriteRules returned an error: %v", err)
	}

	// Read the written file
	testRuleConfig := &RuleConfig{
		Path:  tmpFile.Name(),
		Rules: make(map[string]Rules),
	}
	testRuleConfig.IngestRules()
	// Validate the file content
	if !reflect.DeepEqual(testRuleConfig, ruleConfig) {
		t.Errorf("Expected: %v , received %v", ruleConfig, testRuleConfig)
	}
}

// For the current implementation of our regex generator this test should suffice
// For when we move to a better regex generator this test will need to be updated accordingly
func TestGenerateRegex(t *testing.T) {
	// Test cases
	testCases := []struct {
		input          []string
		expected       string
		shouldMatch    []string
		shouldNotMatch []string
	}{
		{
			input:          []string{"apple", "banana", "cherry"},
			expected:       "apple|banana|cherry",
			shouldMatch:    []string{"apple", "banana", "cherry"},
			shouldNotMatch: []string{"grape", "mango"},
		},
		{
			input:          []string{"123", "456", "789"},
			expected:       "123|456|789",
			shouldMatch:    []string{"123", "456", "789"},
			shouldNotMatch: []string{"111", "222", "333"},
		},
		{
			input:          []string{"one", "two", "three", "four", "five"},
			expected:       "one|t(?:wo|hree)|f(?:our|ive)",
			shouldMatch:    []string{"one", "two", "three"},
			shouldNotMatch: []string{"six", "seven"},
		},
	}

	for _, testCase := range testCases {
		result, err := GenerateRegex(testCase.input)
		if err != nil {
			t.Errorf("Error generating regex: %v", err)
		}

		if result != testCase.expected {
			t.Errorf("Expected %s, but got %s", testCase.expected, result)
		}

		for _, str := range testCase.shouldMatch {
			if !regexp.MustCompile(result).MatchString(str) {
				t.Errorf("Expected %s to match %s", str, result)
			}
		}

		for _, str := range testCase.shouldNotMatch {
			if regexp.MustCompile(result).MatchString(str) {
				t.Errorf("Expected %s not to match %s", str, result)
			}
		}
	}
}
