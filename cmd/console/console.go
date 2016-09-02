package console

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/os/config"
)

const (
	consoleDone = "/run/console-done"
	dockerHome  = "/home/docker"
	gettyCmd    = "/sbin/agetty"
	rancherHome = "/home/rancher"
	startScript = "/opt/rancher/bin/start.sh"
)

type symlink struct {
	oldname, newname string
}

func Main() {
	password := config.GetCmdline("rancher.password")
	cmd := exec.Command("chpasswd")
	cmd.Stdin = strings.NewReader(fmt.Sprint("rancher:", password))
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	cmdline, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		log.Error(err)
	}

	respawnConf := generateRespawnConf(string(cmdline))

	if err = ioutil.WriteFile("/etc/respawn.conf", respawnConf, 0644); err != nil {
		log.Error(err)
	}

	os.Setenv("TERM", "linux")

	respawnBinPath, err := exec.LookPath("respawn")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(syscall.Exec(respawnBinPath, []string{"respawn", "-f", "/etc/respawn.conf"}, os.Environ()))
}

func generateRespawnConf(cmdline string) []byte {
	autologin := strings.Contains(cmdline, "rancher.autologin")

	var respawnConf bytes.Buffer

	for i := 1; i < 7; i++ {
		respawnConf.WriteString(gettyCmd)
		if autologin {
			respawnConf.WriteString(" --autologin rancher")
		}
		respawnConf.WriteString(fmt.Sprintf(" 115200 tty%d\n", i))
	}

	for _, tty := range []string{"ttyS0", "ttyS1", "ttyS2", "ttyS3", "ttyAMA0"} {
		if !strings.Contains(cmdline, fmt.Sprintf("console=%s", tty)) {
			continue
		}

		respawnConf.WriteString(gettyCmd)
		if autologin {
			respawnConf.WriteString(" --autologin rancher")
		}
		respawnConf.WriteString(fmt.Sprintf(" 115200 %s\n", tty))
	}

	return respawnConf.Bytes()
}
