package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("sendMessage not authorized")
		http.Error(w, "sendMessage operation not allowed", http.StatusUnauthorized)
		return
	}

	var mess utilities.Message
	mess.Sender = id
	pMess := &mess

	// Get the photo file path (if an image has been sent)
	if mess.Photo, err = rt.GetPhotoPath(w, r, context); err != nil {
		context.Logger.WithError(err).Error("error during GetPhotoPath sendMessage")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check the correct format for the string
	if mess.Text = r.FormValue("text"); len(mess.Text) == 0 || len(mess.Text) > 250 {
		context.Logger.Error(utilities.ErrTextString)
		http.Error(w, utilities.ErrTextString.Error(), http.StatusBadRequest)
		return
	}

	// Get the conv id where to send the message from the path
	if mess.Conv, err = strconv.ParseUint(params.ByName("convID"), 10, 64); err != nil {
		context.Logger.WithError(err).Error("error in getting convID for sendMessage")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query the database to add the new message
	if err = rt.db.AddMessage(pMess); err != nil {
		context.Logger.WithError(err).Error("error during addMessage db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(mess); err != nil {
		context.Logger.WithError(err).Error("json send message encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
