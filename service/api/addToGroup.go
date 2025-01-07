package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, _, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("addToGroup operation not authorized")
		http.Error(w, "addToGroup operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get info about user we want to add
	var userAdded utilities.User
	if err := json.NewDecoder(r.Body).Decode(&userAdded); err != nil {
		context.Logger.WithError(err).Error("json add to group decode error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if userAdded.ID, err = rt.db.GetIDByUsername(userAdded.Username); err != nil {
		context.Logger.WithError(err).Error("error during getIDByUsername db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get conv id of the group conversation where we want to add the user
	convID, err := strconv.ParseUint(params.ByName("convID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting convID for addToGroup")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query the database to insert a new membership
	if err = rt.db.AddToGroup(convID, userAdded); err != nil {
		context.Logger.WithError(err).Error("error during addToGroup db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send a notification to the client for the success of the operation
	response := utilities.Notification{
		Report: "User added to group successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		context.Logger.WithError(err).Error("json add to group encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
