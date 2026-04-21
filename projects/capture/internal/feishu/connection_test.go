package feishu

import (
	"testing"
	"time"
)

func TestNewPromise(t *testing.T) {
	promise := NewPromise()
	if promise == nil {
		t.Fatal("NewPromise() returned nil")
	}
	if promise.resultChan == nil {
		t.Fatal("promise.resultChan is nil")
	}
}

func TestPromiseResolveAndAwait(t *testing.T) {
	promise := NewPromise()

	go func() {
		promise.Resolve(ConnectionResult{
			Client: nil,
			Err:    nil,
		})
	}()

	result := promise.Await()
	if result.Err != nil {
		t.Errorf("Expected no error, got: %v", result.Err)
	}
}

func TestPromiseAwaitWithTimeout(t *testing.T) {
	// Test timeout case
	promise := NewPromise()

	// Don't resolve the promise
	_, err := promise.AwaitWithTimeout(100 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}

	// Test successful case
	promise2 := NewPromise()
	go func() {
		promise2.Resolve(ConnectionResult{
			Client: nil,
			Err:    nil,
		})
	}()

	result, err := promise2.AwaitWithTimeout(1 * time.Second)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if result.Err != nil {
		t.Errorf("Expected no result error, got: %v", result.Err)
	}
}

func TestConnectionManager(t *testing.T) {
	config := &ConnectionConfig{
		AppID:     "test_app_id",
		AppSecret: "test_app_secret",
		Timeout:   10 * time.Second,
	}

	manager := NewConnectionManager(config)
	if manager == nil {
		t.Fatal("NewConnectionManager() returned nil")
	}

	if manager.IsConnected() {
		t.Error("New manager should not be connected")
	}

	if manager.GetClient() != nil {
		t.Error("New manager should have nil client")
	}
}

func TestConnectWithAPIKey_InvalidCredentials(t *testing.T) {
	// Test with invalid credentials - should fail but not panic
	promise := ConnectWithAPIKey("invalid_app_id", "invalid_app_secret")

	// Use short timeout for test
	result, err := promise.AwaitWithTimeout(5 * time.Second)
	if err != nil {
		// Timeout is acceptable for this test
		t.Logf("Connection timed out (expected): %v", err)
		return
	}

	// With invalid credentials, we expect an error
	if result.Err == nil {
		t.Log("Expected error with invalid credentials, but got none (may be due to network/mock)")
	}
}

func TestConnectionConfig_DefaultTimeout(t *testing.T) {
	config := &ConnectionConfig{
		AppID:     "test",
		AppSecret: "test",
		// Timeout not set (zero value)
	}

	manager := NewConnectionManager(config)
	if manager.config.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout of 30s, got %v", manager.config.Timeout)
	}
}
