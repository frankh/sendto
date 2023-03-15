package db

import (
	"fmt"
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
}

type MemoryDB struct {
	db map[string]Message
}

func NewMemoryDB() DB {
	return &MemoryDB{
		db: map[string]Message{},
	}
}

func (d *MemoryDB) Save(message Message) error {
	d.db[message.ID] = message
	return nil
}

func (d *MemoryDB) Load(ID string) (*Message, error) {
	message, ok := d.db[ID]
	if !ok {
		return nil, fmt.Errorf("Not found")
	}
	return &message, nil
}
