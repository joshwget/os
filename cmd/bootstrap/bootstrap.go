package cloudinitexecute

import (
	"os"
	"os/exec"
	"time"

	"github.com/rancher/os/config"
	"github.com/rancher/os/util"
)

func Main() {
	if err := udevSettle(); err != nil {
		panic(err)
	}

	cfg := config.LoadConfig()
	waitForRoot(cfg)

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

	// TODO autoformat

	if err := udevSettle(); err != nil {
		panic(err)
	}
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
