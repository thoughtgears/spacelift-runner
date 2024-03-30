package main

import (
	"github.com/thoughtgears/spacelift-runner/cli"
	"github.com/thoughtgears/spacelift-runner/internal/log"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		log.Log(log.ErrorLevel, err.Error())
	}
}
