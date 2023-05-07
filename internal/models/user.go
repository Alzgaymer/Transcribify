package models

import "time"

type User struct {
	ID        int
	Email     string `validate:"required"`
	Password  string `validate:"required"`
	Role      string
	CreatedAt time.Time `validate:"required"`
	LastVisit time.Time `validate:"required"`
}
