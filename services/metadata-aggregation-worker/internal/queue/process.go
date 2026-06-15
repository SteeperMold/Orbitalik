package worker

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func (w *Worker) processMessage(ctx context.Context, msg redis.XMessage) error {
	rawID, ok := msg.Values["norad_id"]
	if !ok {
		return fmt.Errorf("missing norad_id")
	}

	noradID, err := parseNoradID(rawID)
	if err != nil {
		return err
	}

	if err := w.aggregator.RebuildSatellite(ctx, noradID); err != nil {
		return err
	}

	if err := w.rdb.XAck(ctx, w.streamName, w.groupName, msg.ID).Err(); err != nil {
		return err
	}

	if err := w.rdb.SRem(ctx, w.dirtySet, strconv.Itoa(noradID)).Err(); err != nil {
		return err
	}

	return nil
}

func parseNoradID(v any) (int, error) {
	switch t := v.(type) {
	case string:
		return strconv.Atoi(t)
	case int:
		return t, nil
	case int64:
		return int(t), nil
	default:
		return 0, fmt.Errorf("invalid norad_id type")
	}
}
