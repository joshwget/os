package switchconsoles

import (
	"os"
	"os/exec"

	"github.com/docker/engine-api/types"
	"github.com/rancher/os/docker"
	"golang.org/x/net/context"
)

func Main() {
	client, err := docker.NewSystemClient()
	if err != nil {
		panic(err)
	}

	/*err = client.ContainerRestart(context.Background(), "console", 10)
	if err != nil {
		panic(err)
	}*/

	err = client.ContainerStop(context.Background(), "console", 10)
	if err != nil {
		panic(err)
	}

	err = client.ContainerRemove(context.Background(), "console", types.ContainerRemoveOptions{})
	if err != nil {
		panic(err)
	}

	err = client.ContainerStart(context.Background(), "console2")
	if err != nil {
		panic(err)
	}

	return

	cmd := exec.Command("ros", "service", "up", "debian-console")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
