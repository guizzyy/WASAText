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
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("forwardMessage not authorized")
		http.Error(w, "forwardMessage operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get the conversation where forward the message to
	var conv utilities.Conversation
	if err = json.NewDecoder(r.Body).Decode(&conv); err != nil {
		context.Logger.WithError(err).Error("json forward message decode error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the message to forward from its id in the path
	messID, err := strconv.ParseUint(params.ByName("messID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting messID for forwardMessage")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msg, err := rt.db.GetMessageInfo(messID)
	if err != nil {
		context.Logger.WithError(err).Error("error during getMessageInfo db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add the forwarded message as a new message in the database
	msg.Sender = id
	msg.Conv = conv.ID
	msg.IsForward = true
	pMess := &msg
	if err = rt.db.AddMessage(pMess); err != nil {
		context.Logger.WithError(err).Error("error during AddMessage db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(pMess); err != nil {
		context.Logger.WithError(err).Error("json forward message encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
