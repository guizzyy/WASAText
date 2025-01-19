package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("setGroupPhoto not authorized")
		http.Error(w, "setGroupPhoto operation not allowed", http.StatusUnauthorized)
		return
	}

	var conv utilities.Conversation

	// Get the group id interested
	if conv.ID, err = strconv.ParseUint(params.ByName("convID"), 10, 64); err != nil {
		context.Logger.WithError(err).Error("error in getting convID for setGroupPhoto")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the previous group photo if there was an existing one
	currPhoto, err := rt.db.GetConvPhoto(conv.ID)
	if err != nil {
		context.Logger.WithError(err).Error("error during get current group photo db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if currPhoto != "" {
		if err = rt.DeletePhotoPath(currPhoto); err != nil {
			context.Logger.WithError(err).Error("error during delete current group photo path")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	//Get the photo from the request body and save the file path; get the group id
	if conv.Photo, err = rt.GetPhotoPath(w, r, context); err != nil {
		context.Logger.WithError(err).Error("error during get photo path setGroupPhoto")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the new group photo in the database
	if err = rt.db.SetGroupPhoto(conv, id); err != nil {
		context.Logger.WithError(err).Error("error during setGroupPhoto db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := utilities.Notification{
		Report: "Group photo updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		context.Logger.WithError(err).Error("json set group photo encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
