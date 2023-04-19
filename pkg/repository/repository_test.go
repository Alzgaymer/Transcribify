package repository

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const (
	PathToRoot = "../../"
)

func TestMain(m *testing.M) {

	err := godotenv.Load(PathToRoot + ".env")
	if err != nil {
		log.Println(err)
		return
	}
	m.Run()
}

func Test_formatQuery(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "With \\n",
			query:    "\nCREATE DATABASE example;",
			expected: " CREATE DATABASE example;",
		},
		{
			name:     "With \\t",
			query:    "\tCREATE DATABASE example;",
			expected: "CREATE DATABASE example;",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := formatQuery(tt.query)

			assert.Equal(t, tt.expected, actual)
		})
	}
}
