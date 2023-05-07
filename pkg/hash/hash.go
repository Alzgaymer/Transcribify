package hash

import (
	"crypto/sha1"
	"fmt"
)

// PasswordHasher provide password hashing for securely store passwords
type PasswordHasher interface {
	Hash(password string) string
}

// SHA1PSHasher uses SHA1 to hash passwords with provided salt
type SHA1PSHasher struct {
	salt string
}

func (s *SHA1PSHasher) Hash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(s.salt)))
}

func NewSHA1PSHasher(salt string) *SHA1PSHasher {
	return &SHA1PSHasher{salt: salt}
}
