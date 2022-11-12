package stores

import (
	"context"
	"event-driven-architecture-in-golang/internal/monolith"
	"event-driven-architecture-in-golang/stores/internal/application"
	"event-driven-architecture-in-golang/stores/internal/grpc"
	"event-driven-architecture-in-golang/stores/internal/logging"
	"event-driven-architecture-in-golang/stores/internal/postgres"
	"event-driven-architecture-in-golang/stores/internal/rest"
)

type Module struct {
}

func (m *Module) Startup(ctx context.Context, mono monolith.Monolith) error {
	stores := postgres.NewStoreRepository("stores.stores", mono.DB())
	participatingStores := postgres.NewParticipatingStoreRepository("stores.stores", mono.DB())
	products := postgres.NewProductRepository("stores.products", mono.DB())

	var app application.App
	app = application.New(stores, participatingStores, products)
	app = logging.LogApplicationAccess(app, mono.Logger())

	if err := grpc.RegisterServer(ctx, app, mono.RPC()); err != nil {
		return err
	}
	if err := rest.RegisterGateway(ctx, mono.Mux(), mono.Config().Rpc.Address()); err != nil {
		return err
	}
	if err := rest.RegisterSwagger(mono.Mux()); err != nil {
		return err
	}

	return nil
}
