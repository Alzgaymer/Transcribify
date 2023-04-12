package models

import "time"

type Token struct {
	T         string
	ExpiresAt time.Time
}
