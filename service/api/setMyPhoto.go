package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}

	var newPhoto Photo
	if err := json.NewDecoder(r.Body).Decode(&newPhoto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err = rt.db.SetPhoto(newPhoto.Photo, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := Notification{
		Outcome:   true,
		Report:    "Profile photo updated successfully",
		ErrorCode: 0,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
