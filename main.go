package main

import (
	"os"

	"github.com/Yoseph-code/go-blockchain/cmd"
)

func main() {
	defer os.Exit(0)
	cli := cmd.CommandLine{}
	cli.Run()
}
