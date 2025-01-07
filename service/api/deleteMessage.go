package api

import (
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, _, err := rt.checkToken(r)
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

	// Query the database to delete the message
	if err = rt.db.RemoveMessage(mess_id); err != nil {
		context.Logger.WithError(err).Error("error in deleting message db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
