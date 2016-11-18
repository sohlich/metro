package main

import "testing"
import "os"

func TestCfgValidation(t *testing.T) {
	cfg := &InputConfig{}
	err := cfg.Validate()
	checkError(t, err, "Host and Port cannot be empty", "Validation for Host and Port failed")
	cfg.Host = "localhost"
	cfg.Port = "22"
	err = cfg.Validate()
	checkError(t, err, "Please provide user and password", "Validation for User and Password failed")
	cfg.User = "user"
	err = cfg.Validate()
	checkError(t, err, "Please provide user and password", "Validation for User and Password failed")
	cfg.Password = "pass"
	err = cfg.Validate()
	checkError(t, err, "No tunel list provided", "Validation for tunel list failed")

}

func TestConfigParse(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"metro",
		"-host=localhost",
		"-port=30",
		"-user=ssh",
		"-password=sshpass",
		"-list=lst.csv",
		"-timeout=25"}
	cfg := parseConfig()

	if cfg.Host != "localhost" {
		t.Error("Parsing host failed")
		t.Fail()
	}

	if cfg.Port != "30" {
		t.Error("Parsing port failed")
		t.Fail()
	}
	if cfg.User != "ssh" {
		t.Error("Parsing user failed")
		t.Fail()
	}
	if cfg.Password != "sshpass" {
		t.Error("Parsing password failed")
		t.Fail()
	}
	if cfg.TunelList != "lst.csv" {
		t.Error("Parsing list failed")
		t.Fail()
	}
	if cfg.SSHTimeout != 25 {
		t.Error("Parsing timeout failed")
		t.Fail()
	}
}

func checkError(t *testing.T, err error, msg, errmsg string) {
	if err == nil || err.Error() != msg {
		t.Error(errmsg)
		t.Fail()
	}
}
