package main

import (
	"github.com/zero-miao/curl-go/cmd"
	"os"
)

func main() {
	if err := cmd.Entry(); err != nil {
		os.Exit(1)
	}
}
