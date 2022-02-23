package src

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/graphql-services/graphql-files/model"
	"github.com/machinebox/graphql"
)

const (
	graphqlSaveFile = `mutation createFile($input: FileCreateInput!) {
		result: createFile(input:$input) {
			id
			name
			size
			contentType
			status
		}
	}`
	graphqlFetchFile = `query file($id: ID!) {
		result: file(id: $id) {
			id
			name
			size
			contentType
			status
		}
	}`
)

// SaveFile ...
func SaveFile(ctx context.Context, f model.UploadResponse, auth string, data map[string]interface{}) (model.UploadResponse, error) {
	var res struct {
		Result model.UploadResponse
	}

	data["id"] = f.ID
	data["name"] = f.Name
	data["size"] = f.Size
	data["contentType"] = f.ContentType
	data["status"] = f.Status
	req := graphql.NewRequest(graphqlSaveFile)
	req.Var("input", data)

	if auth != "" {
		req.Header.Set("authorization", auth)
	}
	err := sendRequest(ctx, req, &res)

	return res.Result, err
}

// FetchFile ...
func FetchFile(ctx context.Context, id, auth string) (*model.UploadResponse, error) {
	var res struct {
		Result *model.UploadResponse
	}

	req := graphql.NewRequest(graphqlFetchFile)
	req.Var("id", id)

	if auth != "" {
		req.Header.Set("authorization", auth)
	}
	err := sendRequest(ctx, req, &res)

	return res.Result, err
}

func sendRequest(ctx context.Context, req *graphql.Request, data interface{}) error {
	URL := os.Getenv("GRAPHQL_URL")

	if URL == "" {
		return fmt.Errorf("Missing required environment variable GRAPHQL_URL")
	}

	client := graphql.NewClient(URL)
	client.Log = func(s string) {
		log.Info(s)
	}

	return client.Run(ctx, req, data)
}
