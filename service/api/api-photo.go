package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

	// Creare la directory "photos" se non esiste
	photosDir := "./photos"
	if _, err := os.Stat(photosDir); os.IsNotExist(err) {
		err := os.Mkdir(photosDir, 0755)
		if err != nil {
			log.Printf("Error creating photos directory")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Printf("Photos directory created")
	} else {
		log.Printf("Photos directory already exists")
	}

	// Ottenere il file dall'input del modulo
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println("Error retrieving file from form data:", err)
		return
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Println("Error closing file:", err)
			return // Ignorare l'errore
		}
	}(file)

	// Generare un nome di file univoco basato sull'ID dell'utente e sul timestamp
	timestamp := time.Now().Format("20060102150405") // Formato timestamp: YYYYMMDDHHmmSS
	filename := fmt.Sprintf("%s_%s.jpg", userID, timestamp)

	// Salvataggio della foto nel sistema di archiviazione locale
	photoFilePath := fmt.Sprintf("./photos/%s", filename)
	photoFile, err := os.Create(photoFilePath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error creating photo file:", err)
		return
	}
	defer func(photoFile *os.File) {
		err := photoFile.Close()
		if err != nil {
			log.Println("Error closing photo file:", err)
			return // Ignorare l'errore
		}
	}(photoFile)

	// Copiare il contenuto del file caricato nel file di archiviazione locale
	_, err = io.Copy(photoFile, file)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error copying file contents:", err)
		return
	}

	// Log per mostrare che il file è stato caricato con successo
	log.Printf("File uploaded successfully: %s\n", photoFilePath)

	// Costruire l'URL della foto utilizzando l'ID dell'utente e il timestamp
	photoURL := fmt.Sprintf("./photos/%s", filename)

	// Inserire l'URL della foto nel database e ottenere l'ID della foto inserita
	photoID, err := ctx.Database.SetPhoto(userID, photoURL, timestamp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error saving photo URL to database:", err)
		return
	}

	// Costruire l'oggetto Photo da restituire come risposta JSON
	photo := database.Photo{
		ID:        int(photoID), // Converto int64 a int
		UserID:    user.ID,      // Utilizzo user.ID come ID dell'utente autenticato
		URL:       photoURL,
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
	photoID := ps.ByName("photoId")

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
	user, err := reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verificare se l'utente possiede la foto
	photo, err := ctx.Database.GetPhotoByID(photoID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if photo.UserID != user.ID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ricavare il nome del file dall'URL della foto
	fileName := filepath.Base(photo.URL)

	// Eliminare la foto dal sistema di archiviazione locale
	err = os.Remove(fmt.Sprintf("./photos/%s", fileName))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Eliminare la foto dal database
	err = ctx.Database.DeletePhoto(photoID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Rispondere con lo stato di successo
	w.WriteHeader(http.StatusOK)
}

// getPhotoHandler ottiene i dettagli di una foto dal database

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Ottenere l'ID dell'utente e l'ID della foto dalla richiesta
	userID := ps.ByName("userId")
	photoID := ps.ByName("photoId")

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

	// Ottenere i dettagli della foto dal database
	photo, err := ctx.Database.GetPhotoByID(photoID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Creare la risposta JSON contenente i dettagli della foto
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// likePhotoHandler aggiunge un like a una foto nel database
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Ottenere l'ID dell'utente e l'ID della foto dalla richiesta
	userID := ps.ByName("userId")
	photoID := ps.ByName("photoId")

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

	// Ottenere l'ID dell'utente, l'ID della foto e l'ID del like dalla richiesta
	userID := ps.ByName("userId")
	photoID := ps.ByName("photoId")
	likeID := ps.ByName("likeId")

	// Verificare che l'ID dell'utente, l'ID della foto e l'ID del like siano validi
	if userID == "" || photoID == "" || likeID == "" {
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

	// Verificare se l'utente è chi ha messo il like
	like, err := ctx.Database.GetLikeByID(likeID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if like.UserID != ctx.User.ID {
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

// commentPhotoHandler aggiunge un commento a una foto nel database
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Ottenere l'ID dell'utente e l'ID della foto dalla richiesta
	userID := ps.ByName("userId")
	photoID := ps.ByName("photoId")

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

	// Ottenere il testo del commento dalla richiesta
	comment := r.FormValue("comment")

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

	// Ottenere l'ID dell'utente, l'ID della foto e l'ID del commento dalla richiesta
	userID := ps.ByName("userId")
	photoID := ps.ByName("photoId")
	commentID := ps.ByName("commentId")

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

	if comment.UserId != ctx.User.ID {
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
	photoID := ps.ByName("photoId")

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
