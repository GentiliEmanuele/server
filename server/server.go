package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"server/services"
)

func main() {
	playOnServer(
		GetOutboundIP().String(),
		os.Getenv("PORT"),
	)
	select {}
}

func playOnServer(ip, port string) {
	var flag bool
	//Wait request from load balancer
	server := rpc.NewServer()
	err := server.RegisterName("Server", services.NewServer())
	if err != nil {
		fmt.Printf("An error occured %s\n", err)
	}
	dotPort := fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp", dotPort)
	go server.Accept(lis)
	registryAddress := os.Getenv("REGISTRY")
	registry, err := rpc.Dial("tcp", registryAddress)
	if err != nil {
		fmt.Printf("An error occured %s\n", err)
	}
	err = registry.Call("Registry.Register", &services.Args{IPAddress: ip, PortNumber: port}, &flag)
	if err != nil {
		fmt.Printf("An error occured %s\n", err)
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}
