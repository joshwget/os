package condition

import (
	"os"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	conditionDir = "/var/lib/rancher/conditions"
)

func WaitFor(condition string) error {
	for {
		if _, err := os.Stat(path.Join(conditionDir, condition)); err == nil {
			return nil
		} else if os.IsNotExist(err) {
			log.Infof("Waiting on condition %s", condition)
			time.Sleep(time.Second)
		} else if err != nil {
			return err
		}
	}
}

func Release(condition string) error {
	f, err := os.Create(path.Join(conditionDir, condition))
	if err != nil {
		return err
	}
	return f.Close()
}
