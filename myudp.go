package main

import (
	"log"
	"net"
)

func udpLocal2(server, target string) {
	// listen to incoming udp packets
	pc, err := net.ListenPacket("udp", ":13532")
	if err != nil {
		log.Fatal(err)
	}
	//simple read
	buffer := make([]byte, 1024)
	for {
		n, srcAddr, err := pc.ReadFrom(buffer)
		if err != nil {
			log.Println("udp ReadFrom:" + err.Error())
		}
		go udpHandler(pc, n, srcAddr, buffer)

	}

}

func udpHandler(pc net.PacketConn, n int, srcAddr net.Addr, buf []byte) {
	remoAddr, err := net.ResolveUDPAddr("udp4", "52.77.232.211:13532")
	if err != nil {

	}
	remoCon, err := net.DialUDP("udp", nil, remoAddr)
	defer remoCon.Close()
	remoCon.Write(buf[len(srcAddr.String()):n])
	log.Println("srcAddr:" + srcAddr.String())
}
