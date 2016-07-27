package control

import (
	"fmt"
	"sort"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	composeConfig "github.com/docker/libcompose/config"
	"github.com/docker/libcompose/project/options"
	"github.com/rancher/os/compose"
	"github.com/rancher/os/config"
	"github.com/rancher/os/util/network"
)

func dockerSubcommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "switch",
			Usage:  "switch console without a reboot",
			Action: dockerSwitch,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "do not prompt for input",
				},
				cli.BoolFlag{
					Name:  "no-pull",
					Usage: "don't pull console image",
				},
			},
		},
		{
			Name:   "list",
			Usage:  "list available consoles",
			Action: dockerList,
		},
	}
}

func dockerSwitch(c *cli.Context) error {
	if len(c.Args()) != 1 {
		log.Fatal("Must specify exactly one Docker to switch to")
	}
	newDocker := c.Args()[0]

	cfg := config.LoadConfig()

	project, err := compose.GetProject(cfg, true, false)
	if err != nil {
		log.Fatal(err)
	}

	if err = project.Stop(context.Background(), 10, "docker"); err != nil {
		log.Fatal(err)
	}

	if err = project.Down(context.Background(), options.Down{}, "docker"); err != nil {
		log.Fatal(err)
	}

	if err = project.Delete(context.Background(), options.Delete{}, "docker"); err != nil {
		log.Fatal(err)
	}

	// TODO: this is crap
	project.ServiceConfigs.Add("docker", &composeConfig.ServiceConfig{})

	if err = compose.LoadService(project, cfg, true, newDocker); err != nil {
		log.Fatal(err)
	}

	fmt.Println(project.ServiceConfigs.Get("docker"))

	if err = project.Up(context.Background(), options.Up{}, "docker"); err != nil {
		log.Fatal(err)
	}

	return nil
}

func dockerList(c *cli.Context) error {
	cfg := config.LoadConfig()

	dockers, err := network.GetDockers(cfg.Rancher.Repositories.ToArray())
	if err != nil {
		return err
	}
	sort.Strings(dockers)

	for _, docker := range dockers {
		fmt.Println(docker)
	}

	return nil
}
