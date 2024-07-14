package provider

import (
	"errors"
	"flag"
	"os"
	"regexp"
	"strings"
)

var duckdnsDomainRE = regexp.MustCompile(`^([a-z][a-z\d-]+\.)*[a-z][a-z\d-]+\.duckdns\.org$`)

type domains []string

func (d *domains) String() string {
	return strings.Join(*d, `,`)
}

func (d *domains) Set(value string) error {
	*d = strings.Split(strings.ToLower(value), `,`)
	for _, s := range *d {
		if !duckdnsDomainRE.MatchString(s) {
			return errors.New("invalid domain")
		}
	}
	return nil
}

type DuckDNSProviderConfig struct {
	DuckdnsToken string
	Domains      domains
}

func (c *DuckDNSProviderConfig) SetFlags() {
	flag.StringVar(&(c.DuckdnsToken), "duckdnstoken", os.Getenv("DUCKDNS_TOKEN"), "duckdns token")
	c.Domains.Set(os.Getenv("DUCKDNS_DOMAINS"))
	flag.Var(&(c.Domains), `domains`, `existing duckdns domains`)
}

func (c *DuckDNSProviderConfig) Validate() bool {
	return len(c.DuckdnsToken) > 0 && len(c.Domains) > 0
}
