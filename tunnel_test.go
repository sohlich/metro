package main

import "testing"

func TestEndpointFromHostPort(t *testing.T) {
	hostPort := "localhost:22"

	endpoint, err := EndpointFromHostPort(hostPort)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if endpoint.Host != "localhost" && endpoint.Port != "22" {
		t.Error("Parsing endpoint failed")
		t.Fail()
	}

	hostPort = "localhost"
	_, err = EndpointFromHostPort(hostPort)
	if err == nil {
		t.Fail()
	}

}
