package main

import (
	"github.com/coveo/credentials-sync/cli"
)

var (
	commit  = "COMMITHASH"
	date    = "1900-01-01"
	version = "main"
)

func main() {
	cli.Execute(commit, date, version)
}
