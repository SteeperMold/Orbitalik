package satnogs

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	baseBackoff = time.Second
	maxBackoff  = 30 * time.Second
)

type Client struct {
	client  *http.Client
	baseURL string
	retries int
	rng     *rand.Rand
}

func NewClient(baseURL string, timeout time.Duration, retries int) *Client {
	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
		retries: retries,
		//nolint:gosec // math/rand is used only for non-security backoff jitter.
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type Transmitter struct {
	SatID       string `json:"sat_id"`
	NoradCatID  int    `json:"norad_cat_id"`
	Description string `json:"description"`
	Mode        string `json:"mode"`
	Alive       bool   `json:"alive"`

	DownlinkLow  float64 `json:"downlink_low"`
	DownlinkHigh float64 `json:"downlink_high"`
	UplinkLow    float64 `json:"uplink_low"`
	UplinkHigh   float64 `json:"uplink_high"`
	Baud         float64 `json:"baud"`
}

func (c *Client) FetchTransmitters(ctx context.Context) ([]Transmitter, error) {
	var out []Transmitter

	err := c.getJSON(ctx, "/transmitters/", &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) getJSON(ctx context.Context, path string, target any) error {
	var lastErr error

	for attempt := 0; attempt <= c.retries; attempt++ {
		err := c.fetchOnce(ctx, path, target)
		if err == nil {
			return nil
		}

		lastErr = err

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(c.backoff(attempt)):
		}
	}

	return fmt.Errorf("satnogs request failed after %d retries: %w", c.retries, lastErr)
}

func (c *Client) fetchOnce(ctx context.Context, path string, target any) error {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "orbitalik-satnogs-client")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(target); err != nil {
		return err
	}

	return nil
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
