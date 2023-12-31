package internal

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	DuckdnsToken string
	DnsName      string
	Debug        bool
	IPv4         bool
	IPv6         bool
}

func NewConfig() *Config {
	var c Config
	flag.StringVar(&(c.DuckdnsToken), "duckdnstoken", os.Getenv("DUCKDNS_TOKEN"), "duckdns token")
	flag.StringVar(&(c.DnsName), "dnsname", os.Getenv("DNSNAME"), "dns name to register")
	flag.BoolVar(&(c.Debug), "debug", boolVarOrDefault("LOGLEVEL_DEBUG", false), "debug")
	flag.BoolVar(&(c.IPv4), "4", boolVarOrDefault("DUCKDNS_IPV4", true), "ipv4")
	flag.BoolVar(&(c.IPv4), "6", boolVarOrDefault("DUCKDNS_IPV6", true), "ipv6")
	flag.Parse()
	return &c
}

func stringVarOrDefault(envVar string, defaultValue string) string {
	result, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultValue
	}
	return result
}

func boolVarOrDefault(envVar string, defaultValue bool) bool {
	result, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultValue
	}
	return strings.EqualFold(result, "true")
}
