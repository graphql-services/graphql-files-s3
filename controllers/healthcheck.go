package controller

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// HealthcheckHandler ...
func HealthcheckHandler(r *mux.Router) error {

	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	return nil
}
