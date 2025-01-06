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
	isAuth, _, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !isAuth {
		context.Logger.WithError(err).Error("startConversation not authorized")
		http.Error(w, "startConversation operation not allowed", http.StatusUnauthorized)
		return
	}

	// Retrieve and check the format of the username in the request body
	var receiver utilities.User
	if err := json.NewDecoder(r.Body).Decode(&receiver); err != nil {
		context.Logger.WithError(err).Error("json start conv decode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if check, err := rt.checkStringFormat(receiver.Username); err != nil {
		context.Logger.WithError(err).Error("error during string format check")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !check {
		context.Logger.WithError(err).Error(utilities.ErrString)
		http.Error(w, "Username format not allowed", http.StatusBadRequest)
		return
	}

	// TODO: find a way to manage the redundancy of private conversations
}
