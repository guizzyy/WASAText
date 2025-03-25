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
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, "Error checking the token", http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("createGroup not authorized")
		http.Error(w, "createGroup operation not allowed", http.StatusUnauthorized)
		return
	}

	var group utilities.Conversation
	pGroup := &group

	// Get the group name from the request body
	if err = json.NewDecoder(r.Body).Decode(&group); err != nil {
		context.Logger.WithError(err).Error("json create group decode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the group name format is correct
	if check, err := rt.checkGroupStringFormat(group.Name); err != nil {
		context.Logger.WithError(err).Error("error during string format check")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !check {
		context.Logger.Error(utilities.ErrGroupNameString)
		http.Error(w, utilities.ErrGroupNameString.Error(), http.StatusBadRequest)
		return
	}

	// Query the database to save the new group
	group.Type = "group"
	if err = rt.db.CreateGroupConv(pGroup, id); err != nil {
		context.Logger.WithError(err).Error("error during createGroup db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Schedule the deletion of conversation if no messages are sent within 3 minutes
	go rt.ScheduleConvDeleting(group.ID, context, w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(group); err != nil {
		context.Logger.WithError(err).Error("json create group encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
