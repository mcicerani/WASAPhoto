package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {
	var requestBody struct {
		Username string `json:"username"`
	}

	// Decodifica il corpo JSON della richiesta
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username := strings.ToLower(requestBody.Username)
	log.Printf("Login attempt with username: %s", username)

	// Verifica se l'utente esiste nel database
	log.Printf("Context database: %+v", ctx.Database)
	user, err := ctx.Database.GetUserByUsername(username)
	if err != nil {
		log.Printf("Error retrieving user from database: %v", err)

		// Se l'errore è dovuto a "no rows in result set", crea un nuovo utente
		if strings.Contains(err.Error(), "no rows in result set") {
			log.Printf("User '%s' does not exist, creating new user", username)

			// Crea il nuovo utente nel database
			err := ctx.Database.SetUser(username)
			if err != nil {
				log.Printf("Error creating new user in database: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Recupera l'utente appena creato
			user, err = ctx.Database.GetUserByUsername(username)
			if err != nil {
				log.Printf("Error retrieving newly created user from database: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else {
			// Altrimenti, gestisci l'errore interno del server
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Genera un token con l'identificatore utente
	token := "Bearer " + strconv.Itoa(user.ID)

	// Restituisci il token nella risposta
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"token": token})
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Login successful for user '%s'", username)
}
