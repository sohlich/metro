package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"sync"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// Endpoint is wrapping
// the host and port of SSHtunnel endpoint
type Endpoint struct {
	Host string
	Port string
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%s", endpoint.Host, endpoint.Port)
}

// EndpointFromHostPort is  helper function
// for creating endpoint from string in format "<host>:<port>"
func EndpointFromHostPort(hostPort string) (*Endpoint, error) {
	remote := strings.Split(hostPort, ":")
	if len(remote) < 2 {
		return nil, fmt.Errorf("Tunnel does not contain a port")
	}
	return &Endpoint{
		remote[0],
		remote[1],
	}, nil
}

// SSHtunnel represents the
// tunneled connection via ssh
type SSHtunnel struct {
	Local    *Endpoint
	Remote   *Endpoint
	Active   *AtomicBool
	stopChan chan bool
	listener net.Listener
}

// Start opens(activates) the
// SSHtunnel tunnel
func (t *SSHtunnel) Start(client *ssh.Client, wait *sync.WaitGroup) error {
	if client == nil {
		return errors.New("metro: ssh client cannot be nil")
	}

	listener, err := net.Listen("tcp", t.Local.String())
	if err != nil {
		return err
	}

	wait.Add(1)
	defer wait.Done()

	t.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		t.Active.Set(true)
		go t.forward(conn, client)
	}
}

// Stop closes the SSH tunnel
func (t *SSHtunnel) Stop() {
	t.listener.Close()
	if t.Active.Get() {
		// closing stopChan causes
		// nil value  while channel is
		// readed
		close(t.stopChan)
	}
}

func (t *SSHtunnel) forward(local net.Conn, client *ssh.Client) {
	remote, err := client.Dial("tcp", t.Remote.String())
	if err != nil {
		log.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConnection(local, remote, t.stopChan)
	closeConnections(local, remote)
}

func copyConnection(local net.Conn, remote net.Conn, stopChan chan bool) {
	localChan := chanFromConn(local)
	remoteChan := chanFromConn(remote)
	for {
		select {
		case <-stopChan:
			return
		case b1 := <-localChan:
			if b1 == nil {
				return
			}
			remote.Write(b1.Data)
		case b2 := <-remoteChan:
			if b2 == nil {
				return
			}
			local.Write(b2.Data)
		}
	}
}

// Data wraps the content
// and amount of data read from
// connection
type Data struct {
	Size int
	Data []byte
}

// chanFromConn creates a channel from a Conn object, and sends everything it
//  Read()s from the socket to the channel.
func chanFromConn(conn net.Conn) chan *Data {
	c := make(chan *Data)
	go func() {
		// If connection closed,
		// EOR err is returned and channel closed
		// Reading from closed channel returns nil
		defer close(c)
		buf1 := make([]byte, 32*1024)
		buf2 := make([]byte, 32*1024)
		for {
			n, err := conn.Read(buf1)
			if n > 0 {
				copy(buf2, buf1[:n])
				c <- &Data{n, buf2[:n]}
			}
			if err != nil {
				c <- nil
				break
			}
		}
	}()
	return c
}

func closeConnections(connections ...net.Conn) {
	for _, conn := range connections {
		if err := conn.Close(); err != nil {
			log.Println("Cannot close connection: %v\n", err)
		}
	}
}
