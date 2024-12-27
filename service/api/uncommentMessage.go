package api

import (
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
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

	// Get the username to delete his reaction from the message
	username, err := rt.db.GetUsernameByID(id)
	if err != nil {
		http.Error(w, "Error getting user by id", http.StatusNotFound)
		return
	}

	// Get the id of the message which we want to delete the reaction from
	messID, err := strconv.ParseUint(params.ByName("messID"), 10, 64)
	if err != nil {
		http.Error(w, "error getting the message id", http.StatusBadRequest)
		return
	}

	// Query the database to delete the reaction with given mess id and username
	if err := rt.db.RemoveReaction(messID, username); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
