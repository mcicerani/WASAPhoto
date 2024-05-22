package database

import (
	"fmt"
	"github.com/google/uuid"
)

// SetUser crea nuovo user nel database per le tabelle UserDetails, Identifier(attribuendo UserId unico stringa) e UserProfile
func (a *appdbimpl) SetUser(name string) error {
	_, err := a.c.Exec(`INSERT INTO user_details (username) VALUES (?)`, name)
	if err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}

	// Generate a unique user ID
	userID := generateUniqueUserID()

	_, err = a.c.Exec(`INSERT INTO identifier (username, user_id, is_new_user) VALUES (?, ?, ?)`, name, userID, true)
	if err != nil {
		return fmt.Errorf("inserting identifier: %w", err)
	}

	_, err = a.c.Exec(`INSERT INTO user_profile (username, user_id, follower_count, followers, following_count, follows, photos, photos_count, banned_user) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, name, "", 0, "", 0, "", "", 0, "")
	if err != nil {
		return fmt.Errorf("inserting user_profile: %w", err)
	}

	return nil
}

// generateUniqueUserID genera un UserId unico per ogni user del database
func generateUniqueUserID() string {
	return uuid.New().String()
}

// GetUser restituisce i dettagli dell'user in user_profile con username=name
func (a *appdbimpl) GetUserByUsername(name string) (UserProfile, error) {
	var user UserProfile
	err := a.c.QueryRow(`SELECT * FROM user_profile WHERE username = ?`, name).Scan(&user.Username, &user.UserId, &user.FollowerCount, &user.Followers, &user.FollowingCount, &user.Follows, &user.PhotosCount, &user.BannedUser)
	if err != nil {
		return user, fmt.Errorf("selecting user: %w", err)
	}

	return user, nil
}

// GetUserById restituisce i dettagli dell'user in user_profile con UserId=id
func (a *appdbimpl) GetUserById(id string) (UserProfile, error) {
	var user UserProfile
	err := a.c.QueryRow(`SELECT * FROM user_profile WHERE user_id = ?`, id).Scan(&user.Username, &user.UserId, &user.FollowerCount, &user.Followers, &user.FollowingCount, &user.Follows, &user.Photos, &user.PhotosCount, &user.BannedUser)
	if err != nil {
		return user, fmt.Errorf("selecting user: %w", err)
	}

	return user, nil
}

// DeleteUser elimina l'user con username e userId dalle tabelle UserDetails, Identifier e UserProfile
func (a *appdbimpl) DeleteUser(name string, id string) error {
	_, err := a.c.Exec(`DELETE FROM user_details WHERE username = ?`, name)
	if err != nil {
		return fmt.Errorf("deleting user: %w", err)
	}

	_, err = a.c.Exec(`DELETE FROM identifier WHERE user_id = ?`, id)
	if err != nil {
		return fmt.Errorf("deleting identifier: %w", err)
	}

	_, err = a.c.Exec(`DELETE FROM user_profile WHERE user_id = ?`, id)
	if err != nil {
		return fmt.Errorf("deleting user_profile: %w", err)
	}

	return nil
}

// cambia username dell'user con username=name e userid=id controllando prima il corrispondente id in identifier e user_profile e infine aggiornando il campo username in user_details
func (a *appdbimpl) UpdateUsername(name string, id string, newname string) error {
	var user Identifier
	err := a.c.QueryRow(`SELECT * FROM identifier WHERE user_id = ?`, id).Scan(&user.Username, &user.UserId, &user.IsNewUser)
	if err != nil {
		return fmt.Errorf("selecting user: %w", err)
	}

	_, err = a.c.Exec(`UPDATE user_details SET username = ? WHERE username = ?`, newname, name)
	if err != nil {
		return fmt.Errorf("updating user_details: %w", err)
	}

	_, err = a.c.Exec(`UPDATE identifier SET username = ? WHERE username = ?`, newname, name)
	if err != nil {
		return fmt.Errorf("updating identifier: %w", err)
	}

	_, err = a.c.Exec(`UPDATE user_profile SET username = ? WHERE username = ?`, newname, name)
	if err != nil {
		return fmt.Errorf("updating user_profile: %w", err)
	}

	return nil
}

// aggiunge alla lista dei follows l'utente da seguire
func (a *appdbimpl) FollowUser(userID string, followedUserID string) error {

	_, err := a.c.Exec("UPDATE user_profile SET follows = array_append(follows, $2) WHERE username = $1")
	if err != nil {
		return fmt.Errorf("adding follow to: %w", err)
	}

	_, err = a.c.Exec("UPDATE user_profile SET followers = array_append(followers, $1) WHERE username = $2")
	if err != nil {
		return fmt.Errorf("adding follower to: %s", err)
	}

	_, err = a.c.Exec("UPDATE user_profile SET following_count = following_count + 1 WHERE username = $1")
	if err != nil {
		return fmt.Errorf("adding follow to following count")
	}

	_, err = a.c.Exec("UPDATE user_profile SET follower_count = follower_count + 1 WHERE username = $2")
	if err != nil {
		return fmt.Errorf("adding follower to followers count")
	}

	return nil
}

// rimuove dalla lista dei follows l'utente
func (a *appdbimpl) UnfollowUser(UserID string, followedUserID string) error {
	_, err := a.c.Exec("UPDATE user_profile SET follows = array_remove(follows, $2) WHERE username = $1")
	if err != nil {
		return fmt.Errorf("removing follow to: %w", err)
	}

	_, err = a.c.Exec("UPDATE user_profile SET followers = array_remove(followers, $1) WHERE username = $2")
	if err != nil {
		return fmt.Errorf("removing follower to: %s", err)
	}

	_, err = a.c.Exec("UPDATE user_profile SET following_count = following_count - 1 WHERE username = $1")
	if err != nil {
		return fmt.Errorf("remove follow from following count")
	}

	_, err = a.c.Exec("UPDATE user_profile SET follower_count = follower_count -1 1 WHERE username = $2")
	if err != nil {
		return fmt.Errorf("removing follower from followers count")
	}

	return nil
}

// GetFollowersByUserID restituisce i dettagli dei followers in user_profile userid=id
func (a *appdbimpl) GetFollowersByUserID(userID string) ([]UserProfile, error) {
	rows, err := a.c.Query(`SELECT * FROM user_profile WHERE user_id = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("selecting followers: %w", err)
	}
	defer rows.Close()

	var followers []UserProfile
	for rows.Next() {
		var follower UserProfile
		err := rows.Scan(&follower.Username, &follower.UserId, &follower.FollowerCount, &follower.Followers, &follower.FollowingCount, &follower.Follows, &follower.Photos, &follower.PhotosCount, &follower.BannedUser)
		if err != nil {
			return nil, fmt.Errorf("scanning followers: %w", err)
		}
		followers = append(followers, follower)
	}
	return followers, nil
}

// GetFollowsByUserID restituisce i dettagli dei follows in user_profile con username=name
func (a *appdbimpl) GetFollowsByUserID(userID string) ([]UserProfile, error) {
	rows, err := a.c.Query(`SELECT * FROM user_profile WHERE user_id = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("selecting follows: %w", err)
	}
	defer rows.Close()

	var follows []UserProfile
	for rows.Next() {
		var follow UserProfile
		err := rows.Scan(&follow.Username, &follow.UserId, &follow.FollowerCount, &follow.Followers, &follow.FollowingCount, &follow.Follows, &follow.Photos, &follow.PhotosCount, &follow.BannedUser)
		if err != nil {
			return nil, fmt.Errorf("scanning follows: %w", err)
		}
		follows = append(follows, follow)
	}
	return follows, nil
}

// aggiunge alla lista dei ban l'utente da seguire
func (a *appdbimpl) BanUser(userID string, bannedUserID string) error {
	_, err := a.c.Exec("UPDATE user_profile SET banned_user = array_append(banned_user, $2) WHERE username = $1")
	if err != nil {
		return fmt.Errorf("adding ban to: %w", err)
	}
	return nil
}

// rimuove dalla lista dei ban l'utente da seguire
func (a *appdbimpl) UnbanUser(userID string, bannedUseriID string) error {
	_, err := a.c.Exec("UPDATE user_profile SET banned_user = array_remove(banned_user, $2) WHERE username = $1")
	if err != nil {
		return fmt.Errorf("removing ban to: %w", err)
	}
	return nil
}




