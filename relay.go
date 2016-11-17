package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh"
)

// Relay represents ssh connection
// which provides the tuneling
type Relay struct {
	Host    string
	Port    string
	Config  *ssh.ClientConfig
	Client  *ssh.Client
	Tunnels []*SSHtunnel
	Active  bool
}

// AddTunnel adds a tunel configuration
// to Relay struct.
// After the struct is activated
// all tunnels, which are sucessfuly activated
// provide relay to given destination.
func (s *Relay) AddTunnel(localPort, remoteHostPort string) error {
	if s.Active {
		return fmt.Errorf("Cannot add tunel to active relay")
	}

	// Split host port
	remote, err := EndpointFromHostPort(remoteHostPort)
	if err != nil {
		return err
	}

	s.Tunnels = append(s.Tunnels, &SSHtunnel{
		&Endpoint{
			"localhost",
			localPort,
		},
		remote,
		false,
	})
	return nil
}

// PrintTunels prints
// all active tunnels
func (s *Relay) PrintActiveTunels() {
	for _, tunnel := range s.Tunnels {
		if tunnel.Active {
			log.Printf("%s -> %s\n", tunnel.Local.Port, tunnel.Remote.String())
		}
	}
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

func (s *Relay) connect() (err error) {
	if len(s.Port) == 0 {
		s.Port = fmt.Sprintf("%d", s.Port)
	}

	hostPort := fmt.Sprintf("%s:%s", s.Host, s.Port)
	s.Client, err = ssh.Dial("tcp", hostPort, s.Config)
	return err
}
