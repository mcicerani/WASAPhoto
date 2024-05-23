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

	err = a.c.QueryRow(`SELECT * FROM users WHERE user_id = ?`, userID).Scan(&user.Username, &user.ID)
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

// ------> UpdateUsername  cambia username dell'user con username=name e userid=id controllando prima il corrispondente id in identifier e user_profile e infine aggiornando il campo username in user_details
func (a *appdbimpl) UpdateUsername(name string, id string, newname string) error {
	var user User
	err := a.c.QueryRow(`SELECT * FROM identifier WHERE user_id = ?`, id).Scan(&user.Username, &user.UserId, &user.IsNewUser)
	if err != nil {
		return fmt.Errorf("selecting user: %w", err)
	}

	return nil
}

// FollowUser aggiunge alla lista dei follows l'utente da seguire
func (a *appdbimpl) FollowUser(userID string, followedUserID string) error {

	_, err := a.c.Exec("UPDATE user_profile SET follows = json_insert(follows, '$[#]', ?) WHERE user_id = ?", followedUserID, userID)
	if err != nil {
		return fmt.Errorf("adding follow to: %w", err)
	}

	_, err = a.c.Exec("UPDATE user_profile SET followers = json_insert(followers, '$[#]', ?) WHERE user_id = ?", userID, followedUserID)
	if err != nil {
		return fmt.Errorf("adding follower to: %s", err)
	}

	_, err = a.c.Exec("UPDATE user_profile SET following_count = following_count + 1 WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("adding follow to following count")
	}

	_, err = a.c.Exec("UPDATE user_profile SET follower_count = follower_count + 1 WHERE user_id = ?", followedUserID)
	if err != nil {
		return fmt.Errorf("adding follower to followers count")
	}

	return nil
}

// UnfollowUser  rimuove dalla lista dei follows l'utente
func (a *appdbimpl) UnfollowUser(UserID string, followedUserID string) error {

	//Rimuove dalla lista dei follows

	// Recupera l'array JSON corrente dalla colonna `follows'
	var followsJSON string
	err := a.c.QueryRow("SELECT follows FROM user_profile WHERE user_id = ?", UserID).Scan(&followsJSON)
	if err != nil {
		return fmt.Errorf("retrieving follows array: %w", err)
	}

	// Parsing dell'array JSON
	var follows []string
	if followsJSON != "" {
		if err := json.Unmarshal([]byte(followsJSON), &follows); err != nil {
			return fmt.Errorf("parsing follows array: %w", err)
		}
	}

	// Rimuovi l'elemento specificato dall'array
	var newFollows []string
	for _, id := range follows {
		if id != followedUserID {
			newFollows = append(newFollows, id)
		}
	}

	// Converti l'array aggiornato di nuovo in JSON
	newFollowsJSON, err := json.Marshal(newFollows)
	if err != nil {
		return fmt.Errorf("serializing follows array: %w", err)
	}

	// Aggiorna la colonna con il nuovo array JSON
	_, err = a.c.Exec("UPDATE user_profile SET follows = ? WHERE user_id = ?", string(newFollowsJSON), UserID)
	if err != nil {
		return fmt.Errorf("updating follows array: %w", err)
	}

	//Rimuove dalla lista dei followers dell'altro utente

	// Recupera l'array JSON corrente dalla colonna `followers`
	var followersJSON string
	err = a.c.QueryRow("SELECT followers FROM user_profile WHERE user_id = ?", followedUserID).Scan(&followersJSON)
	if err != nil {
		return fmt.Errorf("retrieving followers array: %w", err)
	}

	// Parsing dell'array JSON
	var followers []string
	if followersJSON != "" {
		if err := json.Unmarshal([]byte(followersJSON), &followers); err != nil {
			return fmt.Errorf("parsing followers array: %w", err)
		}
	}

	// Rimuovi l'elemento specificato dall'array
	var newFollowers []string
	for _, id := range followers {
		if id != followedUserID {
			newFollowers = append(newFollowers, id)
		}
	}

	// Converti l'array aggiornato di nuovo in JSON
	newFollowersJSON, err := json.Marshal(newFollowers)
	if err != nil {
		return fmt.Errorf("serializing followers array: %w", err)
	}

	// Aggiorna la colonna con il nuovo array JSON
	_, err = a.c.Exec("UPDATE user_profile SET followers = ? WHERE user_id = ?", string(newFollowersJSON), followedUserID)
	if err != nil {
		return fmt.Errorf("updating followers array: %w", err)
	}

	//Aggiorna follow count
	_, err = a.c.Exec("UPDATE user_profile SET following_count = following_count - 1 WHERE user_id = ?", UserID)
	if err != nil {
		return fmt.Errorf("remove follow from following count")
	}

	//Aggiorna follower count
	_, err = a.c.Exec("UPDATE user_profile SET follower_count = follower_count -1 WHERE user_id = ?", followedUserID)
	if err != nil {
		return fmt.Errorf("removing follower from followers count")
	}

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
