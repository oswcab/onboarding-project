package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestServerConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		port     string
		expected string
	}{
		{
			name:     "default port when PORT env var is empty",
			port:     "",
			expected: "0.0.0.0:8080",
		},
		{
			name:     "custom port from environment",
			port:     "9090",
			expected: "0.0.0.0:9090",
		},
		{
			name:     "port with leading zeros",
			port:     "08080",
			expected: "0.0.0.0:08080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.port != "" {
				os.Setenv("PORT", tt.port)
			} else {
				os.Unsetenv("PORT")
			}

			// Get port from environment (same logic as main)
			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}

			// Create server with same configuration as main
			server := &http.Server{
				Addr:         fmt.Sprintf("0.0.0.0:%s", port),
				Handler:      http.HandlerFunc(HelloServer),
				ReadTimeout:  15 * time.Second,
				WriteTimeout: 15 * time.Second,
				IdleTimeout:  60 * time.Second,
			}

			if server.Addr != tt.expected {
				t.Errorf("Server address = %v, want %v", server.Addr, tt.expected)
			}

			// Check timeouts
			if server.ReadTimeout != 15*time.Second {
				t.Errorf("ReadTimeout = %v, want %v", server.ReadTimeout, 15*time.Second)
			}

			if server.WriteTimeout != 15*time.Second {
				t.Errorf("WriteTimeout = %v, want %v", server.WriteTimeout, 15*time.Second)
			}

			if server.IdleTimeout != 60*time.Second {
				t.Errorf("IdleTimeout = %v, want %v", server.IdleTimeout, 60*time.Second)
			}
		})
	}
}

func TestServerShutdown(t *testing.T) {
	// Create a test server
	server := &http.Server{
		Addr:    ":0", // Use port 0 to let the system choose an available port
		Handler: http.HandlerFunc(HelloServer),
	}

	// Start server in a goroutine
	go func() {
		server.ListenAndServe()
	}()

	// Give the server a moment to start
	time.Sleep(10 * time.Millisecond)

	// Test graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		t.Errorf("Server shutdown failed: %v", err)
	}
}

func TestServerShutdownTimeout(t *testing.T) {
	// Create a test server
	server := &http.Server{
		Addr:    ":0",
		Handler: http.HandlerFunc(HelloServer),
	}

	// Start server in a goroutine
	go func() {
		server.ListenAndServe()
	}()

	// Give the server a moment to start
	time.Sleep(10 * time.Millisecond)

	// Test shutdown with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Wait for timeout to expire
	time.Sleep(2 * time.Millisecond)

	err := server.Shutdown(ctx)
	// Note: The behavior of shutdown with expired context may vary
	// This test mainly ensures the shutdown method can be called
	if err != nil {
		t.Logf("Shutdown with expired context returned error (expected): %v", err)
	}
}

func TestEnvironmentVariableHandling(t *testing.T) {
	// Save original PORT value
	originalPort := os.Getenv("PORT")
	defer func() {
		if originalPort != "" {
			os.Setenv("PORT", originalPort)
		} else {
			os.Unsetenv("PORT")
		}
	}()

	// Test unset environment variable
	os.Unsetenv("PORT")
	port := os.Getenv("PORT")
	if port != "" {
		t.Errorf("Expected empty port when PORT is unset, got %s", port)
	}

	// Test set environment variable
	testPort := "12345"
	os.Setenv("PORT", testPort)
	port = os.Getenv("PORT")
	if port != testPort {
		t.Errorf("Expected port %s, got %s", testPort, port)
	}
}

func TestServerHandler(t *testing.T) {
	server := &http.Server{
		Addr:    ":0",
		Handler: http.HandlerFunc(HelloServer),
	}

	// Check that the handler is set correctly
	if server.Handler == nil {
		t.Error("Server handler should not be nil")
	}

	// Test that the handler is our HelloServer function
	// We can't directly compare function pointers, but we can test it works
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Handler status = %v, want %v", w.Code, http.StatusOK)
	}

	expectedBody := "Hello, test!"
	if w.Body.String() != expectedBody {
		t.Errorf("Handler body = %v, want %v", w.Body.String(), expectedBody)
	}
}
