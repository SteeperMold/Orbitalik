package route

import (
	"net/http"
	"time"

	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/transport/adapter"
	"github.com/SteeperMold/Orbitalik/api-gateway-service/internal/transport/http/handler"
	"github.com/gorilla/mux"
)

func NewAuthRoute(router *mux.Router, connections *adapter.Connections, timeout time.Duration) {
	h := handler.NewAuthHandler(connections.Auth, timeout)

	router.HandleFunc("/auth/register", h.Register).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", h.Login).Methods(http.MethodPost)
	router.HandleFunc("/auth/refresh", h.Refresh).Methods(http.MethodPost)
	router.HandleFunc("/auth/logout", h.Logout).Methods(http.MethodPost)
}
