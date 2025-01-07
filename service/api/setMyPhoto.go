package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Get the authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("Error during checking token")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.WithError(err).Error("setMyPhoto not unauthorized")
		http.Error(w, "setMyPhoto operation not allowed", http.StatusUnauthorized)
		return
	}
	var user utilities.User

	// Get the photo from the request body and save the file path
	filePath, err := rt.GetPhotoPath(w, r, context)
	if err != nil {
		context.Logger.WithError(err).Error("Error during get photo path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert/Update the photo path in the database
	user.ID = id
	user.Photo = filePath
	if err = rt.db.SetPhoto(user); err != nil {
		context.Logger.WithError(err).Error("Error during setPhoto db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the client a notification for the success of the operation
	response := utilities.Notification{
		Report: "Profile photo updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		context.Logger.WithError(err).Error("json setPhoto encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
