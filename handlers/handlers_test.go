package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestSuite represents a collection of test cases for Handlers.
type TestSuite struct {
	t *testing.T
}

// NewTestSuite creates a new instance of the test suite.
func NewTestSuite(t *testing.T) *TestSuite {
	return &TestSuite{t: t}
}

func TestIndexHandlerCases(t *testing.T) {
	suite := NewTestSuite(t)
	suite.TestIndexHandler()
	suite.TestTemplateRendering()
}

func TestAsciiArtHandlerCases(t *testing.T) {
	suite := NewTestSuite(t)
	suite.TestAsciiArtHandler()
	suite.TestAsciiArtHandlerInvalidMethod()
}

// Test cases for IndexHandler function.
func (ts *TestSuite) TestIndexHandler() {
	testCases := []struct {
		name              string
		method            string
		path              string
		expectedStatus    int
		expectedSubstring string // Optional: Check for substring in the response body
	}{
		{
			name:           "Root Path GET",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK, // expected status code 200
		},
		{
			name:           "Non-Root Path GET",
			method:         "GET",
			path:           "/about",
			expectedStatus: http.StatusNotFound, // expected status code 404
		},
		{
			name:           "Invalid Method POST",
			method:         "POST",
			path:           "/",
			expectedStatus: http.StatusMethodNotAllowed, // expected status code 405
		},
	}

	for _, tc := range testCases {
		ts.t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				ts.t.Fatal(err)
			}
			// Create a ResponseRecorder to capture the handler's response
			rr := httptest.NewRecorder()
			IndexHandler(rr, req)

			// Check the response status code
			if status := rr.Code; status != tc.expectedStatus {
				ts.t.Errorf("handler returned wrong status code for %s: got %v, want %v", tc.name, status, tc.expectedStatus)
			}

			// Check response body for substring
			if tc.expectedSubstring != "" {
				if !strings.Contains(rr.Body.String(), tc.expectedSubstring) {
					ts.t.Errorf("handler response body does not contain expected substring for %s", tc.name)
				}
			}
		})
	}
}

func (ts *TestSuite) TestTemplateRendering() {
	// Create a request to simulate an HTTP GET to the root path
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		ts.t.Fatal(err)
	}

	// Create a ResponseRecorder to capture the handler's response
	rr := httptest.NewRecorder()
	IndexHandler(rr, req)

	// Check the content type (expecting text/html)
	if contentType := rr.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		ts.t.Errorf("handler did not set the Content-Type header correctly: got %v, want %v", contentType, "text/html; charset=utf-8")
	}
}

func (ts *TestSuite) TestAsciiArtHandler() {
	testCases := []struct {
		name               string
		text               string
		banner             string
		expectedStatusCode int
	}{
		{
			name:               "Valid POST Request",
			text:               "Hello",
			banner:             "shadow.txt",
			expectedStatusCode: http.StatusOK, // expected code 200
		},
		{
			name:               "Empty Text",
			text:               "",
			banner:             "shadow.txt",
			expectedStatusCode: http.StatusBadRequest, // expected code 400
		},
		{
			name:               "Empty Banner",
			text:               "Hello",
			banner:             "",
			expectedStatusCode: http.StatusBadRequest, // expected code 400
		},
		{
			name:               "Invalid Banner",
			text:               "Hello",
			banner:             "nonexistent.txt",
			expectedStatusCode: http.StatusInternalServerError, // expected code 500
		},
	}

	for _, tc := range testCases {
		ts.t.Run(tc.name, func(t *testing.T) {
			// Prepare request
			formData := strings.NewReader("text=" + tc.text + "&banner=" + tc.banner)
			req, err := http.NewRequest("POST", "/ascii-art", formData)
			if err != nil {
				ts.t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Prepare response recorder
			rr := httptest.NewRecorder()

			// Call AsciiArtHandler
			AsciiArtHandler(rr, req)

			// Check the response status code
			if status := rr.Code; status != tc.expectedStatusCode {
				ts.t.Errorf("handler returned wrong status code for %s: got %v, want %v", tc.name, status, tc.expectedStatusCode)
			}
		})
	}
}

func (ts *TestSuite) TestAsciiArtHandlerInvalidMethod() {
	// Prepare a GET request (invalid method)
	req, err := http.NewRequest("GET", "/ascii-art", nil)
	if err != nil {
		ts.t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	AsciiArtHandler(rr, req)

	// Check the response status code (expecting 405 Method Not Allowed)
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		ts.t.Errorf("handler returned wrong status code for invalid method: got %v, want %v", status, http.StatusMethodNotAllowed)
	}
}
