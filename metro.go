package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"fmt"

	"flag"

	"bufio"

	"strings"

	"golang.org/x/crypto/ssh"
)

// InputConfig is used to parse and
// verify arguments passed
// to application.
type InputConfig struct {
	Host      string
	Port      string
	User      string
	Password  string
	TunelList string
	PKFile    string
}

func (cfg *InputConfig) Validate() error {

	if isEmpty(cfg.Host) && isEmpty(cfg.Port) {
		return fmt.Errorf("Host and Port cannot be empty")
	}

	if isEmpty(cfg.PKFile) && (isEmpty(cfg.User) || isEmpty(cfg.Password)) {
		return fmt.Errorf("Please provide user and password")
	}

	if isEmpty(cfg.TunelList) {
		return fmt.Errorf("No tunel list provided")
	}

	return nil
}

func isEmpty(str string) bool {
	return len(str) == 0
}

func main() {

	// Read args from cmd
	sshHost := flag.String("host", "", "Host for SSH")
	sshPort := flag.String("port", "22", "Port for SSH")
	sshUser := flag.String("user", "", "User for SSH")
	sshPass := flag.String("password", "", "Password for SSH")
	tunelList := flag.String("list", "", "CSV list of tunnels")
	flag.Parse()

	// Validate input

	cfg := &ssh.ClientConfig{
		User: *sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(*sshPass),
		},
	}

	s := &Relay{
		Host:   *sshHost,
		Port:   *sshPort,
		Config: cfg,
	}

	log.Println(*tunelList)

	f, err := os.Open(*tunelList)
	if err != nil {
		log.Fatalln("Cannot read config file")
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tunnel := scanner.Text()
		cfgTunnel := strings.Split(tunnel, ";")
		s.AddTunnel(cfgTunnel[0], cfgTunnel[1], cfgTunnel[2])
	}

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
