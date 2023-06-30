package main

import (
	"bytes"
	"io"
	"net/http"
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
	_, err = regexp.Compile(result.Method)
	if err != nil {
		t.Errorf("Failed to compile generated regex: %v", err)
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
	_, err = regexp.Compile(result.Headers.Key)
	if err != nil {
		t.Errorf("Failed to compile generated headers key regex: %v", err)
	}
	_, err = regexp.Compile(result.Headers.Value)
	if err != nil {
		t.Errorf("Failed to compile generated headers value regex: %v", err)
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
