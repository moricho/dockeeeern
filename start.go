package dockeeeern

import (
	"fmt"
)

func (c *Container) Start() error {
	if err := c.loadState(); err != nil {
		return fmt.Errorf("failed to load: %s", err)
	}

	return c.startChild()
}

func (c *Container) startChild() error {
	return c.writePipe(child)
}
