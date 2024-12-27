package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get the username to insert in the reaction table
	username, err := rt.db.GetUsernameByID(id)
	if err != nil {
		http.Error(w, "Error getting user by id", http.StatusNotFound)
		return
	}

	// Get the reaction react from the request body
	var react utilities.Reaction
	react.User = username
	if err := json.NewDecoder(r.Body).Decode(&react); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the message id from the parameters in the path
	idMess, err := strconv.ParseUint(params.ByName("messID"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add the reaction info in the database
	if err := rt.db.AddReaction(react, idMess); err != nil {
		http.Error(w, "Error adding reaction", http.StatusInternalServerError)
		return
	}

	// Return the content to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(react); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
