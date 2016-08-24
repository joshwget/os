package console

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/os/config"
	"github.com/rancher/os/util"
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

	respawnConf.WriteString("/usr/sbin/sshd -D")

	return respawnConf.Bytes()
}

func modifySshdConfig() error {
	sshdConfig, err := ioutil.ReadFile("/etc/ssh/sshd_config")
	if err != nil {
		return err
	}
	sshdConfigString := string(sshdConfig)

	for _, item := range []string{
		"UseDNS no",
		"PermitRootLogin no",
		"ServerKeyBits 2048",
		"AllowGroups docker",
	} {
		match, err := regexp.Match("^"+item, sshdConfig)
		if err != nil {
			return err
		}
		if !match {
			sshdConfigString += fmt.Sprintf("%s\n", item)
		}
	}

	return ioutil.WriteFile("/etc/ssh/sshd_config", []byte(sshdConfigString), 0644)
}

func writeOsRelease() error {
	idLike := "busybox"
	if osRelease, err := ioutil.ReadFile("/etc/os-release"); err == nil {
		for _, line := range strings.Split(string(osRelease), "\n") {
			if strings.HasPrefix(line, "ID_LIKE") {
				split := strings.Split(line, "ID_LIKE")
				if len(split) > 1 {
					idLike = split[1]
				}
			}
		}
	}

	return ioutil.WriteFile("/etc/os-release", []byte(fmt.Sprintf(`
NAME="RancherOS"
VERSION=%s
ID=rancheros
ID_LIKE=%s
VERSION_ID=%s
PRETTY_NAME="RancherOS %s"
HOME_URL=
SUPPORT_URL=
BUG_REPORT_URL=
BUILD_ID=
`, config.VERSION, idLike, config.VERSION, config.VERSION)), 0644)
}

func setupSSH(cfg *config.CloudConfig) error {
	for _, keyType := range []string{"rsa", "dsa", "ecdsa", "ed25519"} {
		outputFile := fmt.Sprintf("/etc/ssh/ssh_host_%s_key", keyType)
		outputFilePub := fmt.Sprintf("/etc/ssh/ssh_host_%s_key.pub", keyType)

		if _, err := os.Stat(outputFile); err == nil {
			continue
		}

		saved, savedExists := cfg.Rancher.Ssh.Keys[keyType]
		pub, pubExists := cfg.Rancher.Ssh.Keys[keyType+"-pub"]

		if savedExists && pubExists {
			// TODO check permissions
			if err := util.WriteFileAtomic(outputFile, []byte(saved), 0600); err != nil {
				return err
			}
			if err := util.WriteFileAtomic(outputFilePub, []byte(pub), 0600); err != nil {
				return err
			}
			continue
		}

		cmd := exec.Command("bash", "-c", fmt.Sprintf("ssh-keygen -f %s -N '' -t %s", outputFile, keyType))
		if err := cmd.Run(); err != nil {
			return err
		}

		savedBytes, err := ioutil.ReadFile(outputFile)
		if err != nil {
			return err
		}

		pubBytes, err := ioutil.ReadFile(outputFilePub)
		if err != nil {
			return err
		}

		config.Set(fmt.Sprintf("rancher.ssh.keys.%s", keyType), string(savedBytes))
		config.Set(fmt.Sprintf("rancher.ssh.keys.%s-pub", keyType), string(pubBytes))
	}

	return os.MkdirAll("/var/run/sshd", 0644)
}
