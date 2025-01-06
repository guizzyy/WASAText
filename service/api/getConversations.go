package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getConversations(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check the authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.WithError(err).Error("getConversations not authorized")
		http.Error(w, "getConversations operation not allowed", http.StatusUnauthorized)
		return
	}

	// TODO: modify the function because conversations must be sorted
	convs, err := rt.db.GetConversations(id)
	if err != nil {
		context.Logger.WithError(err).Error("error during GetConversations db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(convs); err != nil {
		context.Logger.WithError(err).Error("json getConversations encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
