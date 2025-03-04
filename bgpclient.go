package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sort"
	"strings"

	api "github.com/osrg/gobgp/v3/api"
	"github.com/osrg/gobgp/v3/pkg/apiutil"

	"github.com/osrg/gobgp/v3/pkg/packet/bgp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client api.GobgpApiClient
	ctx    context.Context
)

var subOpts struct {
	AddressFamily string `short:"a" long:"address-family" description:"specifying an address family"`
	BatchSize     uint64 `short:"b" long:"batch-size" description:"Size of the temporary buffer in the server memory. Zero is unlimited (default)"`
}

// addr2AddressFamily - получить тип Family заполненный
//
// Example use:
//
//	     def := addr2AddressFamily(net.ParseIP(name))
//		family, err := checkAddressFamily(def)
//		if err != nil {
//			return err
//		}
func addr2AddressFamily(a net.IP) *api.Family {
	if a.To4() != nil {
		return &api.Family{
			Afi:  api.Family_AFI_IP,
			Safi: api.Family_SAFI_UNICAST,
		}
	} else if a.To16() != nil {
		return &api.Family{
			Afi:  api.Family_AFI_IP6,
			Safi: api.Family_SAFI_UNICAST,
		}
	}
	return nil
}

func getNeighbors(address string, enableAdv bool) ([]*api.Peer, error) {
	stream, err := client.ListPeer(ctx, &api.ListPeerRequest{
		Address:          address,
		EnableAdvertised: enableAdv,
	})
	if err != nil {
		return nil, err
	}

	l := make([]*api.Peer, 0, 1024)
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		l = append(l, r.Peer)
	}
	if address != "" && len(l) == 0 {
		return l, fmt.Errorf("not found neighbor %s", address)
	}
	return l, err
}

func getNeigboorIPs(ctx context.Context, srvAddress string) ([]string, error) {

	conn, err := grpc.NewClient(srvAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return []string{}, fmt.Errorf("did not connect: %s", err)
	}
	defer conn.Close()

	client := api.NewGobgpApiClient(conn)
	query := &api.ListPeerRequest{}

	// Prepare
	stream, err := client.ListPeer(ctx, query)
	if err != nil {
		return []string{}, fmt.Errorf("error when calling client.ListPeer: %s", err)
	}

	done := make(chan bool)
	var neigbrips []string
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true //means stream is finished
				return
			}
			if err != nil {
				log.Fatalf("cannot receive %v", err)
				done <- true //means stream is finished
				return
			}
			// Сохраним в переменную программы в виде массива, ибо неигборов может быть много.
			neigbrips = append(neigbrips, resp.Peer.Conf.NeighborAddress)
		}
	}()

	<-done //we will wait until all response is received

	return neigbrips, nil
}

// PrintNeibror - get neiboor
// Examle use: PrintNeibror("192.168.1.11:50051","172.24.1.1")
func PrintNeibror(ctx context.Context, srvAddress string, neibrIp string) {

	//var conn *grpc.ClientConn
	conn, err := grpc.NewClient(srvAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := api.NewGobgpApiClient(conn)

	var query *api.ListPeerRequest

	if len(neibrIp) > 0 {
		query = &api.ListPeerRequest{Address: neibrIp} // extended, by IP
	} else {
		query = &api.ListPeerRequest{} // simple - get all peeer's

	}

	stream, err := client.ListPeer(ctx, query)
	if err != nil {
		log.Fatalf("Error when calling client.ListPolicy: %s", err)
	}

	// https://www.freecodecamp.org/news/grpc-server-side-streaming-with-go/
	done := make(chan bool)

	var neip []string
	go func() {
		var _ipneigbr string
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true //means stream is finished
				return
			}
			if err != nil {
				log.Fatalf("cannot receive %v", err)
			}
			log.Printf("Respond received: <--\n%s\n-->", resp.String())

			log.Printf("Local ASn: %d\n", resp.Peer.Conf.LocalAsn)
			log.Printf("Peer ASn: %d\n", resp.Peer.Conf.PeerAsn)

			_ipneigbr = resp.Peer.Conf.NeighborAddress
			log.Printf("IP: %s", _ipneigbr)

			// Сохраним в переменную программы в виде массива, ибо неигборов может быть много.
			neip = append(neip, _ipneigbr)
		}
	}()

	<-done //we will wait until all response is received
	log.Printf("finished")

	if len(neip) > 0 {
		fmt.Println("Neigbors IP:", neip)
	}

}

func parseCIDRorIP(str string) (net.IP, *net.IPNet, error) {
	ip, n, err := net.ParseCIDR(str)
	if err == nil {
		return ip, n, nil
	}
	ip = net.ParseIP(str)
	if ip == nil {
		return ip, nil, fmt.Errorf("invalid CIDR/IP")
	}
	return ip, nil, nil
}

func ShowRibPathByIp(serverApi string, neigbrIp string, target string) (string, error) {

	ctx = context.Background()

	var r string = "global"
	var def *api.Family
	var family *api.Family
	var err error
	var parsedNeigbrIp net.IP
	var neigbrIps []string

	//if len(neigbrIp) == 0 {
	//	PrintNeibror(ctx, serverApi, "")
	//}

	// Если указали IP пира, то используем только его
	if len(neigbrIp) > 0 {
		neigbrIps = append(neigbrIps, neigbrIp)
	} else {
		// Иначе запросим список у самого GoBGP
		neigbrIps, err = getNeigboorIPs(ctx, serverApi)
		if err != nil {
			return "", err
		}
	}

	for _, ni := range neigbrIps {

		//PrintNeibror(ctx, serverApi, "")
		parsedNeigbrIp = net.ParseIP(ni)
		if parsedNeigbrIp == nil {
			log.Fatalf("Error parsing neigbor ip\n")
		}

		def = addr2AddressFamily(parsedNeigbrIp)

		family, err = checkAddressFamily(def)
		if err != nil {
			return "", err
		}

		// Parse target for IP or CIDR
		if _, _, err = parseCIDRorIP(target); err != nil {
			return "", err
		}

		var (
			option         api.TableLookupPrefix_Type
			tableType      api.TableType
			enableFiltered bool
			conn           *grpc.ClientConn
		)

		rd := ""
		tableType = api.TableType_ADJ_IN

		filter := []*api.TableLookupPrefix{{
			Prefix: target, // "101.0.0.9"
			Rd:     rd,
			Type:   option,
		},
		}

		// Создаем клиента к API GoBGP - сначала коннекцию
		conn, err = grpc.NewClient(serverApi, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()

		// Теперь самого клиента.
		client = api.NewGobgpApiClient(conn)

		// Получаем поток от GoBGP в виде stream
		stream, err := client.ListPath(ctx, &api.ListPathRequest{
			TableType:      tableType,
			Family:         family,
			Name:           ni,
			Prefixes:       filter,
			SortType:       api.ListPathRequest_PREFIX,
			EnableFiltered: enableFiltered,
			BatchSize:      subOpts.BatchSize,
		})
		if err != nil {
			return "", err
		}

		rib := make([]*api.Destination, 0)
		for {
			r, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				return "", err
			}
			rib = append(rib, r.Destination)
		}
		if len(rib) == 0 {
			l, err := getNeighbors(ni, false)
			if err != nil {
				return "", err
			}
			if l[0].State.SessionState != api.PeerState_ESTABLISHED {
				return "", fmt.Errorf("neighbor %v's BGP session is not established", neigbrIp)
			}
		}

		routeFamily := apiutil.ToRouteFamily(family)

		// show RIB
		var dsts []*api.Destination
		// Перебираем
		switch routeFamily {
		case bgp.RF_IPv4_UC, bgp.RF_IPv6_UC:
			type d struct {
				prefix net.IP
				dst    *api.Destination
			}
			l := make([]*d, 0, len(rib))
			for _, dst := range rib {
				prefix := dst.Prefix
				if tableType == api.TableType_VRF {
					// extract prefix from original which is RD(AS:VRF):IPv4 or IPv6 address
					s := strings.SplitN(prefix, ":", 3)
					prefix = s[len(s)-1]
				}
				_, p, _ := net.ParseCIDR(prefix)
				l = append(l, &d{prefix: p.IP, dst: dst})
			}

			sort.Slice(l, func(i, j int) bool {
				return bytes.Compare(l[i].prefix, l[j].prefix) < 0
			})

			dsts = make([]*api.Destination, 0, len(rib))
			for _, s := range l {
				dsts = append(dsts, s.dst)
			}
		default:
			dsts = append(dsts, rib...)
		}

		for _, d := range dsts {
			if enableFiltered {
				showFiltered := r == cmdRejected
				l := make([]*api.Path, 0, len(d.Paths))
				for _, p := range d.GetPaths() {
					if p.Filtered == showFiltered {
						l = append(l, p)
					}
				}
				d.Paths = l
			}
		}
		//var ret string
		if len(dsts) > 0 {
			// Print result:
			//showRoute(dsts, showAge, showBest, showLabel, showMUP, showSendMaxFiltered, showIdentifier)
			//var tpref = showTargetPrefix(dsts)
			// Принт префикс
			for _, dpref := range dsts {
				fmt.Printf("Prefix: %s ", dpref.GetPrefix())
				//fmt.Println("Patch", dpref.GetPaths()[0].)
				//appa := dpref.Paths
				for _, p := range dpref.Paths {
					attrs, _ := apiutil.GetNativePathAttributes(p)
					// Next Hop
					nexthop := "fictitious"
					if n := getNextHopFromPathAttributes(attrs); n != nil {
						nexthop = n.String()
					}

					fmt.Printf("NeigBorIP: %s, Next Hop: %s ", p.NeighborIp, nexthop)
				}
			}

			fmt.Printf("Path: --> %s\n", showAsRoute(dsts))

		} else {
			return "", fmt.Errorf("network not in table")
			//fmt.Println("Network not in table")
		}

		/*
			// TODO: Output as JSON
				if globalOpts.Json {
					d := make(map[string]*apiutil.Destination)
					for _, dst := range rib {
						d[dst.Prefix] = apiutil.NewDestination(dst)
					}
					j, _ := json.Marshal(d)
					fmt.Println(string(j))
					return nil
				}
		*/
	}

	return "", nil
}

func showAsRoute(dsts []*api.Destination) string {
	var attrs []bgp.PathAttributeInterface
	//pathStrs := make([][]interface{}, 0, len(dsts))
	//now := time.Now()
	for _, dst := range dsts {
		for _, p := range dst.Paths {
			//pathStrs = append(pathStrs, makeShowRouteArgs(p, idx, now, showAge, showBest, showLabel, showMUP, showSendMaxFiltered, showIdentifier))
			//fmt.Printf("p: %v\n", p)

			atti, _ := apiutil.GetNativePathAttributes(p)
			attrs = append(attrs, atti...)
			// Test data print
			//fmt.Printf("attrs: %v\n", attrs)
		}
	}
	// AS_PATH
	aspathstr := func() string {
		for _, attr := range attrs {
			switch a := attr.(type) {
			case *bgp.PathAttributeAsPath:
				return bgp.AsPathString(a)
			}
		}
		return ""
	}()
	// debug out
	//fmt.Println(aspathstr)

	return aspathstr
}

func getNextHopFromPathAttributes(attrs []bgp.PathAttributeInterface) net.IP {
	for _, attr := range attrs {
		switch a := attr.(type) {
		case *bgp.PathAttributeNextHop:
			return a.Value
		case *bgp.PathAttributeMpReachNLRI:
			return a.Nexthop
		}
	}
	return nil
}
