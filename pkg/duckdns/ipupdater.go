package duckdns

import (
	"net"

	"nvalembois/duckdns/webhook/pkg/utils"

	"github.com/sirupsen/logrus"
)

type IPUpdater struct {
	currentIPv4 net.IP
	currentIPv6 net.IP
	Provider    *DuckdnsProviderSolver
	Config      *Config
}

func (i *IPUpdater) CheckIPv4() {
	ip := utils.GuessMyIPv4()
	if nil == ip {
		logrus.Warnln("pas d'adresse ipv4 détectée")
		return
	}

	if ip.Equal(i.currentIPv4) {
		logrus.Infof("l'adresse IPv4 publique n'a pas changé")
		return
	}

	if i.updateIP(ip) {
		i.currentIPv4 = ip
	}
}

func (i *IPUpdater) CheckIPv6() {
	ip := utils.GuessMyIPv6()
	if nil == ip {
		logrus.Warnln("pas d'adresse ipv6 détectée")
		return
	}

	if ip.Equal(i.currentIPv6) {
		logrus.Infof("l'adresse IPv6 publique n'a pas changé")
		return
	}

	if i.updateIP(ip) {
		i.currentIPv6 = ip
	}
}

func (i *IPUpdater) updateIP(ip net.IP) bool {
	success := true
	for _, domain := range i.Config.Domains {
		if !i.updateIPRecord(domain, ip) {
			success = false
		}
	}
	return success
}

func (i *IPUpdater) resolveDomain(domain string) []net.IP {

	dnsIps, err := net.LookupIP(domain)
	if err != nil {
		logrus.Errorln("erreur lors de la résolution DNS : ", err)
		return []net.IP{}
	}

	if logrus.GetLevel() == logrus.DebugLevel {
		for _, ip := range dnsIps {
			logrus.Debugf("résolution de %s : %s", domain, ip.String())
		}
	}

	return dnsIps
}

func (i *IPUpdater) updateIPRecord(domain string, ip net.IP) bool {
	if ipListContains(i.resolveDomain(domain), ip) {
		logrus.Infof("l'enregistrement %s pointe déjà vers l'adress %s", domain, ip)
		return true
	}
	if err := i.Provider.UpdateIP(domain, ip); err != nil {
		logrus.Errorf("erreur lors de la mise à jour DuckDNS (%s IN A %s) : %s", domain, ip, err)
		return false
	}
	return true
}

func ipListContains(ipList []net.IP, ip net.IP) bool {
	for _, dnsIp := range ipList {
		if ip.Equal(dnsIp) {
			return true
		}
	}
	return false
}
