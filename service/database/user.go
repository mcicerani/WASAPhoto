package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
)

// SetUser crea un nuovo utente nel database.
func (a *appdbimpl) SetUser(name string) error {
	_, err := a.c.Exec(`INSERT INTO users (username) VALUES (?)`, name)
	if err != nil {
		// Controlla se l'errore è dovuto alla violazione di unicità del nome utente.
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("username already exists: %w", err)
		}

		// Altrimenti, gestisci l'errore generico di inserimento utente.
		return fmt.Errorf("inserting user: %w", err)
	}

	// Se l'inserimento ha avuto successo, ritorna nil.
	return nil
}

// GetUserByUsername restituisce i dettagli dell'user in users con username=name
func (a *appdbimpl) GetUserByUsername(name string) (User, error) {
	var user User
	err := a.c.QueryRow(`SELECT * FROM users WHERE username = ?`, name).Scan(&user.ID, &user.Username)
	if err != nil {
		return user, fmt.Errorf("selecting user: %w", err)
	}

	return user, nil
}

// GetUserById restituisce i dettagli dell'utente con l'ID specificato
func (a *appdbimpl) GetUserById(iD string) (User, error) {
	var user User

	userID, err := strconv.Atoi(iD)
	if err != nil {
		return user, fmt.Errorf("converting user ID to integer: %w", err)
	}

	log.Printf("User ID estratto: %d\n", userID) // Log dell'ID estratto

	err = a.c.QueryRow(`SELECT * FROM users WHERE ID = ?`, userID).Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Nessuna riga trovata per l'ID specificato")
		} else {
			log.Printf("Errore nella query SQL: %v\n", err) // Log dell'errore SQL
		}
		return user, fmt.Errorf("selecting user: %w", err)
	}

	return user, nil
}

// DeleteUser elimina l'user con username e userId dalle tabelle UserDetails, Identifier e User
func (a *appdbimpl) DeleteUser(username string) error {
	_, err := a.c.Exec(`DELETE FROM users WHERE username = ?`, username)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}
	return nil
}

// UpdateUsername cambia username dell'user con username=newname controllando prima il corrispondente id in users
func (a *appdbimpl) UpdateUsername(iD string, newname string) error {

	userID, err := strconv.Atoi(iD)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)
	}

	_, err = a.c.Exec(`UPDATE users SET username = ? WHERE ID = ?`, newname, userID)
	if err != nil {
		return fmt.Errorf("updating username: %w", err)
	}

	return nil
}

// FollowUser crea nella tabella la relazione followed/follower
func (a *appdbimpl) FollowUser(userID string, followedUserID string) error {
	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)
	}

	FollowedUserID, err := strconv.Atoi(followedUserID)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)
	}

	// Log dei valori di userID e followedUserID per debug
	log.Printf("FollowUser: userID = %d, followedUserID = %d", UserID, FollowedUserID)

	_, err = a.c.Exec(`INSERT into followers (follower_id, followed_id) VALUES(?, ?)`, UserID, FollowedUserID)
	if err != nil {
		return fmt.Errorf("following user: %w", err)
	}

	// Log per vedere quando viene eseguito il follow tra due utenti
	log.Printf("User %d followed user %d successfully", UserID, FollowedUserID)

	return nil
}

// UnfollowUser cancella dalla tabella followers la relazione tra i 2 account
func (a *appdbimpl) UnfollowUser(userID string, followedUserID string) error {

	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)
	}

	FollowedUserID, err := strconv.Atoi(followedUserID)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)
	}

	_, err = a.c.Exec(`DELETE FROM followers WHERE follower_id = ? AND followed_id = ?`, UserID, FollowedUserID)
	if err != nil {
		return fmt.Errorf("unfollowing user: %w", err)
	}

	return nil
}

// GetFollowers restituisce lista dei followers per followed_id=userID contando i follower da followers
func (a *appdbimpl) GetFollowers(userID string) ([]User, error) {
	var followers []User

	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return followers, fmt.Errorf("converting user ID to integer: %w", err)
	}

	rows, err := a.c.Query(`SELECT follower_id FROM followers WHERE followed_id = ?`, UserID)
	if err != nil {
		return followers, fmt.Errorf("selecting followers: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows:", err)
			return
		}
	}(rows) // Ensure rows are closed after function returns

	for rows.Next() {
		var followerID int
		err = rows.Scan(&followerID)
		if err != nil {
			return followers, fmt.Errorf("scanning follower ID: %w", err)
		}

		var follower User
		err = a.c.QueryRow(`SELECT ID, username FROM users WHERE ID = ?`, followerID).Scan(&follower.ID, &follower.Username)
		if err != nil {
			return followers, fmt.Errorf("selecting follower: %w", err)
		}

		followers = append(followers, follower)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return followers, nil
}

// GetFollows restituisce i dettagli dei follows in user_profile con username=name
func (a *appdbimpl) GetFollows(userID string) ([]User, error) {
	var follows []User

	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return follows, fmt.Errorf("converting user ID to integer: %w", err)
	}

	log.Printf("Getting follows for user ID: %s", userID)

	rows, err := a.c.Query(`SELECT followed_id FROM followers WHERE follower_id = ?`, UserID)
	if err != nil {
		return follows, fmt.Errorf("selecting follows: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows:", err)
		}
	}(rows) // Ensure rows are closed after function returns

	for rows.Next() {
		var followedID int
		err = rows.Scan(&followedID)
		if err != nil {
			return follows, fmt.Errorf("scanning followed ID: %w", err)
		}

		var followed User
		// Query per ottenere il nome utente dell'utente seguito
		err = a.c.QueryRow(`SELECT ID, username FROM users WHERE ID = ?`, followedID).Scan(&followed.ID, &followed.Username)
		if err != nil {
			return follows, fmt.Errorf("selecting followed: %w", err)
		}

		// Log dopo aver aggiunto ciascun utente seguito alla lista follows
		log.Printf("Added followed user: ID = %d, Username = %s", followed.ID, followed.Username)

		follows = append(follows, followed)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return follows, nil
}

// IsFollowed controlla se l'utente segue un altro utente
func (a *appdbimpl) IsFollowed(userID string, otherUserID string) (bool, error) {
	// Converti userID e otherUserID in interi
	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return false, fmt.Errorf("converting user ID to integer: %w", err)
	}

	OtherUserID, err := strconv.Atoi(otherUserID)
	if err != nil {
		return false, fmt.Errorf("converting other user ID to integer: %w", err)
	}

	// Esegui la query per verificare se l'utente segue l'altro utente
	var exists bool
	err = a.c.QueryRow(`SELECT EXISTS (SELECT 1 FROM followers WHERE followed_id = ? AND follower_id = ?)`, OtherUserID, UserID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking follow: %w", err)
	}

	// Se exists è true, l'utente segue l'altro utente; altrimenti, non lo segue
	return exists, nil
}

// BanUser aggiunge alla lista dei ban l'utente da seguire
func (a *appdbimpl) BanUser(userID string, bannedUserID string) error {

	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)
	}

	BannedUserID, err := strconv.Atoi(bannedUserID)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)
	}

	_, err = a.c.Exec(`INSERT into bans (user_id, banned_id) VALUES(?, ?)`, UserID, BannedUserID)
	if err != nil {
		return fmt.Errorf("banning user: %w", err)
	}

	return nil
}

// UnbanUser rimuove dalla lista dei ban l'utente da seguire
func (a *appdbimpl) UnbanUser(userID string, bannedUserID string) error {

	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)
	}

	BannedUserID, err := strconv.Atoi(bannedUserID)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)
	}

	_, err = a.c.Exec(`DELETE FROM bans WHERE user_id = ? AND banned_id = ?`, UserID, BannedUserID)
	if err != nil {
		return fmt.Errorf("unbanning user: %w", err)
	}

	return nil
}

// IsBanned controlla se l'utente è stato bannato da un altro utente specifico e restituisce true o false
func (a *appdbimpl) IsBanned(userID string, otherUserID string) (bool, error) {
	// Converti userID e otherUserID in interi
	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return false, fmt.Errorf("converting user ID to integer: %w", err)
	}

	OtherUserID, err := strconv.Atoi(otherUserID)
	if err != nil {
		return false, fmt.Errorf("converting other user ID to integer: %w", err)
	}

	// Esegui la query per verificare se l'utente è bannato
	var exists bool
	err = a.c.QueryRow(`SELECT EXISTS (SELECT 1 FROM bans WHERE user_id = ? AND banned_id = ?)`, OtherUserID, UserID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking ban: %w", err)
	}

	// Se exists è true, l'utente è bannato; altrimenti, non è bannato
	return exists, nil
}

// GetBans restituisce la lista degli utenti bannati da un determinato utente
func (a *appdbimpl) GetBans(userID string) ([]User, error) {
	var bans []User

	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return bans, fmt.Errorf("converting user ID to integer: %w", err)
	}

	rows, err := a.c.Query(`SELECT banned_id FROM bans WHERE user_id = ?`, UserID)
	if err != nil {
		return bans, fmt.Errorf("selecting bans: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows:", err)
		}
	}(rows) // Ensure rows are closed after function returns

	for rows.Next() {
		var bannedID int
		err = rows.Scan(&bannedID)
		if err != nil {
			return bans, fmt.Errorf("scanning banned ID: %w", err)
		}

		var banned User
		err = a.c.QueryRow(`SELECT ID, username FROM users WHERE ID = ?`, bannedID).Scan(&banned.ID, &banned.Username)
		if err != nil {
			return bans, fmt.Errorf("selecting banned: %w", err)
		}

		bans = append(bans, banned)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return bans, nil
}

// CountFollowersByUserID restituisce il numero di followers di un utente
func (a *appdbimpl) CountFollowersByUserID(userID string) (int, error) {

	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return 0, fmt.Errorf("converting user ID to integer: %w", err)
	}

	rows, err := a.c.Query(`SELECT COUNT(*) FROM followers WHERE followed_id = ?`, UserID)
	if err != nil {
		return 0, fmt.Errorf("selecting followers: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows:", err)
			return
		}
	}(rows) // Ensure rows are closed after function returns

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("scanning follower: %w", err)
		}
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return 0, fmt.Errorf("iterating rows: %w", err)
	}

	return count, nil
}

// CountFollowsByUserID restituisce il numero di follows di un utente

func (a *appdbimpl) CountFollowsByUserID(userID string) (int, error) {

	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return 0, fmt.Errorf("converting user ID to integer: %w", err)
	}

	rows, err := a.c.Query(`SELECT COUNT(*) FROM followers WHERE follower_id = ?`, UserID)
	if err != nil {
		return 0, fmt.Errorf("selecting follows: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Error closing rows:", err)
			return
		}
	}(rows) // Ensure rows are closed after function returns

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("scanning follow: %w", err)
		}
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return 0, fmt.Errorf("iterating rows: %w", err)
	}

	return count, nil
}
