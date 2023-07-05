package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestNewProxy(t *testing.T) {
	urlRaw := "http://example.com"
	timeout := 5 * time.Second
	monitor := true

	proxy, err := NewProxy(urlRaw, timeout, monitor)
	if err != nil {
		t.Errorf("Failed to create proxy: %v", err)
	}

	// Check if the proxy is initialized correctly
	if proxy.Proxy == nil {
		t.Error("Proxy instance is nil")
	}

	parsedURL, err := url.Parse(urlRaw)
	if err != nil {
		t.Errorf("Failed to parse URL: %v", err)
	}

	// Check if the proxy's Director function sets the correct target URL
	request := httptest.NewRequest("GET", "http://example.com", nil)
	proxy.Proxy.Director(request)

	if request.URL.String() != parsedURL.String()+"/" {
		t.Errorf("Proxy target URL mismatch. Expected: %s, Got: %s", parsedURL, request.URL)
	}

	if proxy.Timeout != timeout {
		t.Errorf("Timeout mismatch. Expected: %s, Got: %s", timeout, proxy.Timeout)
	}

	if proxy.MonitoringModeEnabled != monitor {
		t.Errorf("Monitoring mode mismatch. Expected: %v, Got: %v", monitor, proxy.MonitoringModeEnabled)
	}

	// Check if the proxy's transport is set correctly
	transport := proxy.Proxy.Transport.(*http.Transport)

	if transport.DialContext == nil {
		t.Error("DialContext is not set")
	}

	if transport.MaxIdleConns != 100 {
		t.Errorf("MaxIdleConns mismatch. Expected: %d, Got: %d", 100, transport.MaxIdleConns)
	}

	if transport.IdleConnTimeout != 90*time.Second {
		t.Errorf("IdleConnTimeout mismatch. Expected: %s, Got: %s", 90*time.Second, transport.IdleConnTimeout)
	}

	if transport.TLSHandshakeTimeout != 10*time.Second {
		t.Errorf("TLSHandshakeTimeout mismatch. Expected: %s, Got: %s", 10*time.Second, transport.TLSHandshakeTimeout)
	}

	if transport.ExpectContinueTimeout != timeout {
		t.Errorf("ExpectContinueTimeout mismatch. Expected: %s, Got: %s", timeout, transport.ExpectContinueTimeout)
	}
}

// func TestSimpleProxy_ServeHTTP(t *testing.T) {
// 	// Create a mock HTTP request
// 	request := httptest.NewRequest("GET", "http://example.com", nil)

// 	// Create a mock HTTP response recorder
// 	recorder := httptest.NewRecorder()

// 	// Create a SimpleProxy instance with test configurations
// 	proxy := &SimpleProxy{
// 		Proxy:                 nil, // Mock the reverse proxy for testing purposes
// 		Timeout:               5 * time.Second,
// 		MonitoringModeEnabled: true,
// 	}

// 	// Call the ServeHTTP method with the mock request and response recorder
// 	proxy.ServeHTTP(recorder, request)

// 	// TODO: Add assertions based on your specific test requirements

// 	// Example assertions for verifying logging functionality
// 	expectedLogMessage := "127.0.0.1: GET /"
// 	if !logContains(expectedLogMessage) {
// 		t.Errorf("Expected log message not found: %s", expectedLogMessage)
// 	}

// 	// Example assertions for verifying context timeout setting
// 	expectedTimeout := 5 * time.Second
// 	ctx := request.Context()
// 	deadline, ok := ctx.Deadline()
// 	if !ok || deadline.Sub(time.Now()) != expectedTimeout {
// 		t.Errorf("Context timeout not set correctly. Expected: %s, Got: %s", expectedTimeout, deadline.Sub(time.Now()))
// 	}

// 	// Example assertions for verifying the routing logic
// 	if !isMonitoringRequest(recorder) {
// 		t.Error("Monitoring request not detected")
// 	}

// 	// TODO: Add more assertions based on your specific test requirements
// }

// // Helper function to check if the log contains the specified message
// func logContains(message string) bool {
// 	// TODO: Implement the log checking logic
// 	return false
// }

// // Helper function to check if the request is a monitoring request
// func isMonitoringRequest(recorder *httptest.ResponseRecorder) bool {
// 	// TODO: Implement the monitoring request detection logic
// 	return false
// }
