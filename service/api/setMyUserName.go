package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	_ "git.guizzyy.it/WASAText/service/database"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
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

	// Get the new username wanted form the request body
	var user utilities.User
	user.ID = id
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the username provided has the correct format
	if check, err := rt.checkStringFormat(user.Username); err != nil {
		http.Error(w, "Error while checking the username", http.StatusInternalServerError)
		return
	} else if !check {
		http.Error(w, "Invalid username proposed", http.StatusBadRequest)
		return
	}

	// Set the new username in the database
	if err = rt.db.SetUsername(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the client a notification for the success of the operation
	response := utilities.Notification{
		Report: "Username updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
