package repository

import "github.com/SteeperMold/Orbitalik/tle-ingestion-service/internal/models"

type tleCopySource struct {
	tles []*models.TLE
	idx  int
}

func newTLECopySource(tles []*models.TLE) *tleCopySource {
	return &tleCopySource{
		tles: tles,
		idx:  0,
	}
}

func (s *tleCopySource) Next() bool {
	if s.idx >= len(s.tles) {
		return false
	}
	s.idx++
	return true
}

func (s *tleCopySource) Values() ([]any, error) {
	t := s.tles[s.idx-1]

	return []any{
		t.NoradID,
		t.SatelliteName,
		t.Line1,
		t.Line2,
		t.Epoch,
	}, nil
}

func (s *tleCopySource) Err() error {
	return nil
}
