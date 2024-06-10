package main

import (
	"os"

	"github.com/greenplum-db/gpdb/gpservice/internal/cli"
)

func main() {
	root := cli.RootCommand()

	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
