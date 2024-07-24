package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestErrorHandlerCases(t *testing.T) {
	suite := NewTestSuite(t)
	suite.TestErrorHandler()
	suite.TestTemplateNotFound()
	suite.TestTemplateExecutionError()
	suite.TestTemplateSyntaxError()
	suite.TestVeryLongErrorMessage()
}

func (ts *TestSuite) TestErrorHandler() {
	longMessage := make([]byte, 10)
	for i := range longMessage {
		longMessage[i] = 'a'
	}
	testCases := []struct {
		name               string
		str                string
		code               int
		expectedStatusCode int
	}{
		{
			name:               "Test Valid Error",
			str:                "Bad Request",
			code:               400,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Test Invalid Error",
			str:                "Error 500: Internal server error",
			code:               4000,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "Test Empty Message",
			str:                "",
			code:               500,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range testCases {
		ts.t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ErrorHandler(w, tt.str, tt.code)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("ErrorHandler() status code = %v, want %v", resp.StatusCode, tt.expectedStatusCode)
			}
			body := w.Body.String()
			if !strings.Contains(body, tt.str) {
				t.Errorf("ErrorHandler() body = %v, want %v", body, tt.str)
			}
		})
	}
}

// Mock template parsing to induce errors
func ErrorTemplate() (*template.Template, error) {
	return nil, fmt.Errorf("failed to parse template")
}

func TestErrorHandlerTemplateError4(t *testing.T) {
	w := httptest.NewRecorder()
	ErrorHandler(w, "Error 500: Internal server error", http.StatusInternalServerError)
	resp := w.Result()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("ErrorHandler() status code = %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Error 500: Internal server error") {
		t.Errorf("ErrorHandler() body = %v, want %v", body, "Error 500: Internal server error")
	}
}

// Mocktemplate parsing to induce an error
var templateParseFiles = template.ParseFiles

func mockTemplateParseFilesError() {
	templateParseFiles = func(filenames ...string) (*template.Template, error) {
		return nil, fmt.Errorf("mock template parse error")
	}
}

func restoreTemplateParseFiles() {
	templateParseFiles = template.ParseFiles
}

// Mock ResponseWriter to induce errors
type mockResponseWriter struct {
	httptest.ResponseRecorder
	forceError bool
}

func (mw *mockResponseWriter) Write(data []byte) (int, error) {
	if mw.forceError {
		return 0, fmt.Errorf("mock write error")
	}
	return mw.ResponseRecorder.Write(data)
}

func (ts *TestSuite) TestTemplateNotFound() {
	restoreTemplateParseFiles()
	w := httptest.NewRecorder()
	ErrorHandler(w, "Test error", http.StatusInternalServerError)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		ts.t.Errorf("Expected ErrorHandler() status code = %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

func (ts *TestSuite) TestTemplateExecutionError() {
	mockTemplateParseFilesError()
	defer restoreTemplateParseFiles()

	w := &mockResponseWriter{forceError: true}
	ErrorHandler(w, "Test error", http.StatusInternalServerError)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		ts.t.Errorf("Expected ErrorHandler() status code = %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

func (ts *TestSuite) TestTemplateSyntaxError() {
	w := httptest.NewRecorder()
	ErrorHandler(w, "Test error", http.StatusInternalServerError)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		ts.t.Errorf("Expected ErrorHandler() status code = %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

func (ts *TestSuite) TestVeryLongErrorMessage() {
	longMessage := make([]byte, 10000)
	for i := range longMessage {
		longMessage[i] = 'a'
	}
	w := httptest.NewRecorder()
	ErrorHandler(w, string(longMessage), http.StatusInternalServerError)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		ts.t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}
}
