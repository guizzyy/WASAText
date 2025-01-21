/*
 In this file GO, there is the definition of various function used to manipulate
 and control the data, in order to ensure the integrity of the application.
*/

package api

import (
	"errors"
	"fmt"
	"git.guizzyy.it/WASAText/service/api/reqcontext"
	"git.guizzyy.it/WASAText/service/utilities"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
	if isIn, err := rt.db.IsUserInDatabase(token); err != nil || !isIn {
		return false, 0, err
	}
	return true, token, nil
}

// Check that the strings provided for the username respect the pattern defined
func (rt *_router) checkStringFormat(name string) (bool, error) {
	pattern := `^.*?$`

	if len(name) < 3 || len(name) > 16 {
		return false, utilities.ErrUsernameString
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, errors.New("error compiling regex: " + err.Error())
	}
	if re.MatchString(name) {
		return true, nil
	} else {
		return false, utilities.ErrUsernameString
	}
}

// Check that the strings provided for the group name respect the pattern defined
func (rt *_router) checkGroupStringFormat(grName string) (bool, error) {
	pattern := `^.*?$`

	if len(grName) < 3 || len(grName) > 25 {
		return false, utilities.ErrGroupNameString
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, errors.New("error compiling regex: " + err.Error())
	}
	if re.MatchString(grName) {
		return true, nil
	} else {
		return false, utilities.ErrGroupNameString
	}
}

// Check that the reaction provided respect the pattern (is an emoji)
func (rt *_router) checkEmojiFormat(emoji string) (bool, error) {
	pattern := `^[\\u1F600-\\u1F64F\\u1F300-\\u1F5FF\\u1F680-\\u1F6FF\\u1F700-\\u1F77F\\u1F780-\\u1F7FF\\u1F800-\\u1F8FF\\u1F900-\\u1F9FF\\u1FA00-\\u1FA6F\\u1FA70-\\u1FAFF\\u2600-\\u26FF\\u2700-\\u27BF\\u2300-\\u23FF\\u2B50\\u23F0\\u231A\\u25AA\\u25FE\\u2B06\\u2194\\u2B05\\u2195\\u21A9\\u21AA\\u2753\\u2754\\u2755\\u274C\\u274E\\u2199]$`

	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, errors.New("error compiling regex: " + err.Error())
	}
	if re.MatchString(emoji) {
		return true, nil
	} else {
		return false, errors.New("error is not emoji")
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

func (rt *_router) GetPhotoPath(w http.ResponseWriter, r *http.Request, context reqcontext.RequestContext) (string, error) {
	// Set the dimension of the request body
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		context.Logger.WithError(err).Error("error during ParseMultipartForm")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	// Get the file from the request body (if missing, return empty string)
	file, handler, err := r.FormFile("photo")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			file = nil
		} else {
			context.Logger.WithError(err).Error("Error during file upload")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return "", err
		}
	}
	if file == nil {
		return "", nil
	}
	defer file.Close()

	// Check if the file is an image file
	if isImage, err := rt.checkFileFormat(file); err != nil {
		context.Logger.WithError(err).Error("Error during check file format")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	} else if !isImage {
		context.Logger.Error("File is not an image")
		http.Error(w, "Not a file image uploaded", http.StatusBadRequest)
		return "", err
	}

	// Create a file in the folder (unique name) and copy the image in it
	uniqueFile := fmt.Sprintf("%s_%s", uuid.New().String(), filepath.Ext(handler.Filename))
	filePath := filepath.Join("service/api/photos", uniqueFile)
	fileLocal, err := os.Create(filePath)
	if err != nil {
		context.Logger.WithError(err).Error("Error during file creation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	defer fileLocal.Close()
	_, err = io.Copy(fileLocal, file)
	if err != nil {
		context.Logger.WithError(err).Error("Error during file copy")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	return filePath, nil
}

func (rt *_router) DeletePhotoPath(oldPhoto string) error {
	if err := os.Remove(oldPhoto); err != nil {
		return err
	}
	return nil
}
