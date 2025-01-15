package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, _, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("setGroupName not authorized")
		http.Error(w, "setGroupName operation not allowed", http.StatusUnauthorized)
		return
	}

	var conv utilities.Conversation

	// Get conv id and the new name to update from request body and path params
	if conv.ID, err = strconv.ParseUint(params.ByName("convID"), 10, 64); err != nil {
		context.Logger.WithError(err).Error("error in getting convID for setGroupName")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewDecoder(r.Body).Decode(&conv); err != nil {
		context.Logger.WithError(err).Error("json set group name decode error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the new name has the correct format
	if check, err := rt.checkStringFormat(conv.Name); err != nil {
		context.Logger.WithError(err).Error("error during string format check")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !check {
		context.Logger.Error(utilities.ErrNameString)
		http.Error(w, utilities.ErrNameString.Error(), http.StatusBadRequest)
		return
	}

	// Set the new group name in the database
	if err = rt.db.SetGroupName(conv); err != nil {
		context.Logger.WithError(err).Error("error during set group name db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the client a notification for the success of the operation
	response := utilities.Notification{
		Report: "Group name update successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		context.Logger.WithError(err).Error("json set group name encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
