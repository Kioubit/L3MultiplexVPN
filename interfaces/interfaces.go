package interfaces

import (
	"log"
	"net"
	"net/netip"
	"strconv"
)
import "github.com/songgao/water"

var interfaceTable [255]*water.Interface
var PeerIP []byte

func Startup(l3idList []uint8, peerIP netip.Addr) {
	if peerIP.Is4() {
		t := peerIP.As4()
		PeerIP = t[:]
	} else {
		t := peerIP.As16()
		PeerIP = t[:]
	}

	outPacketChan := make(chan *[]byte, 10)
	inPacketChan := make(chan *[]byte, 10)
	for i := 0; i < len(l3idList)*4; i++ {
		go txProcessor(outPacketChan)
	}
	for i := 0; i < len(l3idList)*4; i++ {
		go rxProcessor(inPacketChan)
	}
	for _, u := range l3idList {
		go createInterface("layerd"+strconv.Itoa(int(u)), u, outPacketChan)
	}
	rxListener(inPacketChan)

}

func createInterface(interfaceName string, l3id uint8, packetChan chan *[]byte) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = interfaceName
	iFace, err := water.New(config)
	if err != nil {
		log.Fatal("Is the 'tun' device available? Failed creating TUN interface ", interfaceName, " - ", err)
	}
	interfaceTable[l3id] = iFace
	for {
		packet := make([]byte, 2000)
		packet[0] = l3id
		n, err := iFace.Read(packet[1:])
		if err != nil {
			log.Fatal("Failed reading from tun interface ", err)
		}
		packet = packet[:n+1]
		packetChan <- &packet
	}
}

func txProcessor(packetChan chan *[]byte) {
	for {
		pkt := <-packetChan
		destinationIP := &net.UDPAddr{
			IP:   net.IP(PeerIP[:]),
			Port: 3643,
		}
		conn, err := net.DialUDP("udp", nil, destinationIP)
		if err != nil {
			continue
		}
		_, _ = conn.Write(*pkt)
	}
}

func rxListener(packetChan chan *[]byte) {
	addr := net.UDPAddr{
		IP:   nil,
		Port: 3643,
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	for {
		packetBuf := make([]byte, 2000)
		readLen, _, err := ser.ReadFromUDP(packetBuf)
		if err != nil {
			continue
		}
		raw := packetBuf[:readLen]
		packetChan <- &raw
	}
}

func rxProcessor(packetChan chan *[]byte) {
	for {
		pkt := <-packetChan
		packet := (*pkt)[1:]
		l3id := (*pkt)[0]
		dest := interfaceTable[l3id]
		if dest == nil {
			continue
		}
		_, _ = dest.Write(packet)
	}
}
