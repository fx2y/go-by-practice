package domain

import "context"

type StoreRepository interface {
	Save(ctx context.Context, store *Store) error
	Update(ctx context.Context, store *Store) error
	Delete(ctx context.Context, id string) error
	Find(ctx context.Context, id string) (*Store, error)
	FindAll(ctx context.Context) ([]*Store, error)
}
