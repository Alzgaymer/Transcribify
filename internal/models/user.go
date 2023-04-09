package models

import "time"

type User struct {
	Name      string    `validate:"required"`
	Password  string    `validate:"required"`
	Email     string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	LastVisit time.Time `validate:"required"`
}

type SignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
