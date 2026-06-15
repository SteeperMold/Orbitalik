package domain

import (
	"context"
)

type IngestionService interface {
	IngestMetadata(ctx context.Context) error
}
