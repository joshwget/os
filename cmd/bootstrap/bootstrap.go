package cloudinitexecute

import (
	"os"
	"os/exec"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/os/config"
	"github.com/rancher/os/util"
)

func Main() {
	log.Info("1")
	if err := udevSettle(); err != nil {
		panic(err)
	}

	log.Info("2")
	cfg := config.LoadConfig()
	log.Info("3")

	if cfg.Rancher.State.Dev != "" && cfg.Rancher.State.Wait {
		log.Info("4")
		waitForRoot(cfg)
		log.Info("5")
	}

	log.Info("6")

	if cfg.Rancher.State.MdadmScan {
		cmd := exec.Command("mdadm", "--assemble", "--scan")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	log.Info("7")
	stateScript := cfg.Rancher.State.Script
	if stateScript != "" {
		// TODO stateScript
	}

	log.Info("8")

	// TODO autoformat

	if err := udevSettle(); err != nil {
		panic(err)
	}

	log.Info("9")
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
