# Connect Feishu

## Overview

This task implements a Feishu (飞书) API connection module with a Promise-like async pattern in both **Go** and **Python**.

The implementation provides:
- Async/Promise-based connection pattern
- Connection verification via tenant access token
- Connection manager for lifecycle management
- Timeout support
- Error handling

---

## Task 1: Create Connect Feishu with API Key Function

### ✅ Completed

1. ✅ Created a function to connect Feishu with API key (App ID + App Secret)
2. ✅ The function returns a Promise that resolves to the connected Feishu instance

---

## Implementation Details

### Go Version

**Location**: `projects/capture/internal/feishu/connection.go`

#### Key Components

1. **Promise Pattern** (`Promise` struct)
   - Channel-based Promise implementation
   - `Await()` - blocking wait for result
   - `AwaitWithTimeout()` - wait with timeout

2. **ConnectionResult** struct
   - `Client *lark.Client` - the Feishu client
   - `Err error` - any connection error

3. **ConnectWithAPIKey** function
   ```go
   func ConnectWithAPIKey(appID, appSecret string) *Promise
   ```
   - Returns immediately with a Promise
   - Runs connection in background goroutine
   - Verifies credentials by requesting tenant access token

4. **ConnectionManager** struct
   - Manages connection lifecycle
   - Configuration with timeout
   - Connection testing

#### Usage Example

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/variableway/innate/capture/internal/feishu"
)

func main() {
    // Method 1: Using Promise with timeout
    promise := feishu.ConnectWithAPIKey("cli_xxxxxxxx", "xxxxxxxxxxxxxxxx")
    result, err := promise.AwaitWithTimeout(30 * time.Second)
    if err != nil {
        log.Fatal("Timeout:", err)
    }
    if result.Err != nil {
        log.Fatal("Connection failed:", result.Err)
    }
    client := result.Client
    fmt.Println("Connected!")

    // Method 2: Synchronous wrapper
    client, err = feishu.ConnectWithAPIKeySync("cli_xxxxxxxx", "xxxxxxxxxxxxxxxx")
    if err != nil {
        log.Fatal(err)
    }

    // Method 3: Using ConnectionManager
    manager := feishu.NewConnectionManager(&feishu.ConnectionConfig{
        AppID:     "cli_xxxxxxxx",
        AppSecret: "xxxxxxxxxxxxxxxx",
        Timeout:   30 * time.Second,
    })
    promise = manager.Connect()
    result = promise.Await()
}
```

#### Testing

Run tests:
```bash
cd projects/capture
go test ./internal/feishu/... -v
```

---

### Python Version

**Location**: `projects/capture-py/feishu_connection.py`

#### Key Components

1. **ConnectionPromise** class
   - Async/await based Promise
   - `await_result()` - coroutine to wait for result
   - `await_result_with_timeout()` - wait with timeout

2. **ConnectionResult** dataclass
   - `client: Optional[FeishuClient]`
   - `error: Optional[Exception]`
   - `success: bool` property

3. **connect_with_api_key** async function
   ```python
   async def connect_with_api_key(app_id: str, app_secret: str) -> ConnectionPromise
   ```
   - Returns a ConnectionPromise immediately
   - Runs connection verification in background task

4. **ConnectionManager** class
   - Async context manager support
   - Connection state tracking
   - Automatic cleanup

#### Usage Example

```python
import asyncio
from feishu_connection import (
    connect_with_api_key,
    connect_with_api_key_sync,
    ConnectionManager,
    create_feishu_client,
    FeishuAPIError
)

async def main():
    # Method 1: Using Promise with timeout
    promise = await connect_with_api_key("cli_xxxxxxxx", "xxxxxxxxxxxxxxxx")
    result = await promise.await_result_with_timeout(timeout=30.0)
    
    if result.success:
        client = result.client
        print("Connected!")
        await client.close()
    else:
        print(f"Connection failed: {result.error}")

    # Method 2: Direct sync-style call
    result = await connect_with_api_key_sync("cli_xxxxxxxx", "xxxxxxxxxxxxxxxx")
    if result.success:
        client = result.client

    # Method 3: Using ConnectionManager
    manager = ConnectionManager("cli_xxxxxxxx", "xxxxxxxxxxxxxxxx", timeout=30.0)
    promise = await manager.connect()
    result = await promise.await_result()
    
    if manager.is_connected:
        print("Connected via manager!")
        await manager.disconnect()

    # Method 4: Direct client creation (raises on failure)
    try:
        client = await create_feishu_client("cli_xxxxxxxx", "xxxxxxxxxxxxxxxx")
        print(f"Client created: {client}")
        await client.close()
    except FeishuAPIError as e:
        print(f"Failed: {e}")

if __name__ == "__main__":
    asyncio.run(main())
```

#### Installation & Testing

```bash
cd projects/capture-py

# Install dependencies
pip install -r requirements.txt

# Run tests
pytest test_feishu_connection.py -v

# Run example
python feishu_connection.py
```

---

## API Reference

### Go API

| Function/Type | Description |
|--------------|-------------|
| `Promise` | Channel-based promise for async results |
| `ConnectionResult` | Result container with Client and Err |
| `ConnectWithAPIKey(appID, appSecret)` | Returns Promise, connects in background |
| `ConnectWithAPIKeySync(appID, appSecret)` | Blocks until connected, returns (client, error) |
| `ConnectionManager` | Lifecycle manager for connections |
| `ConnectionConfig` | Configuration struct with timeout |

### Python API

| Function/Class | Description |
|---------------|-------------|
| `ConnectionPromise` | Async/await promise for connection results |
| `ConnectionResult` | Dataclass with client, error, success properties |
| `connect_with_api_key(app_id, app_secret)` | Async function returning ConnectionPromise |
| `connect_with_api_key_sync(app_id, app_secret)` | Async function returning ConnectionResult directly |
| `create_feishu_client(app_id, app_secret)` | Async function returning client or raising exception |
| `ConnectionManager` | Async connection lifecycle manager |
| `FeishuClient` | HTTP client with auth token management |
| `FeishuAPIError` | Exception class for API errors |

---

## Environment Variables

Both implementations use the same environment variables:

```bash
# Required
FEISHU_APP_ID=cli_xxxxxxxxxx
FEISHU_APP_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# Optional (for other features)
FEISHU_BITABLE_APP_TOKEN=xxx
FEISHU_BITABLE_TABLE_ID=xxx
```

---

## File Structure

```
projects/
├── capture/
│   └── internal/feishu/
│       ├── client.go              # Existing: Basic client
│       ├── connection.go          # NEW: Promise-based connection
│       └── connection_test.go     # NEW: Unit tests
│
└── capture-py/
    ├── feishu_connection.py       # NEW: Python implementation
    ├── test_feishu_connection.py  # NEW: Python tests
    └── requirements.txt           # NEW: Dependencies
```

---

## Next Steps

1. **Integration**: Integrate with existing bot and bitable modules
2. **Retry Logic**: Add exponential backoff for connection failures
3. **Caching**: Cache access tokens until expiry
4. **WebSocket**: Add WebSocket connection support for real-time events

---

## Notes

- Go version uses `larksuite/oapi-sdk-go/v3` for Feishu API
- Python version uses `httpx` for async HTTP requests
- Both verify connection by requesting tenant access token
- Timeout defaults: Go (30s), Python (30s)
- Both support custom timeout configuration
