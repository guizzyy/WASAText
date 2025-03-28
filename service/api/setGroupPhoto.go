package api

import (
	"encoding/json"
	"fmt"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
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

	var group utilities.Conversation

	// Get the group id interested
	if group.ID, err = strconv.ParseUint(params.ByName("convID"), 10, 64); err != nil {
		context.Logger.WithError(err).Error("error in getting convID for setGroupPhoto")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete the previous group photo if there was an existing one
	currPhoto, err := rt.db.GetGroupPhoto(group.ID)
	if err != nil {
		context.Logger.WithError(err).Error("error during get current group photo db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if currPhoto != "" {
		if err = rt.DeleteGroupPhoto(currPhoto); err != nil {
			context.Logger.WithError(err).Error("error during delete current group photo path")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Get the photo from the request body and save the file path; get the group id
	gPhoto, file, err := rt.GetFilePath(w, r, context)
	if err != nil {
		context.Logger.WithError(err).Error("error during GetFilePath for group photo")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filePath := fmt.Sprintf("./uploads/%s/%s", "groups", gPhoto)
	dst, err := os.Create(filePath)
	if err != nil {
		context.Logger.WithError(err).Error("Error during create file path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err = io.Copy(dst, file); err != nil {
		context.Logger.WithError(err).Error("Error during copy file to path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	group.Photo = filePath

	// Set the new group photo in the database
	if err = rt.db.SetGroupPhoto(group, id); err != nil {
		context.Logger.WithError(err).Error("error during setGroupPhoto db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiPhoto, err := rt.GetFile(group.Photo)
	if err != nil {
		context.Logger.WithError(err).Error("error during GetFile for group photo")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := utilities.PhotoResponse{
		Report: "Group photo updated successfully",
		Photo:  apiPhoto,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		context.Logger.WithError(err).Error("json set group photo encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
