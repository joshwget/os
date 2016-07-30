package integration

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *QemuSuite) TestRebootWithContainerRunning(c *C) {
	s.RunQemu(c, "--cloud-config", "./tests/assets/test_03/cloud-config.yml")

	s.CheckCall(c, fmt.Sprintf(`
docker run -d --restart=always %s`, NginxImage))

	s.Reboot(c)

	s.CheckCall(c, "docker ps -f status=running | grep nginx")
}
