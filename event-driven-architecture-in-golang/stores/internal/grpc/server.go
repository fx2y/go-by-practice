package grpc

import (
	"context"
	"event-driven-architecture-in-golang/stores/internal/application"
	"event-driven-architecture-in-golang/stores/internal/application/commands"
	"event-driven-architecture-in-golang/stores/internal/application/queries"
	"event-driven-architecture-in-golang/stores/internal/domain"
	"event-driven-architecture-in-golang/stores/storespb"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type server struct {
	app application.App
	storespb.UnimplementedStoresServiceServer
}

var _ storespb.StoresServiceServer = (*server)(nil)

func RegisterServer(_ context.Context, app application.App, registrar grpc.ServiceRegistrar) error {
	storespb.RegisterStoresServiceServer(registrar, &server{app: app})
	return nil
}

func (s *server) CreateStore(ctx context.Context, req *storespb.CreateStoreRequest) (*storespb.CreateStoreResponse, error) {
	storeID := uuid.New().String()
	err := s.app.CreateStore(ctx, commands.CreateStore{
		ID:       storeID,
		Name:     req.GetName(),
		Location: req.GetLocation(),
	})
	if err != nil {
		return nil, err
	}
	return &storespb.CreateStoreResponse{
		Id: storeID,
	}, nil
}

func (s *server) GetStore(ctx context.Context, req *storespb.GetStoreRequest) (*storespb.GetStoreResponse, error) {
	store, err := s.app.GetStore(ctx, queries.GetStore{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, err
	}
	return &storespb.GetStoreResponse{
		Store: s.storeFromDomain(store),
	}, nil
}

func (s *server) GetStores(ctx context.Context, req *storespb.GetStoresRequest) (*storespb.GetStoresResponse, error) {
	stores, err := s.app.GetStores(ctx, queries.GetStores{})
	if err != nil {
		return nil, err
	}
	var protoStores []*storespb.Store
	for _, store := range stores {
		protoStores = append(protoStores, s.storeFromDomain(store))
	}
	return &storespb.GetStoresResponse{
		Stores: protoStores,
	}, nil
}

func (s *server) EnableParticipation(ctx context.Context, req *storespb.EnableParticipationRequest) (*storespb.EnableParticipationResponse, error) {
	err := s.app.EnableParticipation(ctx, commands.EnableParticipation{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, err
	}
	return &storespb.EnableParticipationResponse{}, nil
}

func (s *server) DisableParticipation(ctx context.Context, req *storespb.DisableParticipationRequest) (*storespb.DisableParticipationResponse, error) {
	err := s.app.DisableParticipation(ctx, commands.DisableParticipation{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, err
	}
	return &storespb.DisableParticipationResponse{}, nil
}

func (s *server) GetParticipatingStores(ctx context.Context, req *storespb.GetParticipatingStoresRequest) (*storespb.GetParticipatingStoresResponse, error) {
	stores, err := s.app.GetParticipatingStores(ctx, queries.GetParticipatingStores{})
	if err != nil {
		return nil, err
	}
	var protoStores []*storespb.Store
	for _, store := range stores {
		protoStores = append(protoStores, s.storeFromDomain(store))
	}
	return &storespb.GetParticipatingStoresResponse{
		Stores: protoStores,
	}, nil
}

func (s *server) AddProduct(ctx context.Context, req *storespb.AddProductRequest) (*storespb.AddProductResponse, error) {
	id := uuid.New().String()
	err := s.app.AddProduct(ctx, commands.AddProduct{
		ID:          id,
		StoreID:     req.GetStoreId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		SKU:         req.GetSku(),
		Price:       req.GetPrice(),
	})
	if err != nil {
		return nil, err
	}
	return &storespb.AddProductResponse{Id: id}, nil
}

func (s *server) RemoveProduct(ctx context.Context, req *storespb.RemoveProductRequest) (*storespb.RemoveProductResponse, error) {
	err := s.app.RemoveProduct(ctx, commands.RemoveProduct{
		ID: req.GetId(),
	})
	return &storespb.RemoveProductResponse{}, err
}

func (s *server) GetCatalog(ctx context.Context, req *storespb.GetCatalogRequest) (*storespb.GetCatalogResponse, error) {
	products, err := s.app.GetCatalog(ctx, queries.GetCatalog{
		StoreID: req.GetStoreId(),
	})
	if err != nil {
		return nil, err
	}
	var protoProducts []*storespb.Product
	for _, product := range products {
		protoProducts = append(protoProducts, s.productFromDomain(product))
	}
	return &storespb.GetCatalogResponse{
		Products: protoProducts,
	}, nil
}

func (s *server) GetProduct(ctx context.Context, req *storespb.GetProductRequest) (*storespb.GetProductResponse, error) {
	product, err := s.app.GetProduct(ctx, queries.GetProduct{
		ID: req.GetId(),
	})
	if err != nil {
		return nil, err
	}
	return &storespb.GetProductResponse{
		Product: s.productFromDomain(product),
	}, nil
}

func (s *server) storeFromDomain(store *domain.Store) *storespb.Store {
	return &storespb.Store{
		Id:            store.ID,
		Name:          store.Name,
		Location:      store.Location,
		Participating: store.Participating,
	}
}

func (s *server) productFromDomain(product *domain.Product) *storespb.Product {
	return &storespb.Product{
		Id:          product.ID,
		StoreId:     product.StoreID,
		Name:        product.Name,
		Description: product.Description,
		Sku:         product.SKU,
		Price:       product.Price,
	}
}
