package dockerinit

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/os/util"
)

const (
	consoleDone = "/run/console-done"
	dockerConf  = "/var/lib/rancher/conf/docker"
)

func Main() {
	if os.Getenv("DOCKER_CONF_SOURCED") == "" {
		if err := shellProxy(os.Args); err != nil {
			log.Warnf("Failed to source %s: %v", dockerConf, err)
		}
	}

	for {
		if _, err := os.Stat(consoleDone); err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	dockerBin := "/usr/bin/docker"
	for _, binPath := range []string{
		"/opt/bin",
		"/usr/local/bin",
		"/var/lib/rancher/docker",
	} {
		if util.ExistsAndExecutable(path.Join(binPath, "dockerd")) {
			dockerBin = path.Join(binPath, "dockerd")
			break
		}
		if util.ExistsAndExecutable(path.Join(binPath, "docker")) {
			dockerBin = path.Join(binPath, "docker")
			break
		}
	}

	if err := syscall.Mount("", "/", "", syscall.MS_SHARED|syscall.MS_REC, ""); err != nil {
		log.Error(err)
	}
	if err := syscall.Mount("", "/run", "", syscall.MS_SHARED|syscall.MS_REC, ""); err != nil {
		log.Error(err)
	}

	mountInfo, err := ioutil.ReadFile("/proc/self/mountinfo")
	if err != nil {
		log.Fatal(err)
	}

	for _, mount := range strings.Split(string(mountInfo), "\n") {
		if strings.Contains(mount, "/var/lib/docker /var/lib/docker") && strings.Contains(mount, "rootfs") {
			os.Setenv("DOCKER_RAMDISK", "1")
		}
	}

	args := []string{
		"dockerlaunch",
		dockerBin,
	}

	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
	}

	if os.Getenv("DOCKER_OPTS") != "" {
		args = append(args, os.Getenv("DOCKER_OPTS"))
	}

	log.Fatal(syscall.Exec("/usr/bin/dockerlaunch", args, os.Environ()))
}

func shellProxy(args []string) error {
	args = append([]string{"source-docker-conf"}, args...)
	env := os.Environ()
	env = append(env, "DOCKER_CONF_SOURCED=1")
	path, err := exec.LookPath("source-docker-conf")
	if err != nil {
		return err
	}
	return syscall.Exec(path, args, env)
}
