package main

import (
	"log/slog"
	"os"
	"testing"
)

// Test utilities and helpers

// setupTestLogger creates a test logger that writes to a buffer
func setupTestLogger(t *testing.T) *slog.Logger {
	t.Helper()
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

// cleanupEnvVar restores an environment variable to its original value
func cleanupEnvVar(t *testing.T, key, originalValue string) {
	t.Helper()
	if originalValue != "" {
		os.Setenv(key, originalValue)
	} else {
		os.Unsetenv(key)
	}
}

// TestLoggerSetup tests the logger configuration
func TestLoggerSetup(t *testing.T) {
	logger := setupTestLogger(t)

	if logger == nil {
		t.Error("Logger should not be nil")
	}

	// Test that we can log without errors
	logger.Info("Test log message", "test", true)
}

// TestEnvironmentCleanup tests environment variable cleanup
func TestEnvironmentCleanup(t *testing.T) {
	// Save original PORT value
	originalPort := os.Getenv("PORT")
	defer cleanupEnvVar(t, "PORT", originalPort)

	// Test setting and cleaning up environment variable
	os.Setenv("PORT", "test-port")
	if os.Getenv("PORT") != "test-port" {
		t.Error("Failed to set environment variable")
	}

	cleanupEnvVar(t, "PORT", originalPort)
	if os.Getenv("PORT") != originalPort {
		t.Errorf("Failed to restore environment variable. Expected %s, got %s", originalPort, os.Getenv("PORT"))
	}
}

// TestMain provides setup and teardown for all tests
func TestMain(m *testing.M) {
	// Setup
	code := m.Run()

	// Teardown - cleanup any global state if needed
	os.Exit(code)
}
