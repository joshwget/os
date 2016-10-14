package control

import (
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/codegangsta/cli"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/os/config"
	"github.com/rancher/os/util"
)

func bootstrapAction(c *cli.Context) error {
	if err := udevSettle(); err != nil {
		panic(err)
	}

	cfg := config.LoadConfig()

	if cfg.Rancher.State.Dev != "" && cfg.Rancher.State.Wait {
		waitForRoot(cfg)
	}

	if cfg.Rancher.State.MdadmScan {
		cmd := exec.Command("mdadm", "--assemble", "--scan")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	stateScript := cfg.Rancher.State.Script
	if stateScript != "" {
		// TODO stateScript
	}

	log.Info("1")

	// TODO autoformat
	cmd := exec.Command("/usr/sbin/auto-format2.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{
		"AUTOFORMAT=" + strings.Join(cfg.Rancher.State.Autoformat, " "),
	}
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := udevSettle(); err != nil {
		return err
	}

	return nil
}

func udevSettle() error {
	cmd := exec.Command("udevd", "--daemon")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("udevadm", "trigger", "--action=add")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmd = exec.Command("udevadm", "settle")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	return nil
}

func waitForRoot(cfg *config.CloudConfig) {
	var dev string
	for i := 0; i < 30; i++ {
		dev = util.ResolveDevice(cfg.Rancher.State.Dev)
		if dev != "" {
			break
		}
		time.Sleep(time.Millisecond * 1000)
	}
	if dev == "" {
		return
	}
	for i := 0; i < 30; i++ {
		if _, err := os.Stat(dev); err == nil {
			break
		}
		time.Sleep(time.Millisecond * 1000)
	}
}
