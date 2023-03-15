package db

import (
	"database/sql"
	_ "embed"
	"errors"
	"io/fs"

	_ "modernc.org/sqlite"
)

type SqliteDB struct {
	db *sql.DB
}

//go:embed schema.sql
var sqlSchema string

func NewSqliteDB(filepath string) DB {
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}

	if _, err = db.Exec(sqlSchema); err != nil {
		panic(err)
	}

	return &SqliteDB{
		db: db,
	}
}

func (d *SqliteDB) Save(message Message) error {
	_, err := d.db.Exec(
		"INSERT INTO Messages (ID, FromUser, ToUser, CipherText) VALUES (?, ?, ?, ?)",
		message.ID, message.From, message.To, message.CipherText)
	if err != nil {
		return err
	}
	return nil
}

func (d *SqliteDB) Load(ID string) (*Message, error) {
	var message Message
	row := d.db.QueryRow("SELECT ID, FromUser, ToUser, CipherText FROM Messages WHERE ID = ?1 LIMIT 1", ID)
	err := row.Scan(&message.ID, &message.From, &message.To, &message.CipherText)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = fs.ErrNotExist
		}
		return nil, err
	}
	return &message, nil
}
