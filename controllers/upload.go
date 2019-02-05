package controller

import (
	"net/http"

	"github.com/gorilla/mux"
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

		err = src.UploadToS3(src.UploadToS3Config{
			Body:        file,
			Size:        header.Size,
			ContentType: header.Header.Get("Content-Type"),
			Bucket:      bucket,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		// w.Header().Set("content-type", "application/json")
		// w.Write(data)
	}).Methods("POST")

	return nil
}
