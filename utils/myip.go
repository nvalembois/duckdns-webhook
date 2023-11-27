package utils

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const (
	myIpApiUrl string = "https://api.myip.com"
)

type myIpResult struct {
	Ip string `json:"ip"`
}

// Returns IPv4 and IPv6 addresses
func GuessMyIPs() (net.IP, net.IP) {
	return guessMyIp("tcp4").To4(), guessMyIp("tcp6").To16()
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
