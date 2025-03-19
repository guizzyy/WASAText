package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/database"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, "Error checking the token", http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("commentMessage not authorized")
		http.Error(w, "commentMessage operation not allowed", http.StatusUnauthorized)
		return
	}

	var react utilities.Reaction

	// Get the reaction react from the request body
	react.User, err = rt.db.GetUserByID(id)
	if err != nil {
		context.Logger.WithError(err).Error("error during GetUserByID db")
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}
	if err = json.NewDecoder(r.Body).Decode(&react); err != nil {
		context.Logger.WithError(err).Error("json comment message decode error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the conversation id from the parameters in the path
	idConv, err := strconv.ParseUint(params.ByName("convID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting convID for comment message")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the message id from the parameters in the path
	idMess, err := strconv.ParseUint(params.ByName("messID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting messID for commentMessage")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the user is in the conversation in order to comment a message
	if isIn, err := rt.db.IsUserInConv(idConv, id); err != nil {
		context.Logger.WithError(err).Error("error in checking if user is in convID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !isIn {
		context.Logger.Error("user not in convID")
		http.Error(w, "user not in convID", http.StatusUnauthorized)
		return
	}

	// Check if the message is in the conversation
	if isIn, err := rt.db.IsMessageInConv(idMess, idConv); err != nil {
		context.Logger.WithError(err).Error("error in checking message in convID")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !isIn {
		context.Logger.Error("message not found in conversation")
		http.Error(w, database.ErrMessageNotFound.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the comment is an emoji
	if isEmoji, err := rt.checkEmojiFormat(react.Emoji); err != nil {
		context.Logger.WithError(err).Error("error in checking emoji format")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !isEmoji {
		context.Logger.Error("emoji format not supported")
		http.Error(w, "emoji format not supported", http.StatusBadRequest)
		return
	}

	// Add the reaction info in the database
	if err = rt.db.AddReaction(react, idMess); err != nil {
		context.Logger.WithError(err).Error("error in adding reaction db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the content to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(react); err != nil {
		context.Logger.WithError(err).Error("json comment message encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
