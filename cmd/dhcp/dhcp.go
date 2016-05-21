package dhcp

import (
	"io/ioutil"
	"os"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/os/config"
	"github.com/rancher/os/hostname"
)

var (
	defaultDhcpArgs = []string{"-MA4b", "eth0"}
)

func Main() {
	log.Infoln("Running dhcp")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	nameservers := cfg.Rancher.DefaultNetwork.Dns.Nameservers
	search := cfg.Rancher.DefaultNetwork.Dns.Search
	/*if _, err := resolvconf.Build("/etc/resolv.conf", nameservers, search, nil); err != nil {
		log.Error(err)
	}*/

	nameservers = cfg.Rancher.Network.Dns.Nameservers
	search = cfg.Rancher.Network.Dns.Search
	userSetDns := len(nameservers) > 0 || len(search) > 0
	userSetHostname := cfg.Hostname != ""

	/*if userSetHostname {
		if err := hostname.SetHostnameFromCloudConfig(cfg); err != nil {
			log.Error(err)
		}
	}*/

	args := defaultDhcpArgs

	if userSetHostname {
		args = append(args, "--nohook", "hostname")
	}

	if userSetDns {
		args = append(args, "--nohook", "resolv.conf")
	}

	conf := []byte(`
# A sample configuration for dhcpcd.
# See dhcpcd.conf(5) for details.

# Allow users of this group to interact with dhcpcd via the control socket.
#controlgroup wheel

# Inform the DHCP server of our hostname for DDNS.
hostname

# Use the hardware address of the interface for the Client ID.
#clientid
# or
# Use the same DUID + IAID as set in DHCPv6 for DHCPv4 ClientID as per RFC4361.
# Some non-RFC compliant DHCP servers do not reply with this set.
# In this case, comment out duid and enable clientid above.
duid

# Persist interface configuration when dhcpcd exits.
persistent

# Rapid commit support.
# Safe to enable by default because it requires the equivalent option set
# on the server to actually work.
option rapid_commit

# A list of options to request from the DHCP server.
option domain_name_servers, domain_name, domain_search, host_name
option classless_static_routes
# Most distributions have NTP support.
option ntp_servers
# Respect the network MTU. This is applied to DHCP routes.
option interface_mtu

# A ServerID is required by RFC2131.
require dhcp_server_identifier

# Generate Stable Private IPv6 Addresses instead of hardware based ones
slaac private
env force_hostname=true
	`)
	err = ioutil.WriteFile("/etc/dhcpcd.conf", conf, 0644)
	if err != nil {
		log.Error(err)
	}

	cmd := exec.Command("dhcpcd", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Error(err)
	}

	if err := hostname.SyncHostname(); err != nil {
		log.Error(err)
	}

	select {}
}
