package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getMembers(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("addToGroup operation not authorized")
		http.Error(w, "addToGroup operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get the conversation id interested from the path
	convID, err := strconv.ParseUint(params.ByName("convID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting convID for getMembers")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Query the database to retrieve all the members information
	members, err := rt.db.GetMembers(convID, id)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting members db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range members {
		members[i].Photo, err = rt.GetFile(members[i].Photo)
		if err != nil {
			context.Logger.WithError(err).Error("error in GetFile in getMembers")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(members); err != nil {
		context.Logger.WithError(err).Error("json get members encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
