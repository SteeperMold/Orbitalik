package service

import (
	"math"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/models"
)

const (
	defaultTrajectoryStepSeconds = 60

	minLatitudeDegrees  = -90.0
	maxLatitudeDegrees  = 90.0
	minLongitudeDegrees = -180.0
	maxLongitudeDegrees = 180.0

	minLatitudeRadians  = -math.Pi / 2
	maxLatitudeRadians  = math.Pi / 2
	minLongitudeRadians = -math.Pi
	maxLongitudeRadians = math.Pi

	maxElevationDegrees = 90.0
	maxElevationRadians = math.Pi / 2
)

func validateFinite(value float64, name string) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return domain.NewValidationError(
			"%s must be a finite number",
			name,
		)
	}

	return nil
}

func validateUnits(units models.UnitSettings) error {
	switch units.DistanceUnit {
	case "meters", "kilometers", "miles":
	default:
		return domain.NewValidationError(
			"distance unit must be one of: meters, kilometers, miles",
		)
	}

	switch units.AngleUnit {
	case "degrees", "radians":
	default:
		return domain.NewValidationError(
			"angle unit must be one of: degrees, radians",
		)
	}

	return nil
}

func validateTimeRange(timeRange models.TimeRange) error {
	if timeRange.Start.IsZero() {
		return domain.NewValidationError("start time is required")
	}

	if timeRange.End.IsZero() {
		return domain.NewValidationError("end time is required")
	}

	if !timeRange.End.After(timeRange.Start) {
		return domain.NewValidationError(
			"end must be later than start",
		)
	}

	return nil
}

func validatePositionParams(params domain.PositionParams) error {
	if params.Datetime.IsZero() {
		return domain.NewValidationError("datetime is required")
	}

	return validateUnits(params.Units)
}

func validateTrajectoryParams(params domain.TrajectoryParams) error {
	if err := validateTimeRange(params.Range); err != nil {
		return err
	}

	if params.Sampling.StepSeconds <= 0 {
		return domain.NewValidationError(
			"step seconds must be a positive integer",
		)
	}

	return validateUnits(params.Units)
}

func validateLookAnglesParams(params domain.LookAnglesParams) error {
	if params.Datetime.IsZero() {
		return domain.NewValidationError("datetime is required")
	}
	if err := validateUnits(params.Units); err != nil {
		return err
	}
	return validateGeodeticInput(params.Observer)
}

func validateObserverTrajectoryParams(params domain.ObserverTrajectoryParams) error {
	trajectoryParams := domain.TrajectoryParams{
		Range:    params.Range,
		Sampling: params.Sampling,
		Units:    params.Units,
	}

	if err := validateTrajectoryParams(trajectoryParams); err != nil {
		return err
	}
	return validateGeodeticInput(params.Observer)
}

func validateGeodeticInput(input domain.GeodeticInput) error {
	if (input.LatDeg == nil) == (input.LatRad == nil) {
		return domain.NewValidationError(
			"specify exactly one latitude representation",
		)
	}

	if (input.LonDeg == nil) == (input.LonRad == nil) {
		return domain.NewValidationError(
			"specify exactly one longitude representation",
		)
	}

	if (input.AltM == nil) == (input.AltKm == nil) {
		return domain.NewValidationError(
			"specify exactly one altitude representation",
		)
	}

	if input.LatDeg != nil {
		if err := validateFinite(*input.LatDeg, "latitude"); err != nil {
			return err
		}

		if *input.LatDeg < minLatitudeDegrees ||
			*input.LatDeg > maxLatitudeDegrees {
			return domain.NewValidationError(
				"latitude must be between -90 and 90 degrees",
			)
		}
	}

	if input.LatRad != nil {
		if err := validateFinite(*input.LatRad, "latitude"); err != nil {
			return err
		}

		if *input.LatRad < minLatitudeRadians ||
			*input.LatRad > maxLatitudeRadians {
			return domain.NewValidationError(
				"latitude must be between -pi/2 and pi/2 radians",
			)
		}
	}

	if input.LonDeg != nil {
		if err := validateFinite(*input.LonDeg, "longitude"); err != nil {
			return err
		}

		if *input.LonDeg < minLongitudeDegrees ||
			*input.LonDeg > maxLongitudeDegrees {
			return domain.NewValidationError(
				"longitude must be between -180 and 180 degrees",
			)
		}
	}

	if input.LonRad != nil {
		if err := validateFinite(*input.LonRad, "longitude"); err != nil {
			return err
		}

		if *input.LonRad < minLongitudeRadians ||
			*input.LonRad > maxLongitudeRadians {
			return domain.NewValidationError(
				"longitude must be between -pi and pi radians",
			)
		}
	}

	if input.AltM != nil {
		if err := validateFinite(*input.AltM, "altitude"); err != nil {
			return err
		}
	}

	if input.AltKm != nil {
		if err := validateFinite(*input.AltKm, "altitude"); err != nil {
			return err
		}
	}

	return nil
}

func validateElevation(
	degrees *float64,
	radians *float64,
	name string,
) error {
	if degrees != nil && radians != nil {
		return domain.NewValidationError(
			"specify either %s in degrees or radians, not both",
			name,
		)
	}

	if degrees != nil {
		if err := validateFinite(*degrees, name); err != nil {
			return err
		}

		if *degrees < 0 || *degrees > maxElevationDegrees {
			return domain.NewValidationError(
				"%s must be between 0 and 90 degrees",
				name,
			)
		}
	}

	if radians != nil {
		if err := validateFinite(*radians, name); err != nil {
			return err
		}

		if *radians < 0 || *radians > maxElevationRadians {
			return domain.NewValidationError(
				"%s must be between 0 and pi/2 radians",
				name,
			)
		}
	}

	return nil
}

func validatePassesParams(params domain.PassesParams) error {
	if err := validateTimeRange(params.Range); err != nil {
		return err
	}

	if err := validateUnits(params.Units); err != nil {
		return err
	}

	if err := validateGeodeticInput(params.Observer); err != nil {
		return err
	}

	if err := validateElevation(
		params.MinPeakElevationDeg,
		params.MinPeakElevationRad,
		"minimum peak elevation",
	); err != nil {
		return err
	}

	if err := validateElevation(
		params.MinElevationDeg,
		params.MinElevationRad,
		"minimum elevation",
	); err != nil {
		return err
	}

	if params.MaxResults != nil && *params.MaxResults == 0 {
		return domain.NewValidationError(
			"max results must be greater than zero",
		)
	}

	return nil
}

func validateNextPassesParams(params domain.NextPassesParams) error {
	if params.Count == 0 {
		return domain.NewValidationError("count must be greater than zero")
	}
	if err := validateUnits(params.Units); err != nil {
		return err
	}
	if err := validateGeodeticInput(params.Observer); err != nil {
		return err
	}
	if err := validateElevation(
		params.MinPeakElevationDeg,
		params.MinPeakElevationRad,
		"minimum peak elevation",
	); err != nil {
		return err
	}
	return validateElevation(params.MinElevationDeg, params.MinElevationRad, "minimum elevation")
}
