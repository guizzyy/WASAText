package api

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
)

func (rt *_router) checkToken(r *http.Request) (bool, uint64, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false, 0, nil
	}
	token, err := strconv.ParseUint(authHeader, 10, 64)
	if err != nil {
		return false, 0, err
	}
	if isIn, err := rt.db.IsInDatabase(token); !isIn {
		return false, 0, err
	}
	return true, token, nil
}

func (rt *_router) checkStringFormat(name string) (bool, error) {
	pattern := `^.*?$`

	if len(name) < 3 || len(name) > 16 {
		return false, errors.New("string length must be between 3 and 16 characters")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, errors.New("error compiling regex: " + err.Error())
	}
	if re.MatchString(name) {
		return true, nil
	} else {
		return false, errors.New("string contains invalid characters")
	}
}
