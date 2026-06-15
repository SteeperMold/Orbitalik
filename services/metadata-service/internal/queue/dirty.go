package queue

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type DirtyMarker struct {
	rdb      *redis.Client
	setKey   string
	queueKey string
}

func NewDirtyMarker(rdb *redis.Client, dirtySetKey, queueKey string) *DirtyMarker {
	return &DirtyMarker{
		rdb:      rdb,
		setKey:   dirtySetKey,
		queueKey: queueKey,
	}
}

func (m *DirtyMarker) MarkDirty(ctx context.Context, noradIDs []int) error {
	seen := make(map[int]struct{})

	for _, id := range noradIDs {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}

		norad := strconv.Itoa(id)

		added, err := m.rdb.SAdd(ctx, m.setKey, norad).Result()
		if err != nil {
			return err
		}

		// already queued
		if added == 0 {
			continue
		}

		err = m.rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: m.queueKey,
			Values: map[string]any{
				"norad_id": id,
			},
		}).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
