package main

import (
	"log"
	"os"

	"github.com/akrylysov/algnhsa"
	"github.com/gorilla/mux"
	controller "github.com/graphql-services/graphql-files/controllers"
	"github.com/rs/cors"
)

func main() {

	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		log.Fatal("Missing S3_BUCKET environment variable")
	}

	r := mux.NewRouter()
	controller.HealthcheckHandler(r)
	controller.UploadHandler(r, bucket)
	controller.FilesHandler(r, bucket)

	handler := cors.AllowAll().Handler(r)

	algnhsa.ListenAndServe(handler, nil)
}
