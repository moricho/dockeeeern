package dockeeeern

import (
	"fmt"
	"io/ioutil"
	"os"
)

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
