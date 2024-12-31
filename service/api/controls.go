/*
 In this file GO, there is the definition of various function used to manipulate
 and control the data, in order to ensure the integrity of the application.
*/

package api

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"
)

// Check the authorization token provided whether is correct, missed, unknown
func (rt *_router) checkToken(r *http.Request) (bool, uint64, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false, 0, nil
	}
	token, err := strconv.ParseUint(authHeader, 10, 64)
	if err != nil {
		return false, 0, err
	}
	if isIn, err := rt.db.IsInDatabase(token); err != nil || !isIn {
		return false, 0, nil
	}
	return true, token, nil
}

// Check that the strings provided for the name respect the pattern defined
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

// Check that the file provided is an image file (with various extension in it)
func (rt *_router) checkFileFormat(file multipart.File) (bool, error) {
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return false, fmt.Errorf("error reading file: %v", err)
	}
	contentType := http.DetectContentType(buffer)
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/apng", "image/bmp", "image/tiff":
		return true, nil
	default:
		return false, nil
	}
}
