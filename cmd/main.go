package main

import (
	"flag"

	"github.com/sirupsen/logrus"

	"nvalembois/duckdns/webhook/internal"
	"nvalembois/duckdns/webhook/pkg/provider"
	"nvalembois/duckdns/webhook/pkg/updater"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/robfig/cron/v3"
)

func main() {
	// initialisation des logs
	logrus.SetFormatter(&logrus.TextFormatter{})

	var config internal.Config
	config.SetFlags()
	flag.Parse()

	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	provider := provider.DuckdnsProviderSolver{Config: config}

	cmd.RunWebhookServer("DuckDNS", &provider)

	ipUpdater := updater.IPUpdater{Provider: provider, Config: config}
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
