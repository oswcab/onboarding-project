package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloServer(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "root path returns Hello World",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello World!",
		},
		{
			name:           "path with name returns personalized greeting",
			path:           "/John",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello, John!",
		},
		{
			name:           "path with multiple segments",
			path:           "/John/Doe",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello, John/Doe!",
		},
		{
			name:           "path with special characters",
			path:           "/user-123",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello, user-123!",
		},
		{
			name:           "empty path segment",
			path:           "//",
			expectedStatus: http.StatusOK,
			expectedBody:   "Hello, /!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			// Create a ResponseRecorder to record the response
			w := httptest.NewRecorder()

			// Call the handler
			HelloServer(w, req)

			// Check the status code
			if w.Code != tt.expectedStatus {
				t.Errorf("HelloServer() status = %v, want %v", w.Code, tt.expectedStatus)
			}

			// Check the response body
			if w.Body.String() != tt.expectedBody {
				t.Errorf("HelloServer() body = %v, want %v", w.Body.String(), tt.expectedBody)
			}

			// Check the content type header
			contentType := w.Header().Get("Content-Type")
			expectedContentType := "text/plain; charset=utf-8"
			if contentType != expectedContentType {
				t.Errorf("HelloServer() Content-Type = %v, want %v", contentType, expectedContentType)
			}
		})
	}
}

func TestHelloServerWithDifferentHTTPMethods(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/test", nil)
			w := httptest.NewRecorder()

			HelloServer(w, req)

			// Should always return 200 OK regardless of HTTP method
			if w.Code != http.StatusOK {
				t.Errorf("HelloServer() with %s method status = %v, want %v", method, w.Code, http.StatusOK)
			}

			// Should return the expected greeting
			expectedBody := "Hello, test!"
			if w.Body.String() != expectedBody {
				t.Errorf("HelloServer() with %s method body = %v, want %v", method, w.Body.String(), expectedBody)
			}
		})
	}
}

func TestHelloServerWithRequestBody(t *testing.T) {
	// Test with a request that has a body
	body := bytes.NewBufferString("test body")
	req := httptest.NewRequest(http.MethodPost, "/test", body)
	w := httptest.NewRecorder()

	HelloServer(w, req)

	// Should still work normally
	if w.Code != http.StatusOK {
		t.Errorf("HelloServer() with body status = %v, want %v", w.Code, http.StatusOK)
	}

	expectedBody := "Hello, test!"
	if w.Body.String() != expectedBody {
		t.Errorf("HelloServer() with body response = %v, want %v", w.Body.String(), expectedBody)
	}
}

func TestHelloServerWithQueryParameters(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test?param1=value1&param2=value2", nil)
	w := httptest.NewRecorder()

	HelloServer(w, req)

	// Should ignore query parameters and only use the path
	if w.Code != http.StatusOK {
		t.Errorf("HelloServer() with query params status = %v, want %v", w.Code, http.StatusOK)
	}

	expectedBody := "Hello, test!"
	if w.Body.String() != expectedBody {
		t.Errorf("HelloServer() with query params response = %v, want %v", w.Body.String(), expectedBody)
	}
}

func TestHelloServerWithHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("Accept", "text/plain")
	req.RemoteAddr = "192.168.1.1:12345"

	w := httptest.NewRecorder()

	HelloServer(w, req)

	// Should work normally with headers
	if w.Code != http.StatusOK {
		t.Errorf("HelloServer() with headers status = %v, want %v", w.Code, http.StatusOK)
	}

	expectedBody := "Hello, test!"
	if w.Body.String() != expectedBody {
		t.Errorf("HelloServer() with headers response = %v, want %v", w.Body.String(), expectedBody)
	}
}

// Benchmark tests
func BenchmarkHelloServer(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		HelloServer(w, req)
	}
}

func BenchmarkHelloServerRoot(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		HelloServer(w, req)
	}
}
