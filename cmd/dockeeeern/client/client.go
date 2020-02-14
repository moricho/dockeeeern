package client

import (
	"encoding/json"
	"os"

	"github.com/moricho/dockeeeern"
	"github.com/spf13/cobra"
)

var (
	c Client
)

var encoder = json.NewEncoder(os.Stdout)

type Client struct {
	*cobra.Command
	container Container
}

type Container interface {
	Create()
}

func New() *Client {
	setupCmd()
	c.container = &dockeeeern.Container{}
	return &c
}

func setupCmd() {
	rootCmd := &cobra.Command{
		Use:   "dockeeeern",
		Short: "",
		Long:  "",
	}

	c = Client{
		Command: rootCmd,
	}

	c.AddCommand(createCmd)
}
