package database

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

/*SetPhoto salva la foto in locale e inserisce dati in photos (id, user_id, url, timestamp) */

func (a *appdbimpl) SetPhoto(userID string, url string) error {

	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)
	}

	_, err = a.c.Exec(`INSERT INTO photos (user_id, url) VALUES (?, ?)`, UserID, url)
	if err != nil {
		return fmt.Errorf("inserting photo: %w", err)
	}

	return nil
}

// GetPhotoByID restituisce i dettagli della foto in photos con photos_id=id
func (a *appdbimpl) GetPhotoByID(photoID string) (Photo, error) {

	var photo Photo

	PhotoID, err := strconv.Atoi(photoID)
	if err != nil {
		return photo, fmt.Errorf("converting photo ID to integer: %w", err)
	}

	err = a.c.QueryRow(`SELECT * FROM photos WHERE id = ?`, PhotoID).Scan(&photo.ID, &photo.UserID, &photo.URL, &photo.Timestamp)
	if err != nil {
		return photo, fmt.Errorf("selecting photo: %w", err)
	}

	return photo, nil
}

// DeletePhoto elimina la foto con photos_id=id dalla tabella photos
func (a *appdbimpl) DeletePhoto(photoID string) error {

	PhotoID, err := strconv.Atoi(photoID)
	if err != nil {
		return fmt.Errorf("converting photo ID to integer: %w", err)
	}

	_, err = a.c.Exec(`DELETE FROM photos WHERE id = ?`, PhotoID)
	if err != nil {
		return fmt.Errorf("deleting photo: %w", err)
	}

	_, err = a.c.Exec(`DELETE FROM comments WHERE photo_id = ?`, PhotoID)
	if err != nil {
		return fmt.Errorf("deleting comments: %w", err)
	}

	_, err = a.c.Exec(`DELETE FROM likes WHERE photo_id = ?`, PhotoID)
	if err != nil {
		return fmt.Errorf("deleting likes: %w", err)
	}

	return nil
}

// SetComment inserisce un nuovo commento nel database nella tabella comment
func (a *appdbimpl) SetComment(userID string, photoID string, comment string) error {

	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)
	}

	PhotoID, err := strconv.Atoi(photoID)
	if err != nil {
		return fmt.Errorf("converting photo ID to integer: %w", err)
	}

	_, err = a.c.Exec(`INSERT INTO comments (user_id, photo_id, comment_text) VALUES (?, ?, ?)`, UserID, PhotoID, comment)
	if err != nil {
		return fmt.Errorf("inserting comment: %w", err)
	}

	return nil
}

// GetCommentByID restituisce i dettagli del commento in comment con comment_id=id
func (a *appdbimpl) GetCommentByID(commentID string) (Comment, error) {

	var comment Comment

	CommentID, err := strconv.Atoi(commentID)
	if err != nil {
		return comment, fmt.Errorf("converting comment ID to integer: %w", err)
	}

	err = a.c.QueryRow(`SELECT * FROM comments WHERE id = ?`, CommentID).Scan(&comment.ID, &comment.UserId, &comment.PhotoId, &comment.CommentText, &comment.Timestamp)
	if err != nil {
		return comment, fmt.Errorf("selecting comment: %w", err)
	}

	return comment, nil
}

// DeleteComment elimina il commento con comment_id=id dalla tabella comment
func (a *appdbimpl) DeleteComment(commentID string) error {

	CommentID, err := strconv.Atoi(commentID)
	if err != nil {
		return fmt.Errorf("converting comment ID to integer: %w", err)
	}

	_, err = a.c.Exec(`DELETE FROM comments WHERE id = ?`, CommentID)
	if err != nil {
		return fmt.Errorf("deleting comment: %w", err)
	}

	return nil
}

// GetCommentsByPhotoID restituisce i dettagli dei commenti in comment con photos_id=id
func (a *appdbimpl) GetCommentsByPhotoID(photoID string) ([]Comment, error) {

	PhotoID, err := strconv.Atoi(photoID)
	if err != nil {
		return nil, fmt.Errorf("converting photo ID to integer: %w", err)
	}

	rows, err := a.c.Query(`SELECT * FROM comments WHERE photo_id = ?`, PhotoID)
	if err != nil {
		return nil, fmt.Errorf("selecting comments: %w", err)
	}

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.ID, &comment.UserId, &comment.PhotoId, &comment.CommentText, &comment.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("scanning comment: %w", err)
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

// GetPhotosByUserID restituisce i dettagli delle foto in photos con user_id=id
func (a *appdbimpl) GetPhotosByUserID(userId string) ([]Photo, error) {

	UserID, err := strconv.Atoi(userId)
	if err != nil {
		return nil, fmt.Errorf("converting user ID to integer: %w", err)
	}

	rows, err := a.c.Query(`SELECT * FROM photos WHERE user_id = ?`, UserID)
	if err != nil {
		return nil, fmt.Errorf("selecting photos: %w", err)
	}

	var photos []Photo

	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.ID, &photo.UserID, &photo.URL, &photo.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("scanning photo: %w", err)
		}

		photos = append(photos, photo)
	}

	return photos, nil
}

// SetLike incrementa il numero di like di una foto
func (a *appdbimpl) SetLike(userId string, photoID string) error {

	UserID, err := strconv.Atoi(userId)
	if err != nil {
		return fmt.Errorf("converting user ID to integer: %w", err)

	}

	PhotoID, err := strconv.Atoi(photoID)
	if err != nil {
		return fmt.Errorf("converting photo ID to integer: %w", err)
	}

	_, err = a.c.Exec(`INSERT INTO likes (user_id, photo_id) VALUES (?, ?)`, UserID, PhotoID)
	if err != nil {
		return fmt.Errorf("inserting like: %w", err)
	}

	return nil
}

// DeleteLike decrementa il numero di like di una foto
func (a *appdbimpl) DeleteLike(likeID string) error {

	LikeID, err := strconv.Atoi(likeID)
	if err != nil {
		return fmt.Errorf("converting like ID to integer: %w", err)
	}

	_, err = a.c.Exec(`DELETE FROM likes WHERE id = ?`, LikeID)
	if err != nil {
		return fmt.Errorf("deleting like: %w", err)
	}

	return nil
}

// GetLikesByPhotoID restituisce il numero di like di una foto
func (a *appdbimpl) GetLikesByPhotoID(photoID string) ([]Like, error) {

	PhotoID, err := strconv.Atoi(photoID)
	if err != nil {
		return nil, fmt.Errorf("converting photo ID to integer: %w", err)
	}

	rows, err := a.c.Query(`SELECT * FROM likes WHERE photo_id = ?`, PhotoID)
	if err != nil {
		return nil, fmt.Errorf("selecting likes: %w", err)
	}

	var likes []Like
	for rows.Next() {
		var like Like
		err = rows.Scan(&like.ID, &like.UserID, &like.PhotoID)
		if err != nil {
			return nil, fmt.Errorf("scanning like: %w", err)
		}

		likes = append(likes, like)
	}

	return likes, nil
}

// GetPhotosStreamByUserID restituisce lista foto in ordine cronologico inverso di tutti account seguiti da userID
func (a *appdbimpl) GetPhotosStreamByUserID(userID string) ([]Photo, error) {

	follows, err := a.GetFollows(userID)

	if err != nil {
		return nil, fmt.Errorf("getting follows: %w", err)
	}

	var photos []Photo
	for _, follow := range follows {
		rows, err := a.c.Query(`SELECT * FROM photos WHERE user_id = ?`, follow.ID)
		if err != nil {
			return nil, fmt.Errorf("selecting photos: %w", err)
		}

		for rows.Next() {
			var photo Photo
			err = rows.Scan(&photo.ID, &photo.UserID, &photo.URL, &photo.Timestamp)
			if err != nil {
				return nil, fmt.Errorf("scanning photo: %w", err)
			}

			photos = append(photos, photo)
		}

		sort.Slice(photos, func(i, j int) bool {
			t1, _ := time.Parse(time.RFC3339, photos[i].Timestamp)
			t2, _ := time.Parse(time.RFC3339, photos[j].Timestamp)
			return t2.Before(t1)
		})

	}

	return photos, nil
}

//CountLikesByPhotoID restituisce il numero di like di una foto

func (a *appdbimpl) CountLikesByPhotoID(photoID string) (int, error) {

	PhotoID, err := strconv.Atoi(photoID)
	if err != nil {
		return 0, fmt.Errorf("converting photo ID to integer: %w", err)
	}

	rows, err := a.c.Query(`SELECT COUNT(*) FROM likes WHERE photo_id = ?`, PhotoID)
	if err != nil {
		return 0, fmt.Errorf("selecting likes: %w", err)
	}

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("scanning like: %w", err)
		}
	}

	return count, nil
}

//CountCommentsByPhotoID restituisce il numero di commenti di una foto

func (a *appdbimpl) CountCommentsByPhotoID(photoID string) (int, error) {
	
	PhotoID, err := strconv.Atoi(photoID)
	if err != nil {
		return 0, fmt.Errorf("converting photo ID to integer: %w", err)
	}

	rows, err := a.c.Query(`SELECT COUNT(*) FROM comments WHERE photo_id = ?`, PhotoID)
	if err != nil {
		return 0, fmt.Errorf("selecting comments: %w", err)
	}

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("scanning comment: %w", err)
		}
	}

	return count, nil
}


// CountPhotosByUserID restituisce il numero di foto di un utente

func (a *appdbimpl) CountPhotosByUserID(userID string) (int, error) {
	
	UserID, err := strconv.Atoi(userID)
	if err != nil {
		return 0, fmt.Errorf("converting user ID to integer: %w", err)
	}

	rows, err := a.c.Query(`SELECT COUNT(*) FROM photos WHERE user_id = ?`, UserID)
	if err != nil {
		return 0, fmt.Errorf("selecting photos: %w", err)
	}

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, fmt.Errorf("scanning photo: %w", err)
		}
	}

	return count, nil
}