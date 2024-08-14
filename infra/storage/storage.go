package storage

import (
	"context"
	"time"
)

type Storage interface {
	RegisterAccess(ctx context.Context, key string, value string, limit int) (bool, int64, error)
	Search(ctx context.Context, key string, value string) (*time.Time, error)
	Block(ctx context.Context, key string, value string, block int) (*time.Time, error)
}
