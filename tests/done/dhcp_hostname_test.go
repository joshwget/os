package integration

import . "gopkg.in/check.v1"

func (s *QemuSuite) TestDhcpHostname(c *C) {
	err := s.RunQemu("--cloud-config", "./tests/assets/test_12/cloud-config.yml")
	c.Assert(err, IsNil)

	s.CheckCall(c, "hostname | grep rancher-dev")
	s.CheckCall(c, "cat /etc/hosts | grep rancher-dev")
}
