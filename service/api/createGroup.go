package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
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

	var group utilities.Conversation
	pGroup := &group

	// Get the group name from the request body
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the group name format is correct
	if check, err := rt.checkStringFormat(group.Name); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if !check {
		http.Error(w, "Group name is invalid", http.StatusBadRequest)
		return
	}

	// Query the database to save the new group
	group.Type = "group"
	if err := rt.db.CreateGroupConv(pGroup, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(group); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
