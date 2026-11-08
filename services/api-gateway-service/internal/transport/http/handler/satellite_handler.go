package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/domain"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/models"
	applog "github.com/SteeperMold/Orbitalik/common/go/log"
)

type SatelliteHandler struct {
	satellites domain.SatelliteService
	logger     applog.Logger
	timeout    time.Duration
}

func NewSatelliteHandler(
	svc domain.SatelliteService,
	logger applog.Logger,
	timeout time.Duration,
) *SatelliteHandler {
	return &SatelliteHandler{
		satellites: svc,
		logger:     logger,
		timeout:    timeout,
	}
}

func (h *SatelliteHandler) List(w http.ResponseWriter, r *http.Request) {
	pageSize, pageToken, err := queryPageParams(r)
	if err != nil {
		badQuery(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	satellites, nextToken, err := h.satellites.List(
		ctx,
		pageSize,
		pageToken,
	)
	if err != nil {
		h.logger.Error("failed to list satellites")
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"nextPageToken": nextToken,
		"satellites":    satellites,
	})
}

func (h *SatelliteHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseSatelliteIdentifier(r)
	if err != nil {
		writeInvalidSatelliteID(w)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	satellite, err := h.satellites.Get(ctx, id)
	if err != nil {
		h.logger.Error("failed to get satellite")
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, satellite)
}

func (h *SatelliteHandler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	id, err := parseSatelliteIdentifier(r)
	if err != nil {
		writeInvalidSatelliteID(w)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	metadata, err := h.satellites.GetMetadata(ctx, id)
	if err != nil {
		h.logger.Error("failed to get satellite metadata")
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, metadata)
}

func (h *SatelliteHandler) GetFrequencies(w http.ResponseWriter, r *http.Request) {
	id, err := parseSatelliteIdentifier(r)
	if err != nil {
		writeInvalidSatelliteID(w)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	frequencies, err := h.satellites.GetFrequencies(ctx, id)
	if err != nil {
		h.logger.Error("failed to get satellite frequencies")
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, frequencies)
}

func (h *SatelliteHandler) GetSources(w http.ResponseWriter, r *http.Request) {
	id, err := parseSatelliteIdentifier(r)
	if err != nil {
		writeInvalidSatelliteID(w)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	sources, err := h.satellites.GetSources(ctx, id)
	if err != nil {
		h.logger.Error("failed to get satellite sources")
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, sources)
}

func (h *SatelliteHandler) GetTle(w http.ResponseWriter, r *http.Request) {
	id, err := parseSatelliteIdentifier(r)
	if err != nil {
		writeInvalidSatelliteID(w)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	tle, err := h.satellites.GetTle(ctx, id)
	if err != nil {
		h.logger.Error("failed to get satellite TLE")
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, tle)
}

func (h *SatelliteHandler) ListTles(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	tles, err := h.satellites.ListTles(ctx)
	if err != nil {
		h.logger.Error("failed to list TLEs")
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"tles": tles})
}

func (h *SatelliteHandler) GetPosition(w http.ResponseWriter, r *http.Request) {
	id, err := parseSatelliteIdentifier(r)
	if err != nil {
		writeInvalidSatelliteID(w)
		return
	}

	datetime, err := queryTime(r, "datetime")
	if err != nil {
		badQuery(w, err)
		return
	}

	params := domain.PositionParams{
		Datetime: datetime,
		Fields:   queryFields(r),
		Units:    queryUnits(r),
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	position, err := h.satellites.GetPosition(ctx, id, params)
	if err != nil {
		h.logger.Error("failed to get satellite position")
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, position)
}

func (h *SatelliteHandler) GetLookAngles(w http.ResponseWriter, r *http.Request) {
	id, err := parseSatelliteIdentifier(r)
	if err != nil {
		writeInvalidSatelliteID(w)
		return
	}

	datetime, err := queryTime(r, "datetime")
	if err != nil {
		badQuery(w, err)
		return
	}

	observer, err := queryObserver(r)
	if err != nil {
		badQuery(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	params := domain.LookAnglesParams{
		Datetime: datetime,
		Observer: observer,
		Fields:   queryFields(r),
		Units:    queryUnits(r),
	}

	result, err := h.satellites.GetLookAngles(ctx, id, params)
	if err != nil {
		h.logger.Error("failed to get look angles")
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *SatelliteHandler) GetTrajectory(w http.ResponseWriter, r *http.Request) {
	id, err := parseSatelliteIdentifier(r)
	if err != nil {
		writeInvalidSatelliteID(w)
		return
	}

	timeRange, err := queryRange(r)
	if err != nil {
		badQuery(w, err)
		return
	}

	step, err := queryStepSeconds(r)
	if err != nil {
		badQuery(w, err)
		return
	}

	params := domain.TrajectoryParams{
		Range: timeRange,
		Sampling: models.SamplingOptions{
			StepSeconds: step,
		},
		Fields: queryFields(r),
		Units:  queryUnits(r),
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	trajectory, err := h.satellites.GetTrajectory(ctx, id, params)
	if err != nil {
		h.logger.Error("failed to get satellite trajectory")
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, trajectory)
}

func (h *SatelliteHandler) GetObserverTrajectory(w http.ResponseWriter, r *http.Request) {
	id, err := parseSatelliteIdentifier(r)
	if err != nil {
		writeInvalidSatelliteID(w)
		return
	}

	timeRange, err := queryRange(r)
	if err != nil {
		badQuery(w, err)
		return
	}

	step, err := queryStepSeconds(r)
	if err != nil {
		badQuery(w, err)
		return
	}

	observer, err := queryObserver(r)
	if err != nil {
		badQuery(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	params := domain.ObserverTrajectoryParams{
		Range: timeRange,
		Sampling: models.SamplingOptions{
			StepSeconds: step,
		},
		Observer: observer,
		Fields:   queryFields(r),
		Units:    queryUnits(r),
	}

	result, err := h.satellites.GetObserverTrajectory(ctx, id, params)
	if err != nil {
		h.logger.Error("failed to get observer trajectory")
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *SatelliteHandler) GetPasses(w http.ResponseWriter, r *http.Request) {
	id, err := parseSatelliteIdentifier(r)
	if err != nil {
		writeInvalidSatelliteID(w)
		return
	}

	timeRange, err := queryRange(r)
	if err != nil {
		badQuery(w, err)
		return
	}

	observer, err := queryObserver(r)
	if err != nil {
		badQuery(w, err)
		return
	}

	minPeakDeg, minPeakRad, err := queryPassElevation(
		r,
		"min_peak_elevation_deg",
		"min_peak_elevation_rad",
	)
	if err != nil {
		badQuery(w, err)
		return
	}

	minElevationDeg, minElevationRad, err := queryPassElevation(
		r,
		"min_elevation_deg",
		"min_elevation_rad",
	)
	if err != nil {
		badQuery(w, err)
		return
	}

	maxResults, err := queryOptionalUint32(r, "max_results")
	if err != nil {
		badQuery(w, err)
		return
	}

	params := domain.PassesParams{
		Range:               timeRange,
		Observer:            observer,
		Units:               queryUnits(r),
		MinPeakElevationDeg: minPeakDeg,
		MinPeakElevationRad: minPeakRad,
		MinElevationDeg:     minElevationDeg,
		MinElevationRad:     minElevationRad,
		MaxResults:          maxResults,
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	passes, err := h.satellites.GetPasses(ctx, id, params)
	if err != nil {
		h.logger.Error("failed to get satellite passes")
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, passes)
}

func (h *SatelliteHandler) GetNextPasses(w http.ResponseWriter, r *http.Request) {
	id, err := parseSatelliteIdentifier(r)
	if err != nil {
		writeInvalidSatelliteID(w)
		return
	}

	observer, err := queryObserver(r)
	if err != nil {
		badQuery(w, err)
		return
	}

	count, err := queryRequiredUint32(r, "count")
	if err != nil {
		badQuery(w, err)
		return
	}

	minPeakDeg, minPeakRad, err := queryPassElevation(
		r,
		"min_peak_elevation_deg",
		"min_peak_elevation_rad",
	)
	if err != nil {
		badQuery(w, err)
		return
	}

	minElevationDeg, minElevationRad, err := queryPassElevation(
		r,
		"min_elevation_deg",
		"min_elevation_rad",
	)
	if err != nil {
		badQuery(w, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	params := domain.NextPassesParams{
		Observer:            observer,
		Units:               queryUnits(r),
		Count:               count,
		MinPeakElevationDeg: minPeakDeg,
		MinPeakElevationRad: minPeakRad,
		MinElevationDeg:     minElevationDeg,
		MinElevationRad:     minElevationRad,
	}

	result, err := h.satellites.GetNextPasses(ctx, id, params)
	if err != nil {
		h.logger.Error("failed to get next passes")
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}
