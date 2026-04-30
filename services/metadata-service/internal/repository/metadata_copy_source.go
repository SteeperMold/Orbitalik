package repository

import (
	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/models"
	"github.com/jackc/pgx/v5"
)

type rawCopySource struct {
	data []*models.SatelliteIngestRecord
	idx  int
}

func newRawMetadataCopySource(data []*models.SatelliteIngestRecord) pgx.CopyFromSource {
	return &rawCopySource{data: data}
}

func (s *rawCopySource) Next() bool {
	return s.idx < len(s.data)
}

func (s *rawCopySource) Values() ([]any, error) {
	r := s.data[s.idx]
	s.idx++

	payload, err := r.Payload.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return []any{
		r.NoradID,
		r.CosparID,
		r.Source,
		payload,
		r.FetchedAt,
	}, nil
}

func (s *rawCopySource) Err() error {
	return nil
}
