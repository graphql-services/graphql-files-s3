package controller

import (
	"net/http"

	"github.com/gorilla/mux"
)

// UploadHandler ...
func FilesHandler(r *mux.Router, bucket string) error {

	r.HandleFunc("/files/{uid}", func(w http.ResponseWriter, r *http.Request) {
		uid := mux.Vars(r)["uid"]
		w.Write([]byte(uid))
	}).Methods("GET")

	return nil
}
