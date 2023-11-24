package utils

import (
	"encoding/json"
	"fmt"
	"net/netip"

	"github.com/go-resty/resty/v2"
)

const (
	myIpApiUrl string = "https://api.myip.com"
)

type myIpResult struct {
	Ip string `json:"ip"`
}

func FetchMyIp() (netip.Addr, error) {
	var ip netip.Addr

	client := resty.New()
	getResponse, err := client.R().Get(myIpApiUrl)
	if err != nil {
		return ip, fmt.Errorf("erreur lors de la requete à %s : %s", myIpApiUrl, err)
	}

	var ipResult myIpResult
	if !getResponse.IsSuccess() {
		return ip, fmt.Errorf("la requête GET a échoué, code de statut : %d", getResponse.StatusCode())
	}

	err = json.Unmarshal(getResponse.Body(), &ipResult)
	if err != nil {
		return ip, fmt.Errorf("erreur lors du décodage JSON : %s", err)
	}

	ip, err = netip.ParseAddr(ipResult.Ip)
	if err != nil {
		return ip, fmt.Errorf("l'adresse IP n'est pas valide: %s", ipResult.Ip)
	}

	return ip, nil
}
