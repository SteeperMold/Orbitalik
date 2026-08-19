package ucs

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"strings"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Stream(r io.Reader, fn func(ingestion.Row) error) error {
	reader := csv.NewReader(r)
	reader.Comma = '\t'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	idx := buildIndex(header)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			return nil
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		if err != nil {
			continue
		}

		rec := mapRow(idx, row)
		if rec == nil {
			continue
		}

		if err := fn(rec); err != nil {
			return err
		}
	}
}

func buildIndex(header []string) map[string]int {
	m := make(map[string]int)

	for i, h := range header {
		key := normalizeHeader(h)
		m[key] = i
	}

	return m
}

func normalizeHeader(h string) string {
	return strings.ToLower(strings.TrimSpace(h))
}

func mapRow(idx map[string]int, row []string) ingestion.Row {
	get := func(key string) string {
		i, ok := idx[key]
		if !ok || i >= len(row) {
			return ""
		}
		return clean(row[i])
	}

	norad := get("norad number")
	if norad == "" {
		return nil
	}

	return ingestion.Row{
		"norad_id":       norad,
		"name":           get("current official name of satellite"),
		"aliases":        get("name of satellite, alternate names"),
		"operator":       get("operator/owner"),
		"owner":          get("country of operator/owner"),
		"users":          get("users"),
		"purpose":        get("purpose"),
		"orbit_class":    get("class of orbit"),
		"launch_date":    get("date of launch"),
		"launch_site":    get("launch site"),
		"launch_vehicle": get("launch vehicle"),
		"cospar":         get("cospar number"),
	}
}

func clean(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, `"`)
	return s
}
