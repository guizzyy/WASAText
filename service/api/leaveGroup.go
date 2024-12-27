package api

import (
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get conv id of the group to leave
	convID, err := strconv.ParseUint(params.ByName("convID"), 10, 64)
	if err != nil {
		http.Error(w, "Error converting convID to uint64", http.StatusBadRequest)
		return
	}

	// Query the database to delete the membership of the group conversation
	if err := rt.db.LeaveGroup(convID, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
