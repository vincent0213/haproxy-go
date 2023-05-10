// Implementation of a UDP proxy

package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Information maintained for each client/server connection
type UDPCap struct {
	ClientAddr *net.UDPAddr // Address of the client
	ServerConn *net.UDPConn // UDP connection to server
}

type UDPCapsLock struct {
	UDPCaps map[string]*UDPCap
	Lock    sync.Mutex
}

// Generate a new connection by opening a UDP connection to the server
func NewConnection(srvAddr, cliAddr *net.UDPAddr) *UDPCap {
	conn := new(UDPCap)
	conn.ClientAddr = cliAddr
	srvudp, err := net.DialUDP("udp", nil, srvAddr)
	if checkreport(1, err) {
		return nil
	}
	conn.ServerConn = srvudp
	return conn
}

func setup(hostport string, port int) {
	// Set up Proxy
	saddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if checkreport(1, err) {
	}
	localConn, err := net.ListenUDP("udp", saddr)
	if checkreport(1, err) {
	}
	Vlogf(2, "UDP Proxy serving on port %d\n", port)

	// Get server address
	ServerAddr, err := net.ResolveUDPAddr("udp", hostport)
	if checkreport(1, err) {
	}
	Vlogf(2, "UDP Connected to server at %s\n", hostport)
	RunProxy(localConn, ServerAddr)
}

// Go routine which manages connection from server to single client
func RunConnection(conn *UDPCap, udpcapsLock UDPCapsLock, localConn *net.UDPConn) {
	var buffer [1500]byte
	for {
		// Read from server
		conn.ServerConn.SetReadDeadline(time.Now().Add(time.Minute * 5))
		n, err := conn.ServerConn.Read(buffer[0:])
		if checkreport2(1, err) == 1 {
			//删除
			log.Println("udp delete timeout router")
			udpcapsLock.Lock.Lock()
			conn.ServerConn.Close()
			delete(udpcapsLock.UDPCaps, conn.ClientAddr.String())
			udpcapsLock.Lock.Unlock()
			break
		}
		// Relay it to client
		_, err = localConn.WriteToUDP(buffer[0:n], conn.ClientAddr)

	}
}

// Routine to handle inputs to Proxy port
func RunProxy(localConn *net.UDPConn, ServerAddr *net.UDPAddr) {
	var buffer [1500]byte
	log.Println("UDP start")
	var uncapsLock UDPCapsLock
	uncapsLock.UDPCaps = map[string]*UDPCap{}
	for {
		n, cliaddr, err := localConn.ReadFromUDP(buffer[0:])
		if checkreport(1, err) {
			continue
		}
		//Vlogf(3, "UDP Read sth. from client %s, routers length %d \n", cliaddr.String(), len(uncapsLock.UDPCaps))
		saddr := cliaddr.String()
		udpcap, found := uncapsLock.UDPCaps[saddr]
		log.Println("remote udp addr:", saddr)
		if !found {
			udpcap = NewConnection(ServerAddr, cliaddr)
			if udpcap == nil {
				continue
			}
			uncapsLock.Lock.Lock()
			uncapsLock.UDPCaps[saddr] = udpcap
			uncapsLock.Lock.Unlock()
			Vlogf(2, "UDP Created new connection for client %s  routers length %d \n", saddr, len(uncapsLock.UDPCaps))
			// Fire up routine to manage new connection
			go RunConnection(udpcap, uncapsLock, localConn)
		}

		// Relay to server
		_, err = udpcap.ServerConn.Write(buffer[0:n])
		if checkreport(1, err) {
			continue
		}
	}
}

// Log result if verbosity level high enough
func Vlogf(level int, format string, v ...interface{}) {
	//if level <= verbosity {
	log.Printf(format, v...)
	//}
}

// Handle errors
func checkreport(level int, err error) bool {
	if err == nil {
		return false
	}
	Vlogf(level, "Error: %s", err.Error())
	return true
}

// Handle errors
func checkreport2(level int, err error) int {
	if err == nil {
		return 0
	}
	Vlogf(level, "Error: %s", err.Error())
	if strings.Contains(err.Error(), "timeout") {
		return 1
	}

	if strings.Contains(err.Error(), "EOF") {
		return 1
	}

	if strings.Contains(err.Error(), "connection refused") {
		return 1
	}

	return 1
}

func udpLocal(localPort int, serverAddr string, serverPort int) {
	hostport := fmt.Sprintf("%s:%d", serverAddr, serverPort)
	Vlogf(3, "UDP Proxy port = %d, Server address = %s\n", localPort, hostport)
	setup(hostport, localPort)
}
