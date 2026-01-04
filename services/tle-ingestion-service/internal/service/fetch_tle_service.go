package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/models"
)

type FetchTLEService struct {
	repository domain.TLERepository
	client     *http.Client
	sourceURL  string
	maxRetries int
}

func NewFetchTLEService(r domain.TLERepository, sourceURL string, requestTimeout time.Duration, maxRetries int) *FetchTLEService {
	return &FetchTLEService{
		repository: r,
		client:     &http.Client{Timeout: requestTimeout},
		sourceURL:  sourceURL,
		maxRetries: maxRetries,
	}
}

func (s *FetchTLEService) FetchTLE(ctx context.Context) error {
	var lastErr error

	for attempt := 1; attempt <= s.maxRetries; attempt++ {
		if err := s.fetchOnce(ctx); err != nil {
			lastErr = err

			if ctx.Err() != nil {
				return ctx.Err()
			}

			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}

		return nil
	}

	return fmt.Errorf("fetch TLE failed after %d attempts: %w", s.maxRetries, lastErr)
}

func (s *FetchTLEService) fetchOnce(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.sourceURL, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("source server responded with: %s", resp.Status)
	}

	tles, err := streamParseTLE(resp.Body)
	if err != nil {
		return err
	}

	if len(tles) == 0 {
		return fmt.Errorf("no TLEs parsed from source")
	}

	if err := s.repository.SaveBatch(ctx, tles); err != nil {
		return fmt.Errorf("failed to save TLEs: %w", err)
	}

	return nil
}

func streamParseTLE(r io.Reader) ([]*models.TLE, error) {
	scanner := bufio.NewScanner(r)

	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	var (
		result []*models.TLE
		lines  [3]string
		i      int
	)

	for scanner.Scan() {
		lines[i] = strings.TrimSpace(scanner.Text())
		i++

		if i < 3 {
			continue
		}

		tle, err := parseTLE(lines[0], lines[1], lines[2])
		if err != nil {
			return nil, err
		}

		result = append(result, tle)
		i = 0
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if i != 0 {
		return nil, fmt.Errorf("incomplete TLE record at end of stream")
	}

	return result, nil
}

func parseTLE(name, line1, line2 string) (*models.TLE, error) {
	if len(line1) < 32 {
		return nil, fmt.Errorf("invalid TLE line1 length")
	}

	// norad id = chars 2-7 of line1
	noradID, err := strconv.Atoi(strings.TrimSpace(line1[2:7]))
	if err != nil {
		return nil, err
	}

	// epoch = chars 18-32 of line1
	epochStr := strings.TrimSpace(line1[18:32])
	epoch, err := parseEpoch(epochStr)
	if err != nil {
		return nil, err
	}

	return &models.TLE{
		NoradID:       noradID,
		SatelliteName: name,
		Line1:         line1,
		Line2:         line2,
		Epoch:         epoch,
	}, nil
}

func parseEpoch(epochStr string) (time.Time, error) {
	year, err := strconv.Atoi(epochStr[:2])
	if err != nil {
		return time.Time{}, err
	}

	doy, err := strconv.ParseFloat(epochStr[2:], 64)
	if err != nil {
		return time.Time{}, err
	}

	if year < 57 {
		year += 2000
	} else {
		year += 1900
	}

	base := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	days := int(doy)
	frac := doy - float64(days)

	return base.AddDate(0, 0, days-1).Add(time.Duration(frac * 24 * float64(time.Hour))), nil
}
