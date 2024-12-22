package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	isAuth, _, err := rt.checkToken(r)
	if err != nil {
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}

	var newUsername string
	if err := json.NewDecoder(r.Body).Decode(&newUsername); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	group_id, err := strconv.ParseUint(params.ByName("grID"), 10, 64)
	if err != nil {
		http.Error(w, "Error parsing group id", http.StatusBadRequest)
		return
	}

	if err := rt.db.SetGroupName(newUsername, group_id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	response := Notification{
		Outcome:   true,
		Report:    "Group name updated successfully",
		ErrorCode: 0,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)

}
