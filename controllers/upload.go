package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/graphql-services/graphql-files/model"
	"github.com/graphql-services/graphql-files/src"
)

type UploadInput struct {
	Filename    string
	Size        int64
	ContentType string
	Status      string
}

// UploadHandler ...
func UploadHandler(r *mux.Router, bucket string) error {

	r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("authorization")

		var input UploadInput
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ID := uuid.Must(uuid.NewV4()).String()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f := model.UploadResponse{
			ID:          ID,
			Name:        input.Filename,
			Size:        input.Size,
			ContentType: input.ContentType,
			Status:      input.Status,
		}

		presignedURL, err := src.PutObjectPresignedURL(bucket, ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Save in GraphQL
		additionalData := map[string]interface{}{}
		for key, value := range r.URL.Query() {
			additionalData[key] = value[0]
		}
		ctx := context.Background()
		response, err := src.SaveFile(ctx, f, auth, additionalData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.UploadURL = presignedURL

		json.NewEncoder(w).Encode(response)
	}).Methods("POST")

	return nil
}
