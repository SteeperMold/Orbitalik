package satnogs

import (
	"context"
	"strconv"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
)

type Source struct {
	client    *Client
	mapper    *Mapper
	batchSize int
}

func NewSource(client *Client, mapper *Mapper, batchSize int) *Source {
	return &Source{
		client:    client,
		mapper:    mapper,
		batchSize: batchSize,
	}
}

func (s *Source) StreamBatch(
	ctx context.Context,
	fn func([]*ingestion.SatelliteSourceRecord) error,
) error {

	tx, err := s.client.FetchTransmitters(ctx)
	if err != nil {
		return err
	}

	now := time.Now()
	batch := make([]*ingestion.SatelliteSourceRecord, 0, s.batchSize)

	for _, t := range tx {
		if !t.Alive {
			continue
		}

		norad := strconv.Itoa(t.NoradCatID)

		row := ingestion.Row{
			"norad_id":      norad,
			"mode":          t.Mode,
			"downlink_low":  strconv.FormatFloat(t.DownlinkLow, 'f', 2, 64),
			"downlink_high": strconv.FormatFloat(t.DownlinkHigh, 'f', 2, 64),
			"uplink_low":    strconv.FormatFloat(t.UplinkLow, 'f', 2, 64),
			"uplink_high":   strconv.FormatFloat(t.UplinkHigh, 'f', 2, 64),
			"description":   t.Description,
			"baud":          strconv.FormatFloat(t.Baud, 'f', 2, 64),
		}

		raw, err := s.mapper.Map(row)
		if err != nil {
			continue
		}

		rec := &ingestion.SatelliteSourceRecord{
			Source:    models.SourceSatNOGS,
			FetchedAt: now,
			Raw:       raw,
		}

		batch = append(batch, rec)

		if len(batch) >= s.batchSize {
			if err := fn(batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		return fn(batch)
	}

	return nil
}

func (s *Source) Name() string {
	return string(models.SourceSatNOGS)
}
