package hash

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher provide password hashing for securely store passwords
type PasswordHasher interface {
	Hash(password string) string
	Compare(password, hashed string) error
}

type BCHasher struct {
	salt string
	cost int
}

func NewBCHasher(cost int) *BCHasher {
	return &BCHasher{cost: cost}
}

func (b *BCHasher) Hash(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	return string(hash)
}
func (b *BCHasher) Compare(password, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
