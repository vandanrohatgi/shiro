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

	// Check if the generated regex is valid
	expectedMethod := "GET,POST"
	result.MethodRegex = *regexp.MustCompile(result.Method)
	if !result.MethodRegex.MatchString(expectedMethod) {
		t.Errorf("Expected updated rule method regex to match pattern '%s'", expectedMethod)
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

	// Check if the generated regexes are valid
	expectedHeaderKey := "User-Agent,Content-Type"
	result.Headers.KeyRegex = *regexp.MustCompile(result.Headers.Key)
	if !result.Headers.KeyRegex.MatchString(expectedHeaderKey) {
		t.Errorf("Expected updated rule header key regex to match pattern '%s'", expectedHeaderKey)
	}

	expectedHeaderValue := "Chrome,Mozilla/5.0"
	result.Headers.ValueRegex = *regexp.MustCompile(result.Headers.Value)
	if !result.Headers.ValueRegex.MatchString(expectedHeaderValue) {
		t.Errorf("Expected updated rule header value regex to match pattern '%s'", expectedHeaderValue)
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
	// Verify the generated regex
	expectedmatch := "Sample body,initial regex"
	updatedRule.BodyRegex = *regexp.MustCompile(updatedRule.Body)
	if !updatedRule.BodyRegex.MatchString(expectedmatch) {
		t.Errorf("Expected updated rule body regex to match pattern '%s'",
			expectedmatch)
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
