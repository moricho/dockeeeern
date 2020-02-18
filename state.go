package dockeeeern

import (
	"fmt"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func (c *Container) State() (*specs.State, error) {
	if err := c.loadState(); err != nil {
		return nil, fmt.Errorf("failed to load: %s", err)
	}

	return &c.state, nil
}
