package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestBlockRequest(t *testing.T) {
	// Create a mock response writer
	response := httptest.NewRecorder()
	w := http.ResponseWriter(response)

	// Create a mock request
	request := httptest.NewRequest("GET", "/example", nil)

	// Call the blockRequest function
	blockRequest(&w, request)

	// Check the response status code
	if response.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, but got %d", http.StatusForbidden, response.Code)
	}

	// Check the response body
	expectedBody := "Forbidden"
	actualBody := strings.TrimSpace(response.Body.String())
	if actualBody != expectedBody {
		t.Errorf("Expected response body '%s', but got '%s'", expectedBody, actualBody)
	}
}

func TestCheckHeaders(t *testing.T) {
	// Create a new HTTP request for testing
	req := httptest.NewRequest("GET", "http://example.com", nil)

	// Create a test rule for checking headers
	rule := Rules{
		Headers: Headers{
			Key:        "Content-Type",
			Value:      "application/json",
			KeyRegex:   *regexp.MustCompile(`^Content-Type$`),
			ValueRegex: *regexp.MustCompile(`^application/json$`),
		},
	}

	// Set a valid header in the request
	req.Header.Set("Content-Type", "application/json")

	// Call the checkHeaders function with the test request and rule
	violation, err := checkHeaders(req, rule)

	// Check if a violation was found
	if violation {
		t.Error("Unexpected violation found for valid header")
	}

	// Check if an error occurred
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	// Set an invalid header in the request
	req.Header.Set("Content-Type", "text/plain")

	// Call the checkHeaders function with the test request and rule
	violation, err = checkHeaders(req, rule)

	// Check if a violation was found
	if !violation {
		t.Error("Expected violation not found for invalid header")
	}

	// Check if the error message is correct
	expectedErrorMessage := "request header [text/plain] violates defined header application/json"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
	}
}

func TestCheckBody(t *testing.T) {
	// Create a test request with a specific body
	reqBody := []byte("Test body content")
	req := httptest.NewRequest("POST", "http://example.com", bytes.NewBuffer(reqBody))

	// Define the rule for testing
	rule := Rules{
		BodyRegex: *regexp.MustCompile("^Test"),
		Body:      "Test body",
	}

	// Call the function being tested
	result, err := checkBody(req, rule)

	// Assert the expected behavior
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result {
		t.Error("Expected false, got true")
	}
}

func TestCheckMethod(t *testing.T) {
	// Create a test request with a specific method
	req := httptest.NewRequest("GET", "http://example.com", nil)
	// Define the rule for testing
	rule := Rules{
		MethodRegex: *regexp.MustCompile("^(GET|POST)$"),
		Method:      "GET",
	}

	// Call the function being tested
	result, err := checkMethod(req, rule)

	// Assert the expected behavior
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result {
		t.Error("Expected false, got true")
	}
}

func TestIsRequestBlocked(t *testing.T) {
	// Create a test request with a specific method, body, and headers
	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("Accept", "application/json")
	// Define the rule for testing
	rule := Rules{
		MethodRegex: *regexp.MustCompile("GET"),
		Method:      "GET",
		BodyRegex:   *regexp.MustCompile(""),
		Body:        "",
		Headers: Headers{
			Key:        "Accept",
			KeyRegex:   *regexp.MustCompile("Accept"),
			Value:      "application/json",
			ValueRegex: *regexp.MustCompile("application/json"),
		},
	}

	// Call the function being tested
	result, err := IsRequestBlocked(req, rule)

	// Assert the expected behavior
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result {
		t.Error("Expected false, got true")
	}
}
