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
	promise := NewPromise()
	_, err := promise.AwaitWithTimeout(100 * time.Millisecond)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}

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
