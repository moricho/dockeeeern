package client

import (
	"github.com/spf13/cobra"
)

func (c *Client) create(args []string) error {
	// id, specPath := args[0], args[1]
	// return c.container.Create("init", id, specPath)
	return nil
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a container",
	Run: func(cmd *cobra.Command, args []string) {
		c.create(args)
	},
}
