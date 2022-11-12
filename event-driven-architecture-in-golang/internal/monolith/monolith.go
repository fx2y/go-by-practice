package monolith

import (
	"context"
	"database/sql"
	"event-driven-architecture-in-golang/internal/config"
	"event-driven-architecture-in-golang/internal/waiter"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Monolith interface {
	Config() config.AppConfig
	DB() *sql.DB
	Logger() zerolog.Logger
	Mux() *chi.Mux
	RPC() *grpc.Server
	Waiter() waiter.Waiter
}

type Module interface {
	Startup(ctx context.Context, m Monolith) error
}
