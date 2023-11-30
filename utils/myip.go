package utils

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

const (
	myIpApiUrl string = "https://api.myip.com"
)

type myIpResult struct {
	Ip string `json:"ip"`
}

// Returns IPv4 and IPv6 addresses
func GuessMyIPs() (net.IP, net.IP) {
	return GuessMyIPv4(), GuessMyIPv6()
}

func GuessMyIPv4() (net.IP) {
	return guessMyIp("tcp4").To4()
}

func GuessMyIPv6() (net.IP) {
	return guessMyIp("tcp6").To16()
}

func guessMyIp(networkString string) net.IP {
	// crée un client pour le réseau souhaité
	var dialer net.Dialer
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = func(ctx context.Context, network string, addr string) (net.Conn, error) {
		return dialer.DialContext(ctx, networkString, addr)
	}
	client := resty.New()
	client.SetTransport(transport)

	// interroge le serveur
	getResponse, err := client.R().Get(myIpApiUrl)
	if err != nil {
		logrus.Errorf("erreur lors de la requete à %s : %s", myIpApiUrl, err)
		return nil
	}
	if !getResponse.IsSuccess() {
		logrus.Errorln("la requête GET a échoué, code de statut : ", getResponse.StatusCode())
		return nil
	}

	// Extraction de l'adresse IP de la réponse
	var ipResult myIpResult
	err = json.Unmarshal(getResponse.Body(), &ipResult)
	if err != nil {
		logrus.Errorln("erreur lors du décodage JSON : ", err)
		return nil
	}

	// resultat
	logrus.Debugf("GuessMyIP(%s) -> %s", networkString, ipResult.Ip)
	return net.ParseIP(ipResult.Ip)
}

func HasIPv6() (bool) {
	interfaces, err := net.Interfaces()
	if err != nil {
		logrus.Errorln("erreur lors de la récupération des interfaces réseau: ", err)
		return false
	}

	for _, iface := range interfaces {
		logrus.Debugf("interface: %s (%s)", iface.Name, iface.Flags.String())
		if iface.Name == "lo" {
			logrus.Debugln("skip interface lo")
			continue
		}
		if ( iface.Flags & net.FlagRunning ) == 0 {
			logrus.Debugln("skip down interface interface:", iface.Name)
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			logrus.Errorf("erreur lors de la récupération des adresses pour l'interface %s : %s", iface.Name, err)
			continue
		}

		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				logrus.Errorf("erreur lors de la conversion de l'adresse CIDR %s de %s: ", addr.String(), iface.Name, err)
				continue
			}

			if ip.To4() == nil {
				logrus.Debugf("interface %s has IPv6: %s", iface.Name, ip)
				return true
			}
		}
	}
	return false
}

func HasIPv6Connectivity() (bool) {
	ipList, err := net.LookupIP("ipv6.google.com")
	if err != nil {
		logrus.Errorln("erreur lors de la résolution DNS de 'ipv6.google.com': ", err)
		return false
	}

	for _, ip := range ipList {
		if ip.To4() != nil {
			logrus.Debugln("la résolution de 'ipv6.google.com' a renvoyé une adresse ipv4: ", ip.String())
			continue
		}
		routes, err := netlink.RouteGet(ip)
		if err != nil {
			logrus.Debugf("pas de connectivité réseau pour %s: %s", ip.String(), err)
			return false
		}
		if len(routes) > 0 {
			return true
		}
		return false
	}
	logrus.Debugln("pas de résolution ipv6 pour 'ipv6.google.com'")
	return false
}