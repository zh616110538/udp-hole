package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

type udpHoleClient struct {
	handle *net.UDPConn
	peers  []net.Addr
	server net.Addr
	key    string
}

func (c *udpHoleClient) recvMessage(info chan []byte) {
	for {
		buf := make([]byte, 1024)
		count, peer_addr, err := c.handle.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		switch peer_addr.String() {
		case c.server.String():
			info <- buf[:count]
		default:
			fmt.Println(string(buf[:count]))
		}
	}
}

func (c *udpHoleClient) sendMessage(info chan []byte) {
	for {
		select {
		case clients := <-info:
			c.addToPeers(clients)
		case <-time.Tick(time.Second):
			c.keepAlive()
			c.sayHelloToPeers()
		}
	}
}

func (c *udpHoleClient) addToPeers(clients []byte) {
	r := bufio.NewReader(bytes.NewReader(clients))
	for {
		l, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		addr, err := net.ResolveUDPAddr("udp", string(l))
		if err != nil {
			panic(err)
		}
		c.peers = append(c.peers, addr)
	}
}
func (c *udpHoleClient) keepAlive() {
	c.handle.WriteTo([]byte("key:"+c.key), c.server)
}

func (c *udpHoleClient) sayHelloToPeers() {
	for _, peer := range c.peers {
		c.handle.WriteTo([]byte("hello "+peer.String()), peer)
	}
}
