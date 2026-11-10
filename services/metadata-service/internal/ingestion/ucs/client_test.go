package ucs

import (
	"compress/gzip"
	"context"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := NewClient(
		"https://example.com/data",
		5*time.Second,
		3,
	)

	require.NotNil(t, client)
	require.NotNil(t, client.client)
	require.NotNil(t, client.rng)

	assert.Equal(t, "https://example.com/data", client.sourceURL)
	assert.Equal(t, 3, client.retries)
	assert.Equal(t, 5*time.Second, client.client.Timeout)
}

func TestClient_Fetch_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "gzip", r.Header.Get("Accept-Encoding"))
		assert.Equal(t, "text/plain, */*;q=0.9", r.Header.Get("Accept"))
		assert.Equal(t, "en-US,en;q=0.9", r.Header.Get("Accept-Language"))
		assert.Equal(t, "Mozilla/5.0 ...", r.Header.Get("User-Agent"))
		assert.Equal(t, "keep-alive", r.Header.Get("Connection"))
		assert.Equal(t, "no-cache", r.Header.Get("Cache-Control"))

		_, err := w.Write([]byte("UCS data"))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	body, err := client.Fetch(context.Background())

	require.NoError(t, err)
	require.NotNil(t, body)
	t.Cleanup(func() {
		_ = body.Close()
	})

	data, err := io.ReadAll(body)
	require.NoError(t, err)

	assert.Equal(t, "UCS data", string(data))
}

func TestClient_fetchOnce_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("hello"))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

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

func TestClient_fetchOnce_Gzip(t *testing.T) {
	const expected = "compressed UCS data"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		_, err := gz.Write([]byte(expected))
		require.NoError(t, err)
		require.NoError(t, gz.Close())
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	body, err := client.fetchOnce(context.Background())

	require.NoError(t, err)
	require.NotNil(t, body)
	t.Cleanup(func() {
		_ = body.Close()
	})

	data, err := io.ReadAll(body)
	require.NoError(t, err)

	assert.Equal(t, expected, string(data))
}

func TestClient_fetchOnce_InvalidGzip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		_, err := w.Write([]byte("this is not gzip"))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	body, err := client.fetchOnce(context.Background())

	assert.Nil(t, body)
	require.Error(t, err)
}

func TestClient_fetchOnce_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadGateway)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	body, err := client.fetchOnce(context.Background())

	assert.Nil(t, body)
	require.Error(t, err)
	assert.EqualError(t, err, "unexpected status: 502")
}

func TestClient_Fetch_InvalidURL(t *testing.T) {
	client := NewClient("://invalid-url", time.Second, 0)

	body, err := client.Fetch(context.Background())

	assert.Nil(t, body)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ucs fetch failed after 0 retries")
}

func TestClient_fetchOnce_InvalidURL(t *testing.T) {
	client := NewClient("://invalid-url", time.Second, 0)

	body, err := client.fetchOnce(context.Background())

	assert.Nil(t, body)
	require.Error(t, err)
}

func TestClient_Fetch_ConnectionError(t *testing.T) {
	client := NewClient(
		"http://127.0.0.1:1",
		100*time.Millisecond,
		0,
	)

	body, err := client.Fetch(context.Background())

	assert.Nil(t, body)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ucs fetch failed after 0 retries")
}

func TestClient_Fetch_RetriesAfterFailure(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := attempts.Add(1)

		if attempt == 1 {
			http.Error(w, "temporary failure", http.StatusServiceUnavailable)
			return
		}

		_, err := w.Write([]byte("success after retry"))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 1)

	//nolint:gosec // G404: deterministic rng is intentional for testing
	client.rng = rand.New(rand.NewSource(1))

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
}

func TestClient_Fetch_ExhaustsRetries(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		http.Error(w, "temporary failure", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 2)

	//nolint:gosec // G404: deterministic rng is intentional for testing
	client.rng = rand.New(rand.NewSource(1))

	body, err := client.Fetch(context.Background())

	assert.Nil(t, body)
	require.Error(t, err)

	assert.Equal(t, int32(3), attempts.Load())
	assert.Contains(t, err.Error(), "ucs fetch failed after 2 retries")
	assert.Contains(t, err.Error(), "unexpected status: 503")
}

func TestClient_Fetch_ContextAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("request should not reach server")
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 3)

	body, err := client.Fetch(ctx)

	assert.Nil(t, body)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestClient_Fetch_ContextCancelledDuringRetry(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		http.Error(w, "temporary failure", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 3)

	//nolint:gosec // G404: deterministic rng is intentional for testing
	client.rng = rand.New(rand.NewSource(1))

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

func TestClient_Backoff_IsWithinRange(t *testing.T) {
	client := NewClient("", time.Second, 0)

	for attempt := 0; attempt < 10; attempt++ {
		delay := client.backoff(attempt)

		cap := baseBackoff * time.Duration(1<<attempt)
		if cap > maxBackoff {
			cap = maxBackoff
		}

		minDelay := cap / 2

		assert.GreaterOrEqual(t, delay, minDelay)
		assert.LessOrEqual(t, delay, cap)
	}
}

func TestClient_Backoff_IsCapped(t *testing.T) {
	client := NewClient("", time.Second, 0)

	for attempt := 10; attempt < 20; attempt++ {
		delay := client.backoff(attempt)

		assert.GreaterOrEqual(t, delay, maxBackoff/2)
		assert.LessOrEqual(t, delay, maxBackoff)
	}
}

func TestClient_Fetch_GzipViaFetch(t *testing.T) {
	const expected = "gzip response"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		_, err := gz.Write([]byte(expected))
		require.NoError(t, err)
		require.NoError(t, gz.Close())
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	body, err := client.Fetch(context.Background())

	require.NoError(t, err)
	require.NotNil(t, body)
	t.Cleanup(func() {
		_ = body.Close()
	})

	data, err := io.ReadAll(body)
	require.NoError(t, err)

	assert.Equal(t, expected, string(data))
}

func TestClient_Fetch_SuccessAfterGzipFailure(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := attempts.Add(1)

		if attempt == 1 {
			w.Header().Set("Content-Encoding", "gzip")
			_, err := w.Write([]byte("invalid gzip"))
			require.NoError(t, err)
			return
		}

		_, err := w.Write([]byte("valid response"))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 1)

	body, err := client.Fetch(context.Background())

	require.NoError(t, err)
	require.NotNil(t, body)
	t.Cleanup(func() {
		_ = body.Close()
	})

	data, err := io.ReadAll(body)
	require.NoError(t, err)

	assert.Equal(t, "valid response", string(data))
	assert.Equal(t, int32(2), attempts.Load())
}

func TestClient_fetchOnce_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	body, err := client.fetchOnce(ctx)

	assert.Nil(t, body)
	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}
