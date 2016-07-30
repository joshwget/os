package integration

import . "gopkg.in/check.v1"

func (s *QemuSuite) TestCloudConfigMounts(c *C) {
	s.RunQemu(c, "--cloud-config", "./tests/assets/test_16/cloud-config.yml")

	s.CheckCall(c, "cat /home/rancher/test | grep test")
}
