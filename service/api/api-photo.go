package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Ottenere l'ID dell'utente dalla richiesta
	userID := ps.ByName("userId")

	// Log per mostrare che la richiesta di upload è stata ricevuta
	log.Printf("Upload request received for user ID: %s\n", userID)

	token, err := reqcontext.ExtractBearerToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Error extracting Bearer token:", err)
		return
	}

	// Autenticare l'utente utilizzando il token
	user, err := reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Println("Error authenticating user:", err)
		return
	}

	// Log per mostrare che l'autenticazione è avvenuta con successo
	log.Printf("User authenticated: %s (ID: %d)\n", user.Username, user.ID)

	// Verifica che il metodo di richiesta sia POST e che il contenuto sia di tipo multipart/form-data
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err = r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error parsing multipart form:", err)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error retrieving file from form data:", err)
		return
	}
	defer file.Close()

	// Leggi i dati del file in un byte slice
	imageData, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error reading file contents:", err)
		return
	}

	log.Printf("Received image data length: %d", len(imageData))

	// Salvataggio dell'immagine nel database e ottenimento dell'ID della foto
	timestamp := time.Now().Format("20060102150405") // Formato timestamp: YYYYMMDDHHmmSS
	photoID, err := ctx.Database.SetPhoto(userID, imageData, timestamp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error saving photo URL to database:", err)
		return
	}

	// Costruisci l'oggetto Photo da restituire come risposta JSON
	photo := database.Photo{
		ID:        int(photoID),
		UserID:    user.ID, // Utilizzo user.ID come ID dell'utente autenticato
		ImageData: imageData,
		Timestamp: timestamp,
	}

	// Creare la risposta JSON contenente i dettagli della foto
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error encoding JSON response:", err)
		return
	}

	// Log per mostrare che la risposta JSON è stata inviata correttamente
	log.Println("JSON response sent successfully")
}

// deletePhotoHandler elimina una foto dal sistema di archiviazione locale e dal database

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Ottenere l'ID dell'utente e l'ID della foto dalla richiesta
	userID := ps.ByName("userId")
	photoID := ps.ByName("photosId")

	// Log per mostrare i parametri ricevuti
	log.Printf("Deleting photo with userId: %s, photosId: %s\n", userID, photoID)

	// Verificare che l'ID dell'utente e l'ID della foto siano validi
	if userID == "" || photoID == "" {
		log.Println("Bad Request: Empty userId or photosId")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	token, err := reqcontext.ExtractBearerToken(r)
	if err != nil {
		log.Println("Unauthorized: Failed to extract token")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Autentica l'utente utilizzando il token
	user, err := reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		log.Println("Unauthorized: Authentication failed")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verificare se l'utente possiede la foto
	photo, err := ctx.Database.GetPhotoByID(photoID)
	if err != nil {
		log.Printf("Internal Server Error: Failed to retrieve photo with photoId: %s\n", photoID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if photo.UserID != user.ID {
		log.Println("Unauthorized: User does not own the photo")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Eliminare la foto dal database
	err = ctx.Database.DeletePhoto(photoID)
	if err != nil {
		log.Printf("Internal Server Error: Failed to delete photo with photoId: %s\n", photoID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Rispondere con lo stato di successo
	w.WriteHeader(http.StatusOK)
	log.Printf("Photo deleted successfully with photoId: %s\n", photoID)
}

// getPhotoHandler ottiene i dettagli di una foto dal database
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Ottenere l'ID dell'utente e l'ID della foto dalla richiesta
	userID := ps.ByName("userId")
	photoID := ps.ByName("photosId")

	// Verificare che l'ID dell'utente e l'ID della foto siano validi
	if userID == "" || photoID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Bad Request: Empty userID or photosID")
		return
	}

	token, err := reqcontext.ExtractBearerToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Unauthorized: %v\n", err)
		return
	}

	// Autentica l'utente utilizzando il token
	_, err = reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Unauthorized: %v\n", err)
		return
	}

	// Ottenere i dettagli della foto dal database
	photo, err := ctx.Database.GetPhotoByID(photoID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Internal Server Error: %v\n", err)
		return
	}

	// Log per indicare il successo nel recupero dei dettagli della foto
	log.Printf("Photo details fetched successfully: %v\n", photo)

	// Creare la risposta JSON contenente i dettagli della foto
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error encoding JSON response: %v\n", err)
		return
	}
}

// likePhotoHandler aggiunge un like a una foto nel database
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := reqcontext.ExtractBearerToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Autentica l'utente utilizzando il token
	user, err := reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ID della foto dalla richiesta
	photoID := ps.ByName("photosId")
	userID := strconv.Itoa(user.ID)

	// Verificare che l'ID dell'utente e l'ID della foto siano validi
	if userID == "" || photoID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Aggiungere un like alla foto nel database
	err = ctx.Database.SetLike(userID, photoID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Rispondere con lo stato di successo
	w.WriteHeader(http.StatusOK)
}

// unlikePhotoHandler rimuove un like da una foto nel database
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := reqcontext.ExtractBearerToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Autentica l'utente utilizzando il token
	user, err := reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	photoID := ps.ByName("photosId")
	likeID := ps.ByName("likesId")
	userID := strconv.Itoa(user.ID)

	// Verificare che l'ID dell'utente, l'ID della foto e l'ID del like siano validi
	if userID == "" || photoID == "" || likeID == "" {
		log.Printf("userID: %s, photoID: %s, likeID: %s", userID, photoID, likeID)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Verificare se l'utente è chi ha messo il like
	like, err := ctx.Database.GetLikeByID(likeID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("likeID: %d, userID: %d, photoID: %d", like.ID, like.UserID, like.PhotoID)

	if like.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Rimuovere il like dalla foto nel database
	err = ctx.Database.DeleteLike(likeID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Rispondere con lo stato di successo
	w.WriteHeader(http.StatusOK)
}

// GetPhotoLike ritorna lista likes a una photo
func (rt *_router) getPhotoLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Ottenere l'ID dell'utente e l'ID della foto dalla richiesta
	userID := ps.ByName("userId")
	photoID := ps.ByName("photosId")

	// Verificare che l'ID dell'utente e l'ID della foto siano validi
	if userID == "" || photoID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	token, err := reqcontext.ExtractBearerToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Autentica l'utente utilizzando il token
	_, err = reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ottenere i likes della foto dal database
	likes, err := ctx.Database.GetLikesByPhotoID(photoID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Creare la risposta JSON contenente i commenti della foto
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(likes)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// commentPhotoHandler aggiunge un commento a una foto nel database
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := reqcontext.ExtractBearerToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Autentica l'utente utilizzando il token
	user, err := reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ID della foto dalla richiesta
	photoID := ps.ByName("photosId")
	userID := strconv.Itoa(user.ID)
	log.Printf("userID: %s, photoId: %s", userID, photoID)

	// Verificare che l'ID dell'utente e l'ID della foto siano validi
	if userID == "" || photoID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Ottenere il testo del commento dalla richiesta
	comment := r.FormValue("comment")
	log.Printf("comment: %s", comment)

	// Aggiungere il commento alla foto nel database
	err = ctx.Database.SetComment(userID, photoID, comment)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Rispondere con lo stato di successo
	w.WriteHeader(http.StatusOK)
}

// uncommentPhotoHandler rimuove un commento da una foto nel database
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	token, err := reqcontext.ExtractBearerToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Autentica l'utente utilizzando il token
	user, err := reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ottenere l'ID dell'utente, l'ID della foto e l'ID del commento dalla richiesta
	userID := strconv.Itoa(user.ID)
	photoID := ps.ByName("photosId")
	commentID := ps.ByName("commentsId")

	// Verificare che l'ID dell'utente, l'ID della foto e l'ID del commento siano validi
	if userID == "" || photoID == "" || commentID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Verificare se l'utente è chi ha scritto il commento
	comment, err := ctx.Database.GetCommentByID(commentID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if comment.UserId != user.ID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Rimuovere il commento dalla foto nel database
	err = ctx.Database.DeleteComment(commentID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Rispondere con lo stato di successo
	w.WriteHeader(http.StatusOK)
}

// getPhotoCommentsHandler ottiene i commenti di una foto dal database
func (rt *_router) getPhotoComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Ottenere l'ID dell'utente e l'ID della foto dalla richiesta
	userID := ps.ByName("userId")
	photoID := ps.ByName("photosId")

	// Verificare che l'ID dell'utente e l'ID della foto siano validi
	if userID == "" || photoID == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	token, err := reqcontext.ExtractBearerToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Autentica l'utente utilizzando il token
	_, err = reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ottenere i commenti della foto dal database
	comments, err := ctx.Database.GetCommentsByPhotoID(photoID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Creare la risposta JSON contenente i commenti della foto
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
