package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) startConversation(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check the authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("startConversation not authorized")
		http.Error(w, "startConversation operation not allowed", http.StatusUnauthorized)
		return
	}

	// Retrieve and check the format of the username in the request body
	var receiver utilities.User
	pReceiver := &receiver
	if err = json.NewDecoder(r.Body).Decode(&receiver); err != nil {
		context.Logger.WithError(err).Error("json start conv decode error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if check, err := rt.checkStringFormat(receiver.Username); err != nil {
		context.Logger.WithError(err).Error("error during string format check")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !check {
		context.Logger.Error(utilities.ErrNameString)
		http.Error(w, utilities.ErrNameString.Error(), http.StatusBadRequest)
		return
	}

	// Query the database for receiver info
	if err = rt.db.GetUserByUsername(pReceiver); err != nil {
		context.Logger.WithError(err).Error("error during getUserByUsername db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the new conversation in the database
	conv, err := rt.db.CreatePrivConv(id, receiver)
	if err != nil {
		context.Logger.WithError(err).Error("error during CreatePrivConv db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(conv); err != nil {
		context.Logger.WithError(err).Error("json start conv encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
