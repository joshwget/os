package integration

import . "gopkg.in/check.v1"

func (s *QemuSuite) TestUp(c *C) {
	err := s.RunQemu("--cloud-config", "./tests/integration/assets/test_03/cloud-config.yml")
	c.Assert(err, IsNil)

	s.CheckCall("apt-get --version")
}
