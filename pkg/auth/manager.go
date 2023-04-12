package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"time"
	"transcribify/internal/models"
)

const (
	Token   = 15 * time.Minute
	Refresh = 24 * time.Hour
)

type TokenManager interface {
	NewJWT(user *models.User, ttl time.Duration) (models.Token, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey string
}

func (m *Manager) NewJWT(user *models.User, ttl time.Duration) (models.Token, error) {

	expires := time.Now().Add(ttl)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = (*user).ID
	claims["exp"] = expires.Unix()
	claims["role"] = user.Role

	signed, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return models.Token{}, err
	}

	return models.Token{
		T:         signed,
		ExpiresAt: expires,
	}, nil
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}
