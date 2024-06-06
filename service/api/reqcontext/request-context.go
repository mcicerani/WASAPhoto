/*
Package reqcontext contains the request context. Each request will have its own instance of RequestContext filled by the
middleware code in the api-context-wrapper.go (parent package).

Each value here should be assumed valid only per request only, with some exceptions like the logger.
*/
package reqcontext

import (
	"errors"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// RequestContext is the context of the request, for request-dependent parameters
type RequestContext struct {
	// ReqUUID is the request unique ID
	ReqUUID uuid.UUID

	// Database is the instance of database.AppDatabase where data is saved
	Database database.AppDatabase

	// Logger is a custom field logger for the request
	Logger logrus.FieldLogger

	// User is the user logged in that made the request
	User database.User

	// BearerToken is the bearer token used to authenticate the user
	Token string
}

// ExtractBearerToken extracts the bearer token from the Authorization header
func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}

// AuthenticateUser authenticates the user using the bearer token
func AuthenticateUser(token string, db database.AppDatabase) (database.User, error) {

	// Check if the token is a valid integer
	user, err := db.GetUserById(token)
	if err != nil {
		return database.User{}, errors.New("user not found")
	}

	return user, nil
}
