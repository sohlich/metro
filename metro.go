package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"time"

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

// Validate validates config struct
// if any of missiong arguments
// error is returned
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

var config *InputConfig

func parseConfig() *InputConfig {
	cfg := &InputConfig{}
	flag.StringVar(&cfg.Host, "host", "", "Host for SSH")
	flag.StringVar(&cfg.Port, "port", "22", "Port for SSH")
	flag.StringVar(&cfg.User, "user", "", "User for SSH")
	flag.StringVar(&cfg.Password, "password", "", "Password for SSH")
	flag.StringVar(&cfg.TunelList, "list", "", "CSV list of tunnels")
	flag.Parse()
	return cfg
}

func main() {

	config = parseConfig()
	if err := config.Validate(); err != nil {
		log.Fatal(err.Error())
	}
	// Read args from cmd

	cfg := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		Timeout: 20 * time.Second,
	}

	s := &Relay{
		Host:   config.Host,
		Port:   config.Port,
		Config: cfg,
	}

	loadTunnelsFromFile(s, config.TunelList)

	log.Println("Connecting to ssh endpoint ...")
	if err := s.Activate(); err != nil {
		log.Println("Cannot connect to SSH : \n" + err.Error())
		return
	}

	s.PrintActiveTunels()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Press Ctrl+C to close.")
	<-signalChan
	log.Println("Bye bye")

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
			log.Printf("Cannot add tunnel: %s reason: %s\n", tunnel, err.Error())
		} else {
			count++
		}
	}
	log.Printf("Loaded %d tunnels\n", count)
}
