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
		context.Logger.Error(utilities.ErrUsernameString)
		http.Error(w, utilities.ErrUsernameString.Error(), http.StatusBadRequest)
		return
	}

	// Query the database for receiver info
	if err = rt.db.GetUserByUsername(pReceiver); err != nil {
		context.Logger.WithError(err).Error("error during getUserByUsername db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query the database for the user info
	user, err := rt.db.GetUserByID(id)
	if err != nil {
		context.Logger.WithError(err).Error("error during getUserByID db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the new conversation in the database
	conv, err := rt.db.CreatePrivConv(user, receiver)
	if err != nil {
		context.Logger.WithError(err).Error("error during CreatePrivConv db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if conv.Photo, err = rt.GetFile(conv.Photo); err != nil {
		context.Logger.WithError(err).Error("error during GetFile in GetConversations")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Schedule the deletion of conversation if no messages are sent within 5 minutes
	go rt.ScheduleConvDeleting(conv.ID, context, w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(conv); err != nil {
		context.Logger.WithError(err).Error("json start conv encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
