package api

import (
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, token, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("leaveGroup not authorized")
		http.Error(w, "leaveGroup operation not allowed", http.StatusUnauthorized)
		return
	}

	// Check if the id from the path and the token correspond
	leaverID, err := strconv.ParseUint(params.ByName("uID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting leaverID for leaveGroup")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if leaverID != token {
		context.Logger.Error("Security error:")
		http.Error(w, "You can't delete another user from the group", http.StatusUnauthorized)
		return
	}

	// Get conv id of the group to leave
	convID, err := strconv.ParseUint(params.ByName("convID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting convID for leaveGroup")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Query the database to delete the membership of the group conversation
	if err = rt.db.LeaveGroup(convID, leaverID); err != nil {
		context.Logger.WithError(err).Error("error during leaveGroup db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
