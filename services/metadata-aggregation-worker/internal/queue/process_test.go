package worker

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestParseNoradID(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    int
		wantErr bool
	}{
		{
			name:  "string",
			input: "25544",
			want:  25544,
		},
		{
			name:  "int",
			input: 25544,
			want:  25544,
		},
		{
			name:  "int64",
			input: int64(25544),
			want:  25544,
		},
		{
			name:    "invalid string",
			input:   "not-a-number",
			wantErr: true,
		},
		{
			name:    "unsupported type",
			input:   float64(25544),
			wantErr: true,
		},
		{
			name:    "nil",
			input:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseNoradID(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWorker_ProcessMessage_MissingNoradID(t *testing.T) {
	w := newTestWorker(nil, nil, nil)

	msg := redis.XMessage{
		ID:     "1-0",
		Values: map[string]any{},
	}

	err := w.processMessage(context.Background(), msg)

	require.EqualError(t, err, "missing norad_id")
}

func TestWorker_ProcessMessage_InvalidNoradID(t *testing.T) {
	w := newTestWorker(nil, nil, nil)

	msg := redis.XMessage{
		ID: "1-0",
		Values: map[string]any{
			"norad_id": "not-a-number",
		},
	}

	err := w.processMessage(context.Background(), msg)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid syntax")
}

func TestWorker_ProcessMessage_AggregatorError(t *testing.T) {
	rdb := new(mockRedis)
	aggregator := new(mockAggregator)

	aggregationErr := errors.New("aggregation failed")

	aggregator.
		On("RebuildSatellite", mock.Anything, 25544).
		Return(aggregationErr).
		Once()

	w := newTestWorker(rdb, aggregator, nil)

	msg := redis.XMessage{
		ID: "1-0",
		Values: map[string]any{
			"norad_id": "25544",
		},
	}

	err := w.processMessage(context.Background(), msg)

	assert.ErrorIs(t, err, aggregationErr)

	aggregator.AssertExpectations(t)
	rdb.AssertNotCalled(t, "XAck", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	rdb.AssertNotCalled(t, "SRem", mock.Anything, mock.Anything, mock.Anything)
}

func TestWorker_ProcessMessage_XAckError(t *testing.T) {
	rdb := new(mockRedis)
	aggregator := new(mockAggregator)

	ackErr := errors.New("ack failed")

	aggregator.
		On("RebuildSatellite", mock.Anything, 25544).
		Return(nil).
		Once()

	ackCmd := redis.NewIntCmd(context.Background())
	ackCmd.SetErr(ackErr)

	rdb.
		On(
			"XAck",
			mock.Anything,
			"satellite-events",
			"metadata-workers",
			[]string{"1-0"},
		).
		Return(ackCmd).
		Once()

	w := newTestWorker(rdb, aggregator, nil)

	msg := redis.XMessage{
		ID: "1-0",
		Values: map[string]any{
			"norad_id": int64(25544),
		},
	}

	err := w.processMessage(context.Background(), msg)

	assert.ErrorIs(t, err, ackErr)

	aggregator.AssertExpectations(t)
	rdb.AssertExpectations(t)

	rdb.AssertNotCalled(t, "SRem", mock.Anything, mock.Anything, mock.Anything)
}

func TestWorker_ProcessMessage_SRemError(t *testing.T) {
	rdb := new(mockRedis)
	aggregator := new(mockAggregator)

	sremErr := errors.New("remove dirty entry failed")

	aggregator.
		On("RebuildSatellite", mock.Anything, 25544).
		Return(nil).
		Once()

	ackCmd := redis.NewIntCmd(context.Background())
	ackCmd.SetVal(1)

	sremCmd := redis.NewIntCmd(context.Background())
	sremCmd.SetErr(sremErr)

	rdb.
		On(
			"XAck",
			mock.Anything,
			"satellite-events",
			"metadata-workers",
			[]string{"1-0"},
		).
		Return(ackCmd).
		Once()

	rdb.
		On(
			"SRem",
			mock.Anything,
			"dirty-satellites",
			[]any{"25544"},
		).
		Return(sremCmd).
		Once()

	w := newTestWorker(rdb, aggregator, nil)

	msg := redis.XMessage{
		ID: "1-0",
		Values: map[string]any{
			"norad_id": 25544,
		},
	}

	err := w.processMessage(context.Background(), msg)

	assert.ErrorIs(t, err, sremErr)

	aggregator.AssertExpectations(t)
	rdb.AssertExpectations(t)
}

func TestWorker_ProcessMessage_Success(t *testing.T) {
	rdb := new(mockRedis)
	aggregator := new(mockAggregator)

	aggregator.
		On("RebuildSatellite", mock.Anything, 25544).
		Return(nil).
		Once()

	ackCmd := redis.NewIntCmd(context.Background())
	ackCmd.SetVal(1)

	sremCmd := redis.NewIntCmd(context.Background())
	sremCmd.SetVal(1)

	rdb.
		On(
			"XAck",
			mock.Anything,
			"satellite-events",
			"metadata-workers",
			[]string{"123-0"},
		).
		Return(ackCmd).
		Once()

	rdb.
		On(
			"SRem",
			mock.Anything,
			"dirty-satellites",
			[]any{"25544"},
		).
		Return(sremCmd).
		Once()

	w := newTestWorker(rdb, aggregator, nil)

	msg := redis.XMessage{
		ID: "123-0",
		Values: map[string]any{
			"norad_id": "25544",
		},
	}

	err := w.processMessage(context.Background(), msg)

	require.NoError(t, err)

	aggregator.AssertExpectations(t)
	rdb.AssertExpectations(t)
}
