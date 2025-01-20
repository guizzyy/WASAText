package api

import (
	"encoding/json"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params, context reqcontext.RequestContext) {
	// Get the username for log in the request body
	var userLog utilities.User
	if err := json.NewDecoder(r.Body).Decode(&userLog); err != nil {
		context.Logger.WithError(err).Error("json login decode error")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pUser := &userLog

	// Check if the username provided has the correct format
	if check, err := rt.checkStringFormat(userLog.Username); err != nil {
		context.Logger.WithError(err).Error("error during string format check")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !check {
		context.Logger.Error(utilities.ErrUsernameString)
		http.Error(w, utilities.ErrUsernameString.Error(), http.StatusBadRequest)
		return
	}

	// Ask the database if it is a new/existing user and get their ID and photo
	isNew, err := rt.db.LogUser(pUser)
	if err != nil {
		context.Logger.WithError(err).Error("error during logUser db")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := utilities.LoginResponse{
		UserLogged: userLog,
	}

	w.Header().Set("Content-Type", "application/json")
	if !isNew {
		// The user already exists
		w.WriteHeader(http.StatusOK)
		response.Message = "Login successful"
		if err = json.NewEncoder(w).Encode(response); err != nil {
			context.Logger.WithError(err).Error("json login encode error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// The user is new
		w.WriteHeader(http.StatusCreated)
		response.Message = "User created successfully"
		if err = json.NewEncoder(w).Encode(response); err != nil {
			context.Logger.WithError(err).Error("json login encode error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
