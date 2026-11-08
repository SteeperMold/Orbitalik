package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/gen/authpb"
	"google.golang.org/grpc"
)

type AuthHandler struct {
	client  authpb.AuthServiceClient
	timeout time.Duration
}

func NewAuthHandler(conn *grpc.ClientConn, timeout time.Duration) *AuthHandler {
	return &AuthHandler{
		client:  authpb.NewAuthServiceClient(conn),
		timeout: timeout,
	}
}

type credentialsRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()

	return decoder.Decode(target)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body credentialsRequest

	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, 400, "invalid_request", "Request body must be valid JSON")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	response, err := h.client.Register(ctx, &authpb.RegisterRequest{
		Email:    strings.TrimSpace(body.Email),
		Username: strings.TrimSpace(body.Username),
		Password: body.Password,
	})

	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body credentialsRequest

	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, 400, "invalid_request", "Request body must be valid JSON")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	response, err := h.client.Login(ctx, &authpb.LoginRequest{
		Email:    strings.TrimSpace(body.Email),
		Password: body.Password,
	})

	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body refreshRequest

	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, 400, "invalid_request", "Request body must be valid JSON")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	response, err := h.client.RefreshToken(ctx, &authpb.RefreshTokenRequest{
		RefreshToken: body.RefreshToken,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var body refreshRequest

	if err := decodeJSON(w, r, &body); err != nil {
		writeError(w, 400, "invalid_request", "Request body must be valid JSON")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	_, err := h.client.Logout(ctx, &authpb.LogoutRequest{
		RefreshToken: body.RefreshToken,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
