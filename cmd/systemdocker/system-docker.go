package systemdocker

import (
	"os"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/os/config"
)

func Main() {
	var newEnv []string
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "DOCKER_HOST=") {
			newEnv = append(newEnv, env)
		}
	}

	newEnv = append(newEnv, "DOCKER_HOST="+config.DOCKER_SYSTEM_HOST)

	if os.Geteuid() != 0 {
		log.Fatalf("%s: Need to be root", os.Args[0])
	}

	os.Args[0] = config.DOCKER_DIST_BIN
	if err := syscall.Exec(os.Args[0], os.Args, newEnv); err != nil {
		log.Fatal(err)
	}
}
