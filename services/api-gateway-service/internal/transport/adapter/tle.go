package adapter

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/gen/tlepb"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/models"
	"google.golang.org/grpc"
)

type TLEAdapter struct {
	client tlepb.TleServiceClient
}

func NewTLEAdapter(conn *grpc.ClientConn) *TLEAdapter {
	return &TLEAdapter{
		client: tlepb.NewTleServiceClient(conn),
	}
}

func (a *TLEAdapter) GetTle(ctx context.Context, id domain.SatelliteIdentifier) (*models.Tle, error) {
	identifier, err := satelliteIdentifierToTLEProto(id)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.GetTle(ctx, &tlepb.GetTleRequest{
		Identifier: identifier,
	})
	if err != nil {
		return nil, fmt.Errorf("TLE GetTle: %w", err)
	}
	if resp == nil || resp.GetTle() == nil {
		return nil, fmt.Errorf("TLE GetTle: empty response")
	}

	fetchedAt := time.Now()

	return toDomainTLE(resp.GetTle(), fetchedAt), nil
}

func (a *TLEAdapter) ListTles(ctx context.Context) ([]models.Tle, error) {
	resp, err := a.client.ListTles(ctx, &tlepb.ListTlesRequest{})
	if err != nil {
		return nil, fmt.Errorf("TLE ListTles: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("TLE ListTles: empty response")
	}

	fetchedAt := time.Now()

	orbits := make([]models.Tle, 0, len(resp.GetTles()))
	for _, tle := range resp.GetTles() {
		if tle == nil {
			continue
		}
		orbits = append(orbits, *toDomainTLE(tle, fetchedAt))
	}

	return orbits, nil
}

func satelliteIdentifierToTLEProto(id domain.SatelliteIdentifier) (*tlepb.SatelliteIdentifier, error) {
	switch {
	case id.NORADID != nil && id.Name == nil:
		if *id.NORADID <= 0 || uint64(*id.NORADID) > uint64(math.MaxUint32) {
			return nil, fmt.Errorf("invalid NORAD ID: %d", *id.NORADID)
		}

		return &tlepb.SatelliteIdentifier{
			Kind: &tlepb.SatelliteIdentifier_NoradId{
				NoradId: uint32(*id.NORADID),
			},
		}, nil

	case id.Name != nil && id.NORADID == nil:
		name := strings.TrimSpace(*id.Name)
		if name == "" {
			return nil, fmt.Errorf("satellite name cannot be empty")
		}

		return &tlepb.SatelliteIdentifier{
			Kind: &tlepb.SatelliteIdentifier_SatelliteName{
				SatelliteName: name,
			},
		}, nil

	case id.NORADID != nil && id.Name != nil:
		return nil, fmt.Errorf(
			"provide either NORAD ID or satellite name, not both",
		)

	default:
		return nil, fmt.Errorf("satellite NORAD ID or name is required")
	}
}

func toDomainTLE(tle *tlepb.Tle, fetchedAt time.Time) *models.Tle {
	if tle == nil {
		return nil
	}

	return &models.Tle{
		NoradID:       tle.NoradId,
		SatelliteName: tle.SatelliteName,
		Line1:         tle.Line1,
		Line2:         tle.Line2,
		Epoch:         timestampToTime(tle.Epoch),
		FetchedAt:     fetchedAt,
	}
}
