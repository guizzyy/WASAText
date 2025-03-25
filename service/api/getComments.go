package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
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

	// Get mess id interested
	messID, err := strconv.ParseUint(params.ByName("messID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting messID for getComments")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the conversation id from the path
	convID, err := strconv.ParseUint(params.ByName("convID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting convID for getComments")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the user is in the conversation to see the reactions
	if isIn, err := rt.db.IsUserInConv(convID, id); err != nil {
		context.Logger.WithError(err).Error("error in checking if user is in conversation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !isIn {
		context.Logger.Error("user is not in conversation")
		http.Error(w, database.ErrUserNotFound.Error(), http.StatusBadRequest)
		return
	}

	// Check if the message is in the conversation
	if isIn, err := rt.db.IsMessageInConv(messID, convID); err != nil {
		context.Logger.WithError(err).Error("error in checking if message is in conversation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !isIn {
		context.Logger.Error("message is not in conversation")
		http.Error(w, database.ErrMessageNotFound.Error(), http.StatusBadRequest)
		return
	}

	// Query the database to retrieve a list of reactions in the message
	reactions, err := rt.db.GetReactions(messID)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting reactions db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range reactions {
		reactions[i].User.Photo, err = rt.GetFile(reactions[i].User.Photo)
		if err != nil {
			context.Logger.WithError(err).Error("error in GetFile in GetComments")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(reactions); err != nil {
		context.Logger.WithError(err).Error("json get comments encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
