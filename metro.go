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

	s := &Server{
		Host:   "192.168.0.104",
		Port:   "22",
		Config: cfg,
	}

	s.AddTunnel("7777", "seznam.cz", "80")

	fmt.Println("Trying to connect")
	if err := s.Connect(); err != nil {
		log.Panic("Cannot connect to SSH" + err.Error())
	}
	s.StartAllTunnels()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Press Ctrl+C to close.")
	<-signalChan
	fmt.Println("Bye bye")

}
