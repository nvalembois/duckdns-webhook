package provider

import (
	"errors"
	"fmt"
	"net"
	"slices"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
)

// func (b DuckDNS) AdjustEndpoints(endpoints []*endpoint.Endpoint) ([]*endpoint.Endpoint, error) {
// 	return endpoints, nil
// }

// func (b DuckDNS) GetDomainFilter() endpoint.DomainFilter {
// 	return endpoint.DomainFilter{}
// }

// type contextKey struct {
// 	name string
// }

// func (k *contextKey) String() string { return "provider context value " + k.name }

// // RecordsContextKey is a context key. It can be used during ApplyChanges
// // to access previously cached records. The associated value will be of
// // type []*endpoint.Endpoint.
// var RecordsContextKey = &contextKey{"records"}

// // EnsureTrailingDot ensures that the hostname receives a trailing dot if it hasn't already.
// func EnsureTrailingDot(hostname string) string {
// 	if net.ParseIP(hostname) != nil {
// 		return hostname
// 	}

// 	return strings.TrimSuffix(hostname, ".") + "."
// }

// // Difference tells which entries need to be respectively
// // added, removed, or left untouched for "current" to be transformed to "desired"
// func Difference(current, desired []string) ([]string, []string, []string) {
// 	add, remove, leave := []string{}, []string{}, []string{}
// 	index := make(map[string]struct{}, len(current))
// 	for _, x := range current {
// 		index[x] = struct{}{}
// 	}
// 	for _, x := range desired {
// 		if _, found := index[x]; found {
// 			leave = append(leave, x)
// 			delete(index, x)
// 		} else {
// 			add = append(add, x)
// 			delete(index, x)
// 		}
// 	}
// 	for x := range index {
// 		remove = append(remove, x)
// 	}
// 	return add, remove, leave
// }

const duckDNSApiUrl string = "https://www.duckdns.org/update"

// duckdnsProvider implements the provider-specific logic needed to
// 'present' an ACME challenge TXT record for your own DNS provider.
// To do so, it must implement the `github.com/cert-manager/cert-manager/pkg/acme/webhook.Solver`
// interface.
type DuckdnsProviderSolver struct {
	config DuckDNSProviderConfig
}

// Name is used as the name for this DNS solver when referencing it on the ACME
// Issuer resource.
func (c *DuckdnsProviderSolver) Name() string {
	return "duckdns"
}

// Initialize will be called when the webhook first starts.
// This method can be used to instantiate the webhook, i.e. initialising
// connections or warming up caches.
// Typically, the kubeClientConfig parameter is used to build a Kubernetes
// client that can be used to fetch resources from the Kubernetes API, e.g.
// Secret resources containing credentials used to authenticate with DNS
// provider accounts.
// The stopCh can be used to handle early termination of the webhook, in cases
// where a SIGTERM or similar signal is sent to the webhook process.
func (c *DuckdnsProviderSolver) Initialize(stopCh <-chan struct{}) error {
	return nil
}

func (c *DuckdnsProviderSolver) UpdateIP(domain string, ip net.IP) error {
	if err := c.assertValidDomain(domain); err != nil {
		return err
	}

	req := resty.New().R().
		SetQueryParam("domains", domain).
		SetQueryParam("token", c.config.DuckdnsToken).
		SetQueryParam("verbose", "true")
	if ip.To4() != nil {
		req = req.SetQueryParam("ip", ip.String())
	} else {
		req = req.SetQueryParam("ipv6", ip.String())
	}
	req.Method = resty.MethodGet
	req.URL = duckDNSApiUrl

	logrus.Debugln("Query : ", req.RawRequest)
	getResponse, err := req.Send()
	if err != nil {
		return fmt.Errorf("erreur lors de la requête GET : %s", err)
	}
	if !getResponse.IsSuccess() {
		return fmt.Errorf("la requête GET a échoué, code de statut : %d", getResponse.StatusCode())
	}
	if !strings.HasPrefix(getResponse.String(), "OK") {
		return fmt.Errorf("echec de la reuête de mise à jour DuckDNS: %s", getResponse.String())
	}

	logrus.Debugln("succès de la reuête de mise à jour DuckDNS : ", getResponse.String())
	// Success
	return nil
}

// Present is responsible for actually presenting the DNS record with the
// DNS provider.
// This method should tolerate being called multiple times with the same value.
// cert-manager itself will later perform a self check to ensure that the
// solver has correctly configured the DNS provider.
func (c *DuckdnsProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	if err := c.assertValidDomain(ch.DNSName); err != nil {
		return err
	}

	req := resty.New().R().
		SetQueryParam("domains", ch.DNSName).
		SetQueryParam("token", c.config.DuckdnsToken).
		SetQueryParam("verbose", "true").
		SetQueryParam("txt", ch.Key)
	req.Method = resty.MethodGet
	req.URL = duckDNSApiUrl

	logrus.Debugln("Query : ", req.RawRequest)
	getResponse, err := req.Send()
	if err != nil {
		return fmt.Errorf("erreur lors de la requête TXT : %s", err)
	}
	if !getResponse.IsSuccess() {
		return fmt.Errorf("la requête TXT a échoué, code de statut : %d", getResponse.StatusCode())
	}
	if !strings.HasPrefix(getResponse.String(), "OK") {
		return fmt.Errorf("echec de la reuête de mise à jour DuckDNS : %s", getResponse.String())
	}

	logrus.Debugln("succès de la reuête de mise à jour DuckDNS : ", getResponse.String())
	// Success

	return nil
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
// If multiple TXT records exist with the same record name (e.g.
// _acme-challenge.example.com) then **only** the record with the same `key`
// value provided on the ChallengeRequest should be cleaned up.
// This is in order to facilitate multiple DNS validations for the same domain
// concurrently.
func (c *DuckdnsProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	if err := c.assertValidDomain(ch.DNSName); err != nil {
		return err
	}

	req := resty.New().R().
		SetQueryParam("domains", ch.DNSName).
		SetQueryParam("token", c.config.DuckdnsToken).
		SetQueryParam("verbose", "true").
		SetQueryParam("txt", ch.Key).
		SetQueryParam("clear", "true")

	req.Method = resty.MethodGet
	req.URL = duckDNSApiUrl

	logrus.Debugln("Query : ", req.RawRequest)
	getResponse, err := req.Send()
	if err != nil {
		return fmt.Errorf("erreur lors de la requête TXT : %s", err)
	}
	if !getResponse.IsSuccess() {
		return fmt.Errorf("la requête TXT a échoué, code de statut : %d", getResponse.StatusCode())
	}
	if !strings.HasPrefix(getResponse.String(), "OK") {
		return fmt.Errorf("echec de la reuête de mise à jour DuckDNS : %s", getResponse.String())
	}

	logrus.Debugln("succès de la reuête de mise à jour DuckDNS : ", getResponse.String())
	// Success

	return nil
}

func (c *DuckdnsProviderSolver) assertValidDomain(domain string) error {
	if !slices.Contains(c.config.Domains, strings.ToLower(domain)) {
		return errors.New("invalid domain")
	}
	return nil
}
