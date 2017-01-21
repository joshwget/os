package athena

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	athenaDir     = "/usr/bin/athena"
	dockerDir     = "/usr/bin/docker"
	dockerRuncDir = "/usr/bin/docker-runc"
)

func Main() {
	fmt.Println("athena")

	/*if err := os.Remove(dockerDir); err != nil {
		//panic(err)
		_ = err
	}
	if err := os.Symlink(athenaDir, dockerDir); err != nil {
		panic(err)
	}
	if err := os.Remove(dockerRuncDir); err != nil {
		//panic(err)
		_ = err
	}
	if err := os.Symlink(athenaDir, dockerRuncDir); err != nil {
		panic(err)
	}*/

	cmd := exec.Command(dockerDir, "daemon")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	cmd = exec.Command(dockerDir, strings.Split("run -d --privileged -v /var/run/docker.sock:/var/run/docker.sock -v /var/lib/rancher:/var/lib/rancher rancher/agent:v1.1.3 http://cb9e163b.ngrok.io/v1/scripts/B3ED45CA99884D1645CE:1484985600000:TrlrDH7mFxfCtE5ote4wcQHgCkk", " ")...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	select {}
}
