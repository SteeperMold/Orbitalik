package satnogs

import (
	"errors"
	"io"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// Stream is not used for streaming file input
// but kept for interface consistency with other sources
func (p *Parser) Stream(r io.Reader, fn func(ingestion.Row) error) error {
	return errors.New("satnogs does not use file streaming parser")
}
