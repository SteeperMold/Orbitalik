//go:build integration
// +build integration

package queue

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()

	mr := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	t.Cleanup(func() {
		_ = rdb.Close()
	})

	return mr, rdb
}

func TestDirtyMarker_MarkDirty_AddsNewIDs(t *testing.T) {
	_, rdb := newTestRedis(t)

	marker := NewDirtyMarker(
		rdb,
		"dirty-satellites",
		"satellite-events",
	)

	err := marker.MarkDirty(context.Background(), []int{25544, 12345})

	require.NoError(t, err)

	members, err := rdb.SMembers(
		context.Background(),
		"dirty-satellites",
	).Result()

	require.NoError(t, err)

	assert.ElementsMatch(
		t,
		[]string{"25544", "12345"},
		members,
	)

	streams, err := rdb.XRange(
		context.Background(),
		"satellite-events",
		"-",
		"+",
	).Result()

	require.NoError(t, err)
	require.Len(t, streams, 2)

	assert.Equal(t, "25544", streams[0].Values["norad_id"])
	assert.Equal(t, "12345", streams[1].Values["norad_id"])
}

func TestDirtyMarker_MarkDirty_DeduplicatesInput(t *testing.T) {
	_, rdb := newTestRedis(t)

	marker := NewDirtyMarker(
		rdb,
		"dirty-satellites",
		"satellite-events",
	)

	err := marker.MarkDirty(
		context.Background(),
		[]int{25544, 25544, 25544, 12345, 12345},
	)

	require.NoError(t, err)

	members, err := rdb.SMembers(
		context.Background(),
		"dirty-satellites",
	).Result()

	require.NoError(t, err)

	assert.ElementsMatch(
		t,
		[]string{"25544", "12345"},
		members,
	)

	streams, err := rdb.XRange(
		context.Background(),
		"satellite-events",
		"-",
		"+",
	).Result()

	require.NoError(t, err)
	require.Len(t, streams, 2)

	assert.Equal(t, "25544", streams[0].Values["norad_id"])
	assert.Equal(t, "12345", streams[1].Values["norad_id"])
}

func TestDirtyMarker_MarkDirty_SkipsAlreadyQueuedIDs(t *testing.T) {
	_, rdb := newTestRedis(t)

	ctx := context.Background()

	require.NoError(
		t,
		rdb.SAdd(ctx, "dirty-satellites", "25544").Err(),
	)

	marker := NewDirtyMarker(
		rdb,
		"dirty-satellites",
		"satellite-events",
	)

	err := marker.MarkDirty(ctx, []int{25544})

	require.NoError(t, err)

	streams, err := rdb.XRange(
		ctx,
		"satellite-events",
		"-",
		"+",
	).Result()

	require.NoError(t, err)
	assert.Empty(t, streams)

	members, err := rdb.SMembers(
		ctx,
		"dirty-satellites",
	).Result()

	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"25544"}, members)
}

func TestDirtyMarker_MarkDirty_MixedExistingAndNewIDs(t *testing.T) {
	_, rdb := newTestRedis(t)

	ctx := context.Background()

	require.NoError(
		t,
		rdb.SAdd(ctx, "dirty-satellites", "25544").Err(),
	)

	marker := NewDirtyMarker(
		rdb,
		"dirty-satellites",
		"satellite-events",
	)

	err := marker.MarkDirty(ctx, []int{25544, 12345})

	require.NoError(t, err)

	members, err := rdb.SMembers(
		ctx,
		"dirty-satellites",
	).Result()

	require.NoError(t, err)

	assert.ElementsMatch(
		t,
		[]string{"25544", "12345"},
		members,
	)

	streams, err := rdb.XRange(
		ctx,
		"satellite-events",
		"-",
		"+",
	).Result()

	require.NoError(t, err)
	require.Len(t, streams, 1)

	assert.Equal(t, "12345", streams[0].Values["norad_id"])
}

func TestDirtyMarker_MarkDirty_EmptyInput(t *testing.T) {
	_, rdb := newTestRedis(t)

	marker := NewDirtyMarker(
		rdb,
		"dirty-satellites",
		"satellite-events",
	)

	err := marker.MarkDirty(
		context.Background(),
		nil,
	)

	require.NoError(t, err)

	members, err := rdb.SMembers(
		context.Background(),
		"dirty-satellites",
	).Result()

	require.NoError(t, err)
	assert.Empty(t, members)

	streams, err := rdb.XRange(
		context.Background(),
		"satellite-events",
		"-",
		"+",
	).Result()

	require.NoError(t, err)
	assert.Empty(t, streams)
}

func TestDirtyMarker_MarkDirty_ReturnsSAddError(t *testing.T) {
	mr, rdb := newTestRedis(t)

	ctx := context.Background()

	require.NoError(
		t,
		rdb.Set(ctx, "dirty-satellites", "not-a-set", 0).Err(),
	)

	marker := NewDirtyMarker(
		rdb,
		"dirty-satellites",
		"satellite-events",
	)

	err := marker.MarkDirty(ctx, []int{25544})

	require.Error(t, err)

	mr.FlushAll()
}

func TestDirtyMarker_MarkDirty_ContextCancellation(t *testing.T) {
	_, rdb := newTestRedis(t)

	marker := NewDirtyMarker(
		rdb,
		"dirty-satellites",
		"satellite-events",
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := marker.MarkDirty(ctx, []int{25544})

	require.Error(t, err)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestDirtyMarker_MarkDirty_PreservesInputOrderInStream(t *testing.T) {
	_, rdb := newTestRedis(t)

	marker := NewDirtyMarker(
		rdb,
		"dirty-satellites",
		"satellite-events",
	)

	err := marker.MarkDirty(
		context.Background(),
		[]int{3, 1, 2},
	)

	require.NoError(t, err)

	streams, err := rdb.XRange(
		context.Background(),
		"satellite-events",
		"-",
		"+",
	).Result()

	require.NoError(t, err)
	require.Len(t, streams, 3)

	assert.Equal(t, "3", streams[0].Values["norad_id"])
	assert.Equal(t, "1", streams[1].Values["norad_id"])
	assert.Equal(t, "2", streams[2].Values["norad_id"])
}
