package main

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"nvalembois/duckdns/webhook/pkg/duckdns"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/robfig/cron/v3"
)

func main() {
	// initialisation des logs
	logrus.SetFormatter(&logrus.TextFormatter{})
	if strings.EqualFold(os.Getenv(`LOGLEVEL_DEBUG`), `true`) {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var config duckdns.Config
	config.Init()

	provider := duckdns.DuckdnsProviderSolver{Config: &config}

	cmd.RunWebhookServer("DuckDNS", &provider)

	ipUpdater := duckdns.IPUpdater{Provider: &provider, Config: &config}
	c := cron.New()
	c.AddFunc("5 * * * *", func() {
		if config.IPv4 {
			ipUpdater.CheckIPv4()
		}
		if config.IPv6 {
			ipUpdater.CheckIPv6()
		}
	})
	c.Start()
}
