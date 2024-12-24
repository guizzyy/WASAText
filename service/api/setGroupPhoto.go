package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	isAuth, _, err := rt.checkToken(r)
	if err != nil {
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}

	var photo utilities.Photo
	if err = json.NewDecoder(r.Body).Decode(&photo); err != nil {
		http.Error(w, "error decoding the photo", http.StatusBadRequest)
		return
	}
	gr_id, err := strconv.ParseUint(params.ByName("grID"), 10, 64)
	if err != nil {
		http.Error(w, "error decoding the id photo", http.StatusBadRequest)
		return
	}

	if err := rt.db.SetGroupPhoto(gr_id, photo.Photo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := utilities.Notification{
		Outcome:   true,
		Report:    "Group photo updated successfully",
		ErrorCode: 0,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
