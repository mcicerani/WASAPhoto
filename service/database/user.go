package database

import (
	"fmt"
	"github.com/google/uuid"
)


//SetUser crea nuovo user nel database per le tabelle UserDetails, Identifier(attribuendo UserId unico stringa) e UserProfile
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

	_, err = a.c.Exec(`INSERT INTO user_profile (username, follower_count, followers, following_count, follows, photos_count, banned_user) VALUES (?, ?, ?, ?, ?, ?, ?)`, name, 0, "", 0, "", 0, "")
	if err != nil {
		return fmt.Errorf("inserting user_profile: %w", err)
	}

	return nil
}


//generateUniqueUserID genera un UserId unico per ogni user del database
func generateUniqueUserID() string {
    return uuid.New().String()
}


//GetUser restituisce i dettagli dell'user in user_profile con username=name
func (a *appdbimpl) GetUserByUsername(name string) (UserProfile, error) {
	var user UserProfile
	err := a.c.QueryRow(`SELECT * FROM user_profile WHERE username = ?`, name).Scan(&user.Username, &user.UserId, &user.FollowerCount, &user.Followers, &user.FollowingCount, &user.Follows, &user.PhotosCount, &user.BannedUser)
	if err != nil {
		return user, fmt.Errorf("selecting user: %w", err)
	}

	return user, nil
}

//GetUserById restituisce i dettagli dell'user in user_profile con UserId=id
func (a *appdbimpl) GetUserById(id string) (UserProfile, error) {
	var user UserProfile
	err := a.c.QueryRow(`SELECT * FROM user_profile WHERE user_id = ?`, id).Scan(&user.Username, &user.UserId, &user.FollowerCount, &user.Followers, &user.FollowingCount, &user.Follows, &user.PhotosCount, &user.BannedUser)
	if err != nil {
		return user, fmt.Errorf("selecting user: %w", err)
	}

	return user, nil
}

//DeleteUser elimina l'user con username e userId dalle tabelle UserDetails, Identifier e UserProfile
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


//cambia username dell'user con username=name e userid=id controllando prima il corrispondente id in identifier e user_profile e infine aggiornando il campo username in user_details
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
