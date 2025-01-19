package api

import (
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, "Error checking the token", http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("deleteMessage not authorized")
		http.Error(w, "deleteMessage operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get the mess ID to delete from the parameters in the path
	mess_id, err := strconv.ParseUint(params.ByName("messID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting messID for deleteMessage")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the conv ID in where the message to delete is
	conv_id, err := strconv.ParseUint(params.ByName("convID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting convID for deleteMessage")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the message is in the conversation and if the conversation exists
	if isIn, err := rt.db.IsMessageInConv(mess_id, conv_id); err != nil {
		context.Logger.WithError(err).Error("error in checking if the message is in conversation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !isIn {
		context.Logger.Error("the message is not in conversation")
		http.Error(w, database.ErrMessageNotFound.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user is in the conversation
	if check, err := rt.db.IsUserInConv(conv_id, id); err != nil {
		context.Logger.WithError(err).Error("error in checking if the user is in conversation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !check {
		context.Logger.Error("the user is not in conversation")
		http.Error(w, database.ErrUserNotFound.Error(), http.StatusBadRequest)
		return
	}

	// Query the database to delete the message
	if err = rt.db.RemoveMessage(mess_id, id); err != nil {
		context.Logger.WithError(err).Error("error in deleting message db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
