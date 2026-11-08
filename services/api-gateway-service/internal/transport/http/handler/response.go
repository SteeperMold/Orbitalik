package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type APIError struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, APIError{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func writeServiceError(w http.ResponseWriter, err error) {
	if grpcStatus, ok := status.FromError(err); ok {
		switch grpcStatus.Code() {
		case codes.InvalidArgument:
			writeError(w, http.StatusBadRequest, "invalid_request", "Request parameters are invalid")
			return
		case codes.NotFound:
			writeError(w, http.StatusNotFound, "not_found", "The requested resource was not found")
			return
		case codes.AlreadyExists:
			writeError(w, http.StatusConflict, "already_exists", "The resource already exists")
			return
		case codes.Unauthenticated:
			writeError(w, http.StatusUnauthorized, "unauthenticated", "Authentication failed")
			return
		case codes.PermissionDenied:
			writeError(w, http.StatusForbidden, "forbidden", "You are not allowed to access this resource")
			return
		case codes.DeadlineExceeded:
			writeError(w, http.StatusGatewayTimeout, "upstream_timeout", "An upstream service timed out")
			return
		case codes.Unavailable:
			writeError(w, http.StatusBadGateway, "upstream_error", "An upstream service is unavailable")
			return
		}
	}
	switch {
	case errors.Is(err, domain.ErrSatelliteNotFound):
		writeError(
			w,
			http.StatusNotFound,
			"satellite_not_found",
			"Satellite was not found",
		)

	case errors.Is(err, domain.ErrInvalidArgument):
		writeError(
			w,
			http.StatusBadRequest,
			"invalid_request",
			"Request parameters are invalid",
		)

	case errors.Is(err, domain.ErrUpstreamTimeout):
		writeError(
			w,
			http.StatusGatewayTimeout,
			"upstream_timeout",
			"An upstream service timed out",
		)

	case errors.Is(err, domain.ErrUpstreamUnavailable):
		writeError(
			w,
			http.StatusBadGateway,
			"upstream_error",
			"An upstream service is unavailable",
		)

	default:
		writeError(
			w,
			http.StatusInternalServerError,
			"internal_error",
			"An unexpected error occurred",
		)
	}
}

func writeInvalidSatelliteID(w http.ResponseWriter) {
	writeError(
		w,
		http.StatusBadRequest,
		"invalid_satellite_id",
		"Satellite ID must be a positive integer",
	)
}
