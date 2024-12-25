package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check the authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get the username wanted from the query and check if it's valid
	username := r.URL.Query().Get("username")
	if check, err := rt.checkStringFormat(username); err != nil {
		http.Error(w, "Error checking the user format", http.StatusBadRequest)
		return
	} else if !check {
		http.Error(w, "username not valid", http.StatusBadRequest)
		return
	}

	// Query the database in order to retrieve name and photo from users
	users, err := rt.db.GetUsers(username, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
