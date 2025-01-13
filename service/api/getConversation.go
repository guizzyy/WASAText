package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, "Error checking the token", http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("getConversation not authorized")
		http.Error(w, "getConversation operation not allowed", http.StatusUnauthorized)
		return
	}

	var convID uint64
	// Get the conv id from the path
	if convID, err = strconv.ParseUint(params.ByName("convID"), 10, 64); err != nil {
		context.Logger.WithError(err).Error("error in getting convID for getConversation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query the database to retrieve all messages for the conversation
	// TODO: finish function in the database
	messages, err := rt.db.GetConversation(convID, id)
	if err != nil {
		context.Logger.WithError(err).Error("error during getConversation db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(messages); err != nil {
		context.Logger.WithError(err).Error("json get conversation encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
