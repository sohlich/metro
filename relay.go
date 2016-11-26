package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh"
)

const (
	// DefaultPort is common
	// used  SSH port
	DefaultPort = 22
)

// Relay represents ssh connection
// which provides the tuneling
type Relay struct {
	Host    string
	Port    string
	Config  *ssh.ClientConfig
	Client  *ssh.Client
	Tunnels []*SSHtunnel
	Active  *AtomicBool
	Wait    *sync.WaitGroup
}

// AddTunnel adds a tunel configuration
// to Relay struct.
// After the struct is activated
// all tunnels, which are sucessfuly activated
// provide relay to given destination.
func (s *Relay) AddTunnel(localPort, remoteHostPort string) error {
	if s.Active.Get() {
		return fmt.Errorf("Cannot add tunel to active relay")
	}

	// Split host port
	remote, err := EndpointFromHostPort(remoteHostPort)
	if err != nil {
		return err
	}

	s.Tunnels = append(s.Tunnels, &SSHtunnel{
		Local: &Endpoint{
			"localhost",
			localPort,
		},
		Remote:   remote,
		Active:   NewAtomicBool(),
		stopChan: make(chan bool),
	})
	return nil
}

// PrintActiveTunels prints
// all active tunnels
func (s *Relay) PrintActiveTunels() {
	for _, tunnel := range s.Tunnels {
		if tunnel.Active.Get() {
			log.Printf("%s -> %s\n", tunnel.Local.Port, tunnel.Remote.String())
		}
	}
}

// Activate connects the SSH connection and
// activates all tunnels from SSHtunnel list.
func (s *Relay) Activate() error {
	if s.Active.Get() {
		return fmt.Errorf("Relay already activated")
	}
	if err := s.connect(); err != nil {
		return errors.Wrap(err, "Cannot estabilish ssh connection")
	}
	for _, tunnel := range s.Tunnels {
		go tunnel.Start(s.Client, s.Wait)
	}
	s.Active.Set(true)
	return nil
}

// Stop disables all active tunnels
// and closes the ssh connection
func (s *Relay) Stop() {
	for _, tun := range s.Tunnels {
		tun.Stop()
	}

	log.Println("Waiting to stop relay")
	s.Wait.Wait()
	s.Client.Close()
}

func (s *Relay) connect() (err error) {
	if len(s.Port) == 0 {
		s.Port = fmt.Sprintf("%d", DefaultPort)
	}

	hostPort := fmt.Sprintf("%s:%s", s.Host, s.Port)
	s.Client, err = ssh.Dial("tcp", hostPort, s.Config)
	return err
}
