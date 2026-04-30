package celestrak

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion/filesource"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Stream(r io.Reader, fn func(filesource.Row) error) error {
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

func parseLine(line string) filesource.Row {
	// SATCAT lines are ~130 chars, guard early
	if len(line) < 50 {
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

	noradID, err := strconv.Atoi(noradStr)
	if err != nil {
		return nil
	}

	return filesource.Row{
		"norad_id":    strconv.Itoa(noradID),
		"object_type": slice(7, 11),
		"status":      slice(21, 23),
		"name":        slice(23, 47),
	}
}
