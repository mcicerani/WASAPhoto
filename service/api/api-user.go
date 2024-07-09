package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// SearchUserHandler ritorna l'utente cercato, se l'utente che effettua la ricerca è bannato da quell'utente, non verrà visualizzato
func (rt *_router) searchUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

	username := r.FormValue("username")
	var loggedUser = ctx.User

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

	userByUsername, err := ctx.Database.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	usernameId := strconv.Itoa(userByUsername.ID)

	isBanned, err := ctx.Database.IsBanned(strconv.Itoa(loggedUser.ID), usernameId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isBanned {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user, err := ctx.Database.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// SetMyUserNameHandler modifica il nome utente
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	type UsernameUpdateRequest struct {
		Username string `json:"username"`
	}

	// Estrae il token dall'header Authorization
	token, err := reqcontext.ExtractBearerToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Autentica l'utente utilizzando il token JWT
	user, err := reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Verifica che l'utente autenticato sia autorizzato a modificare l'username
	userID := ps.ByName("userId")
	if strconv.Itoa(user.ID) != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Decodifica il corpo JSON della richiesta
	var reqBody UsernameUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Effettua l'aggiornamento dell'username nel database
	err = ctx.Database.UpdateUsername(userID, reqBody.Username)
	if err != nil {
		if err.Error() == "username already exists" {
			http.Error(w, "Username già esistente", http.StatusConflict)
			return
		}
		log.Printf("Errore durante l'aggiornamento dell'username: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Risponde con successo
	w.WriteHeader(http.StatusOK)
}

// GetUserProfileHandler ritorna il profilo utente
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")
	log.Printf("Getting profile for user ID: %s", userId)

	var loggedUser = ctx.User

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

	isBanned, err := ctx.Database.IsBanned(strconv.Itoa(loggedUser.ID), userId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isBanned {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user, err := ctx.Database.GetUserById(userId)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	followers, err := ctx.Database.GetFollowers(userId)
	if err != nil {
		log.Printf("Error retrieving followers: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	numFollowers, err := ctx.Database.CountFollowersByUserID(userId)
	if err != nil {
		log.Printf("Error counting followers: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	follows, err := ctx.Database.GetFollows(userId)
	if err != nil {
		log.Printf("Error retrieving follows: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	numFollows, err := ctx.Database.CountFollowsByUserID(userId)
	if err != nil {
		log.Printf("Error counting follows: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	photos, err := ctx.Database.GetPhotosByUserID(userId)
	if err != nil {
		log.Printf("Error retrieving photos: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	numPhotos, err := ctx.Database.CountPhotosByUserID(userId)
	if err != nil {
		log.Printf("Error counting photos: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	bans, err := ctx.Database.GetBans(userId)
	if err != nil {
		log.Printf("Error retrieving bans: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Costruisci il profilo utente con tutte le informazioni
	userProfile := struct {
		User         database.User    `json:"user"`
		Followers    []database.User  `json:"followers"`
		NumFollowers int              `json:"numFollowers"`
		Follows      []database.User  `json:"follows"`
		NumFollowing int              `json:"numFollowing"`
		Photos       []database.Photo `json:"Photos"`
		NumPhotos    int              `json:"numPhotos"`
		Bans         []database.User  `json:"bans"`
	}{
		User:         user,
		Followers:    followers,
		NumFollowers: numFollowers,
		Follows:      follows,
		NumFollowing: numFollows,
		Photos:       photos,
		NumPhotos:    numPhotos,
		Bans:         bans,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userProfile)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Profile response: %+v", userProfile)
}

// getMyStream ritorna lo stream dell'utente cliccando su tasto stream
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai il token dall'header Authorization
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

	photos, err := ctx.Database.GetPhotosStreamByUserID(strconv.Itoa(user.ID))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Costruisci una struttura temporanea con le informazioni di likes e comments
	var userStream struct {
		Photos []struct {
			database.Photo
			Likes    int `json:"likes"`
			Comments int `json:"comments"`
		} `json:"Photos"`
	}

	// Itera su ogni foto per aggiungere le informazioni di likes e comments
	for _, photo := range photos {
		photoID := strconv.Itoa(photo.ID)

		likes, err := ctx.Database.CountLikesByPhotoID(photoID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		comments, err := ctx.Database.CountCommentsByPhotoID(photoID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Aggiungi la foto con le informazioni di likes e comments alla struttura temporanea
		userStream.Photos = append(userStream.Photos, struct {
			database.Photo
			Likes    int `json:"likes"`
			Comments int `json:"comments"`
		}{
			Photo:    photo,
			Likes:    likes,
			Comments: comments,
		})
	}

	// Serializza la risposta in JSON e inviala al client
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userStream)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// followUserHandler segue un utente
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	isBanned, err := ctx.Database.IsBanned(strconv.Itoa(user.ID), ps.ByName("followedId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isBanned {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	followedID := ps.ByName("followedId")
	log.Printf("Before calling FollowUser: userID = %d, followedID = %s", user.ID, followedID)
	err = ctx.Database.FollowUser(strconv.Itoa(user.ID), followedID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Log per vedere su quale utente viene eseguito il follow con successo
	log.Printf("User %d followed user %s successfully", user.ID, followedID)
	w.WriteHeader(http.StatusOK)
}

// unfollowUserHandler smette di seguire un utente
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	isBanned, err := ctx.Database.IsBanned(strconv.Itoa(user.ID), ps.ByName("followedId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isBanned {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	followedID := ps.ByName("followedId")
	err = ctx.Database.UnfollowUser(strconv.Itoa(user.ID), followedID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// banUserHandler banna un utente
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	err = ctx.Database.BanUser(strconv.Itoa(user.ID), ps.ByName("bannedId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// rimuove dai follow e followers

	err = ctx.Database.UnfollowUser(strconv.Itoa(user.ID), ps.ByName("bannedId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ctx.Database.UnfollowUser(ps.ByName("bannedId"), strconv.Itoa(user.ID))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// unbanUserHandler rimuove il ban a un utente
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	err = ctx.Database.UnbanUser(strconv.Itoa(user.ID), ps.ByName("bannedId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// getIsBanned verifica se l'utente è bannato da un altro utente specifico
func (rt *_router) getIsBanned(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrai il token JWT dall'header Authorization
	token, err := reqcontext.ExtractBearerToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Autentica l'utente utilizzando il token JWT
	_, err = reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Estrai i parametri dall'URL
	userID := ps.ByName("userId")
	bannedID := ps.ByName("bannedId")

	// Verifica se l'utente è bannato
	isBanned, err := ctx.Database.IsBanned(userID, bannedID)
	if err != nil {
		// Se l'errore è diverso da "nessuna riga nel set di risultati", restituisci Internal Server Error
		if err.Error() != "sql: no rows in result set" {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		// Se non c'è ban, restituisci false
		isBanned = false
	}

	// Costruisci la risposta JSON
	response := struct {
		IsBanned bool `json:"isBanned"`
	}{
		IsBanned: isBanned,
	}

	// Serializza la risposta JSON e scrivi nella risposta HTTP
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Imposta l'intestazione Content-Type e scrivi la risposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonBytes); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

// getIsFollwed verifica se l'utente segue un altro utente
func (rt *_router) getIsFollowed(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Estrai i parametri dall'URL
	userID := ps.ByName("userId")
	followedID := ps.ByName("followedId")

	// Verifica se l'utente è seguito dall'utente specificato
	isFollowed, err := ctx.Database.IsFollowed(userID, followedID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Costruisci la risposta JSON
	response := struct {
		IsFollowed bool `json:"isFollowed"`
	}{
		IsFollowed: isFollowed,
	}

	// Serializza la risposta JSON e scrivi nella risposta HTTP
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Imposta l'intestazione Content-Type e scrivi la risposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonBytes); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
