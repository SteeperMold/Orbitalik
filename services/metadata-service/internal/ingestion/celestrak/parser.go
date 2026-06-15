package celestrak

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Stream(r io.Reader, fn func(ingestion.Row) error) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		rec := parseLine(line)
		if rec == nil {
			continue
		}

		if err := fn(rec); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func parseLine(line string) ingestion.Row {
	if len(line) < 80 {
		return nil
	}

	slice := func(start, end int) string {
		if start >= len(line) {
			return ""
		}
		if end > len(line) {
			end = len(line)
		}
		return strings.TrimSpace(line[start:end])
	}

	noradStr := slice(13, 18)
	if noradStr == "" {
		return nil
	}

	if _, err := strconv.Atoi(noradStr); err != nil {
		return nil
	}

	return ingestion.Row{
		"norad_id":    noradStr,
		"cospar_id":   slice(0, 11),
		"name":        slice(23, 47),
		"flags":       slice(21, 23), // *, D, etc
		"owner":       slice(49, 54),
		"launch_date": slice(56, 66),
		"launch_site": slice(68, 73),
		"decay_date":  slice(75, 85),
		"period":      slice(87, 94),
		"inclination": slice(96, 101),
		"apogee":      slice(103, 109),
		"perigee":     slice(111, 117),
	}
}
