package celestrak

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Fetch_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "text/plain, */*;q=0.9", r.Header.Get("Accept"))
		assert.Equal(t, "en-US,en;q=0.9", r.Header.Get("Accept-Language"))
		assert.Equal(t, "Mozilla/5.0 ...", r.Header.Get("User-Agent"))
		assert.Equal(t, "gzip", r.Header.Get("Accept-Encoding"))
		assert.Equal(t, "keep-alive", r.Header.Get("Connection"))
		assert.Equal(t, "no-cache", r.Header.Get("Cache-Control"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("celestrak data"))
	}))
	defer server.Close()

	client := NewClient(
		server.URL,
		time.Second,
		1,
	)

	body, err := client.Fetch(context.Background())

	require.NoError(t, err)
	require.NotNil(t, body)
	t.Cleanup(func() {
		_ = body.Close()
	})

	data, err := io.ReadAll(body)

	require.NoError(t, err)
	assert.Equal(t, "celestrak data", string(data))
}

func TestClient_Fetch_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(
		server.URL,
		time.Second,
		1,
	)

	body, err := client.Fetch(context.Background())

	assert.Nil(t, body)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "celestrak fetch failed after 1 retries")
	assert.Contains(t, err.Error(), "unexpected status: 404")
}

func TestClient_Fetch_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(
		server.URL,
		time.Second,
		1,
	)

	body, err := client.Fetch(context.Background())

	assert.Nil(t, body)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status: 500")
}

func TestClient_Fetch_InvalidURL(t *testing.T) {
	client := NewClient(
		"://invalid-url",
		time.Second,
		1,
	)

	body, err := client.Fetch(context.Background())

	assert.Nil(t, body)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "celestrak fetch failed after 1 retries")
}

func TestClient_Fetch_ContextAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("request should not have reached the server")
	}))
	defer server.Close()

	client := NewClient(
		server.URL,
		time.Second,
		3,
	)

	body, err := client.Fetch(ctx)

	assert.Nil(t, body)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestClient_Fetch_RetriesAfterFailure(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := attempts.Add(1)

		if attempt == 1 {
			http.Error(w, "temporary failure", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success after retry"))
	}))
	defer server.Close()

	client := NewClient(
		server.URL,
		time.Second,
		1,
	)

	start := time.Now()

	body, err := client.Fetch(context.Background())

	require.NoError(t, err)
	require.NotNil(t, body)
	t.Cleanup(func() {
		_ = body.Close()
	})

	data, err := io.ReadAll(body)

	require.NoError(t, err)
	assert.Equal(t, "success after retry", string(data))
	assert.Equal(t, int32(2), attempts.Load())

	assert.GreaterOrEqual(t, time.Since(start), 500*time.Millisecond)
}

func TestClient_Fetch_ExhaustsRetries(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		http.Error(w, "temporary failure", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(
		server.URL,
		time.Second,
		2,
	)

	body, err := client.Fetch(context.Background())

	assert.Nil(t, body)
	require.Error(t, err)

	assert.Equal(t, int32(3), attempts.Load())
	assert.Contains(t, err.Error(), "celestrak fetch failed after 2 retries")
	assert.Contains(t, err.Error(), "unexpected status: 503")
}

func TestClient_Fetch_ContextCancelledDuringRetry(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		http.Error(w, "temporary failure", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(
		server.URL,
		time.Second,
		3,
	)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	body, err := client.Fetch(ctx)

	assert.Nil(t, body)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)

	assert.Equal(t, int32(1), attempts.Load())
}

func TestClient_fetchOnce_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	}))
	defer server.Close()

	client := NewClient(
		server.URL,
		time.Second,
		1,
	)

	body, err := client.fetchOnce(context.Background())

	require.NoError(t, err)
	require.NotNil(t, body)
	t.Cleanup(func() {
		_ = body.Close()
	})

	data, err := io.ReadAll(body)

	require.NoError(t, err)
	assert.Equal(t, "hello", string(data))
}

func TestClient_fetchOnce_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	client := NewClient(
		server.URL,
		time.Second,
		1,
	)

	body, err := client.fetchOnce(context.Background())

	assert.Nil(t, body)
	require.Error(t, err)
	assert.EqualError(t, err, "unexpected status: 502")
}

func TestClient_fetchOnce_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	client := NewClient(
		server.URL,
		time.Second,
		1,
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	body, err := client.fetchOnce(ctx)

	assert.Nil(t, body)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestClient_NewClient(t *testing.T) {
	client := NewClient(
		"https://example.com/data",
		5*time.Second,
		3,
	)

	require.NotNil(t, client)
	require.NotNil(t, client.client)

	assert.Equal(t, "https://example.com/data", client.sourceURL)
	assert.Equal(t, 3, client.retries)
	assert.Equal(t, 5*time.Second, client.client.Timeout)
}

func TestClient_Fetch_ConnectionError(t *testing.T) {
	client := NewClient(
		"http://127.0.0.1:1",
		100*time.Millisecond,
		1,
	)

	body, err := client.Fetch(context.Background())

	assert.Nil(t, body)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "celestrak fetch failed after 1 retries")
}

func TestClient_Fetch_ResponseBodyIsClosedOnNonOK(t *testing.T) {
	var requested atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requested.Add(1)
		http.Error(w, "bad request", http.StatusBadRequest)
	}))
	defer server.Close()

	client := NewClient(
		server.URL,
		time.Second,
		1,
	)

	_, err := client.Fetch(context.Background())

	require.Error(t, err)
	assert.Equal(t, int32(2), requested.Load())
}

func TestClient_fetchOnce_InvalidURL(t *testing.T) {
	client := NewClient(
		":://bad",
		time.Second,
		1,
	)

	body, err := client.fetchOnce(context.Background())

	assert.Nil(t, body)
	require.Error(t, err)
}
