package api

import (
	"encoding/json"
	"log"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// homeHandler gestisce richiesta GET alla route principale "/"
func (rt *_router) homeHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

	// Struttura della risposta
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Benvenuto in WASAPhoto",
	}

	// Concerti la struttura in JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Imposta l'intestazione Content-Type e scrivi la risposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responseJSON); err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
