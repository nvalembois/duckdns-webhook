package duckdns

import (
	"errors"
	"flag"
	"os"
	"regexp"
	"slices"
	"strings"
)

var duckdnsDomainRE = regexp.MustCompile(`^([a-z][a-z\d-]+\.)*[a-z][a-z\d-]+\.duckdns\.org$`)

type Config struct {
	IPv4         bool
	IPv6         bool
	Domains      domains
	DuckdnsToken string
}

func (c *Config) SetFlags() {
	flag.BoolVar(&(c.IPv4), `4`,
		boolVarOrDefault(`DUCKDNS_IPV4`, true),
		`ipv4`)
	flag.BoolVar(&(c.IPv4), `6`,
		boolVarOrDefault(`DUCKDNS_IPV6`, false),
		`ipv6`)
	flag.StringVar(&(c.DuckdnsToken), `duckdnstoken`,
		os.Getenv(`DUCKDNS_TOKEN`),
		`duckdns token`)
	c.Domains.Set(os.Getenv(`DUCKDNS_DOMAINS`))
	flag.Var(&(c.Domains), `domains`,
		`existing duckdns domains`)
}

func (c *Config) Init() {
	c.IPv4 = boolVarOrDefault(`DUCKDNS_IPV4`, true)
	c.IPv4 = boolVarOrDefault(`DUCKDNS_IPV6`, false)
	c.DuckdnsToken = os.Getenv(`DUCKDNS_TOKEN`)
	c.Domains.Set(os.Getenv(`DUCKDNS_DOMAINS`))
}

type domains []string

func (d *domains) String() string {
	return strings.Join(*d, `,`)
}

func (d *domains) Set(value string) error {
	*d = strings.Split(strings.ToLower(value), `,`)
	for _, s := range *d {
		if !duckdnsDomainRE.MatchString(s) {
			return errors.New(`invalid domain`)
		}
	}
	return nil
}

func (c *Config) HasDomain(domain string) bool {
	return slices.Contains(c.Domains, strings.ToLower(domain))
}

func boolVarOrDefault(envVar string, defaultValue bool) bool {
	result, ok := os.LookupEnv(envVar)
	if !ok {
		return defaultValue
	}
	return strings.EqualFold(result, `true`)
}
