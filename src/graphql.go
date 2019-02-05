package src

import (
	"context"
	"fmt"
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
func SaveFile(ctx context.Context, f model.File) (model.File, error) {
	var res struct {
		Result model.File
	}

	req := graphql.NewRequest(graphqlSaveFile)
	req.Var("input", f)
	err := sendRequest(ctx, req, &res)

	return res.Result, err
}

func sendRequest(ctx context.Context, req *graphql.Request, data interface{}) error {
	URL := os.Getenv("GRAPHQL_URL")

	if URL == "" {
		return fmt.Errorf("Missing required environment variable GRAPHQL_URL")
	}

	client := graphql.NewClient(URL)
	// client.Log = func(s string) {
	// 	glog.Info(s)
	// }

	return client.Run(ctx, req, data)
}
