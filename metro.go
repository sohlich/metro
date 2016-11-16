package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"fmt"

	"golang.org/x/crypto/ssh"
)

func main() {
	cfg := &ssh.ClientConfig{
		User: "sshuser",
		Auth: []ssh.AuthMethod{
			ssh.Password("sshuser"),
		},
	}

	s := &Relay{
		Host:   "192.168.0.100",
		Port:   "22",
		Config: cfg,
	}

	s.AddTunnel("7777", "seznam.cz", "80")
	s.AddTunnel("8888", "google.com", "80")

	fmt.Println("Trying to connect")
	if err := s.Activate(); err != nil {
		log.Panic("Cannot connect to SSH" + err.Error())
	}

	for _, tunnel := range s.Tunnels {
		if tunnel.Active {
			fmt.Printf("%s -> %s\n", tunnel.Local.Port, tunnel.Remote.String())
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Press Ctrl+C to close.")
	<-signalChan
	fmt.Println("Bye bye")

}
