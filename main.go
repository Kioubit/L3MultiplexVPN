package main

import (
	"L3MultiplexVPN/interfaces"
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: <PeerIP> <Local network id list>")
		fmt.Println("Example: 192.168.5.5 0,2")
		os.Exit(2)
	}
	addr, err := netip.ParseAddr(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	final := make([]uint8, 0)
	s := strings.Split(os.Args[2], ",")
	for _, id := range s {
		intID, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		final = append(final, uint8(intID))
	}

	interfaces.Startup(final, addr)
}
