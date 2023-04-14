package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"transcribify/internal/models"
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

// GetUserByLogin searching in database for user with provided login.
// If success returns user`s id.
// If failed - returns "-1", error. "-1" means there`re no users with provided login
func (u *UserRepository) GetUserByLogin(ctx context.Context, login string) (string, error) {
	query := `SELECT id FROM users WHERE email = $1;`
	var id int

	if err := u.client.QueryRow(ctx, query, login).Scan(&id); err != nil {
		return convCode(NotFound), err
	}

	return fmt.Sprintf("%d", id), nil
}
func convCode(code int) string {
	return fmt.Sprintf("%d", code)
}

// PutUser firstly search for existing user.
// If found - returns user`s id, error.
// If not - inserting. On failure returns InternalRepositoryError, err.
// On success returns user`s id, nil
func (u *UserRepository) PutUser(ctx context.Context, user *models.User) (string, error) {

	userID, err := u.GetUserByLogin(ctx, user.Email)
	if userID != convCode(NotFound) {
		return userID, err
	}

	query := "INSERT INTO users (email, password) VALUES ($1,$2) returning id;"
	var id int

	if err = u.client.QueryRow(ctx, query, user.Email, user.Password).Scan(&id); err != nil {
		return convCode(InternalRepositoryError), err
	}

	return fmt.Sprintf("%d", id), nil
}

func (u *UserRepository) GetUserByLoginPassword(ctx context.Context, user *models.User) (string, error) {
	query := `SELECT id FROM users WHERE email = $1 and password = $2;`
	var id int

	err := u.client.QueryRow(ctx, query, user.Email, user.Password).Scan(&id)
	switch {
	case err == pgx.ErrNoRows:
		return convCode(NotFound), fmt.Errorf("user(%s) doesn`t exist", user.Email)
	case err != nil:
		return convCode(InternalRepositoryError), err
	default:
		return fmt.Sprintf("%d", id), nil
	}
}

func (u *UserRepository) GetUserId(ctx context.Context, login string) (string, error) {
	var (
		rawQuery = `select id from users as u 
					where u.email = $1;
					`
		query = formatQuery(rawQuery)
		id    int
	)

	err := u.client.QueryRow(ctx, query, login).Scan(&id)
	if err != nil {

		return convCode(NotFound), err
	}

	return fmt.Sprintf("%d", id), nil
}

func NewUserRepository(client *pgx.Conn) *UserRepository {
	return &UserRepository{client: client}
}
