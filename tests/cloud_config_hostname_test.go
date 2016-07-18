package integration

import . "gopkg.in/check.v1"

func (s *QemuSuite) TestUp(c *C) {
	err := s.RunQemu("--cloud-config", "./tests/assets/test_13/cloud-config.yml")
	c.Assert(err, IsNil)

	s.CheckCall(c, "hostname | grep rancher-test")
	s.CheckCall(c, "cat /etc/hosts | grep rancher-test")
}
