package integration

import . "gopkg.in/check.v1"

func (s *QemuSuite) TestDhcpHostname(c *C) {
	s.RunQemu(c)

	s.CheckCall(c, "hostname | grep rancher-dev")
	s.CheckCall(c, "cat /etc/hosts | grep rancher-dev")
}
