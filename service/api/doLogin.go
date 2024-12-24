package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	var userLog utilities.User

	if err := json.NewDecoder(r.Body).Decode(&userLog); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pUser := &userLog

	if check, err := rt.checkStringFormat(userLog.Username); err != nil || !check {
		http.Error(w, "Failed username format validation", http.StatusBadRequest)
		return
	}

	isNew, err := rt.db.LogUser(pUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := utilities.LoginResponse{
		UserLogged: userLog,
	}

	w.Header().Set("Content-Type", "application/json")
	if !isNew {
		w.WriteHeader(http.StatusOK)
		response.Message = "Login successful"
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed encoding the response", http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusCreated)
		response.Message = "User created successfully"
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed encoding the response", http.StatusInternalServerError)
			return
		}
	}
}
