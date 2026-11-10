package filesource

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
)

type Parser interface {
	Stream(r io.Reader, fn func(row ingestion.Row) error) error
}

type Mapper interface {
	Map(row ingestion.Row) (json.RawMessage, error)
}

type Source struct {
	name      string
	path      string
	parser    Parser
	mapper    Mapper
	batchSize int
}

func NewSource(name, path string, parser Parser, mapper Mapper, batchSize int) *Source {
	return &Source{
		name:      name,
		path:      path,
		parser:    parser,
		mapper:    mapper,
		batchSize: batchSize,
	}
}

func (s *Source) Name() string {
	return s.name
}

func (s *Source) StreamBatch(ctx context.Context, fn func([]*ingestion.SatelliteSourceRecord) error) error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
	}(f)

	now := time.Now()
	batch := make([]*ingestion.SatelliteSourceRecord, 0, s.batchSize)

	err = s.parser.Stream(f, func(row ingestion.Row) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		raw, err := s.mapper.Map(row)
		if err != nil {
			//nolint:nilerr // invalid rows are intentionally skipped
			return nil
		}

		rec := &ingestion.SatelliteSourceRecord{
			Source:    models.Source(s.Name()),
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
		return fn(batch)
	}

	return nil
}
