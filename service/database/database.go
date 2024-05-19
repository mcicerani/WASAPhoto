/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	//User
	SetUser(name string) error
	UpdateUsername(name string, id string, newname string) error
	GetUserByUsername(name string) (UserProfile, error)
	GetUserById(id string) (UserProfile, error)
	DeleteUser(name string, id string) error

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if the database is reachable
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// Create the tables if they don't exist
	//User table
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS user_details (
		username TEXT
	)`)
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	//Identifier table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS identifier (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		is_new_user BOOLEAN
	)`)
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	//UserProfile table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_profile (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		follower_count INTEGER,
		followers TEXT,
		following_count INTEGER,
		follows TEXT,
		photos_count INTEGER,
		banned_user TEXT
	)`)
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	//Photo table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS photos (
		photos_id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT,
		timestamp TEXT,
		likes_number INTEGER
	)`)
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	//Comment table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS comment (
		comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		photos_id INTEGER,
		text_comment TEXT
	)`)
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
