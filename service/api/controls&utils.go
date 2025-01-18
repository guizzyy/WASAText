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

// Check that the strings provided for the name respect the pattern defined
func (rt *_router) checkStringFormat(name string) (bool, error) {
	pattern := `^.*?$`

	if len(name) < 3 || len(name) > 16 {
		return false, utilities.ErrNameString
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, errors.New("error compiling regex: " + err.Error())
	}
	if re.MatchString(name) {
		return true, nil
	} else {
		return false, utilities.ErrNameString
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
		context.Logger.WithError(err).Error("error during ParseMultipartForm sendMessage")
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
		defer file.Close()
		return "", err
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
