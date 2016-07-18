package integration

import (
        "fmt"

        . "gopkg.in/check.v1"
)

func (s *QemuSuite) TestInstall(c *C) {
	err := s.RunQemu("--no-format")
	c.Assert(err, IsNil)

	s.CheckHostCall(c, fmt.Sprintf(`
sudo sh -c "docker save rancher/os:%s%s | system-docker load"`, Version, Suffix))

	s.CheckCall(c, fmt.Sprintf(`
sudo mkfs.ext4 /dev/vda
sudo ros install -f --no-reboot -d /dev/vda -i rancher/os:%s%s`, Version, Suffix))
}
