package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Get the authorization for the operation
	isAuth, token, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("Error during checking token")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("setMyPhoto not unauthorized")
		http.Error(w, "setMyPhoto operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get the user token from the path and check if it matches the auth token
	loggedID, err := strconv.ParseUint(params.ByName("uID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting loggedID from the path")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if loggedID != token {
		context.Logger.WithError(err).Error("Security error")
		http.Error(w, "setMyPhoto operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get the photo from the request body and save the file path
	filePath, err := rt.GetPhotoPath(w, r, context)
	if err != nil {
		context.Logger.WithError(err).Error("Error during get photo path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the previous user photo if there was an existing one
	user, err := rt.db.GetUserByID(loggedID)
	if err != nil {
		context.Logger.WithError(err).Error("Error during get current user photo db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user.Photo != "" {
		if err = rt.DeletePhotoPath(user.Photo); err != nil {
			context.Logger.WithError(err).Error("Error during delete current user photo path")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Insert/Update the photo path in the database
	user.ID = loggedID
	user.Photo = "http://localhost:3000/" + filePath
	if err = rt.db.SetPhoto(user); err != nil {
		context.Logger.WithError(err).Error("Error during setPhoto db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the client a notification for the success of the operation
	response := utilities.PhotoResponse{
		Message: "Profile photo updated successfully",
		Photo:   filePath,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		context.Logger.WithError(err).Error("json setPhoto encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
