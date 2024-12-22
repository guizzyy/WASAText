package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}

	var text string
	if err := json.NewDecoder(r.Body).Decode(&text); err != nil {
		http.Error(w, "Error decoding body", http.StatusBadRequest)
		return
	}
	convParam := params.ByName("convID")
	conv_id, err := strconv.ParseUint(convParam, 10, 64)
	if err != nil {
		http.Error(w, "Error decoding convID", http.StatusBadRequest)
		return
	}

	if err := rt.db.AddMessage(text, conv_id, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mess := Message{
		Text:   text,
		Sender: id,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mess); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
