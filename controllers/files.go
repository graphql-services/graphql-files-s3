package controller

import (
	"net/http"

	"github.com/gorilla/mux"
)

// FilesHandler ...
func FilesHandler(r *mux.Router, bucket string) error {

	r.HandleFunc("/files/{uid}", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("authorization")
		uid := mux.Vars(r)["uid"]
		w.Write([]byte(uid + auth))
	}).Methods("GET")

	return nil
}
