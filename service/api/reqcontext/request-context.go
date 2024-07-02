/*
Package reqcontext contains the request context. Each request will have its own instance of RequestContext filled by the
middleware code in the api-context-wrapper.go (parent package).

Each value here should be assumed valid only per request only, with some exceptions like the logger.
*/
package reqcontext

import (
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"fmt"
	"errors"
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

// ExtractBearerToken estrae il token bearer dall'header Authorization
func ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}

	token := parts[1]
	fmt.Printf("Token estratto: %s\n", token) // Aggiungi log per vedere il token estratto
	return token, nil
}

// AuthenticateUser autentica l'utente utilizzando il token JWT
func AuthenticateUser(token string, db database.AppDatabase) (database.User, error) {
    // Verifica se il token contiene l'ID dell'utente (esempio)
    userIDStr := token

    // Log per verificare l'ID estratto dal token
    fmt.Printf("User ID estratto dal token: %s\n", userIDStr)

    // Ottieni l'utente dal database usando l'ID estratto dal token
    user, err := db.GetUserById(userIDStr)
    if err != nil {
        // Log dell'errore nel recupero dell'utente
        fmt.Printf("Errore nel recupero dell'utente: %v\n", err)
        return database.User{}, fmt.Errorf("user not found: %w", err)
    }

    // Log dell'utente autenticato
    fmt.Printf("User autenticato: %+v\n", user)

    return user, nil
}