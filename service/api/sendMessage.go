package api

import (
	"encoding/json"
	"fmt"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
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
	mess.Sender, err = rt.db.GetUserByID(id)
	if err != nil {
		context.Logger.WithError(err).Error("error during GetUserByID db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pMess := &mess

	// Get the photo file path (if an image has been sent)
	mPhoto, file, err := rt.GetFilePath(w, r, context)
	if err != nil {
		context.Logger.WithError(err).Error("error during GetFilePath for message photo")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if mPhoto != "" {
		uDir := strconv.FormatUint(mess.Sender.ID, 10)
		if _, err := os.Stat("./tmp/uploads/" + uDir + "/sent"); os.IsNotExist(err) {
			err = os.MkdirAll("./tmp/uploads/"+uDir+"/sent", 0755)
			if err != nil {
				context.Logger.WithError(err).Error("can't create the folder")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if err != nil {
			context.Logger.WithError(err).Error("can't get the folder")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		filePath := fmt.Sprintf("./tmp/uploads/%s/%s/%s", uDir, "sent", mPhoto)
		dst, err := os.Create(filePath)
		if err != nil {
			context.Logger.WithError(err).Error("Error during create file path")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()
		if _, err = io.Copy(dst, file); err != nil {
			context.Logger.WithError(err).Error("Error during copy file to path")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mess.Photo = filePath
	}

	// Check the correct format for the string
	mess.Text = r.FormValue("text")
	if mess.Text == "" && mess.Photo == "" {
		context.Logger.Error("message is completely empty")
		http.Error(w, "message is completely empty", http.StatusBadRequest)
		return
	}
	if len(mess.Text) > 250 {
		context.Logger.Error(utilities.ErrTextString)
		http.Error(w, utilities.ErrTextString.Error(), http.StatusBadRequest)
		return
	}

	reply := r.FormValue("reply")
	if reply != "" {
		replyID, err := strconv.ParseUint(reply, 10, 64)
		if err != nil {
			context.Logger.WithError(err).Error("error in getting message reply id for sendMessage")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messReply, err := rt.db.GetMessageInfo(replyID)
		if err != nil {
			context.Logger.WithError(err).Error("error during GetMessageInfo for reply sendMessage")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mess.ReplyID = replyID
		mess.ReplyText = messReply.Text
		mess.Photo = messReply.Photo
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

	// Get the photo sent for the API
	if mess.Photo != "" {
		if mess.Photo, err = rt.GetFile(mess.Photo); err != nil {
			context.Logger.WithError(err).Error("error during GetFile for mess photo")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(mess); err != nil {
		context.Logger.WithError(err).Error("json send message encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
