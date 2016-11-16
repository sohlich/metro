package main

import (
	"fmt"

	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh"
)

type Relay struct {
	Host    string
	Port    string
	Config  *ssh.ClientConfig
	Client  *ssh.Client
	Tunnels []*SSHtunnel
	Active  bool
}

func (s *Relay) connect() (err error) {
	if len(s.Port) == 0 {
		s.Port = fmt.Sprintf("%d", s.Port)
	}

	hostPort := fmt.Sprintf("%s:%s", s.Host, s.Port)
	s.Client, err = ssh.Dial("tcp", hostPort, s.Config)
	return err
}

func (s *Relay) AddTunnel(localPort, remoteHost, remotePort string) error {
	if s.Active {
		return fmt.Errorf("Cannot add tunel to active relay")
	}
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

func (s *Relay) Activate() error {
	if err := s.connect(); err != nil {
		return errors.Wrap(err, "Cannot estabilish ssh connection")
	}
	for _, tunnel := range s.Tunnels {
		go tunnel.Start(s.Client)
		tunnel.Active = true
	}
	s.Active = true
	return nil
}
