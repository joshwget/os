package integration

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

func init() {
	Suite(&QemuSuite{
		runCommand: "../scripts/run",
		sshCommand: "../scripts/ssh",
	})
}

var (
	BusyboxImage = map[string]string{
		"amd64": "busybox",
		"arm":   "armhf/busybox",
		"arm64": "aarch64/busybox",
	}[runtime.GOARCH]
	NginxImage = map[string]string{
		"amd64": "nginx",
		"arm":   "armhfbuild/nginx",
		"arm64": "armhfbuild/nginx",
	}[runtime.GOARCH]
        DockerUrl = "https://experimental.docker.com/builds/Linux/x86_64/docker-1.10.0-dev"
        Version = os.Getenv("VERSION")
        Suffix = os.Getenv("SUFFIX")
)

type QemuSuite struct {
	runCommand string
	sshCommand string
	qemuCmd    *exec.Cmd
}

func (s *QemuSuite) TearDownTest(c *C) {
	c.Assert(s.qemuCmd.Process.Kill(), IsNil)
	time.Sleep(time.Millisecond * 1000)
}

func (s *QemuSuite) RunQemu(additionalArgs ...string) error {
	runArgs := []string{
		"--qemu",
		"--no-rebuild",
		"--no-rm-usr",
		"--fresh",
	}
	runArgs = append(runArgs, additionalArgs...)

	s.qemuCmd = exec.Command(s.runCommand, runArgs...)
	s.qemuCmd.Stdout = os.Stdout
	s.qemuCmd.Stderr = os.Stderr
	if err := s.qemuCmd.Start(); err != nil {
		return err
	}

	return s.WaitForSSH()
}

func (s *QemuSuite) RestartQemu(additionalArgs ...string) error {
	s.qemuCmd.Process.Kill()
	time.Sleep(time.Millisecond * 1000)

	runArgs := []string{
		"--qemu",
		"--no-rebuild",
		"--no-rm-usr",
	}
	runArgs = append(runArgs, additionalArgs...)

	s.qemuCmd = exec.Command(s.runCommand, runArgs...)
	s.qemuCmd.Stdout = os.Stdout
	s.qemuCmd.Stderr = os.Stderr
	if err := s.qemuCmd.Start(); err != nil {
		return err
	}

	return s.WaitForSSH()
}

func (s *QemuSuite) WaitForSSH() error {
	sshArgs := []string{
		"--qemu",
		"docker",
		"version",
		">/dev/null",
		"2>&1",
	}

	var err error
	for i := 0; i < 300; i++ {
		cmd := exec.Command(s.sshCommand, sshArgs...)
		if err = cmd.Run(); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("Failed to connect to SSH: %v", err)
}

func (s *QemuSuite) MakeCall(additionalArgs ...string) error {
	sshArgs := []string{
		"--qemu",
	}
	sshArgs = append(sshArgs, additionalArgs...)

	cmd := exec.Command(s.sshCommand, sshArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
        return cmd.Run()
}

func (s *QemuSuite) MakeHostCall(arg string) error {
	cmd := exec.Command(arg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
        return cmd.Run()
}

func (s *QemuSuite) Reboot() {
        s.MakeCall("sudo reboot")
        time.Sleep(3000 * time.Millisecond)
}

func (s *QemuSuite) CheckCall(c *C, additionalArgs ...string) {
        err := s.MakeCall(additionalArgs...)
	c.Assert(err, IsNil)
}

func (s *QemuSuite) CheckHostCall(c *C, arg string) {
        err := s.MakeHostCall(arg)
	c.Assert(err, IsNil)
}
