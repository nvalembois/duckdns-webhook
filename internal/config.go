package internal

import (
	"flag"
	"os"
)

type Config struct {
	DuckdnsToken string
	DnsName      string
}

func NewConfig() *Config {
	var c Config
	flag.StringVar(&(c.DuckdnsToken), "duckdnstoken", os.Getenv("DUCKDNS_TOKEN"), "duckdns token")
	flag.StringVar(&(c.DnsName), "dnsname", os.Getenv("DNSNAME"), "dns name to register")
	flag.Parse()
	return &c
}
