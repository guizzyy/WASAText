package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Check authorization for the operation
	isAuth, _, err := rt.checkToken(r)
	if err != nil {
		http.Error(w, "Error checking the token", http.StatusUnauthorized)
		return
	}
	if !isAuth {
		http.Error(w, "Operation not allowed", http.StatusUnauthorized)
		return
	}

	// Get conv id and the new name to update from request body and path params
	convID, err := strconv.ParseUint(params.ByName("convID"), 10, 64)
	if err != nil {
		http.Error(w, "Error converting convID to uint64", http.StatusBadRequest)
		return
	}
	var conv utilities.Conversation
	conv.ID = convID
	if err := json.NewDecoder(r.Body).Decode(&conv); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err :=
}
