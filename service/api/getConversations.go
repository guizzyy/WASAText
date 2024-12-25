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
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}

	// TODO: modify the function because conversations must be sorted
	convs, err := rt.db.GetConversations(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(convs); err != nil {
		http.Error(w, "error during the response encode", http.StatusInternalServerError)
		return
	}
}
