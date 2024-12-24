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
	isAuth, _, err := rt.checkToken(r)
	if err != nil {
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}

	var username utilities.Username
	if err := json.NewDecoder(r.Body).Decode(&username); err != nil {
		http.Error(w, "error decoding the body", http.StatusBadRequest)
	}
	gr_id, err := strconv.ParseUint(params.ByName("group_id"), 10, 64)
	if err != nil {
		http.Error(w, "error retrieving the group ID", http.StatusBadRequest)
	}

	if err := rt.db.AddMembership(user_id, gr_id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := utilities.Notification{
		Outcome:   true,
		Report:    "Enter the group successfully",
		ErrorCode: 0,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
