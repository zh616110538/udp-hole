package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"strconv"
)

func main() {
	var connect_to_ip = flag.String("c", "", "connect to ip")
	var is_server = flag.Bool("s", false, "server mode")
	// var host_ip = flag.String("h", "", "server ip")
	var key = flag.String("k", "", "key")
	flag.Parse()
	if *connect_to_ip != "" {
		client(*connect_to_ip, *key)
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
	var server_addr, err = net.ResolveUDPAddr("udp", host_ip+":9000")
	if err != nil {
		panic(err)
	}
	client := udpHoleClient{
		handle: startUdpServer(0),
		server: server_addr,
		peers:  make([]net.Addr, 0),
		key:    key,
	}
	defer client.handle.Close()
	c := make(chan []byte)
	go func() {
		client.recvMessage(c)
	}()
	client.sendMessage(c)
}

func server() {
	var udp = startUdpServer(9000)
	var key_to_netAddr = make(map[string]map[string]bool)
	defer udp.Close()
	var buffer = make([]byte, 1024)
	for {
		var n, client_addr, err = udp.ReadFrom(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Println("receive message")
		key, err := resolveKey(buffer[:n])
		if err != nil {
			fmt.Println(err)
			continue
		}
		if _, ok := key_to_netAddr[key]; !ok {
			key_to_netAddr[key] = make(map[string]bool)
		}

		// if key_to_netAddr[key][client_addr.String()] {
		// 	// TODO:已经建立过了连接，相当于保持nat，让客户端发慢点就行了
		// 	continue
		// }
		key_to_netAddr[key][client_addr.String()] = true

		// TODO:如果这里有几千个客户端咋整
		str := ""
		for addr := range key_to_netAddr[key] {
			if client_addr.String() != addr {
				str += addr + "\n"
			}
		}
		if str != "" {
			udp.WriteTo([]byte(str[:len(str)-1]), client_addr)
			fmt.Println("**********", str[:len(str)-1])
		}
	}
}

func resolveKey(data []byte) (string, error) {
	var key_str = string(data)
	if key_str[0:4] == "key:" {
		return string(key_str[4:]), nil
	}
	return "", errors.New("error key")
}
