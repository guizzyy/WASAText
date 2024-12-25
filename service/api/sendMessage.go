package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get the text of the message from the request body
	var mess utilities.Message
	pMess := &mess
	if err := json.NewDecoder(r.Body).Decode(&mess); err != nil {
		http.Error(w, "Error decoding body", http.StatusBadRequest)
		return
	}
	mess.Sender = id
	mess.Conv, err = strconv.ParseUint(params.ByName("convID"), 10, 64)
	if err != nil {
		http.Error(w, "Error decoding convID", http.StatusBadRequest)
		return
	}

	if err := rt.db.AddMessage(pMess); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(mess); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
