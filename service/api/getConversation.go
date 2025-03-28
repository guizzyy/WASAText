package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, "Error checking the token", http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("getConversation not authorized")
		http.Error(w, "getConversation operation not allowed", http.StatusUnauthorized)
		return
	}

	var convID uint64
	// Get the conv id from the path
	if convID, err = strconv.ParseUint(params.ByName("convID"), 10, 64); err != nil {
		context.Logger.WithError(err).Error("error in getting convID for getConversation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conv, err := rt.db.GetConvByID(convID, id)
	if err != nil {
		context.Logger.WithError(err).Error("error getting conv from db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if conv.Photo, err = rt.GetFile(conv.Photo); err != nil {
		context.Logger.WithError(err).Error("error getting conv photo for getConversation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastMessageID, err := strconv.ParseUint(r.URL.Query().Get("lastID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting lastMessageID for getConversation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query the database to retrieve all messages for the conversation
	messages, err := rt.db.GetConversation(convID, id, lastMessageID)
	if err != nil {
		context.Logger.WithError(err).Error("error during getConversation db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query the database to retrieve all members in the conversations
	members, err := rt.db.GetMembers(convID, id)
	if err != nil {
		context.Logger.WithError(err).Error("error during getMembers db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range messages {
		msg := &messages[i]
		if msg.Status, err = rt.db.CheckStatus(msg.ID, msg.Sender.ID); err != nil {
			context.Logger.WithError(err).Error("error during checkStatus")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if msg.Photo != "" {
			if msg.Photo, err = rt.GetFile(msg.Photo); err != nil {
				context.Logger.WithError(err).Error("error during GetFile")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if msg.Sender.Photo, err = rt.GetFile(msg.Sender.Photo); err != nil {
			context.Logger.WithError(err).Error("error during GetFile")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	for i := range members {
		mem := &members[i]
		if mem.Photo, err = rt.GetFile(mem.Photo); err != nil {
			context.Logger.WithError(err).Error("error during GetFile")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	response := utilities.ConvResponse{
		Type:     conv.Type,
		Name:     conv.Name,
		Photo:    conv.Photo,
		Messages: messages,
		Members:  members,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		context.Logger.WithError(err).Error("json get conversation encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
