package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

			region := os.Getenv("S3_REGION")
			session, err := session.NewSession(&aws.Config{Region: aws.String(region)})
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			r := mux.NewRouter()
			controller.HealthcheckHandler(r)
			controller.UploadHandler(r, session, bucket)
			controller.FilesHandler(r, session, bucket)

			handler := cors.AllowAll().Handler(r)

			fmt.Println("starting on port: " + port)
			if err := http.ListenAndServe(":"+port, handler); err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			return nil
		},
	}
}
