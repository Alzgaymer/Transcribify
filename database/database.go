package database

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// Client has signature of pgx.Tx interface
type Client interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repository interface {
	Create(context.Context, ...any) error
	Read(context.Context, ...any) error
	Update(context.Context, ...any) error
	Delete(context.Context, ...any) error
}
