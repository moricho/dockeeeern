package client

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	c Client
)

type Client struct {
	*cobra.Command
}

func New() *Client {
	setup()
	return &c
}

func setup() {
	rootCmd := &cobra.Command{
		Use:   "dockeeeern",
		Short: "",
		Long:  "",
	}

	c = Client{
		rootCmd,
	}

	// c.AddCommand(createCmd)
}
