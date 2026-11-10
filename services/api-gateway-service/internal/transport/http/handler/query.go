package handler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/models"
)

const (
	pageSizeKey     = "pageSize"
	pageTokenKey    = "pageToken"
	distanceUnitKey = "distanceUnit"
	angleUnitKey    = "angleUnit"
	fieldsKey       = "fields"

	timeRangeStartKey = "start"
	timeRangeEndKey   = "end"
	stepSecondsKey    = "stepSeconds"

	observerLatitudeDegreesKey    = "observerLatDeg"
	observerLatitudeRadianKey     = "observerLatRad"
	observerLongitudeDegreesKey   = "observerLonDeg"
	observerLongitudeRadianKey    = "observerLonRad"
	observerAltitudeMetersKey     = "observerAltM"
	observerAltitudeKilometersKey = "observerAltKM"
)

func badQuery(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func queryPageParams(
	r *http.Request,
) (pageSize uint32, pageToken string, err error) {
	rawPageSize := r.URL.Query().Get(pageSizeKey)
	if rawPageSize == "" {
		return 0, "", fmt.Errorf(
			"query parameter %q is required",
			pageSizeKey,
		)
	}

	pageSize64, err := strconv.ParseUint(rawPageSize, 10, 32)
	if err != nil {
		return 0, "", fmt.Errorf(
			"query parameter %q must be an unsigned integer",
			pageSizeKey,
		)
	}

	return uint32(pageSize64), r.URL.Query().Get(pageTokenKey), nil
}

func queryTime(r *http.Request, key string) (time.Time, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return time.Time{}, fmt.Errorf(
			"query parameter %q is required",
			key,
		)
	}

	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"query parameter %q must be an RFC3339 timestamp",
			key,
		)
	}

	return t.UTC(), nil
}

func queryUnits(r *http.Request) models.UnitSettings {
	q := r.URL.Query()

	return models.UnitSettings{
		DistanceUnit: q.Get(distanceUnitKey),
		AngleUnit:    q.Get(angleUnitKey),
	}
}

func queryFields(r *http.Request) []string {
	raw := r.URL.Query().Get(fieldsKey)
	if raw == "" {
		return nil
	}

	var fields []string

	for _, field := range strings.Split(raw, ",") {
		field = strings.TrimSpace(field)
		if field != "" {
			fields = append(fields, field)
		}
	}

	return fields
}

func queryRange(r *http.Request) (models.TimeRange, error) {
	start, err := queryTime(r, timeRangeStartKey)
	if err != nil {
		return models.TimeRange{}, err
	}

	end, err := queryTime(r, timeRangeEndKey)
	if err != nil {
		return models.TimeRange{}, err
	}

	return models.TimeRange{
		Start: start,
		End:   end,
	}, nil
}

func queryStepSeconds(r *http.Request) (int, error) {
	raw := r.URL.Query().Get(stepSecondsKey)
	if raw == "" {
		return 0, nil
	}

	seconds, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || int64(int(seconds)) != seconds {
		return 0, fmt.Errorf(
			"%s must be an integer",
			stepSecondsKey,
		)
	}

	return int(seconds), nil
}

func queryOptionalFloat(
	r *http.Request,
	key string,
) (*float64, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil, nil
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return nil, fmt.Errorf(
			"query parameter %q must be a number",
			key,
		)
	}

	return &value, nil
}

func queryOptionalUint32(
	r *http.Request,
	key string,
) (*uint32, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil, nil
	}

	value, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return nil, fmt.Errorf(
			"query parameter %q must be an unsigned integer not greater than %d",
			key,
			uint64(math.MaxUint32),
		)
	}

	result := uint32(value)
	return &result, nil
}

func queryRequiredUint32(r *http.Request, key string) (uint32, error) {
	value, err := queryOptionalUint32(r, key)
	if err != nil {
		return 0, err
	}
	if value == nil {
		return 0, fmt.Errorf("query parameter %q is required", key)
	}
	return *value, nil
}

func queryPassElevation(
	r *http.Request,
	degKey string,
	radKey string,
) (*float64, *float64, error) {
	deg, err := queryOptionalFloat(r, degKey)
	if err != nil {
		return nil, nil, err
	}

	rad, err := queryOptionalFloat(r, radKey)
	if err != nil {
		return nil, nil, err
	}

	return deg, rad, nil
}

func queryObserver(
	r *http.Request,
) (domain.GeodeticInput, error) {
	latDeg, err := queryOptionalFloat(r, observerLatitudeDegreesKey)
	if err != nil {
		return domain.GeodeticInput{}, err
	}

	latRad, err := queryOptionalFloat(r, observerLatitudeRadianKey)
	if err != nil {
		return domain.GeodeticInput{}, err
	}

	lonDeg, err := queryOptionalFloat(r, observerLongitudeDegreesKey)
	if err != nil {
		return domain.GeodeticInput{}, err
	}

	lonRad, err := queryOptionalFloat(r, observerLongitudeRadianKey)
	if err != nil {
		return domain.GeodeticInput{}, err
	}

	altM, err := queryOptionalFloat(r, observerAltitudeMetersKey)
	if err != nil {
		return domain.GeodeticInput{}, err
	}

	altKm, err := queryOptionalFloat(r, observerAltitudeKilometersKey)
	if err != nil {
		return domain.GeodeticInput{}, err
	}

	return domain.GeodeticInput{
		LatDeg: latDeg,
		LatRad: latRad,
		LonDeg: lonDeg,
		LonRad: lonRad,
		AltM:   altM,
		AltKm:  altKm,
	}, nil
}
