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
			uid
			size
			contentType
			url
		}
	}`
	graphqlFetchFile = `query file($uid: ID) {
		result: file(filter: { uid: $uid }) {
			id
			uid
			name
			size
			contentType
		}
	}`
)

// SaveFile ...
func SaveFile(ctx context.Context, f model.UploadResponse, auth string, data map[string]interface{}) (model.UploadResponse, error) {
	var res struct {
		Result model.UploadResponse
	}

	data["name"] = f.Name
	data["size"] = f.Size
	data["uid"] = f.UID
	data["url"] = f.URL
	data["contentType"] = f.ContentType
	req := graphql.NewRequest(graphqlSaveFile)
	req.Var("input", data)

	if auth != "" {
		req.Header.Set("authorization", auth)
	}
	err := sendRequest(ctx, req, &res)

	return res.Result, err
}

// FetchFile ...
func FetchFile(ctx context.Context, uid, auth string) (*model.UploadResponse, error) {
	var res struct {
		Result *model.UploadResponse
	}

	req := graphql.NewRequest(graphqlFetchFile)
	req.Var("uid", uid)

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
