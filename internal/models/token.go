package models

import "time"

type Token struct {
	Key       string
	T         string
	ExpiresAt time.Time
}
