package celestrak

import (
	"errors"
	"strings"
	"testing"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeDefaultCelestrakLine() string {
	line := make([]byte, 117)
	for i := range line {
		line[i] = ' '
	}

	put := func(start int, value string) {
		copy(line[start:], value)
	}

	put(0, "1998-067A")
	put(13, "25544")
	put(21, "D")
	put(23, "ISS")
	put(49, "US")
	put(56, "1998-11-20")
	put(68, "KSC")
	put(75, "2025-01-01")
	put(87, "92.6")
	put(96, "51.6")
	put(103, "420")
	put(111, "410")

	return string(line)
}

func TestParseLine(t *testing.T) {
	line := makeDefaultCelestrakLine()

	got := parseLine(line)

	require.NotNil(t, got)

	assert.Equal(t, "25544", got["norad_id"])
	assert.Equal(t, "1998-067A", got["cospar_id"])
	assert.Equal(t, "ISS", got["name"])
	assert.Equal(t, "D", got["flags"])
	assert.Equal(t, "US", got["owner"])
	assert.Equal(t, "1998-11-20", got["launch_date"])
	assert.Equal(t, "KSC", got["launch_site"])
	assert.Equal(t, "2025-01-01", got["decay_date"])
	assert.Equal(t, "92.6", got["period"])
	assert.Equal(t, "51.6", got["inclination"])
	assert.Equal(t, "420", got["apogee"])
	assert.Equal(t, "410", got["perigee"])
}

func TestParseLine_TooShort(t *testing.T) {
	tests := []string{
		"",
		"short",
		strings.Repeat("x", 79),
	}

	for _, line := range tests {
		t.Run("invalid length", func(t *testing.T) {
			assert.Nil(t, parseLine(line))
		})
	}
}

func TestParseLine_MissingNoradID(t *testing.T) {
	line := make([]byte, 117)

	for i := range line {
		line[i] = ' '
	}

	assert.Nil(t, parseLine(string(line)))
}

func TestParseLine_InvalidNoradID(t *testing.T) {
	line := make([]byte, 117)

	for i := range line {
		line[i] = ' '
	}

	copy(line[13:], "ABCDE")

	assert.Nil(t, parseLine(string(line)))
}

func TestParseLine_NoradIDBoundary(t *testing.T) {
	line := make([]byte, 117)

	for i := range line {
		line[i] = ' '
	}

	copy(line[13:], "12345")

	got := parseLine(string(line))

	require.NotNil(t, got)
	assert.Equal(t, "12345", got["norad_id"])
}

func TestParser_Stream(t *testing.T) {
	parser := NewParser()

	line1 := makeDefaultCelestrakLine()

	line2Bytes := []byte(makeDefaultCelestrakLine())
	copy(line2Bytes[13:], "12345")
	line2 := string(line2Bytes)

	input := strings.Join([]string{
		line1,
		"short invalid line",
		line2,
	}, "\n")

	var rows []ingestion.Row

	err := parser.Stream(
		strings.NewReader(input),
		func(row ingestion.Row) error {
			rows = append(rows, row)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, rows, 2)

	assert.Equal(t, "25544", rows[0]["norad_id"])
	assert.Equal(t, "12345", rows[1]["norad_id"])
}

func TestParser_Stream_EmptyInput(t *testing.T) {
	parser := NewParser()

	called := false

	err := parser.Stream(
		strings.NewReader(""),
		func(row ingestion.Row) error {
			called = true
			return nil
		},
	)

	require.NoError(t, err)
	assert.False(t, called)
}

func TestParser_Stream_CallbackError(t *testing.T) {
	parser := NewParser()
	expectedErr := errors.New("callback failed")

	err := parser.Stream(
		strings.NewReader(makeDefaultCelestrakLine()),
		func(row ingestion.Row) error {
			return expectedErr
		},
	)

	assert.ErrorIs(t, err, expectedErr)
}

func TestParser_Stream_SkipsInvalidLines(t *testing.T) {
	parser := NewParser()

	valid := makeDefaultCelestrakLine()

	var rows []ingestion.Row

	input := strings.Join([]string{
		"too short",
		valid,
		strings.Repeat("x", 100),
	}, "\n")

	err := parser.Stream(
		strings.NewReader(input),
		func(row ingestion.Row) error {
			rows = append(rows, row)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, rows, 1)
	assert.Equal(t, "25544", rows[0]["norad_id"])
}
