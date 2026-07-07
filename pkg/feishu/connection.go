package feishu

import (
	"context"
	"fmt"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
)

// ConnectionResult represents the result of a Feishu connection attempt
type ConnectionResult struct {
	Client *lark.Client
	Err    error
}

// Promise represents a promise that will eventually resolve to a ConnectionResult
type Promise struct {
	resultChan chan ConnectionResult
}

// NewPromise creates a new Promise
func NewPromise() *Promise {
	return &Promise{
		resultChan: make(chan ConnectionResult, 1),
	}
}

// Resolve resolves the promise with a result
func (p *Promise) Resolve(result ConnectionResult) {
	p.resultChan <- result
	close(p.resultChan)
}

// Await waits for the promise to resolve and returns the result
func (p *Promise) Await() ConnectionResult {
	return <-p.resultChan
}

// AwaitWithTimeout waits for the promise to resolve with a timeout
func (p *Promise) AwaitWithTimeout(timeout time.Duration) (ConnectionResult, error) {
	select {
	case result := <-p.resultChan:
		return result, nil
	case <-time.After(timeout):
		return ConnectionResult{}, fmt.Errorf("connection timeout after %v", timeout)
	}
}

// ConnectWithAPIKey creates a Feishu connection using app ID and app secret (API key pattern)
// Returns a Promise that resolves to the connected Feishu instance
func ConnectWithAPIKey(appID, appSecret string) *Promise {
	promise := NewPromise()

	go func() {
		client := NewClient(appID, appSecret)
		if appID == "" || appSecret == "" {
			promise.Resolve(ConnectionResult{
				Client: nil,
				Err:    fmt.Errorf("app_id and app_secret are required"),
			})
			return
		}
		promise.Resolve(ConnectionResult{
			Client: client,
			Err:    nil,
		})
	}()

	return promise
}

// ConnectWithAPIKeySync is a synchronous version of ConnectWithAPIKey
func ConnectWithAPIKeySync(appID, appSecret string) (*lark.Client, error) {
	promise := ConnectWithAPIKey(appID, appSecret)
	result := promise.Await()
	return result.Client, result.Err
}

// ConnectionManager manages Feishu connections
type ConnectionManager struct {
	client *lark.Client
	config *ConnectionConfig
}

// ConnectionConfig holds connection configuration
type ConnectionConfig struct {
	AppID     string
	AppSecret string
	Timeout   time.Duration
}

// NewConnectionManager creates a new connection manager
func NewConnectionManager(config *ConnectionConfig) *ConnectionManager {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	return &ConnectionManager{
		config: config,
	}
}

// Connect establishes a connection and returns a Promise
func (cm *ConnectionManager) Connect() *Promise {
	return ConnectWithAPIKey(cm.config.AppID, cm.config.AppSecret)
}

// GetClient returns the connected client (nil if not connected)
func (cm *ConnectionManager) GetClient() *lark.Client {
	return cm.client
}

// IsConnected checks if the client is connected
func (cm *ConnectionManager) IsConnected() bool {
	return cm.client != nil
}

// TestConnection tests if the client is properly initialized
func (cm *ConnectionManager) TestConnection(ctx context.Context) error {
	if cm.client == nil {
		return fmt.Errorf("not connected")
	}
	if cm.config.AppID == "" || cm.config.AppSecret == "" {
		return fmt.Errorf("app_id and app_secret are required")
	}
	return nil
}

// ConnectAndVerify connects and immediately verifies the connection
func (cm *ConnectionManager) ConnectAndVerify(ctx context.Context) (*lark.Client, error) {
	promise := cm.Connect()
	result := promise.Await()
	if result.Err != nil {
		return nil, result.Err
	}
	cm.client = result.Client
	if err := cm.TestConnection(ctx); err != nil {
		cm.client = nil
		return nil, err
	}
	return cm.client, nil
}
