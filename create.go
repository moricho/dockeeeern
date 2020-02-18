package dockeeeern

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/unix"
)

func (c *Container) Create(id string, args []string) error {
	c.ID = id

	if err := c.initWorkspace(); err != nil {
		return fmt.Errorf("failed to init workspace: %w", err)
	}

	if err := c.saveState(); err != nil {
		return fmt.Errorf("failed to save state: %w", err)
	}

	select {
	case err := <-c.clone(args...):
		if err != nil {
			return err
		}
	}

	if err := c.loadState(); err != nil {
		return fmt.Errorf("failed to load: %s", err)
	}
	if err := c.saveState(); err != nil {
		return fmt.Errorf("failed to save: %s", err)
	}

	return nil

}

func (c *Container) clone(args ...string) chan error {
	ch := make(chan error)
	go func() {
		defer close(ch)

		cmd := buildCloneCmd(args...)
		if err := cmd.Start(); err != nil {
			ch <- err
		}
		c.Pid = cmd.Process.Pid

		ch <- cmd.Wait()
	}()

	return ch
}

func buildCloneCmd(args ...string) *exec.Cmd {
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWUSER |
			unix.CLONE_NEWNET |
			unix.CLONE_NEWPID |
			unix.CLONE_NEWIPC |
			unix.CLONE_NEWUTS |
			unix.CLONE_NEWNS,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd
}

func (c *Container) initWorkspace() error {
	// ワークディレクトリの作成
	path := c.workDirName()
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	err := os.MkdirAll(path, 0744)
	if err != nil {
		return fmt.Errorf("failed to init directoery %s: %w", path, err)
	}

	// コンテナのstateファイルの作成
	stateFilePath := c.stateFileName()
	f, err := os.Create(stateFilePath)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	// FIFOの作成
	fifoPath := c.fifoName()
	if err := syscall.Mkfifo(fifoPath, 700); err != nil {
		return err
	}
	pipe, err := unix.Open(fifoPath, os.O_RDONLY|unix.O_NONBLOCK, 700)
	if err != nil {
		return err
	}
	c.Pipe = pipe

	return nil
}
