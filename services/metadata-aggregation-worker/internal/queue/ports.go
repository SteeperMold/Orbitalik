package worker

import "context"

type Aggregator interface {
	RebuildSatellite(ctx context.Context, noradID int) error
}
