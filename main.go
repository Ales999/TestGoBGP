package main

import "fmt"

func main() {
	//printNeibr("192.168.1.11:50051", "172.24.1.1")

	//ShowRibPathByIp("172.24.1.1", "101.0.0.9")
	// 104.0.0.1
	asPath, err := ShowRibPathByIp("192.168.1.11:50051", "172.24.1.1", "104.0.0.1")
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
