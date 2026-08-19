package satnogs

import (
	"context"
	"encoding/json"
	"io"
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
		"https://example.com/api",
		5*time.Second,
		3,
	)

	require.NotNil(t, client)
	require.NotNil(t, client.client)
	require.NotNil(t, client.rng)

	assert.Equal(t, "https://example.com/api", client.baseURL)
	assert.Equal(t, 3, client.retries)
	assert.Equal(t, 5*time.Second, client.client.Timeout)
}

func TestClient_FetchTransmitters_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/transmitters/", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "orbitalik-satnogs-client", r.Header.Get("User-Agent"))

		w.Header().Set("Content-Type", "application/json")

		_, err := w.Write([]byte(`[
			{
				"sat_id": "ISS",
				"norad_cat_id": 25544,
				"description": "ISS transmitter",
				"mode": "VHF",
				"alive": true,
				"downlink_low": 145.8,
				"downlink_high": 145.8,
				"uplink_low": 145.2,
				"uplink_high": 145.2,
				"baud": 9600
			}
		]`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	result, err := client.FetchTransmitters(context.Background())

	require.NoError(t, err)
	require.Len(t, result, 1)

	assert.Equal(t, "ISS", result[0].SatID)
	assert.Equal(t, 25544, result[0].NoradCatID)
	assert.Equal(t, "ISS transmitter", result[0].Description)
	assert.Equal(t, "VHF", result[0].Mode)
	assert.True(t, result[0].Alive)
	assert.Equal(t, 145.8, result[0].DownlinkLow)
	assert.Equal(t, 9600.0, result[0].Baud)
}

func TestClient_FetchTransmitters_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`[]`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	result, err := client.FetchTransmitters(context.Background())

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestClient_FetchTransmitters_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	result, err := client.FetchTransmitters(context.Background())

	assert.Nil(t, result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected status: 500")
}

func TestClient_FetchTransmitters_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`not valid json`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	result, err := client.FetchTransmitters(context.Background())

	assert.Nil(t, result)
	require.Error(t, err)
}

func TestClient_getJSON_InvalidURL(t *testing.T) {
	client := NewClient("://invalid-url", time.Second, 0)

	var target []Transmitter

	err := client.getJSON(
		context.Background(),
		"/transmitters/",
		&target,
	)

	require.Error(t, err)
}

func TestClient_fetchOnce_InvalidURL(t *testing.T) {
	client := NewClient("://invalid-url", time.Second, 0)

	var target []Transmitter

	err := client.fetchOnce(
		context.Background(),
		"/transmitters/",
		&target,
	)

	require.Error(t, err)
}

func TestClient_getJSON_RetriesAfterFailure(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := attempts.Add(1)

		if attempt == 1 {
			http.Error(w, "temporary failure", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`[]`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 1)

	var target []Transmitter

	err := client.getJSON(
		context.Background(),
		"/transmitters/",
		&target,
	)

	require.NoError(t, err)
	assert.Equal(t, int32(2), attempts.Load())
	assert.Empty(t, target)
}

func TestClient_getJSON_ExhaustsRetries(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		http.Error(w, "temporary failure", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 2)

	var target []Transmitter

	err := client.getJSON(
		context.Background(),
		"/transmitters/",
		&target,
	)

	require.Error(t, err)
	assert.Equal(t, int32(3), attempts.Load())
	assert.Contains(t, err.Error(), "satnogs request failed after 2 retries")
	assert.Contains(t, err.Error(), "unexpected status: 503")
}

func TestClient_getJSON_ContextAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("request should not reach server")
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 3)

	var target []Transmitter

	err := client.getJSON(ctx, "/transmitters/", &target)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestClient_getJSON_ContextCancelledDuringRetry(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		http.Error(w, "temporary failure", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 3)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	var target []Transmitter

	err := client.getJSON(ctx, "/transmitters/", &target)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, int32(1), attempts.Load())
}

func TestClient_fetchOnce_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/custom-path", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"sat_id":"TEST","norad_cat_id":123}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	var target Transmitter

	err := client.fetchOnce(
		context.Background(),
		"/custom-path",
		&target,
	)

	require.NoError(t, err)

	assert.Equal(t, "TEST", target.SatID)
	assert.Equal(t, 123, target.NoradCatID)
}

func TestClient_fetchOnce_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	var target []Transmitter

	err := client.fetchOnce(
		context.Background(),
		"/transmitters/",
		&target,
	)

	require.Error(t, err)
	assert.EqualError(t, err, "unexpected status: 404")
}

func TestClient_fetchOnce_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`{invalid`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	var target Transmitter

	err := client.fetchOnce(
		context.Background(),
		"/transmitters/",
		&target,
	)

	require.Error(t, err)
}

func TestClient_fetchOnce_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var target []Transmitter

	err := client.fetchOnce(ctx, "/transmitters/", &target)

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestClient_Backoff_IsWithinExpectedRange(t *testing.T) {
	client := NewClient("", time.Second, 0)

	for attempt := 0; attempt < 8; attempt++ {
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

		assert.LessOrEqual(t, delay, maxBackoff)
		assert.GreaterOrEqual(t, delay, maxBackoff/2)
	}
}

func TestClient_Backoff_ProducesDifferentValues(t *testing.T) {
	client := NewClient("", time.Second, 0)

	values := make(map[time.Duration]struct{})

	for i := 0; i < 20; i++ {
		values[client.backoff(3)] = struct{}{}
	}

	assert.Greater(t, len(values), 1)
}

func TestTransmitter_JSONTags(t *testing.T) {
	input := `{
		"sat_id": "AO-91",
		"norad_cat_id": 43017,
		"description": "Amateur satellite",
		"mode": "V/U",
		"alive": true,
		"downlink_low": 145.8,
		"downlink_high": 145.8,
		"uplink_low": 435.25,
		"uplink_high": 435.3,
		"baud": 9600
	}`

	var transmitter Transmitter

	require.NoError(t, json.Unmarshal([]byte(input), &transmitter))

	assert.Equal(t, "AO-91", transmitter.SatID)
	assert.Equal(t, 43017, transmitter.NoradCatID)
	assert.Equal(t, "Amateur satellite", transmitter.Description)
	assert.Equal(t, "V/U", transmitter.Mode)
	assert.True(t, transmitter.Alive)
	assert.Equal(t, 145.8, transmitter.DownlinkLow)
	assert.Equal(t, 435.25, transmitter.UplinkLow)
	assert.Equal(t, 9600.0, transmitter.Baud)
}

func TestClient_getJSON_ConnectionError(t *testing.T) {
	client := NewClient(
		"http://127.0.0.1:1",
		100*time.Millisecond,
		0,
	)

	var target []Transmitter

	err := client.getJSON(
		context.Background(),
		"/transmitters/",
		&target,
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "satnogs request failed after 0 retries")
}

func TestClient_fetchOnce_ResponseBodyIsClosed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `[]`)
		require.NoError(t, err)
	}))
	defer server.Close()

	client := NewClient(server.URL, time.Second, 0)

	var target []Transmitter

	err := client.fetchOnce(
		context.Background(),
		"/transmitters/",
		&target,
	)

	require.NoError(t, err)
}
