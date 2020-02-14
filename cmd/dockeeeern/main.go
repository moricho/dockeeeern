package main

import (
	"fmt"
	"os"

	"github.com/moricho/dockeeeern/cmd/dockeeeern/client"
)

func main() {
	c := client.New()
	if err := c.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
