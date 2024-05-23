package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/mattn/go-sqlite3"
)

// SetUser crea nuovo user nel database
func (a *appdbimpl) SetUser(name string) error {

	_, err := a.c.Exec(`INSERT INTO users (username) VALUES (?)`, name)
	if err != nil {
		// Controlla se l'errore è dovuto alla violazione del vincolo UNIQUE
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			return errors.New("username already exists")
		}
		return fmt.Errorf("inserting user: %w", err)
	}
	return nil
}

// GetUserByUsername restituisce i dettagli dell'user in users con username=name
func (a *appdbimpl) GetUserByUsername(name string) (User, error) {
	var user User
	err := a.c.QueryRow(`SELECT * FROM users WHERE username = ?`, name).Scan(&user.Username, &user.ID)
	if err != nil {
		return user, fmt.Errorf("selecting user: %w", err)
	}

	return user, nil
}

// GetUserById restituisce i dettagli dell'user in user_profile con UserId=id
func (a *appdbimpl) GetUserById(id string) (User, error) {
	var user User

	userID, err := strconv.Atoi(id)
	if err != nil {
		return user, fmt.Errorf("converting user ID to integer: %w", err)
	}

	err = a.c.QueryRow(`SELECT * FROM users WHERE ID = ?`, userID).Scan(&user.Username, &user.ID)
	if err != nil {
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

// UpdateUsername  cambia username dell'user con username=newname controllando prima il corrispondente id in users
func (a *appdbimpl) UpdateUsername(id string, newname string) error {

	userID, err := strconv.Atoi(id)
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
func (a *appdbimpl) FollowUser(id string, followedUserID string) error {

	userID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)
	}

	followedID, err := strconv.Atoi(followedUserID)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)
	}

	_, err = a.c.Exec(`INSERT into followers (followerID, followedID) VALUES(?, ?)`, userID, followedID)
	if err != nil {
		return fmt.Errorf("following user: %w", err)

	}
	return nil
}

// --------> UnfollowUser cancella dalla tabella followers la relazione tra i 2 account
func (a *appdbimpl) UnfollowUser(UserID string, followedUserID string) error {

	return nil
}

// GetFollowersByUserID restituisce i dettagli dei followers in user_profile userid=id
func (a *appdbimpl) GetFollowersByUserID(userID string) ([]User, error) {
	rows, err := a.c.Query(`SELECT * FROM user_profile WHERE user_id = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("selecting followers: %w", err)
	}

	var followers []User
	for rows.Next() {
		var follower User
		err := rows.Scan(&follower.Username, &follower.UserId, &follower.FollowerCount, &follower.Followers, &follower.FollowingCount, &follower.Follows, &follower.Photos, &follower.PhotosCount, &follower.BannedUser)
		if err != nil {
			return nil, fmt.Errorf("scanning followers: %w", err)
		}
		followers = append(followers, follower)
	}
	return followers, nil
}

// GetFollowsByUserID restituisce i dettagli dei follows in user_profile con username=name
func (a *appdbimpl) GetFollowsByUserID(userID string) ([]User, error) {
	rows, err := a.c.Query(`SELECT * FROM user_profile WHERE user_id = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("selecting follows: %w", err)
	}

	var follows []User
	for rows.Next() {
		var follow User
		err := rows.Scan(&follow.Username, &follow.UserId, &follow.FollowerCount, &follow.Followers, &follow.FollowingCount, &follow.Follows, &follow.Photos, &follow.PhotosCount, &follow.BannedUser)
		if err != nil {
			return nil, fmt.Errorf("scanning follows: %w", err)
		}
		follows = append(follows, follow)
	}
	return follows, nil
}

// BanUser  aggiunge alla lista dei ban l'utente da seguire
func (a *appdbimpl) BanUser(userID string, bannedUserID string) error {
	_, err := a.c.Exec("UPDATE user_profile SET banned_user = json_insert(banned_user, '$[#]' , ?) WHERE user_id = ?", bannedUserID, userID)
	if err != nil {
		return fmt.Errorf("adding ban to: %w", err)
	}
	return nil
}

// UnbanUser  rimuove dalla lista dei ban l'utente da seguire
func (a *appdbimpl) UnbanUser(userID string, bannedUserID string) error {

	//Ricava la lista banned_user corrente come json
	var bannedUsersJSON string
	err := a.c.QueryRow("SELECT banned_user FROM user_profile WHERE user_id = ?", userID).Scan(&bannedUsersJSON)
	if err != nil {
		return fmt.Errorf("retrieving banned_user array: %w", err)
	}

	// Analizza array
	var bannedUsers []string
	if bannedUsersJSON != "" {
		if err := json.Unmarshal([]byte(bannedUsersJSON), &bannedUsers); err != nil {
			return fmt.Errorf("parsing banned_user array: %w", err)
		}
	}

	// Rimove utente da sbannare
	var newBannedUsers []string
	for _, id := range bannedUsers {
		if id != bannedUserID {
			newBannedUsers = append(newBannedUsers, id)
		}
	}

	// Riconverte in JSON l'array
	newBannedUsersJSON, err := json.Marshal(newBannedUsers)
	if err != nil {
		return fmt.Errorf("serializing banned_user array: %w", err)
	}

	// Update del database con nuovo array
	_, err = a.c.Exec("UPDATE user_profile SET banned_user = ? WHERE user_id = ?", string(newBannedUsersJSON), userID)
	if err != nil {
		return fmt.Errorf("updating banned_user array: %w", err)
	}

	return nil
}
