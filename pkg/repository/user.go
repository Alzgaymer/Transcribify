package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"transcribify/internal/models"
)

type UserRepository struct {
	client *pgx.Conn
}

func (u *UserRepository) GetUserByLogin(ctx context.Context, user *models.User) error {
	return u.client.QueryRow(ctx, "select id, email, password from get_user($1)", user.Email).
		Scan(&user.ID, &user.Email, &user.Password)
}

func (u *UserRepository) PutUser(ctx context.Context, user *models.User) error {
	return u.client.QueryRow(ctx, "select id from put_user($1, $2)", user.Email, user.Password).
		Scan(&user.ID)
}

func (u *UserRepository) GetUserByLoginPassword(ctx context.Context, user *models.User) error {
	return u.client.QueryRow(ctx, "select id from get_user($1, $2)", user.Email, user.Password).
		Scan(&user.ID)
}

func NewUserRepository(client *pgx.Conn) *UserRepository {
	return &UserRepository{client: client}
}
