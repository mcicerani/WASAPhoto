package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// doLogin handles the login request con metodo bearer authentication
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

	username := r.FormValue("username")

	// Verifica se l'utente esiste nel database
	user, err := ctx.Database.GetUserByUsername(username)
	if err != nil {
		// Se si verifica un errore, restituisci un errore 500
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Se l'utente non esiste, crea un nuovo utente
	if user.ID == 0 {
		// Crea il nuovo utente nel database
		err := ctx.Database.SetUser(username)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Recupera l'utente appena creato
		user, err = ctx.Database.GetUserByUsername(username)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Genera un token con l'identificatore utente
	token := strconv.Itoa(user.ID)

	// Restituisci il token nella risposta
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"token": token})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
