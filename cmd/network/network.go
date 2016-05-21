package network

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/docker/libnetwork/resolvconf"
	"github.com/rancher/netconf"
	"github.com/rancher/os/config"
	"github.com/rancher/os/docker"
)

func Main() {
	log.Infoln("Running network")

	client, err := docker.NewSystemClient()
	if err != nil {
		log.Error(err)
	}

	err = client.ContainerRestart(context.Background(), "dhcp", 10)
	if err != nil {
		log.Error(err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	nameservers := cfg.Rancher.Network.Dns.Nameservers
	search := cfg.Rancher.Network.Dns.Search
	if len(nameservers) > 0 || len(search) > 0 {
		if _, err := resolvconf.Build("/etc/resolv.conf", nameservers, search, nil); err != nil {
			log.Error(err)
		}
	}

	if err := netconf.ApplyNetworkConfigs(&cfg.Rancher.Network); err != nil {
		log.Error(err)
	}

	if err := netconf.RunDhcp(&cfg.Rancher.Network); err != nil {
		log.Error(err)
	}

	if err := netconf.WaitForIps(&cfg.Rancher.Network); err != nil {
		log.Error(err)
	}
}
