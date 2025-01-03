package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Get the authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}
	var user utilities.User

	// Limit the dimension of the file to 32MB
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}

	// Get the file from the request body
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error getting file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if the file is an image file
	if isImage, err := rt.checkFileFormat(file); err != nil {
		http.Error(w, "Error checking file format", http.StatusBadRequest)
		return
	} else if !isImage {
		http.Error(w, "File is not an image", http.StatusBadRequest)
		return
	}

	// Create a file in the folder and copy the image in it
	filePath := filepath.Join("photos", handler.Filename)
	fileLocal, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer fileLocal.Close()
	_, err = io.Copy(fileLocal, file)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}

	// Insert/Update the photo path in the database
	user.ID = id
	user.Photo = filePath
	if err = rt.db.SetPhoto(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the client a notification for the success of the operation
	response := utilities.Notification{
		Report: "Profile photo updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
