package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/opencontainers/specs"
)

var generateCommand = cli.Command{
	Name:  "generate",
	Usage: "generate a OCI spec file",
	Action: func(context *cli.Context) {
		spec := specs.LinuxSpec{
			Spec: specs.Spec{
				Version: specs.Version,
				Platform: specs.Platform{
					OS:   runtime.GOOS,
					Arch: runtime.GOARCH,
				},
				Root: specs.Root{
					Path:     "rootfs",
					Readonly: true,
				},
				Process: specs.Process{
					Terminal: true,
					User:     specs.User{},
					Args: []string{
						"sh",
					},
					Env: []string{
						"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
						"TERM=xterm",
					},
				},
				Hostname: "shell",
				Mounts: []specs.MountPoint{
					{
						Name: "proc",
						Path: "/proc",
					},
					{
						Name: "dev",
						Path: "/dev",
					},
					{
						Name: "devpts",
						Path: "/dev/pts",
					},
					{
						Name: "shm",
						Path: "/dev/shm",
					},
					{
						Name: "mqueue",
						Path: "/dev/mqueue",
					},
					{
						Name: "sysfs",
						Path: "/sys",
					},
					{
						Name: "cgroup",
						Path: "/sys/fs/cgroup",
					},
				},
			},
			Linux: specs.Linux{
				Capabilities: []string{
					"AUDIT_WRITE",
					"KILL",
					"NET_BIND_SERVICE",
				},
			},
		}
		rspec := specs.LinuxRuntimeSpec{
			RuntimeSpec: specs.RuntimeSpec{
				Mounts: map[string]specs.Mount{
					"proc": {
						Type:    "proc",
						Source:  "proc",
						Options: nil,
					},
					"dev": {
						Type:    "tmpfs",
						Source:  "tmpfs",
						Options: []string{"nosuid", "strictatime", "mode=755", "size=65536k"},
					},
					"devpts": {
						Type:    "devpts",
						Source:  "devpts",
						Options: []string{"nosuid", "noexec", "newinstance", "ptmxmode=0666", "mode=0620", "gid=5"},
					},
					"shm": {
						Type:    "tmpfs",
						Source:  "shm",
						Options: []string{"nosuid", "noexec", "nodev", "mode=1777", "size=65536k"},
					},
					"mqueue": {
						Type:    "mqueue",
						Source:  "mqueue",
						Options: []string{"nosuid", "noexec", "nodev"},
					},
					"sysfs": {
						Type:    "sysfs",
						Source:  "sysfs",
						Options: []string{"nosuid", "noexec", "nodev"},
					},
					"cgroup": {
						Type:    "cgroup",
						Source:  "cgroup",
						Options: []string{"nosuid", "noexec", "nodev", "relatime", "ro"},
					},
				},
			},
			Linux: specs.LinuxRuntime{
				Namespaces: []specs.Namespace{
					{
						Type: "pid",
					},
					{
						Type: "network",
					},
					{
						Type: "ipc",
					},
					{
						Type: "uts",
					},
					{
						Type: "mount",
					},
				},
				Rlimits: []specs.Rlimit{
					{
						Type: "RLIMIT_NOFILE",
						Hard: uint64(1024),
						Soft: uint64(1024),
					},
				},
				Devices: []specs.Device{
					{
						Type:        'c',
						Path:        "/dev/null",
						Major:       1,
						Minor:       3,
						Permissions: "rwm",
						FileMode:    0666,
						UID:         0,
						GID:         0,
					},
					{
						Type:        'c',
						Path:        "/dev/random",
						Major:       1,
						Minor:       8,
						Permissions: "rwm",
						FileMode:    0666,
						UID:         0,
						GID:         0,
					},
					{
						Type:        'c',
						Path:        "/dev/full",
						Major:       1,
						Minor:       7,
						Permissions: "rwm",
						FileMode:    0666,
						UID:         0,
						GID:         0,
					},
					{
						Type:        'c',
						Path:        "/dev/tty",
						Major:       5,
						Minor:       0,
						Permissions: "rwm",
						FileMode:    0666,
						UID:         0,
						GID:         0,
					},
					{
						Type:        'c',
						Path:        "/dev/zero",
						Major:       1,
						Minor:       5,
						Permissions: "rwm",
						FileMode:    0666,
						UID:         0,
						GID:         0,
					},
					{
						Type:        'c',
						Path:        "/dev/urandom",
						Major:       1,
						Minor:       9,
						Permissions: "rwm",
						FileMode:    0666,
						UID:         0,
						GID:         0,
					},
				},
				Resources: &specs.Resources{
					Memory: specs.Memory{
						Swappiness: -1,
					},
				},
				Seccomp: specs.Seccomp{
					DefaultAction: "SCMP_ACT_ALLOW",
					Syscalls:      []*specs.Syscall{},
				},
			},
		}
		checkNoFile := func(name string) error {
			_, err := os.Stat(name)
			if err == nil {
				return fmt.Errorf("File %s exists. Remove it first", name)
			}
			if !os.IsNotExist(err) {
				return err
			}
			return nil
		}
		cName := "config.json"
		rName := "runtime.json"
		if err := checkNoFile(cName); err != nil {
			logrus.Fatal(err)
		}
		if err := checkNoFile(rName); err != nil {
			logrus.Fatal(err)
		}
		data, err := json.MarshalIndent(&spec, "", "\t")
		if err != nil {
			logrus.Fatal(err)
		}
		if err := ioutil.WriteFile(cName, data, 0666); err != nil {
			logrus.Fatal(err)
		}
		rdata, err := json.MarshalIndent(&rspec, "", "\t")
		if err != nil {
			logrus.Fatal(err)
		}
		if err := ioutil.WriteFile(rName, rdata, 0666); err != nil {
			logrus.Fatal(err)
		}
	},
}
