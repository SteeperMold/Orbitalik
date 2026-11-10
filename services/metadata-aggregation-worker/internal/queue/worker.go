package worker

import (
	"context"
	"errors"
	"strings"
	"time"

	applog "github.com/SteeperMold/Orbitalik/common/go/log"
	"github.com/redis/go-redis/v9"
)

type Worker struct {
	rdb          RedisClient
	aggregator   Aggregator
	logger       applog.Logger
	consumerName string
	streamName   string
	groupName    string
	dirtySet     string
	streamsCount int64
	streamsBlock time.Duration
	retryDelay   time.Duration
}

func New(
	rdb RedisClient,
	aggregator Aggregator,
	logger applog.Logger,
	consumerName,
	queueKey,
	groupName,
	dirtySet string,
	streamsCount int64,
	streamsBlock time.Duration,
	retryDelay time.Duration,
) *Worker {

	return &Worker{
		rdb:          rdb,
		aggregator:   aggregator,
		logger:       logger,
		consumerName: consumerName,
		streamName:   queueKey,
		groupName:    groupName,
		dirtySet:     dirtySet,
		streamsCount: streamsCount,
		streamsBlock: streamsBlock,
		retryDelay:   retryDelay,
	}
}

func (w *Worker) Init(ctx context.Context) error {
	err := w.rdb.XGroupCreateMkStream(ctx, w.streamName, w.groupName, "0").Err()

	if err != nil && !isBusyGroup(err) {
		return err
	}

	return nil
}

func isBusyGroup(err error) bool {
	return err != nil && strings.Contains(err.Error(), "BUSYGROUP")
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		streams, err := w.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    w.groupName,
			Consumer: w.consumerName,
			Streams:  []string{w.streamName, ">"},
			Count:    w.streamsCount,
			Block:    w.streamsBlock,
		}).Result()

		if err != nil {
			if errors.Is(err, redis.Nil) {
				continue
			}

			w.logger.Error(
				"failed to read stream",
				applog.NewErrorField(err),
			)

			if err := sleep(ctx, w.retryDelay); err != nil {
				return err
			}

			continue
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				if err := w.processMessage(ctx, msg); err != nil {
					w.logger.Error(
						"failed to process message",
						applog.NewField("message_id", msg.ID),
						applog.NewErrorField(err),
					)
				}
			}
		}
	}
}

func sleep(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
