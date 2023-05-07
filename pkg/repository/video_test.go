package repository

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"testing"
)

func TestYTVideoRepository(t *testing.T) {
	//repo := NewYTVideoRepository(db)

}
