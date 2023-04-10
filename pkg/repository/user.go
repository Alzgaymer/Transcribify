package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	client *pgx.Conn
}

func (u *UserRepository) GetRefreshTokenByID(ctx context.Context, id string) (string, error) {
	query := "select refresh_token from users where id = $1;"
	var token string

	if err := u.client.QueryRow(ctx, query, id).Scan(&token); err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserRepository) SetRefreshToken(ctx context.Context, login, token string) error {
	query := "UPDATE users SET refresh_token = $2 WHERE email = $1;"

	if err := u.client.Ping(ctx); err != nil {
		return err
	}

	if err := u.client.QueryRow(ctx, query, login, token).Scan(); err != pgx.ErrNoRows {
		return err
	}

	return nil
}

func (u *UserRepository) SignUser(ctx context.Context, login, password string) error {
	query := `SELECT id FROM users WHERE email = $1 AND password = $2;`
	var id int

	err := u.client.QueryRow(ctx, query, login, login, password).Scan(&id)
	if err == nil {
		return fmt.Errorf("signuser: user already exist, its id: %d", id)
	}

	query = "INSERT INTO users (email, password) VALUES ($1,$2);"

	u.client.QueryRow(ctx, query, login, password)

	return nil
}

func (u *UserRepository) GetUserId(ctx context.Context, login string) (int, error) {
	var (
		rawQuery = `select id from users as u 
					where u.email = $1;
					`
		query = formatQuery(rawQuery)
		id    int
	)

	err := u.client.QueryRow(ctx, query, login).Scan(&id)
	if err != nil {
		return NotFound, err
	}

	return id, nil
}

func NewUserRepository(client *pgx.Conn) *UserRepository {
	return &UserRepository{client: client}
}
