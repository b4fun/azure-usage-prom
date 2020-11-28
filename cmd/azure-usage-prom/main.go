package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	azureusageprom "github.com/b4fun/azure-usage-prom"
	"github.com/b4fun/azure-usage-prom/internal/collector"
	"github.com/b4fun/azure-usage-prom/internal/usagelister"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	flagAddr         = flag.String("listen-address", ":8080", "http listen address")
	flagEnvironment  = flag.String("environment", azure.PublicCloud.Name, "azure cloud environment name")
	flagAuthWithCLI  = flag.Bool("use-cli-auth", false, "use azure cli auth?")
	flagQueryTargets = flag.String(
		"query-targets", "",
		"comma separated targets list: <rp>|<subscription-id>|<location>,<rp>|<subscription-id>|<location>",
	)
)

func parseQueryTargets(s string) ([]azureusageprom.QueryTarget, error) {
	if s == "" {
		return nil, nil
	}

	var rv []azureusageprom.QueryTarget
	groups := strings.Split(s, ",")
	for _, group := range groups {
		parts := strings.Split(group, "|")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid group: %s", group)
		}
		rp, err := azureusageprom.ParseResourceProvider(parts[0])
		if err != nil {
			return nil, err
		}

		target := azureusageprom.QueryTarget{
			ResourceProvider: rp,
			SubscriptionID:   parts[1],
			Location:         parts[2],
		}
		rv = append(rv, target)
	}

	return rv, nil
}

func main() {
	flag.Parse()

	var (
		authorizer   autorest.Authorizer
		cloud        azure.Environment
		queryTargets []azureusageprom.QueryTarget
		err          error
	)

	if *flagAuthWithCLI {
		authorizer, err = auth.NewAuthorizerFromCLI()
		if err != nil {
			glog.Fatal(err)
		}
	} else {
		authorizer, err = auth.NewAuthorizerFromEnvironment()
		if err != nil {
			glog.Fatal(err)
		}
	}

	cloud, err = azure.EnvironmentFromName(*flagEnvironment)
	if err != nil {
		glog.Fatal(err)
	}
	usageLister := usagelister.New(cloud, authorizer)

	queryTargets, err = parseQueryTargets(*flagQueryTargets)
	if err != nil {
		glog.Fatal(err)
	}
	if len(queryTargets) < 1 {
		glog.Fatal("no query targets")
	}

	for _, target := range queryTargets {
		prometheus.MustRegister(collector.New(usageLister, target))
	}

	glog.Infof("azure-usage-prom listening at %s", *flagAddr)
	http.Handle("/metrics", promhttp.Handler())
	glog.Fatal(http.ListenAndServe(*flagAddr, nil))
}
