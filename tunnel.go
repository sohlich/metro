package main

import (
	"fmt"
	"io"
	"net"

	"golang.org/x/crypto/ssh"
)

type Endpoint struct {
	Host string
	Port string
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%s", endpoint.Host, endpoint.Port)
}

type SSHtunnel struct {
	Local  *Endpoint
	Remote *Endpoint
	Active bool
}

func (t *SSHtunnel) Start(client *ssh.Client) error {
	listener, err := net.Listen("tcp", t.Local.String())
	if err != nil {
		return err
	}
	t.Active = true
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go t.forward(conn, client)
	}
}

func (t *SSHtunnel) forward(localConn net.Conn, client *ssh.Client) {
	remoteConn, err := client.Dial("tcp", t.Remote.String())
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
	t.Active = true
}
