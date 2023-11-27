package main

import (
	"net"
	"os"

	"github.com/sirupsen/logrus"

	"nvalembois/external-dns/webhook/internal"
	"nvalembois/external-dns/webhook/provider"
	"nvalembois/external-dns/webhook/utils"
)

func main() {
	// initialisation des logs
	logrus.SetFormatter(&logrus.TextFormatter{})

	// traitement des arguments de la ligne de commande
	config := internal.NewConfig()
	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// détection des adresses ipv4 et ipv6 actuelle publiques
	ipv4, ipv6 := utils.GuessMyIPs()

	// récupération de la résolution DNS actuelle
	dnsIps, err := net.LookupIP(config.DnsName)
	if err != nil {
		logrus.Errorln("erreur lors de la résolution DNS : ", err)
		os.Exit(0)
	}
	if logrus.GetLevel() == logrus.DebugLevel {
		for _, ip := range dnsIps {
			logrus.Debugf("résolution de %s : %s", config.DnsName, ip.String())
		}
	}
	retCode := 0

	// vérification pour l'adresse ipv4
	if nil == ipv4 {
		logrus.Warnln("pas d'adresse ipv4 détectée")
	} else if ipListContains(dnsIps, ipv4) {
		logrus.Infof("l'enregistrement %s pointe déjà vers l'adress %s", config.DnsName, ipv4.String())
	} else {
		// mise à jour pour l'adresse ipv4
		err = provider.DuckDNSUpdate(config.DuckdnsToken, config.DnsName, ipv4)
		if err != nil {
			logrus.Errorln("erreur lors de la mise à jour DuckDNS : ", err)
			retCode = 1
		}
	}

	// vérification pour l'adresse ipv6
	if nil == ipv6 {
		logrus.Infoln("pas d'adresse ipv6 détectée")
	} else if ipListContains(dnsIps, ipv6) {
		logrus.Infof("l'enregistrement %s pointe déjà vers l'adress %s", config.DnsName, ipv6.String())
	} else {
		// mise à jour pour l'adresse ipv6
		err = provider.DuckDNSUpdate(config.DuckdnsToken, config.DnsName, ipv4)
		if err != nil {
			logrus.Errorln("erreur lors de la mise à jour DuckDNS : ", err)
			retCode = 1
		}
	}

	os.Exit(retCode)
}

func ipListContains(ipList []net.IP, ip net.IP) bool {
	for _, dnsIp := range ipList {
		if ip.Equal(dnsIp) {
			return true
		}
	}
	return false
}
