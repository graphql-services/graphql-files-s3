package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/graphql-services/graphql-files/model"
	"github.com/graphql-services/graphql-files/src"
)

// FilesHandler ...
func FilesHandler(r *mux.Router, bucket string) error {

	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader == "" {
			token := r.URL.Query().Get("access_token")
			if token != "" {
				authorizationHeader = "Bearer " + token
			}
		}
		contentDisposition := r.URL.Query().Get("content-disposition")
		expire := r.URL.Query().Get("expire")
		id := mux.Vars(r)["id"]

		expiration := 15 * time.Minute
		if expire != "" {
			seconds, err := strconv.Atoi(expire)
			if err != nil {
				http.Error(w, fmt.Sprintf("invalid expire attribute, error: %s", err.Error()), http.StatusBadRequest)
				return
			}
			expiration = time.Second * time.Duration(seconds)
		}

		ctx := context.Background()
		file, err := src.FetchFile(ctx, id, authorizationHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if file == nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		if contentDisposition == "" {
			if strings.HasPrefix(file.ContentType, "image/") || file.ContentType == "application/pdf" {
				contentDisposition = "inline"
			} else {
				contentDisposition = "attachment"
			}
		}
		if contentDisposition != "" {
			contentDisposition += ";filename=" + url.QueryEscape(file.Name)
		}

		presignedURL, err := src.GetObjectPresignedURL(bucket, id, contentDisposition, expiration)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := model.GetFileResponse{URL: presignedURL}

		json.NewEncoder(w).Encode(res)
	}).Methods("GET")

	return nil
}
