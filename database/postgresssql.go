package database

import (
	"context"
)

type Postgres struct {
	client Client
}

func NewPostgres(client Client) Repository {
	return &Postgres{
		client: client,
	}
}

func (p *Postgres) Create(ctx context.Context, a ...any) error {
	//TODO implement me
	panic("implement me")
}

func (p *Postgres) Read(ctx context.Context, a ...any) error {
	//TODO implement me
	panic("implement me")
}

func (p *Postgres) Update(ctx context.Context, a ...any) error {
	//TODO implement me
	panic("implement me")
}

func (p *Postgres) Delete(ctx context.Context, a ...any) error {
	//TODO implement me
	panic("implement me")
}
