package api

import (
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, "Error checking the token", http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("uncommentMessage not authorized")
		http.Error(w, "uncommentMessage operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get the id of the message which we want to delete the reaction from
	messID, err := strconv.ParseUint(params.ByName("messID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting messID for uncommentMessage")
		http.Error(w, "error getting the message id", http.StatusInternalServerError)
		return
	}

	// Query the database to delete the reaction with given mess id and username
	if err = rt.db.RemoveReaction(messID, id); err != nil {
		context.Logger.WithError(err).Error("error removing reaction db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
