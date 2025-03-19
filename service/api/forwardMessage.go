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

	// Get the new conversation where forward the message to
	var newConv utilities.Conversation
	var strConv map[string]string
	if err = json.NewDecoder(r.Body).Decode(&strConv); err != nil {
		context.Logger.WithError(err).Error("json forward message decode error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if newConv.ID, err = strconv.ParseUint(strConv["id"], 10, 64); err != nil {
		context.Logger.WithError(err).Error("error in convert id conv from string to uint64")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the old conversation from the path
	var oldConv utilities.Conversation
	if oldConv.ID, err = strconv.ParseUint(params.ByName("convID"), 10, 64); err != nil {
		context.Logger.WithError(err).Error("error in getting old convID for forward message")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the old conversation exists in the database
	if isIn, err := rt.db.IsConvInDatabase(oldConv.ID); err != nil {
		context.Logger.WithError(err).Error("error in checking if the conversation exists")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !isIn {
		context.Logger.Error("conversation does not exist")
		http.Error(w, database.ErrConversationNotFound.Error(), http.StatusBadRequest)
	}

	// Check if the user is in the old conversation from where it takes the message
	if check, err := rt.db.IsUserInConv(oldConv.ID, id); err != nil {
		context.Logger.WithError(err).Error("error in checking user in conversation db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !check {
		context.Logger.Error("Membership error:")
		http.Error(w, "User not in the old conversation", http.StatusBadRequest)
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

	// Check if the message exists in the conversation interested
	if check, err := rt.db.IsMessageInConv(messID, oldConv.ID); err != nil {
		context.Logger.WithError(err).Error("error during IsMessageInConv db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !check {
		context.Logger.Error("message not in this conversation")
		http.Error(w, database.ErrMessageNotFound.Error(), http.StatusBadRequest)
		return
	}

	// Add the forwarded message as a new message in the database
	msg.Sender, err = rt.db.GetUserByID(id)
	if err != nil {
		context.Logger.WithError(err).Error("error during GetUserByID db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msg.Conv = newConv.ID
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
