package database

import (
	"fmt"
	"time"
	"github.com/google/uuid"
)

// genera id unico per foto e commenti
func generateUniquePhotoID() string {
	return uuid.New().String()
}


func generateUniqueCommentID() string {
	return uuid.New().String()
}

// genera timestamp per foto

func generateTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

/*SetPhoto inserisce una nuova foto nel database nella tabella photos e ne
 inizializza url generando un id unico per la foto*/

func (a *appdbimpl) SetPhoto(userID string, binaryFile string) error {
	photoID := generateUniquePhotoID()
	url := "users/userID/photos/" + photoID
	timestamp := generateTimestamp()
	_, err := a.c.Exec(`INSERT INTO photos (user_id, binary_file, photos_id, url, timestamp, likes_number, comments) VALUES (?, ?, ?, ?, ?, ?)`, userID, binaryFile, photoID, url, timestamp, 0, []string{})
	if err != nil {
		return fmt.Errorf("inserting photo: %w", err)
	}
	return nil
}


// GetPhotoByID restituisce i dettagli della foto in photos con photos_id=id
func (a *appdbimpl) GetPhotoByID(userId string, photoID string) (Photo, error) {
	var photo Photo
	err := a.c.QueryRow(`SELECT * FROM photos WHERE user_id = ? AND photos_id = ?`, userId, photoID).Scan(&photo.UserID, &photo.BinaryFile, &photo.PhotosId, &photo.Url, &photo.Timestamp, &photo.LikesNumber, &photo.Comments)
	if err != nil {
		return photo, fmt.Errorf("selecting photo: %w", err)
	}
	return photo, nil
}


// DeletePhoto elimina la foto con photos_id=id dalla tabella photos
func (a *appdbimpl) DeletePhoto(userId string, photoID string) error {
	_, err := a.c.Exec(`DELETE FROM photos WHERE user_id = ? AND photos_id = ?`, userId, photoID)
	if err != nil {
		return fmt.Errorf("deleting photo: %w", err)
	}
	return nil
}

// SetComment inserisce un nuovo commento nel database nella tabella comment
func (a *appdbimpl) SetComment(userId string, photoID string, text string) error {
	commentID := generateUniqueCommentID()
	_, err := a.c.Exec(`INSERT INTO comment (user_id, photos_id, comment_id, text_comment) VALUES (?, ?, ?, ?)`, userId, photoID, commentID, text)
	if err != nil {
		return fmt.Errorf("inserting comment: %w", err)
	}
	return nil
}

// GetCommentByID restituisce i dettagli del commento in comment con comment_id=id
func (a *appdbimpl) GetCommentByID(userId string, photoID string, commentID string) (Comment, error) {
	var comment Comment
	err := a.c.QueryRow(`SELECT * FROM comment WHERE user_id = ? AND photos_id = ? AND comment_id = ?`, userId, photoID, commentID).Scan(&comment.UserId, &comment.PhotosId, &comment.CommentId, &comment.Text)
	if err != nil {
		return comment, fmt.Errorf("selecting comment: %w", err)
	}
	return comment, nil
}

// DeleteComment elimina il commento con comment_id=id dalla tabella comment
func (a *appdbimpl) DeleteComment(userId string, photoID string, commentID string) error {
	_, err := a.c.Exec(`DELETE FROM comment WHERE user_id = ? AND photos_id = ? AND comment_id = ?`, userId, photoID, commentID)
	if err != nil {
		return fmt.Errorf("deleting comment: %w", err)
	}
	return nil
}


// GetCommentsByPhotoID restituisce i dettagli dei commenti in comment con photos_id=id
func (a *appdbimpl) GetCommentsByPhotoID(userId string, photoID string) ([]Comment, error) {
	rows, err := a.c.Query(`SELECT * FROM comment WHERE user_id = ? AND photos_id = ?`, userId, photoID)
	if err != nil {
		return nil, fmt.Errorf("selecting comments: %w", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.UserId, &comment.PhotosId, &comment.CommentId, &comment.Text)
		if err != nil {
			return nil, fmt.Errorf("scanning comments: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// GetPhotosByUserID restituisce i dettagli delle foto in photos con user_id=id
func (a *appdbimpl) GetPhotosByUserID(userId string) ([]Photo, error) {
	rows, err := a.c.Query(`SELECT * FROM photos WHERE user_id = ?`, userId)
	if err != nil {
		return nil, fmt.Errorf("selecting photos: %w", err)
	}
	defer rows.Close()

	var photos []Photo
	for rows.Next() {
		var photo Photo
		err := rows.Scan(&photo.UserID, &photo.BinaryFile, &photo.PhotosId, &photo.Url, &photo.Timestamp, &photo.LikesNumber, &photo.Comments)
		if err != nil {
			return nil, fmt.Errorf("scanning photos: %w", err)
		}
		photos = append(photos, photo)
	}
	return photos, nil
}


// SetLike incrementa il numero di like di una foto
func (a *appdbimpl) SetLike(userId string, photoID string) error {
	photo, err := a.GetPhotoByID(userId, photoID)
	if err != nil {
		return fmt.Errorf("getting photo: %w", err)
	}
	photo.LikesNumber++
	_, err = a.c.Exec(`UPDATE photos SET likes_number = ? WHERE user_id = ? AND photos_id = ?`, photo.LikesNumber, userId, photoID)
	if err != nil {
		return fmt.Errorf("updating photo: %w", err)
	}
	return nil
}

// DeleteLike decrementa il numero di like di una foto
func (a *appdbimpl) DeleteLike(userId string, photoID string) error {
	photo, err := a.GetPhotoByID(userId, photoID)
	if err != nil {
		return fmt.Errorf("getting photo: %w", err)
	}
	photo.LikesNumber--
	_, err = a.c.Exec(`UPDATE photos SET likes_number = ? WHERE user_id = ? AND photos_id = ?`, photo.LikesNumber, userId, photoID)
	if err != nil {
		return fmt.Errorf("updating photo: %w", err)
	}
	return nil
}

// GetLikesByPhotoID restituisce il numero di like di una foto
func (a *appdbimpl) GetLikesByPhotoID(userId string, photoID string) (int, error) {
	photo, err := a.GetPhotoByID(userId, photoID)
	if err != nil {
		return 0, fmt.Errorf("getting photo: %w", err)
	}
	return photo.LikesNumber, nil
}


// GetPhotosStreamByUserID restituisce lista foto di tutti account seguiti da userID in ordine cronologico inverso
func (a *appdbimpl) GetPhotosStreamByUserID(userID string) ([]Photo, error) {
	//ottengo lista follower da cui poi prendere le foto da mettere nella lista stream
	followers, err := a.GetFollowersByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("getting followers: %w", err)
	}

	// prendo le foto di tutti gli utenti seguiti e le ordino per timestamp
	var photos []Photo
	for _, follower := range followers {
		followerPhotos, err := a.GetPhotosByUserID(follower.UserId)
		if err != nil {
			return nil, fmt.Errorf("getting photos: %w", err)
		}
		photos = append(photos, followerPhotos...)
	}

	// ordino le foto per timestamp
	for i := 0; i < len(photos); i++ {
		for j := i + 1; j < len(photos); j++ {
			if photos[i].Timestamp < photos[j].Timestamp {
				photos[i], photos[j] = photos[j], photos[i]
			}
		}
	}

	return photos, nil
}
