package celestrak

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client    *http.Client
	sourceURL string
	retries   int
}

func NewClient(url string, timeout time.Duration, retries int) *Client {
	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
		sourceURL: url,
		retries:   retries,
	}
}

func (c *Client) Fetch(ctx context.Context) (io.ReadCloser, error) {
	var lastErr error

	for attempt := 1; attempt <= c.retries; attempt++ {
		body, err := c.fetchOnce(ctx)
		if err == nil {
			return body, nil
		}

		lastErr = err

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(attempt) * time.Second):
		}
	}

	return nil, fmt.Errorf("celestrak fetch failed after %d retries: %w", c.retries, lastErr)
}

func (c *Client) fetchOnce(ctx context.Context) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.sourceURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
