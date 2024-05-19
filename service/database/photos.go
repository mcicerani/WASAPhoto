package database

import (
	"github.com/google/uuid"
)

// genera id unico per foto e commenti
func generateUniquePhotoID() string {
	return uuid.New().String()
}

func generateUniqueCommentID() string {
	return uuid.New().String()
}
