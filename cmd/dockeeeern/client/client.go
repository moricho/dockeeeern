package client

import (
	"encoding/json"
	"os"

	"github.com/opencontainers/runtime-spec/specs-go"
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
	Create(id string, args []string) error
	State(id string) (*specs.State, error)
	Start(id string) error
	Delete(id string) error
}

func New(container Container) *Client {
	setupCmd()
	c.container = container
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
