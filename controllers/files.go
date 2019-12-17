package controller

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/graphql-services/graphql-files/src"
)

// FilesHandler ...
func FilesHandler(r *mux.Router, bucket string) error {

	r.HandleFunc("/{uid}", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("access_token")
		contentDisposition := r.URL.Query().Get("content-disposition")
		uid := mux.Vars(r)["uid"]

		ctx := context.Background()
		file, err := src.FetchFile(ctx, uid, "Bearer "+token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if file == nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		if contentDisposition != "" {
			contentDisposition += ";filename=" + file.Name
		}

		presignedURL, err := src.GetObjectPresignedURL(bucket, uid, contentDisposition)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, presignedURL, 301)
	}).Methods("GET")

	return nil
}
