// +build linux
package dockeeeern

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"github.com/opencontainers/runtime-spec/specs-go"
)

var newroot = "rootfs"

func init() {
	mountProc(newroot)
	pivotRoot(newroot)
	cgroup()
}

type Container struct {
	state specs.State
	Pipe  int `json:"pipe"`
}

func (c *Container) workDirName() string {
	return filepath.Join("/exec", "shoten", c.state.ID)
}

func (c *Container) stateFileName() string {
	return filepath.Join(c.workDirName(), "state.json")
}

func (c *Container) fifoName() string {
	return filepath.Join(c.workDirName(), "pipe.fifo")
}

func (c *Container) saveState() error {
	f, err := os.OpenFile(c.stateFileName(), os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		return fmt.Errorf("failed to open spec file: %w", err)
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(c)
}

func (c *Container) loadState() error {
	f, err := os.Open(c.stateFileName())
	if err != nil {
		return fmt.Errorf("failed to open spec file: %s", err)
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(c)
}

func pivotRoot(newroot string) error {
	putold := filepath.Join(newroot, "/.pivot_root")

	if err := syscall.Mount(newroot, newroot, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}

	if err := os.MkdirAll(putold, 0700); err != nil {
		return err
	}

	if err := syscall.PivotRoot(newroot, putold); err != nil {
		return err
	}

	if err := os.Chdir("/"); err != nil {
		return err
	}

	putold = "/.pivot_root"
	if err := syscall.Unmount(putold, syscall.MNT_DETACH); err != nil {
		return err
	}

	if err := os.RemoveAll(putold); err != nil {
		return err
	}

	return nil
}

func mountProc(newroot string) error {
	source := "proc"
	target := filepath.Join(newroot, "/proc")
	fstype := "proc"
	flags := 0
	data := ""

	os.MkdirAll(target, 0755)
	if err := syscall.Mount(source, target, fstype, uintptr(flags), data); err != nil {
		return err
	}

	return nil
}

func cgroup() error {
	if err := os.MkdirAll("/sys/fs/cgroup/cpu/my-container", 0700); err != nil {
		return fmt.Errorf("Cgroups namespace my-container create failed: %w", err)
	}

	if err := ioutil.WriteFile("/sys/fs/cgroup/cpu/my-container/tasks", []byte(fmt.Sprintf("%d\n", os.Getpid())), 0644); err != nil {
		return fmt.Errorf("Cgroups register tasks to my-container namespace failed: %w", err)
	}

	if err := ioutil.WriteFile("/sys/fs/cgroup/cpu/my-container/cpu.cfs_quota_us", []byte("1000\n"), 0644); err != nil {
		return fmt.Errorf("Cgroups add limit cpu.cfs_quota_us to 1000 failed: %w", err)
	}

	return nil
}
