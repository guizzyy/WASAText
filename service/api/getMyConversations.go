package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check the authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.WithError(err).Error("getMyConversations not authorized")
		http.Error(w, "getMyConversations operation not allowed", http.StatusUnauthorized)
		return
	}

	// Query the database to retrieve the conversations of the user
	convs, err := rt.db.GetConversations(id)
	if err != nil {
		context.Logger.WithError(err).Error("error during GetConversations db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the images and last message for each conversation
	for i := range convs {
		convs[i].Photo, err = rt.GetFile(convs[i].Photo)
		if err != nil {
			context.Logger.WithError(err).Error("error during GetFile in GetConversations")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		convs[i].LastMessage, err = rt.db.GetLastMessage(convs[i].ID, id)
		if err != nil {
			context.Logger.WithError(err).Error("error during GetLastMessage in GetConversations")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(convs); err != nil {
		context.Logger.WithError(err).Error("json getMyConversations encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
