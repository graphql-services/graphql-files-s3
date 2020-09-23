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

		UID := uuid.Must(uuid.NewV4()).String()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		f := model.UploadResponse{
			UID:         UID,
			Name:        input.Filename,
			Size:        input.Size,
			ContentType: input.ContentType,
		}

		presignedURL, err := src.PutObjectPresignedURL(bucket, UID)
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

	// r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
	// 	auth := r.Header.Get("authorization")

	// 	r.ParseMultipartForm(10 << 20)
	// 	file, header, err := r.FormFile("file")
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}
	// 	defer file.Close()

	// 	UID := uuid.Must(uuid.NewV4()).String()
	// 	hostURL, err := url.Parse(os.Getenv("HOST_URL"))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	hostURL.Path = path.Join(hostURL.Path, UID)
	// 	f := model.File{
	// 		UID:         UID,
	// 		Name:        header.Filename,
	// 		Size:        header.Size,
	// 		ContentType: header.Header.Get("Content-Type"),
	// 		URL:         hostURL.String(),
	// 	}

	// 	// Upload to S3
	// 	err = src.UploadToS3(src.UploadToS3Config{
	// 		Bucket:      bucket,
	// 		Key:         f.UID,
	// 		Body:        file,
	// 		Size:        f.Size,
	// 		ContentType: f.ContentType,
	// 	})
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Save in GraphQL
	// 	additionalData := map[string]interface{}{}
	// 	for key, value := range r.URL.Query() {
	// 		additionalData[key] = value[0]
	// 	}
	// 	ctx := context.Background()
	// 	response, err := src.SaveFile(ctx, f, auth, additionalData)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Send response
	// 	data, err := json.Marshal(response)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.Header().Set("content-type", "application/json")
	// 	w.Write(data)
	// }).Methods("POST")

	return nil
}
