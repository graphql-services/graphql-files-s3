package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	controller "github.com/graphql-services/graphql-files/controllers"
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

			http.Handle("/", r)

			fmt.Println("starting on port: " + port)
			if err := http.ListenAndServe(":"+port, nil); err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			return nil
		},
	}
}
