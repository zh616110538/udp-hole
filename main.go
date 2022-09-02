package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"
)

func main() {
	var is_client = flag.Bool("c", false, "client mode")
	var is_server = flag.Bool("s", false, "server mode")
	var host_ip = flag.String("h", "", "server ip")
	var key = flag.String("k", "", "key")
	flag.Parse()
	if *is_client {
		if *host_ip == "" {
			flag.PrintDefaults()
			return
		}
		client(*host_ip, *key)
	} else if *is_server {
		server()
	} else {
		flag.PrintDefaults()
	}
}

func startUdpServer(port int) *net.UDPConn {
	var addr, _ = net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
	var udp_server, err = net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("start local udp server:", udp_server.LocalAddr())
	return udp_server
}

func client(host_ip, key string) {
	var udp = startUdpServer(0)
	defer udp.Close()
	var server_addr, err = net.ResolveUDPAddr("udp", host_ip+":9000")
	if err != nil {
		panic(err)
	}
	udp.WriteTo([]byte("key:"+key), server_addr)
	for {
		var buffer = make([]byte, 1024)
		var count, addr, err = udp.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Println(addr, ":", string(buffer[:count]))
	}
}

func server() {
	var udp = startUdpServer(9000)
	var key_to_netAddr = make(map[string]net.Addr)
	defer udp.Close()
	var buffer = make([]byte, 1024)
	for {
		var _, client_addr, err = udp.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}
		go func() {
			var new_udp = startUdpServer(0)
			defer new_udp.Close()
			for i := 0; i < 60; i++ {
				var ret = fmt.Sprintf("hello client:%d", i)
				udp.WriteTo([]byte(ret), client_addr)
				new_udp.WriteTo([]byte(ret), client_addr)
				time.Sleep(time.Second)
			}
		}()
	}
}

func resolveKey(data []byte) (string, error) {
	var key_str = string(data)
	if key_str[0:4] == "key:" {

	}
	return "", errors.New("error key")
}
