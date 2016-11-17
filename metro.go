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

	"io"

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

	loadTunnelsFromFile(s, *tunelList)

	fmt.Println("Connecting to ssh endpoint ...")
	if err := s.Activate(); err != nil {
		log.Panic("Cannot connect to SSH" + err.Error())
	}

	s.PrintTunels()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Press Ctrl+C to close.")
	<-signalChan
	fmt.Println("Bye bye")

}

func loadTunnelsFromFile(s *Relay, filepath string) {
	fmt.Printf("Reading tunel file: %s\n", filepath)
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalln("Cannot read config file")
	}
	defer f.Close()
	loadTunnels(s, f)
}

func loadTunnels(s *Relay, tunels io.Reader) {
	log.Printf("Loading tunnels...")
	count := 0
	scanner := bufio.NewScanner(tunels)
	for scanner.Scan() {
		tunnel := scanner.Text()
		cfgTunnel := strings.Split(tunnel, ";")
		if err := s.AddTunnel(cfgTunnel[0], cfgTunnel[1]); err != nil {
			fmt.Printf("Cannot add tunnel: %s reason: %s\n", tunnel, err.Error())
		}
		count++
	}
	log.Printf("Loaded %d tunnels\n", count)
}
