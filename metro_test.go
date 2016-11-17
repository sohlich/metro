package main

import "testing"

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

func checkError(t *testing.T, err error, msg, errmsg string) {
	if err == nil || err.Error() != msg {
		t.Error(errmsg)
		t.Fail()
	}
}
