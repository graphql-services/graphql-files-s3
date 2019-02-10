package cmd

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	controller "github.com/graphql-services/graphql-files/controllers"
	"github.com/rs/cors"
	"github.com/urfave/cli"
)

// StartCmd ...
func StartCmd() cli.Command {
	return cli.Command{
		Name:        "start",
		Description: "start service",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "bucket",
				EnvVar: "S3_BUCKET",
			},
			cli.StringFlag{
				Name:   "p,port",
				EnvVar: "PORT",
				Value:  "80",
			},
		},
		Action: func(c *cli.Context) error {
			port := c.String("port")
			bucket := c.String("bucket")
			if bucket == "" {
				log.Fatal("Missing S3_BUCKET environment variable")
			}

			r := mux.NewRouter()
			controller.HealthcheckHandler(r)
			controller.UploadHandler(r, bucket)
			controller.FilesHandler(r, bucket)

			handler := cors.AllowAll().Handler(r)

			fmt.Println("starting on port: " + port)
			if err := http.ListenAndServe(":"+port, handler); err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			return nil
		},
	}
}
