package database

import (
	"fmt"
	"strconv"
)

// SetUser crea nuovo user nel database
func (a *appdbimpl) SetUser(name string) error {

	_, err := a.c.Exec(`INSERT INTO users (username) VALUES (?)`, name)
	if err != nil {

		if err.Error() == "UNIQUE constraint failed: users.username" {
			return fmt.Errorf("username already exists: %w", err)
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
func (a *appdbimpl) GetUserById(ID string) (User, error) {
	var user User

	userID, err := strconv.Atoi(ID)
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
func (a *appdbimpl) UpdateUsername(ID string, newname string) error {

	userID, err := strconv.Atoi(ID)
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

	_, err = a.c.Exec(`INSERT into followers (follower_id, followed_id) VALUES(?, ?)`, UserID, FollowedUserID)
	if err != nil {
		return fmt.Errorf("following user: %w", err)

	}
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

	for rows.Next() {
		var followerID int
		err = rows.Scan(&followerID)
		if err != nil {
			return followers, fmt.Errorf("scanning follower ID: %w", err)
		}

		var follower User
		err = a.c.QueryRow(`SELECT username FROM users WHERE ID = ?`, followerID).Scan(&follower.Username)
		if err != nil {
			return followers, fmt.Errorf("selecting follower: %w", err)
		}

		followers = append(followers, follower)
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

	rows, err := a.c.Query(`SELECT followed_id FROM followers WHERE follower_id = ?`, UserID)
	if err != nil {
		return follows, fmt.Errorf("selecting follows: %w", err)
	}

	for rows.Next() {
		var followedID int
		err = rows.Scan(&followedID)
		if err != nil {
			return follows, fmt.Errorf("scanning followed ID: %w", err)
		}

		var followed User
		err = a.c.QueryRow(`SELECT username FROM users WHERE ID = ?`, followedID).Scan(&followed.Username)
		if err != nil {
			return follows, fmt.Errorf("selecting followed: %w", err)
		}

		follows = append(follows, followed)
	}

	return follows, nil
}

// BanUser  aggiunge alla lista dei ban l'utente da seguire
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

// UnbanUser  rimuove dalla lista dei ban l'utente da seguire
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

//CountFollowersByUserID restituisce il numero di followers di un utente

func (a *appdbimpl) CountFollowersByUserID(userID string) (int, error) {

	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return 0, fmt.Errorf("converting user ID to integer: %w", err)
	}

	rows, err := a.c.Query(`SELECT COUNT(*) FROM followers WHERE followed_id = ?`, UserID)
	if err != nil {
		return 0, fmt.Errorf("selecting followers: %w", err)
	}

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("scanning follower: %w", err)
		}
	}

	return count, nil
}

//CountFollowsByUserID restituisce il numero di follows di un utente

func (a *appdbimpl) CountFollowsByUserID(userID string) (int, error) {

	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return 0, fmt.Errorf("converting user ID to integer: %w", err)
	}

	rows, err := a.c.Query(`SELECT COUNT(*) FROM followers WHERE follower_id = ?`, UserID)
	if err != nil {
		return 0, fmt.Errorf("selecting follows: %w", err)
	}

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("scanning follow: %w", err)
		}
	}

	return count, nil
}
