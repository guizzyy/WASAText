package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	_ "git.guizzyy.it/WASAText/service/database"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}

	var newUsername utilities.Username
	if err := json.NewDecoder(r.Body).Decode(&newUsername); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if check, err := rt.checkStringFormat(newUsername.Name); err != nil {
		http.Error(w, "Error while checking the username", http.StatusInternalServerError)
		return
	} else if !check {
		http.Error(w, "Invalid username proposed", http.StatusBadRequest)
		return
	}

	if err = rt.db.SetUsername(newUsername.Name, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := utilities.Notification{
		Report: "Username updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
