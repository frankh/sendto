package db

import (
	"time"
)

type Message struct {
	ID         string
	CipherText string
	ExpiresAt  time.Time
	From       string
	To         string
}

type DB interface {
	Save(message Message) error
	Load(ID string) (*Message, error)
	ClearExpired() error
}
