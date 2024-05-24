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
	UpdateUsername(ID string, newname string) error
	GetUserByUsername(name string) (User, error)
	GetUserById(ID string) (User, error)
	DeleteUser(username string) error
	FollowUser(userID string, followedUserID string) error
	UnfollowUser(userID string, followedUserID string) error
	GetFollowers(userID string) ([]User, error)
	GetFollows(userID string) ([]User, error)
	BanUser(userID string, bannedUserID string) error
	UnbanUser(userID string, bannedUserID string) error
	IsBanned(userID string, otherUserID string) (bool, error)
	CountFollowersByUserID(userID string) (int, error)
	CountFollowsByUserID(userID string) (int, error)
	SetPhoto(userId string, url string) error
	GetPhotoByID(photoID string) (Photo, error)
	DeletePhoto(photoID string) error
	SetComment(userId string, photoID string, comment string) error
	GetCommentByID(commentID string) (Comment, error)
	DeleteComment(commentID string) error
	GetCommentsByPhotoID(photoID string) ([]Comment, error)
	GetPhotosByUserID(userId string) ([]Photo, error)
	SetLike(userId string, photoID string) error
	DeleteLike(likeID string) error
	GetLikesByPhotoID(photoID string) ([]Like, error)
	GetPhotosStreamByUserID(userID string) ([]Photo, error)
	CountCommentsByPhotoID(photoID string) (int, error)
	CountLikesByPhotoID(photoID string) (int, error)
	CountPhotosByUserID(userID string) (int, error)

	//Ping checks if the database is reachable

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
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE
	)`)
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	//Photo table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS photos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		url TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`)
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	//Comment table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		photo_id INTEGER NOT NULL,
		text TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (photo_id) REFERENCES photos(id)
	)`)
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	//Like table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS likes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		photo_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (photo_id) REFERENCES photos(id)
	)`)
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	//followers table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS followers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		follower_id INTEGER NOT NULL,
		followed_id INTEGER NOT NULL,
		FOREIGN KEY (follower_id) REFERENCES users(id),
		FOREIGN KEY (followed_id) REFERENCES users(id)
	)`)
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	//bans table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		banned_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (banned_id) REFERENCES users(id)
	)`)
	if err != nil {
		return nil, fmt.Errorf("creating table: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (a *appdbimpl) Ping() error {
	return a.c.Ping()
}
