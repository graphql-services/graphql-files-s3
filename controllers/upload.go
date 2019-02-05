package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-services/graphql-files/model"
	"github.com/graphql-services/graphql-files/src"
)

// UploadHandler ...
func UploadHandler(r *mux.Router, bucket string) error {

	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(10 << 20)
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		f := model.File{
			Name:        header.Filename,
			Size:        header.Size,
			ContentType: header.Header.Get("Content-Type"),
		}

		// Upload to S3
		err = src.UploadToS3(src.UploadToS3Config{
			Bucket:      bucket,
			Key:         f.Name,
			Body:        file,
			Size:        f.Size,
			ContentType: f.ContentType,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Save in GraphQL
		ctx := context.Background()
		response, err := src.SaveFile(ctx, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send response
		data, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(data)
	}).Methods("POST")

	return nil
}
