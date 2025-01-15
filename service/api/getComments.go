package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check the authorization for the operation
	isAuth, _, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("startConversation not authorized")
		http.Error(w, "startConversation operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get mess id interested
	messID, err := strconv.ParseUint(params.ByName("messID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting messID for getComments")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Query the database to retrieve a list of reactions in the message
	reactions, err := rt.db.GetReactions(messID)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting reactions db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(reactions); err != nil {
		context.Logger.WithError(err).Error("json get comments encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
