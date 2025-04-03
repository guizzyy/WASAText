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

	if _, err := os.Stat("./tmp/uploads/" + strconv.FormatUint(loggedID, 10)); os.IsNotExist(err) {
		err = os.MkdirAll("./tmp/uploads/"+strconv.FormatUint(loggedID, 10), 0755)
		if err != nil {
			context.Logger.WithError(err).Error("can't create the folder")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		context.Logger.WithError(err).Error("can't get the folder")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the photo from the request body and save it in the correct folder
	mainDir := strconv.FormatUint(loggedID, 10)
	fileName, file, err := rt.GetFilePath(w, r, context)
	if err != nil {
		context.Logger.WithError(err).Error("Error during get photo path")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filePath := fmt.Sprintf("./tmp/uploads/%s/%s", mainDir, fileName)
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

	// Delete the previous user photo if there was an existing one
	user, err := rt.db.GetUserByID(loggedID)
	if err != nil {
		context.Logger.WithError(err).Error("Error during get current user photo db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user.Photo != "" {
		if err = rt.DeleteUserPhoto(user.Photo); err != nil {
			context.Logger.WithError(err).Error("Error during delete current user photo path")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Insert/Update the photo path in the database
	user.ID = loggedID
	user.Photo = filePath
	if err = rt.db.SetPhoto(user); err != nil {
		context.Logger.WithError(err).Error("Error during setPhoto db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	apiPhoto, err := rt.GetFile(user.Photo)
	if err != nil {
		context.Logger.WithError(err).Error("Error during GetFile for user photo")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the client a notification for the success of the operation
	response := utilities.PhotoResponse{
		Report: "Profile photo updated successfully",
		Photo:  apiPhoto,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		context.Logger.WithError(err).Error("json setPhoto encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
