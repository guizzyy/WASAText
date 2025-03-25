package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	_ "git.guizzyy.it/WASAText/service/database"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, token, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("Error during checkToken")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("setMyUserName not authorized")
		http.Error(w, "setMyUserName operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get the user id from the path and check if it matches the auth id
	loggedID, err := strconv.ParseUint(params.ByName("uID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting loggedID from the path")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if loggedID != token {
		context.Logger.WithError(err).Error("Security error")
		http.Error(w, "setMyUserName operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get the new username wanted from the request body
	var user utilities.User
	user.ID = loggedID
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		context.Logger.WithError(err).Error("json setUsername decode error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the username provided has the correct format
	if check, err := rt.checkStringFormat(user.Username); err != nil {
		context.Logger.WithError(err).Error("error during string format check")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !check {
		context.Logger.Error(utilities.ErrUsernameString)
		http.Error(w, utilities.ErrUsernameString.Error(), http.StatusBadRequest)
		return
	}

	// Set the new username in the database
	if err = rt.db.SetUsername(user); err != nil {
		context.Logger.WithError(err).Error("error during setUsername db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the client a notification for the success of the operation
	response := utilities.Notification{
		Report: "Username updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		context.Logger.WithError(err).Error("json setUsername encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
