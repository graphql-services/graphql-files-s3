package src

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/graphql-services/graphql-files/model"
	"github.com/machinebox/graphql"
)

const (
	graphqlSaveFile = `mutation createFile($input: FileRawCreateInput!) {
		result: createFile(input:$input) {
			id
			uid
			size
			contentType
			url
		}
	}`
	graphqlFetchFile = `query file($uid: ID) {
		file(filter: { uid: $uid }) {
			id
			uid
			size
			contentType
		}
	}`
)

// SaveFile ...
func SaveFile(ctx context.Context, f model.File, auth string, data map[string]interface{}) (model.File, error) {
	var res struct {
		Result model.File
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
func FetchFile(ctx context.Context, uid, auth string) (*model.File, error) {
	var res struct {
		Result *model.File
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
		log.Print(s)
	}

	return client.Run(ctx, req, data)
}
