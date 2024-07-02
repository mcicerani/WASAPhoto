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

// SetMyUsernameHandler modifica il nome utente
func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	
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

    // Assicurati che il campo 'username' sia presente e non vuoto
    if reqBody.Username == "" {
        http.Error(w, "Bad Request: 'username' field is required", http.StatusBadRequest)
        return
    }

    // Log delle informazioni sull'utente e sull'operazione
    log.Printf("User ID autenticato: %d, Username: %s", user.ID, user.Username)
    log.Printf("UserID da parametro: %s", userID)
    log.Printf("Nuovo username da aggiornare: %s", reqBody.Username)

    // Effettua l'aggiornamento dell'username nel database
    if err := ctx.Database.UpdateUsername(userID, reqBody.Username); err != nil {
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

    // Costruisci il profilo utente con tutte le informazioni
    userProfile := struct {
        User         database.User    `json:"user"`
        Followers    []database.User  `json:"followers"`
        NumFollowers int              `json:"numFollowers"`
        Follows      []database.User  `json:"follows"`
        NumFollowing int              `json:"numFollowing"`
        Photos       []database.Photo `json:"Photos"`
        NumPhotos    int              `json:"numPhotos"`
    }{
        User:         user,
        Followers:    followers,
        NumFollowers: numFollowers,
        Follows:      follows,
        NumFollowing: numFollows,
        Photos:       photos,
        NumPhotos:    numPhotos,
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

// getMyStreamHandle ritorna lo stream dell'utente cliccando su tasto stream
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	Photos, err := ctx.Database.GetPhotosStreamByUserID(ps.ByName("userId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Costruisci lo stream dell'utente
	userStream := struct {
		Photos []database.Photo `json:"Photos"`
	}{
		Photos: Photos,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(userStream)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

// followUserHandler segue un utente
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	isBanned, err := ctx.Database.IsBanned(strconv.Itoa(loggedUser.ID), ps.ByName("followedId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isBanned {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	followedID := ps.ByName("followedId")
	err = ctx.Database.FollowUser(strconv.Itoa(ctx.User.ID), followedID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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
	_, err = reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	isBanned, err := ctx.Database.IsBanned(strconv.Itoa(ctx.User.ID), ps.ByName("followedId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if isBanned {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	followedID := ps.ByName("followedId")
	err = ctx.Database.UnfollowUser(strconv.Itoa(ctx.User.ID), followedID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// getUserFollowsHandler ritorna gli utenti seguiti
func (rt *_router) getUserFollows(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	Follows, err := ctx.Database.GetFollows(ps.ByName("userId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Follows)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// getUserFollowersHandler ritorna gli utenti che seguono
func (rt *_router) getUserFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	Followers, err := ctx.Database.GetFollowers(ps.ByName("userId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Followers)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// banUserHandler banna un utente
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	err = ctx.Database.BanUser(strconv.Itoa(ctx.User.ID), ps.ByName("bannedId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// rimuove dai follow e followers

	err = ctx.Database.UnfollowUser(strconv.Itoa(ctx.User.ID), ps.ByName("bannedId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ctx.Database.UnfollowUser(ps.ByName("bannedId"), strconv.Itoa(ctx.User.ID))
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
	_, err = reqcontext.AuthenticateUser(token, ctx.Database)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = ctx.Database.UnbanUser(strconv.Itoa(ctx.User.ID), ps.ByName("bannedId"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
