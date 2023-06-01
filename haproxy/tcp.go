package haproxy

import (
	"errors"
	"github.com/juju/ratelimit"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// Bucket adding 100KB every second, holding max 100KB
var bucket = ratelimit.NewBucketWithRate(200*1024, 200*1024)

func TcpForward(server string, inPort string, outPort string) {
	l, err := net.Listen("tcp", inPort)
	if err != nil {
		log.Panic(err)
		return
	}
	log.Println("TCP listen " + inPort)
	// 死循环，每当遇到连接时，调用handle
	for {
		client, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		log.Println("remote tcp addr:", client.RemoteAddr())
		go handle(client, server, outPort)
	}

}

func handle(client net.Conn, server string, output string) {
	if client == nil {
		return
	}
	log.Println("TCP handle new connection remoteAddr: ", client.RemoteAddr().String())
	defer client.Close()
	destConn, err := net.DialTimeout("tcp", server+output, time.Second*3)
	if err != nil {
		log.Println(err)
		return
	}
	defer destConn.Close()
	header := make([]byte, 3)
	if _, err := io.ReadAtLeast(client, header, 3); err != nil {
		return
	}
	log.Println("header:", string(header), header)
	destConn.Write(header)

	//将客户端的请求转发至服务端，将服务端的响应转发给客户端。io.Copy为阻塞函数，文件描述符不关闭就不停止
	//go io.Copy(destConn, client)
	//io.Copy(client, destConn)
	relay(client, destConn)
	log.Println("TCP handle close connection.")
}

// relay copies between left and right bidirectionally
func relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
		right.SetReadDeadline(time.Now().Add(wait)) // unblock read on right
	}()

	_, err = io.Copy(left, right)
	left.SetReadDeadline(time.Now().Add(wait)) // unblock read on left
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}

// relay copies between left and right bidirectionally
func relayLimit(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, ratelimit.Reader(left, bucket))
		right.SetReadDeadline(time.Now().Add(wait)) // unblock read on right
	}()

	_, err = io.Copy(left, ratelimit.Reader(right, bucket))
	left.SetReadDeadline(time.Now().Add(wait)) // unblock read on left
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}
