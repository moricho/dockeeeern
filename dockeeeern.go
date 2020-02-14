// +build linux
package dockeeeern

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sys/unix"
)

func Run() {
	clone()
}

func Setup() {
	newroot := "/rootfs"
	pivotRoot(newroot)
	mountProc(newroot)
	cgroup()
}

type Container struct {
}

// clone: clone プロセスのcloneをし,その子プロセスにおいて各Namespaceの分離、UID/GIDのマッピング、必要なファイルシステムのマウントなどを行う
func (c *Container) Create() {
	cmd := exec.Command("/bin/sh")
	cmd.SysProcAttr = &unix.SysProcAttr{
		Cloneflags: unix.CLONE_NEWIPC |
			unix.CLONE_NEWNET |
			unix.CLONE_NEWNS |
			unix.CLONE_NEWPID |
			unix.CLONE_NEWUSER |
			unix.CLONE_NEWUTS,
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

	cmd.Env = []string{"PS1=-[namespace-process]-# "}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
