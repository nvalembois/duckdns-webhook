package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"nvalembois/external-dns/webhook/internal"
	"nvalembois/external-dns/webhook/provider"
	"nvalembois/external-dns/webhook/utils"
)

func main() {
	config := internal.NewConfig()

	ip, err := utils.FetchMyIp()
	if err != nil {
		fmt.Println("Erreur lors de la récupération de l'ip :", err)
		os.Exit(1)
	}

	dnsIps, err := net.LookupIP(config.DnsName)
	if err != nil {
		fmt.Println("Erreur lors de la résolution DNS :", err)
		return
	}
	for _, dnsIp := range dnsIps {
		if strings.EqualFold(ip.String(), dnsIp.String()) {
			fmt.Println("L'enregistrement est déjà valide")
			os.Exit(0)
		}
	}

	err = provider.DuckDNSUpdate(config.DuckdnsToken, config.DnsName, ip)
	if err != nil {
		fmt.Println("Erreur lors de la requête DuckDNS :", err)
		os.Exit(1)
	}

	os.Exit(0)
}
