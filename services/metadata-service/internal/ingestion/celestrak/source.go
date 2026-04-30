package celestrak

import (
	"context"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion/filesource"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
)

type Source struct {
	client    *Client
	parser    *Parser
	mapper    *Mapper
	batchSize int
}

func NewSource(client *Client, parser *Parser, mapper *Mapper, batchSize int) *Source {
	return &Source{
		client:    client,
		parser:    parser,
		mapper:    mapper,
		batchSize: batchSize,
	}
}

func (s *Source) StreamBatch(ctx context.Context, fn func([]*ingestion.SatelliteSourceRecord) error) error {
	reader, err := s.client.Fetch(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	now := time.Now()
	batch := make([]*ingestion.SatelliteSourceRecord, 0, s.batchSize)

	err = s.parser.Stream(reader, func(r filesource.Row) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		raw, err := s.mapper.Map(r)
		if err != nil {
			return err
		}
		if raw == nil {
			return nil
		}

		rec := &ingestion.SatelliteSourceRecord{
			Source:    models.SourceCelestrak,
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

		return nil
	})

	if err != nil {
		return err
	}

	if len(batch) > 0 {
		if err := fn(batch); err != nil {
			return err
		}
	}

	return nil
}

func (s *Source) Name() string {
	return string(models.SourceCelestrak)
}
