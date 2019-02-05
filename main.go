package main

import (
	"os"

	"github.com/graphql-services/graphql-files/cmd"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "graphql-files"
	app.Usage = "..."
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		cmd.StartCmd(),
	}

	app.Run(os.Args)
}
