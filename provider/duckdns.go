package provider

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const (
	allowedDomainSuffix string = ".duckdns.org"
	apiURLScheme        string = "https"
	apiURLHost          string = "www.duckdns.org"
	apiURLPath          string = "/update"
)

func DuckDNSUpdate(token string, name string, ip net.IP) error {
	// validation des arguments
	if !strings.HasSuffix(name, allowedDomainSuffix) {
		return errors.New("invalid suffix")
	}

	// construction de l'URL
	query := url.Values{}
	query.Add("domains", name)
	query.Add("token", token)
	query.Add("verbose", "true")
	if ip.To4() != nil {
		query.Add("ip", ip.String())
	} else {
		query.Add("ipv6", ip.String())
	}
	duckdnsApiUrl := url.URL{
		Scheme:   apiURLScheme,
		Host:     apiURLHost,
		Path:     apiURLPath,
		RawQuery: query.Encode(),
	}

	// Appel HTTP
	logrus.Debugln("Query : ", duckdnsApiUrl.String())
	getResponse, err := resty.New().R().Get(duckdnsApiUrl.String())
	if err != nil {
		return fmt.Errorf("erreur lors de la requête GET : %s", err)
	}
	if !getResponse.IsSuccess() {
		return fmt.Errorf("la requête GET a échoué, code de statut : %d", getResponse.StatusCode())
	}
	if !strings.HasPrefix(getResponse.String(), "OK") {
		return fmt.Errorf("echec de la reuête de mise à jour DuckDNS (%s -> %s) : %s\n", name, ip.String(), getResponse.String())
	}

	logrus.Debugln("succès de la reuête de mise à jour DuckDNS : ", getResponse.String())
	// Success
	return nil
}

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
