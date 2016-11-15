package main

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

type Server struct {
	Host    string
	Port    string
	Config  *ssh.ClientConfig
	Client  *ssh.Client
	Tunnels []*SSHtunnel
	Running bool
}

func (s *Server) Connect() (err error) {
	if len(s.Port) == 0 {
		s.Port = fmt.Sprintf("%d", s.Port)
	}

	hostPort := fmt.Sprintf("%s:%s", s.Host, s.Port)
	s.Client, err = ssh.Dial("tcp", hostPort, s.Config)
	return err
}

func (s *Server) AddTunnel(localPort, remoteHost, remotePort string) error {
	s.Tunnels = append(s.Tunnels, &SSHtunnel{
		&Endpoint{
			"localhost",
			localPort,
		},
		&Endpoint{
			remoteHost,
			remotePort,
		},
		false,
	})

	return nil
}

func (s *Server) StartAllTunnels() {
	for _, tunnel := range s.Tunnels {
		go tunnel.Start(s.Client)
	}
}
