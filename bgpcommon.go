package main

import (
	"fmt"

	api "github.com/osrg/gobgp/v3/api"
)

var (
	ipv4UC = &api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_UNICAST,
	}
	ipv6UC = &api.Family{
		Afi:  api.Family_AFI_IP6,
		Safi: api.Family_SAFI_UNICAST,
	}
	ipv4VPN = &api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_MPLS_VPN,
	}
	ipv6VPN = &api.Family{
		Afi:  api.Family_AFI_IP6,
		Safi: api.Family_SAFI_MPLS_VPN,
	}
	ipv4MPLS = &api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_MPLS_LABEL,
	}
	ipv6MPLS = &api.Family{
		Afi:  api.Family_AFI_IP6,
		Safi: api.Family_SAFI_MPLS_LABEL,
	}
	evpn = &api.Family{
		Afi:  api.Family_AFI_L2VPN,
		Safi: api.Family_SAFI_EVPN,
	}
	l2vpnVPLS = &api.Family{
		Afi:  api.Family_AFI_L2VPN,
		Safi: api.Family_SAFI_VPLS,
	}
	ipv4Encap = &api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_ENCAPSULATION,
	}
	ipv6Encap = &api.Family{
		Afi:  api.Family_AFI_IP6,
		Safi: api.Family_SAFI_ENCAPSULATION,
	}
	rtc = &api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_ROUTE_TARGET_CONSTRAINTS,
	}
	ipv4Flowspec = &api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_FLOW_SPEC_UNICAST,
	}
	ipv6Flowspec = &api.Family{
		Afi:  api.Family_AFI_IP6,
		Safi: api.Family_SAFI_FLOW_SPEC_UNICAST,
	}
	ipv4VPNflowspec = &api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_FLOW_SPEC_VPN,
	}
	ipv6VPNflowspec = &api.Family{
		Afi:  api.Family_AFI_IP6,
		Safi: api.Family_SAFI_FLOW_SPEC_VPN,
	}
	l2VPNflowspec = &api.Family{
		Afi:  api.Family_AFI_L2VPN,
		Safi: api.Family_SAFI_FLOW_SPEC_VPN,
	}
	opaque = &api.Family{
		Afi:  api.Family_AFI_OPAQUE,
		Safi: api.Family_SAFI_KEY_VALUE,
	}
	ls = &api.Family{
		Afi:  api.Family_AFI_LS,
		Safi: api.Family_SAFI_LS,
	}
	ipv4MUP = &api.Family{
		Afi:  api.Family_AFI_IP,
		Safi: api.Family_SAFI_MUP,
	}
	ipv6MUP = &api.Family{
		Afi:  api.Family_AFI_IP6,
		Safi: api.Family_SAFI_MUP,
	}
)

func checkAddressFamily(def *api.Family) (*api.Family, error) {
	var f *api.Family
	var e error
	switch subOpts.AddressFamily {
	case "ipv4", "v4", "4":
		f = ipv4UC
	case "ipv6", "v6", "6":
		f = ipv6UC
	case "ipv4-l3vpn", "vpnv4", "vpn-ipv4":
		f = ipv4VPN
	case "ipv6-l3vpn", "vpnv6", "vpn-ipv6":
		f = ipv6VPN
	case "ipv4-labeled", "ipv4-labelled", "ipv4-mpls":
		f = ipv4MPLS
	case "ipv6-labeled", "ipv6-labelled", "ipv6-mpls":
		f = ipv6MPLS
	case "evpn":
		f = evpn
	case "l2vpn-vpls":
		f = l2vpnVPLS
	case "encap", "ipv4-encap":
		f = ipv4Encap
	case "ipv6-encap":
		f = ipv6Encap
	case "rtc":
		f = rtc
	case "ipv4-flowspec", "ipv4-flow", "flow4":
		f = ipv4Flowspec
	case "ipv6-flowspec", "ipv6-flow", "flow6":
		f = ipv6Flowspec
	case "ipv4-l3vpn-flowspec", "ipv4vpn-flowspec", "flowvpn4":
		f = ipv4VPNflowspec
	case "ipv6-l3vpn-flowspec", "ipv6vpn-flowspec", "flowvpn6":
		f = ipv6VPNflowspec
	case "l2vpn-flowspec":
		f = l2VPNflowspec
	case "opaque":
		f = opaque
	case "ls", "linkstate", "bgpls":
		f = ls
	case "ipv4-mup", "mup-ipv4", "mup4":
		f = ipv4MUP
	case "ipv6-mup", "mup-ipv6", "mup6":
		f = ipv6MUP
	case "":
		f = def
	default:
		e = fmt.Errorf("unsupported address family: %s", subOpts.AddressFamily)
	}
	return f, e
}
