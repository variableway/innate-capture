"""
Tests for Feishu Connection Module
"""

import asyncio
import pytest
from unittest.mock import patch, AsyncMock, MagicMock
from feishu_connection import (
    ConnectionPromise,
    ConnectionResult,
    ConnectionManager,
    FeishuClient,
    FeishuAPIError,
    connect_with_api_key,
    connect_with_api_key_sync,
    create_feishu_client,
)


class TestConnectionPromise:
    """Test ConnectionPromise class."""
    
    @pytest.mark.asyncio
    async def test_resolve_and_await(self):
        """Test resolving and awaiting a promise."""
        promise = ConnectionPromise()
        result = ConnectionResult(client=None, error=None)
        
        promise.resolve(result)
        awaited = await promise.await_result()
        
        assert awaited == result
        assert awaited.success
    
    @pytest.mark.asyncio
    async def test_await_with_timeout_success(self):
        """Test awaiting with timeout - success case."""
        promise = ConnectionPromise()
        result = ConnectionResult(client=None, error=None)
        
        promise.resolve(result)
        awaited = await promise.await_result_with_timeout(timeout=1.0)
        
        assert awaited.success
    
    @pytest.mark.asyncio
    async def test_await_with_timeout_timeout(self):
        """Test awaiting with timeout - timeout case."""
        promise = ConnectionPromise()
        # Don't resolve the promise
        
        awaited = await promise.await_result_with_timeout(timeout=0.1)
        
        assert not awaited.success
        assert isinstance(awaited.error, TimeoutError)
    
    @pytest.mark.asyncio
    async def test_reject(self):
        """Test rejecting a promise."""
        promise = ConnectionPromise()
        error = FeishuAPIError("Test error")
        
        promise.reject(error)
        result = await promise.await_result()
        
        assert not result.success
        assert result.error == error


class TestConnectionManager:
    """Test ConnectionManager class."""
    
    def test_initialization(self):
        """Test manager initialization."""
        manager = ConnectionManager("app_id", "app_secret", timeout=10.0)
        
        assert manager.app_id == "app_id"
        assert manager.app_secret == "app_secret"
        assert manager.timeout == 10.0
        assert not manager.is_connected
        assert manager.client is None
    
    @pytest.mark.asyncio
    async def test_default_timeout(self):
        """Test default timeout value."""
        manager = ConnectionManager("app_id", "app_secret")
        assert manager.timeout == 30.0


class TestConnectionResult:
    """Test ConnectionResult dataclass."""
    
    def test_success_with_client(self):
        """Test success property with client."""
        mock_client = MagicMock(spec=FeishuClient)
        result = ConnectionResult(client=mock_client, error=None)
        
        assert result.success
    
    def test_success_with_error(self):
        """Test success property with error."""
        result = ConnectionResult(client=None, error=FeishuAPIError("Error"))
        
        assert not result.success
    
    def test_success_with_none(self):
        """Test success property with None values."""
        result = ConnectionResult(client=None, error=None)
        
        assert not result.success


class TestConnectWithAPIKey:
    """Test connect_with_api_key function."""
    
    @pytest.mark.asyncio
    @patch('feishu_connection.FeishuClient')
    async def test_successful_connection(self, mock_client_class):
        """Test successful connection."""
        mock_client = AsyncMock(spec=FeishuClient)
        mock_client.verify_connection.return_value = True
        mock_client_class.return_value = mock_client
        
        promise = await connect_with_api_key("app_id", "app_secret")
        result = await promise.await_result()
        
        assert result.success
        assert result.client == mock_client
    
    @pytest.mark.asyncio
    @patch('feishu_connection.FeishuClient')
    async def test_failed_verification(self, mock_client_class):
        """Test connection with failed verification."""
        mock_client = AsyncMock(spec=FeishuClient)
        mock_client.verify_connection.return_value = False
        mock_client_class.return_value = mock_client
        
        promise = await connect_with_api_key("app_id", "app_secret")
        result = await promise.await_result()
        
        assert not result.success
        assert result.error is not None
    
    @pytest.mark.asyncio
    @patch('feishu_connection.FeishuClient')
    async def test_connection_exception(self, mock_client_class):
        """Test connection with exception."""
        mock_client = AsyncMock(spec=FeishuClient)
        mock_client.verify_connection.side_effect = FeishuAPIError("Connection failed")
        mock_client_class.return_value = mock_client
        
        promise = await connect_with_api_key("app_id", "app_secret")
        result = await promise.await_result()
        
        assert not result.success
        assert isinstance(result.error, FeishuAPIError)


class TestConnectWithAPIKeySync:
    """Test connect_with_api_key_sync function."""
    
    @pytest.mark.asyncio
    @patch('feishu_connection.FeishuClient')
    async def test_sync_wrapper(self, mock_client_class):
        """Test synchronous-style wrapper."""
        mock_client = AsyncMock(spec=FeishuClient)
        mock_client.verify_connection.return_value = True
        mock_client_class.return_value = mock_client
        
        result = await connect_with_api_key_sync("app_id", "app_secret")
        
        assert result.success
        assert result.client == mock_client


class TestCreateFeishuClient:
    """Test create_feishu_client function."""
    
    @pytest.mark.asyncio
    @patch('feishu_connection.FeishuClient')
    async def test_successful_creation(self, mock_client_class):
        """Test successful client creation."""
        mock_client = AsyncMock(spec=FeishuClient)
        mock_client.verify_connection.return_value = True
        mock_client_class.return_value = mock_client
        
        client = await create_feishu_client("app_id", "app_secret")
        
        assert client == mock_client
    
    @pytest.mark.asyncio
    @patch('feishu_connection.FeishuClient')
    async def test_failed_creation(self, mock_client_class):
        """Test failed client creation raises exception."""
        mock_client = AsyncMock(spec=FeishuClient)
        mock_client.verify_connection.return_value = False
        mock_client_class.return_value = mock_client
        
        with pytest.raises(FeishuAPIError):
            await create_feishu_client("app_id", "app_secret")


class TestFeishuClient:
    """Test FeishuClient class."""
    
    def test_initialization(self):
        """Test client initialization."""
        client = FeishuClient("app_id", "app_secret")
        
        assert client.app_id == "app_id"
        assert client.app_secret == "app_secret"
        assert client.access_token is None
    
    @pytest.mark.asyncio
    @patch('httpx.AsyncClient.post')
    async def test_get_tenant_access_token(self, mock_post):
        """Test getting tenant access token."""
        mock_response = MagicMock()
        mock_response.json.return_value = {
            "code": 0,
            "tenant_access_token": "test_token",
            "expire": 7200
        }
        mock_response.raise_for_status = MagicMock()
        mock_post.return_value = mock_response
        
        client = FeishuClient("app_id", "app_secret")
        token = await client._get_tenant_access_token()
        
        assert token == "test_token"
        assert client.access_token == "test_token"
    
    @pytest.mark.asyncio
    @patch('httpx.AsyncClient.post')
    async def test_get_token_api_error(self, mock_post):
        """Test API error when getting token."""
        mock_response = MagicMock()
        mock_response.json.return_value = {
            "code": 99991663,
            "msg": "app_id error"
        }
        mock_response.raise_for_status = MagicMock()
        mock_post.return_value = mock_response
        
        client = FeishuClient("app_id", "app_secret")
        
        with pytest.raises(FeishuAPIError) as exc_info:
            await client._get_tenant_access_token()
        
        assert "app_id error" in str(exc_info.value)


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
