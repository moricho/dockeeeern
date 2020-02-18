package dockeeeern

import (
	"fmt"
	"os"
)

type pipeDest int

const (
	self pipeDest = iota
	child
)

func (c *Container) getPipePath(dest pipeDest) string {
	if dest == self {
		return fmt.Sprintf("/proc/self/fd/%d", c.Pipe)
	}

	return fmt.Sprintf("/proc/%d/fd/%d", c.state.Pid, c.Pipe)
}

func (c *Container) writePipe(dest pipeDest) error {
	pipePath := c.getPipePath(dest)
	f, err := os.OpenFile(pipePath, os.O_WRONLY, 0700)
	if err != nil {
		return err
	}

	return f.Close()
}

func (c *Container) readPipe(dest pipeDest) error {
	pipePath := c.getPipePath(dest)
	f, err := os.Open(pipePath)
	if err != nil {
		return err
	}

	return f.Close()
}
