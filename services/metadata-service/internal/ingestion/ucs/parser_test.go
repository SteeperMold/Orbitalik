package ucs

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/SteeperMold/Orbitalik/satellite-metadata-service/internal/ingestion"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_Stream(t *testing.T) {
	input := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite\tName of Satellite, Alternate Names\tOperator/Owner\tCountry of Operator/Owner\tUsers\tPurpose\tClass of Orbit\tDate of Launch\tLaunch Site\tLaunch Vehicle\tCOSPAR Number",
		"25544\tISS\tZARYA\tNASA\tUSA\tGovernment\tSpace Station\tLEO\t11/20/1998\tBaikonur\tProton-K\t1998-067A",
	}, "\n")

	parser := NewParser()

	var rows []ingestion.Row

	err := parser.Stream(
		strings.NewReader(input),
		func(row ingestion.Row) error {
			rows = append(rows, row)
			return nil
		},
	)

	require.NoError(t, err)
	require.Len(t, rows, 1)

	row := rows[0]

	assert.Equal(t, "25544", row["norad_id"])
	assert.Equal(t, "ISS", row["name"])
	assert.Equal(t, "ZARYA", row["aliases"])
	assert.Equal(t, "NASA", row["operator"])
	assert.Equal(t, "USA", row["owner"])
	assert.Equal(t, "Government", row["users"])
	assert.Equal(t, "Space Station", row["purpose"])
	assert.Equal(t, "LEO", row["orbit_class"])
	assert.Equal(t, "11/20/1998", row["launch_date"])
	assert.Equal(t, "Baikonur", row["launch_site"])
	assert.Equal(t, "Proton-K", row["launch_vehicle"])
	assert.Equal(t, "1998-067A", row["cospar"])
}

func TestParser_Stream_WhitespaceInHeaders(t *testing.T) {
	input := " Norad Number \t CURRENT OFFICIAL NAME OF SATELLITE \t Users \n" +
		"25544\tISS\tGovernment\n"

	parser := NewParser()

	var rows []ingestion.Row

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
	assert.Equal(t, "ISS", rows[0]["name"])
	assert.Equal(t, "Government", rows[0]["users"])
}

func TestParser_Stream_SkipsRowsWithoutNoradID(t *testing.T) {
	input := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite",
		"\tISS",
		"25544\tOther Satellite",
	}, "\n")

	parser := NewParser()

	var rows []ingestion.Row

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

func TestParser_Stream_SkipsMalformedRows(t *testing.T) {
	input := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite",
		"25544\tISS",
		`25544\t"broken"`,
		"12345\tSAT-2",
	}, "\n")

	parser := NewParser()

	var rows []ingestion.Row

	err := parser.Stream(
		strings.NewReader(input),
		func(row ingestion.Row) error {
			rows = append(rows, row)
			return nil
		},
	)

	require.NoError(t, err)
	assert.NotEmpty(t, rows)
}

func TestParser_Stream_EmptyInput(t *testing.T) {
	parser := NewParser()

	err := parser.Stream(
		strings.NewReader(""),
		func(row ingestion.Row) error {
			t.Fatal("callback should not be called")
			return nil
		},
	)

	require.NoError(t, err)
}

func TestParser_Stream_CallbackError(t *testing.T) {
	input := strings.Join([]string{
		"Norad Number\tCurrent Official Name of Satellite",
		"25544\tISS",
	}, "\n")

	expectedErr := errors.New("callback failed")
	parser := NewParser()

	err := parser.Stream(
		strings.NewReader(input),
		func(row ingestion.Row) error {
			return expectedErr
		},
	)

	require.ErrorIs(t, err, expectedErr)
}

func TestParser_Stream_ContextDeadlineError(t *testing.T) {
	parser := NewParser()

	reader := &errorReader{
		err: context.DeadlineExceeded,
	}

	err := parser.Stream(
		reader,
		func(row ingestion.Row) error {
			t.Fatal("callback should not be called")
			return nil
		},
	)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestBuildIndex(t *testing.T) {
	header := []string{
		"Norad Number",
		" Current Official Name of Satellite ",
		"USERS",
	}

	got := buildIndex(header)

	assert.Equal(t, 0, got["norad number"])
	assert.Equal(t, 1, got["current official name of satellite"])
	assert.Equal(t, 2, got["users"])
}

func TestNormalizeHeader(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "lowercase",
			input: "NORAD NUMBER",
			want:  "norad number",
		},
		{
			name:  "trims whitespace",
			input: "  Norad Number  ",
			want:  "norad number",
		},
		{
			name:  "lowercase and trims",
			input: "  CURRENT OFFICIAL NAME OF SATELLITE  ",
			want:  "current official name of satellite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, normalizeHeader(tt.input))
		})
	}
}

func TestMapRow(t *testing.T) {
	idx := map[string]int{
		"norad number":                       0,
		"current official name of satellite": 1,
		"name of satellite, alternate names": 2,
		"operator/owner":                     3,
		"country of operator/owner":          4,
		"users":                              5,
		"purpose":                            6,
		"class of orbit":                     7,
		"date of launch":                     8,
		"launch site":                        9,
		"launch vehicle":                     10,
		"cospar number":                      11,
	}

	row := []string{
		`"25544"`,
		`"ISS"`,
		`"ZARYA"`,
		`"NASA"`,
		`"USA"`,
		`"Government"`,
		`"Space Station"`,
		`"LEO"`,
		`"11/20/1998"`,
		`"Baikonur"`,
		`"Proton-K"`,
		`"1998-067A"`,
	}

	got := mapRow(idx, row)

	require.NotNil(t, got)

	assert.Equal(t, "25544", got["norad_id"])
	assert.Equal(t, "ISS", got["name"])
	assert.Equal(t, "ZARYA", got["aliases"])
	assert.Equal(t, "NASA", got["operator"])
	assert.Equal(t, "USA", got["owner"])
	assert.Equal(t, "Government", got["users"])
	assert.Equal(t, "Space Station", got["purpose"])
	assert.Equal(t, "LEO", got["orbit_class"])
	assert.Equal(t, "11/20/1998", got["launch_date"])
	assert.Equal(t, "Baikonur", got["launch_site"])
	assert.Equal(t, "Proton-K", got["launch_vehicle"])
	assert.Equal(t, "1998-067A", got["cospar"])
}

func TestMapRow_MissingNoradID(t *testing.T) {
	idx := map[string]int{
		"norad number": 0,
	}

	row := []string{""}

	assert.Nil(t, mapRow(idx, row))
}

func TestMapRow_MissingOptionalColumns(t *testing.T) {
	idx := map[string]int{
		"norad number": 0,
		"name":         1,
	}

	row := []string{
		"25544",
	}

	got := mapRow(idx, row)

	require.NotNil(t, got)
	assert.Equal(t, "25544", got["norad_id"])
	assert.Equal(t, "", got["name"])
	assert.Equal(t, "", got["operator"])
	assert.Equal(t, "", got["cospar"])
}

func TestMapRow_ColumnBeyondRowLength(t *testing.T) {
	idx := map[string]int{
		"norad number": 0,
		"users":        1,
	}

	row := []string{"25544"}

	got := mapRow(idx, row)

	require.NotNil(t, got)
	assert.Equal(t, "25544", got["norad_id"])
	assert.Equal(t, "", got["users"])
}

func TestClean(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "trims whitespace",
			input: "  hello  ",
			want:  "hello",
		},
		{
			name:  "removes quotes",
			input: `"hello"`,
			want:  "hello",
		},
		{
			name:  "trims and removes quotes",
			input: `  "hello"  `,
			want:  "hello",
		},
		{
			name:  "empty",
			input: "",
			want:  "",
		},
		{
			name:  "spaces only",
			input: "   ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, clean(tt.input))
		})
	}
}

type errorReader struct {
	err error
}

func (r *errorReader) Read([]byte) (int, error) {
	return 0, r.err
}
