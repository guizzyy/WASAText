package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check the authorization for the operation
	isAuth, id, err := rt.checkToken(r)
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

	// Get the username wanted from the query and check if it's valid
	username := r.URL.Query().Get("username")
	if check, err := rt.checkStringFormat(username); err != nil {
		context.Logger.WithError(err).Error("error during string format check")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !check {
		context.Logger.Error(utilities.ErrString)
		http.Error(w, utilities.ErrString.Error(), http.StatusBadRequest)
		return
	}

	// Query the database in order to retrieve name and photo from users
	users, err := rt.db.GetUsers(username, id)
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
