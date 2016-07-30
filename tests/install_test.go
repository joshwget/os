package integration

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (s *QemuSuite) TestInstall(c *C) {
	s.RunQemu(c, "--no-format")

	s.LoadInstallerImage(c)

	s.CheckCall(c, fmt.Sprintf(`
sudo mkfs.ext4 /dev/vda
sudo ros install -f --no-reboot -d /dev/vda -i rancher/os:%s%s`, Version, Suffix))
}
