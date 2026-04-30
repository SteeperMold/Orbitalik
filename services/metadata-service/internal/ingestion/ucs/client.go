package ucs

import (
	"compress/gzip"
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

	return nil, fmt.Errorf("ucs fetch failed after %d retries: %w", c.retries, lastErr)
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
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var reader = resp.Body

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}

		reader = struct {
			io.Reader
			io.Closer
		}{
			Reader: gz,
			Closer: resp.Body,
		}
	}

	return reader, nil
}
