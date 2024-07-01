package api

import(
	"encoding/json"
	"net/http"
)

// homeHandler gestisce richiesta GET alla route principale "/"
func (rt *_router) homeHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

	//Struttura della risposta
	response := struct {
		Message string `json:"message"`
	}{
		Message: "Benvenuto in WASAPhoto",
	}

	// Concerti la struttura in JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Imposta l'intestazione Content-Type e scrivi la risposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}