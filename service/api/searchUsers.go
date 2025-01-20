package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check the authorization for the operation
	isAuth, token, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("searchUsers not authorized")
		http.Error(w, "searchUsers operation not allowed", http.StatusUnauthorized)
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
		http.Error(w, "searchUsers operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get the username wanted from the query and check if it's valid
	username := r.URL.Query().Get("username")
	if username == "" {
		context.Logger.Error("write at least one letter to search users")
		http.Error(w, "searchUsers operation empty string", http.StatusBadRequest)
		return
	}

	// Query the database in order to retrieve name and photo from users
	users, err := rt.db.GetUsers(username, loggedID)
	if err != nil {
		context.Logger.WithError(err).Error("error during getUsers db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(users); err != nil {
		context.Logger.WithError(err).Error("json searchUsers encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
