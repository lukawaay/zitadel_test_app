package config

import (
	"fmt"
	"strconv"
)

type Config struct {
	InstanceDomain string
	InstancePort uint16
	InstanceProtocol string
	InstanceSecure bool
	Key string
	ClientID string
	RedirectURI string
	Port uint16
}

type Source struct {
	InstanceDomain string
	InstancePort string
	InstanceProtocol string
	Key string
	ClientID string
	RedirectURI string
	Port string
}

func Load(source Source) (*Config, error) {
	port, err := strconv.ParseUint(source.Port, 10, 16)
	if err != nil {
		return nil, err
	}

	instancePort, err := strconv.ParseUint(source.InstancePort, 10, 16)
	if err != nil {
		return nil, err
	}

	var instanceSecure bool
	if source.InstanceProtocol == "http" {
		instanceSecure = false
	} else if source.InstanceProtocol == "https" {
		instanceSecure = true
	} else {
		return nil, fmt.Errorf("Invalid protocol: %s", source.InstanceProtocol)
	}

	return &Config {
		InstanceDomain: source.InstanceDomain,
		InstancePort: uint16(instancePort),
		InstanceProtocol: source.InstanceProtocol,
		InstanceSecure: instanceSecure,
		Key: source.Key,
		ClientID: source.ClientID,
		RedirectURI: source.RedirectURI,
		Port: uint16(port),
	}, nil
}

const (
	SDKAuthEndpoint = "/auth"
	AppEndpoint = "/"
)
