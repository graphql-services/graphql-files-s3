package controller

import (
	"context"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-services/graphql-files/src"
)

// FilesHandler ...
func FilesHandler(r *mux.Router, bucket string) error {

	r.HandleFunc("/{uid}", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("authorization")
		uid := mux.Vars(r)["uid"]

		ctx := context.Background()
		file, err := src.FetchFile(ctx, uid, auth)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if file == nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		s3Object, err := src.GetS3Object(src.GetS3ObjectConfig{
			Bucket: bucket,
			Key:    uid,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer s3Object.Body.Close()

		w.Header().Set("content-type", file.ContentType)
		w.Header().Set("content-size", string(file.Size))
		io.Copy(w, s3Object.Body)
	}).Methods("GET")

	return nil
}
