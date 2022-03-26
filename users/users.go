package users

import (
	"database/sql"
	"sync"

	_ "modernc.org/sqlite"
)

var existUser = `
SELECT COUNT(*) FROM USERS WHERE ID==?
`

var add = `
INSERT INTO users (
    id, name, starts, messages, passeds
) VALUES (
    ?, ?, 0, 0, 0
)
`

var upStarts = `
UPDATE users
SET starts = starts+1
WHERE ID==?;
`

var upMessages = `
UPDATE users
SET messages = messages+1
WHERE ID==?;
`

var upPasseds = `
UPDATE users
SET passeds = passeds+1
WHERE ID==?;
`

var schemaSQL = `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    name TEXT,
    starts INTEGER,
    messages INTEGER,
	passeds INTEGER
);

CREATE INDEX IF NOT EXISTS users_Name ON users(name);
CREATE INDEX IF NOT EXISTS users_ID ON users(id);
`

type User struct {
	ID   int64
	Name string
}

// DB - это база данных игроков.
type DB struct {
	mu         sync.Mutex // protects following fields
	sql        *sql.DB
	existUser  *sql.Stmt
	add        *sql.Stmt
	UpStarts   *sql.Stmt
	UpMessages *sql.Stmt
	UpPasseds  *sql.Stmt
}

// NewDB создает Users для.
// Этот API не является потокобезопасным.
func NewDB(dbFile string) (*DB, error) {
	db := DB{}
	var err error

	db.sql, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if _, err = db.sql.Exec(schemaSQL); err != nil {
		return nil, err
	}

	db.existUser, err = db.sql.Prepare(existUser)
	if err != nil {
		return nil, err
	}

	db.add, err = db.sql.Prepare(add)
	if err != nil {
		return nil, err
	}

	db.UpStarts, err = db.sql.Prepare(upStarts)
	if err != nil {
		return nil, err
	}

	db.UpMessages, err = db.sql.Prepare(upMessages)
	if err != nil {
		return nil, err
	}

	db.UpPasseds, err = db.sql.Prepare(upPasseds)
	if err != nil {
		return nil, err
	}
	return &db, nil
}

func (db *DB) Add(ID int64, Name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}
	resp, err := tx.Stmt(db.existUser).Query(ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	var ok bool
	if resp.Scan(&ok); !ok {
		_, err := tx.Stmt(db.add).Exec(ID, Name)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (db *DB) Up(ID int64, stmt *sql.Stmt) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tx, err := db.sql.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Stmt(stmt).Exec(ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Close вносит (посредством Flush) всх игроков в базу данных и закрывает БД.
func (db *DB) Close() {
	db.UpStarts.Close()
	db.UpMessages.Close()
	db.UpPasseds.Close()
	db.sql.Close()
}
