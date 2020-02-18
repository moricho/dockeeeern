package client

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

func (c *Client) create(args []string) error {
	id, bundlePath := args[0], args[1]
	cmd := []string{"init", id, bundlePath}
	return c.container.Create(id, cmd)
}

func (c *Client) state(args []string) error {
	id := args[0]
	state, err := c.container.State()
	if err != nil {
		return err
	}

	return json.NewEncoder(os.Stdout).Encode(state)
}

func (c *Client) start(args []string) error {
	id := args[0]

	return c.container.Start(id)
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a container",
	Run: func(cmd *cobra.Command, args []string) {
		c.create(args)
	},
}
