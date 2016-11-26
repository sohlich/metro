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

func (t *SSHtunnel) forward(localConn net.Conn, client *ssh.Client) {
	remoteConn, err := client.Dial("tcp", t.Remote.String())
	if err != nil {
		log.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConnection(localConn, remoteConn, t.stopChan)
	closeConnections(localConn, remoteConn)
}

func copyConnection(conn1 net.Conn, conn2 net.Conn, stopChan chan bool) {
	chan1 := chanFromConn(conn1)
	chan2 := chanFromConn(conn2)
	for {
		select {
		case <-stopChan:
			return
		case b1 := <-chan1:
			if b1 == nil {
				return
			}
			conn2.Write(b1)
		case b2 := <-chan2:
			if b2 == nil {
				return
			}
			conn1.Write(b2)
		}
	}
}

// chanFromConn creates a channel from a Conn object, and sends everything it
//  Read()s from the socket to the channel.
func chanFromConn(conn net.Conn) chan []byte {
	c := make(chan []byte)
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
				c <- buf2[:n]
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
	log.Println("Closing connection")
	for _, conn := range connections {
		if err := conn.Close(); err != nil {
			log.Println("Cannot close conn")
		}
	}

	log.Println("Connection closed")
}
