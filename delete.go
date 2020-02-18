package dockeeeern

import (
	"fmt"
	"os"
)

func (c *Container) Delete() error {
	if err := c.loadState(); err != nil {
		return fmt.Errorf("failed to load: %W", err)
	}

	return os.RemoveAll(c.workDirName())
}
