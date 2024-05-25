package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

// uploadPhotoHandler carica una foto in locale e salva url nel database

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	
	userId := ps.ByName("userId")
	if userId == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if ctx.User.ID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if strconv.Itoa(ctx.User.ID) != userId {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save file in local storage
	photoUrl, err := ctx.Storage.SavePhoto(file)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	photo := database.Photo{
		UserId: ctx.User.ID,
		Url:    photoUrl,
	}
	




