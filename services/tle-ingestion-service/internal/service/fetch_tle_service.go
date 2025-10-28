package service

import (
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
}

func NewFetchTLEService(r domain.TLERepository, client *http.Client, sourceURL string) *FetchTLEService {
	return &FetchTLEService{
		repository: r,
		client:     client,
		sourceURL:  sourceURL,
	}
}

func (s *FetchTLEService) FetchTLE(ctx context.Context) error {
	resp, err := s.client.Get(s.sourceURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("source server responded with: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	tles, err := parseTLEData(string(body))
	if err != nil {
		return fmt.Errorf("failed to parse TLE data: %w", err)
	}

	if len(tles) == 0 {
		return fmt.Errorf("no TLEs parsed from source")
	}

	err = s.repository.SaveBatch(ctx, tles)
	if err != nil {
		return fmt.Errorf("failed to save TLEs: %w", err)
	}

	return nil
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

func parseTLEData(data string) ([]*models.TLE, error) {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	var result []*models.TLE

	for i := 0; i+2 < len(lines); i += 3 {
		satName := strings.TrimSpace(lines[i])
		line1 := strings.TrimSpace(lines[i+1])
		line2 := strings.TrimSpace(lines[i+2])

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

		result = append(result, &models.TLE{
			NoradID:       noradID,
			SatelliteName: satName,
			Line1:         line1,
			Line2:         line2,
			Epoch:         epoch,
		})
	}

	return result, nil
}
