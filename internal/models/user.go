package models

import "time"

type User struct {
	ID        string
	Email     string `validate:"required"`
	Password  string `validate:"required"`
	Role      string
	CreatedAt time.Time `validate:"required"`
	LastVisit time.Time `validate:"required"`
}

type SignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
