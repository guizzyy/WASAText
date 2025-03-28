package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getGroupInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, id, err := rt.checkToken(r)
	if err != nil {
		context.Logger.WithError(err).Error("error during checkToken")
		http.Error(w, "Error checking the token", http.StatusInternalServerError)
		return
	}
	if !isAuth {
		context.Logger.Error("getGroupInfo not authorized")
		http.Error(w, "getGroupInfo operation not allowed", http.StatusUnauthorized)
		return
	}

	groupID, err := strconv.ParseUint(params.ByName("convID"), 10, 64)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting convID for getGroupInfo")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	group, err := rt.db.GetConvByID(groupID, id)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting conv for getGroupInfo")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	members, err := rt.db.GetMembers(groupID, id)
	if err != nil {
		context.Logger.WithError(err).Error("error in getting group members for getGroupInfo")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i := range members {
		mem := &members[i]
		if mem.Photo, err = rt.GetFile(mem.Photo); err != nil {
			context.Logger.WithError(err).Error("error in getting file for getGroupInfo")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if group.Photo, err = rt.GetFile(group.Photo); err != nil {
		context.Logger.WithError(err).Error("error in getting group photo for getGroupInfo")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := utilities.ConvResponse{
		Type:    group.Type,
		Name:    group.Name,
		Photo:   group.Photo,
		Members: members,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(response); err != nil {
		context.Logger.WithError(err).Error("json get conversation encode error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
