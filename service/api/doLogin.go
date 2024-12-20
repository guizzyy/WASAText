package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	var userLog string

	if err := json.NewDecoder(r.Body).Decode(&userLog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if check, err := rt.checkStringFormat(userLog); err != nil || !check {
		http.Error(w, "Failed username format validation", http.StatusBadRequest)
	}

	dUser, err := rt.db.LogUser(userLog)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dUser)
}
