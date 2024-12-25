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
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}

	// Retrieve and check the format of the username in the request body
	var receiver utilities.User
	if err := json.NewDecoder(r.Body).Decode(&receiver); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if check, err := rt.checkStringFormat(receiver.Username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if !check {
		http.Error(w, "Username format not allowed", http.StatusUnauthorized)
		return
	}

	// TODO: find a way to manage the redundancy of private conversations
	err :=
}
