"""
Feishu Connection Module (Python Version)

This module provides async/await based Feishu API connection functionality
similar to the Go Promise pattern.
"""

import asyncio
from dataclasses import dataclass
from typing import Optional, Callable, Any
import httpx
from datetime import datetime


@dataclass
class ConnectionResult:
    """Result of a Feishu connection attempt."""
    client: Optional['FeishuClient'] = None
    error: Optional[Exception] = None
    
    @property
    def success(self) -> bool:
        """Check if connection was successful."""
        return self.client is not None and self.error is None


class FeishuClient:
    """Feishu API Client."""
    
    BASE_URL = "https://open.feishu.cn/open-apis"
    
    def __init__(self, app_id: str, app_secret: str, access_token: Optional[str] = None):
        self.app_id = app_id
        self.app_secret = app_secret
        self.access_token = access_token
        self._http_client = httpx.AsyncClient(base_url=self.BASE_URL, timeout=30.0)
    
    async def close(self):
        """Close the HTTP client."""
        await self._http_client.aclose()
    
    async def _get_tenant_access_token(self) -> str:
        """Get tenant access token using app credentials."""
        url = "/auth/v3/tenant_access_token/internal"
        payload = {
            "app_id": self.app_id,
            "app_secret": self.app_secret
        }
        
        response = await self._http_client.post(url, json=payload)
        response.raise_for_status()
        
        data = response.json()
        if data.get("code") != 0:
            raise FeishuAPIError(f"Failed to get access token: {data.get('msg')}")
        
        return data["tenant_access_token"]
    
    async def verify_connection(self) -> bool:
        """Verify the connection by getting an access token."""
        try:
            token = await self._get_tenant_access_token()
            self.access_token = token
            return True
        except Exception:
            return False
    
    async def get_auth_headers(self) -> dict:
        """Get authentication headers for API requests."""
        if not self.access_token:
            self.access_token = await self._get_tenant_access_token()
        return {"Authorization": f"Bearer {self.access_token}"}


class FeishuAPIError(Exception):
    """Feishu API Error."""
    pass


class ConnectionPromise:
    """
    A Promise-like object for Feishu connections.
    
    Similar to JavaScript Promise or Go channels, this allows async/await
    style programming for Feishu connections.
    """
    
    def __init__(self):
        self._future: asyncio.Future[ConnectionResult] = asyncio.Future()
    
    def resolve(self, result: ConnectionResult):
        """Resolve the promise with a result."""
        if not self._future.done():
            self._future.set_result(result)
    
    def reject(self, error: Exception):
        """Reject the promise with an error."""
        if not self._future.done():
            self._future.set_result(ConnectionResult(error=error))
    
    async def await_result(self) -> ConnectionResult:
        """Await the promise result."""
        return await self._future
    
    async def await_result_with_timeout(self, timeout: float = 30.0) -> ConnectionResult:
        """Await the promise result with a timeout."""
        try:
            return await asyncio.wait_for(self._future, timeout=timeout)
        except asyncio.TimeoutError:
            return ConnectionResult(error=TimeoutError(f"Connection timeout after {timeout}s"))


async def connect_with_api_key(app_id: str, app_secret: str) -> ConnectionPromise:
    """
    Create a Feishu connection using app ID and app secret.
    
    Args:
        app_id: Feishu App ID
        app_secret: Feishu App Secret
        
    Returns:
        ConnectionPromise that resolves to a ConnectionResult
        
    Example:
        >>> promise = await connect_with_api_key("cli_xxx", "secret_xxx")
        >>> result = await promise.await_result()
        >>> if result.success:
        ...     print("Connected!")
        ... else:
        ...     print(f"Error: {result.error}")
    """
    promise = ConnectionPromise()
    
    async def _connect():
        try:
            client = FeishuClient(app_id, app_secret)
            
            # Verify connection
            is_valid = await client.verify_connection()
            if not is_valid:
                await client.close()
                promise.resolve(ConnectionResult(
                    error=FeishuAPIError("Failed to verify Feishu connection")
                ))
                return
            
            promise.resolve(ConnectionResult(client=client))
            
        except Exception as e:
            promise.resolve(ConnectionResult(error=e))
    
    # Start connection in background
    asyncio.create_task(_connect())
    
    return promise


async def connect_with_api_key_sync(app_id: str, app_secret: str) -> ConnectionResult:
    """
    Synchronous-style connection (still async under the hood).
    
    Args:
        app_id: Feishu App ID
        app_secret: Feishu App Secret
        
    Returns:
        ConnectionResult directly
    """
    promise = await connect_with_api_key(app_id, app_secret)
    return await promise.await_result()


class ConnectionManager:
    """Manages Feishu connections with configuration."""
    
    def __init__(self, app_id: str, app_secret: str, timeout: float = 30.0):
        self.app_id = app_id
        self.app_secret = app_secret
        self.timeout = timeout
        self._client: Optional[FeishuClient] = None
        self._connected = False
    
    @property
    def is_connected(self) -> bool:
        """Check if currently connected."""
        return self._connected and self._client is not None
    
    @property
    def client(self) -> Optional[FeishuClient]:
        """Get the connected client."""
        return self._client
    
    async def connect(self) -> ConnectionPromise:
        """Establish a connection and return a Promise."""
        promise = await connect_with_api_key(self.app_id, self.app_secret)
        
        # Set up callback to store client on success
        async def _on_connect():
            result = await promise.await_result()
            if result.success:
                self._client = result.client
                self._connected = True
        
        asyncio.create_task(_on_connect())
        return promise
    
    async def test_connection(self) -> bool:
        """Test the current connection."""
        if not self.is_connected or not self._client:
            return False
        return await self._client.verify_connection()
    
    async def disconnect(self):
        """Disconnect and cleanup."""
        if self._client:
            await self._client.close()
            self._client = None
            self._connected = False


# Convenience functions for common use cases

async def create_feishu_client(app_id: str, app_secret: str) -> FeishuClient:
    """
    Create and verify a Feishu client.
    
    Raises:
        FeishuAPIError: If connection fails
    """
    result = await connect_with_api_key_sync(app_id, app_secret)
    if not result.success:
        raise result.error or FeishuAPIError("Unknown connection error")
    return result.client


# Example usage and testing
if __name__ == "__main__":
    async def main():
        """Example usage."""
        # Example 1: Using Promise pattern
        print("Example 1: Promise pattern")
        promise = await connect_with_api_key("cli_xxx", "secret_xxx")
        result = await promise.await_result_with_timeout(10.0)
        
        if result.success:
            print(f"Connected! Client: {result.client}")
            await result.client.close()
        else:
            print(f"Connection failed: {result.error}")
        
        # Example 2: Using ConnectionManager
        print("\nExample 2: ConnectionManager")
        manager = ConnectionManager("cli_xxx", "secret_xxx", timeout=10.0)
        promise = await manager.connect()
        result = await promise.await_result()
        
        if result.success:
            print(f"Connected via manager! Is connected: {manager.is_connected}")
            await manager.disconnect()
        else:
            print(f"Manager connection failed: {result.error}")
        
        # Example 3: Direct client creation
        print("\nExample 3: Direct client creation")
        try:
            client = await create_feishu_client("cli_xxx", "secret_xxx")
            print(f"Client created: {client}")
            await client.close()
        except FeishuAPIError as e:
            print(f"Failed to create client: {e}")
    
    # Run examples
    asyncio.run(main())
