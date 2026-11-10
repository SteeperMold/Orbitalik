package celestrak

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const (
	baseBackoff = time.Second
	maxBackoff  = 30 * time.Second
)

type Client struct {
	client    *http.Client
	sourceURL string
	retries   int
	rng       *rand.Rand
}

func NewClient(url string, timeout time.Duration, retries int) *Client {
	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
		sourceURL: url,
		retries:   retries,
		//nolint:gosec // math/rand is used only for non-security backoff jitter.
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (c *Client) Fetch(ctx context.Context) (io.ReadCloser, error) {
	var lastErr error

	for attempt := 0; attempt <= c.retries; attempt++ {
		body, err := c.fetchOnce(ctx)
		if err == nil {
			return body, nil
		}

		lastErr = err

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(c.backoff(attempt)):
		}
	}

	return nil, fmt.Errorf("celestrak fetch failed after %d retries: %w", c.retries, lastErr)
}

func (c *Client) fetchOnce(ctx context.Context) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.sourceURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", "Mozilla/5.0 ...")
	req.Header.Set("Accept", "text/plain, */*;q=0.9")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (c *Client) backoff(attempt int) time.Duration {
	delay := baseBackoff * time.Duration(1<<attempt)

	if delay > maxBackoff {
		delay = maxBackoff
	}

	minDelay := delay / 2
	jitter := time.Duration(c.rng.Int63n(int64(delay-minDelay) + 1))

	return minDelay + jitter
}
