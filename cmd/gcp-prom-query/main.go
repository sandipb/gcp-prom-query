package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sandipb/gcp-prom-query/pkg/prom"
	"gopkg.in/alecthomas/kingpin.v2"
)

const appName = "gcp-prom-query"

var (
	version = "unset"
	commit  = "unset"
	date    = "unset"
)

var (
	showVersion    = kingpin.Flag("version", "Show version").Short('v').Bool()
	debugLevel     = kingpin.Flag("debug", "Debug level logging").Short('d').Bool()
	timeoutSeconds = kingpin.Flag("timeout", "Timeout in seconds for the query").Short('t').Default("10").Int()
	promServer     = kingpin.Flag("prom-api", "URL to API server. Used when gcp-project is not provided.").Short('u').Default("localhost:9090").String()
	project        = kingpin.Flag("gcp-project", "Name of the GCP project. If not given, uses 'prom-api'").Short('p').String()
	token          = kingpin.Flag("gcp-token", "Name of the GCP project. Required if project is given. Can also be provided via env var GCP_ACCESS_TOKEN").
			Short('a').Envar("GCP_ACCESS_TOKEN").String()

	instant = kingpin.Command("instant", "Instant query")
	unixTS  = instant.Flag("now", "Time to run instant query as Unix epoch time").Int64()
	query   = instant.Arg("query", "Promql query").Required().String()
)

const helpText = `Runs query on the gcp prometheus api`

func setup() string {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.Help = helpText
	cmd := kingpin.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	if *debugLevel {
		log.Logger = log.Logger.Level(zerolog.DebugLevel)
	} else {
		log.Logger = log.Logger.Level(zerolog.InfoLevel)
	}

	if *project != "" && *token == "" {
		kingpin.Fatalf("Token is required if project is specified")
	}

	if *project != "" {
		*promServer = fmt.Sprintf("https://monitoring.googleapis.com/v1/projects/%s/location/global/prometheus", *project)
	} else {
		if !strings.HasPrefix(*promServer, "http") {
			*promServer = "http://" + *promServer
		}
	}

	if *unixTS == 0 {
		*unixTS = int64(time.Now().Unix())
	}
	return cmd
}

func printVersion() {
	fmt.Printf("%s %s, commit %s, built %s\n", appName, version, commit, date)
	os.Exit(0)
}

func main() {
	cmd := setup()
	if *showVersion {
		printVersion()
	}
	log.Debug().Msgf("Using prometheus server: %#v", *promServer)

	client, err := prom.GetAPIClient(*promServer, *token)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create prometheus client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeoutSeconds)*time.Second)
	defer cancel()

	switch cmd {
	case "instant":
		err = prom.PrintInstant(ctx, client, *query, time.Unix(*unixTS, 0))
	}
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
