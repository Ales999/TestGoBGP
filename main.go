package main

import (
	"flag"
	"fmt"
)

func main() {

	flagServerApi := flag.String("h", "127.0.0.1:50051", "GoBgp host:port, Example: 192.168.1.10:50051")
	flagNeigbrIp := flag.String("n", "", "Neigbror, Example: 172.24.1.1")
	flagTarget := flag.String("t", "104.0.0.1", "find target, Example: 104.0.0.1")

	flag.Parse()

	//printNeibr("192.168.1.11:50051", "172.24.1.1")

	//ShowRibPathByIp("172.24.1.1", "101.0.0.9")
	// 104.0.0.1
	asPath, err := ShowRibPathByIp(*flagServerApi, *flagNeigbrIp, *flagTarget)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println(asPath)
}

// Example Output:
//  .\TestGoBGP.exe
//  attrs: [{Origin: i} {AsPath: 64503 64502 64500 64504} {Nexthop: 172.24.1.1}]
//  64503 64502 64500 64504
