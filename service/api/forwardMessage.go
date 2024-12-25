package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
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

	// Get the conversation where forward the message to
	var conv utilities.Conversation
	if err = json.NewDecoder(r.Body).Decode(&conv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the message to forward thanks to it's id
	messID, err := strconv.ParseUint(params.ByName("messID"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg, err := rt.db.GetMessageInfo(messID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add the forwarded message as a new message in the database
	msg.Sender = id
	msg.Conv = conv.ID
	pMess := &msg
	if err := rt.db.AddMessage(pMess); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pMess); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
