package main

import (
	"github.com/zero-miao/curl-go/entry"
	"os"
)

func main() {
	if err := entry.Entry(); err != nil {
		os.Exit(1)
	}
}
