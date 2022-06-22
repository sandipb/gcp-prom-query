package prom

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/config"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog/log"
)

func GetAPIClient(url string, jwt string) (v1.API, error) {
	cfg := api.Config{
		Address: url,
	}

	if jwt != "" {
		cfg.RoundTripper = config.NewAuthorizationCredentialsRoundTripper("Bearer", config.Secret(jwt), api.DefaultRoundTripper)
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return v1.NewAPI(client), nil
}

type MetricEntry struct {
	Metric    string
	Value     float64
	Timestamp int64
}

func PrintInstant(ctx context.Context, client v1.API, query string, t time.Time) error {
	result, warnings, err := client.Query(ctx, query, t)
	if err != nil {
		if v1err, ok := err.(*v1.Error); ok {
			return fmt.Errorf("Query error: %v: msg='%s', detail='%s' ", v1err.Type, v1err.Msg, v1err.Detail)
		}
		return fmt.Errorf("Query error: %w", err)
	}
	for _, w := range warnings {
		log.Warn().Msgf("Query returned warning: %s", w)
	}
	if result.Type() != model.ValVector {
		return fmt.Errorf("Unexpected non-vector result type %s received for query: %s", result.Type(), query)
	}
	results := result.(model.Vector)
	log.Debug().Msgf("%d metrics received", len(results))
	out := []MetricEntry{}
	for _, entry := range results {
		out = append(out, MetricEntry{
			Metric:    entry.Metric.String(),
			Value:     float64(entry.Value),
			Timestamp: entry.Timestamp.Unix(),
		})
	}

	table := tablewriter.NewWriter(os.Stdout)
	// table.SetColMinWidth(0, 40)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"metric", "Value", "Timestamp"})
	for _, e := range out {
		table.Append([]string{e.Metric, strconv.FormatFloat(e.Value, 'f', 1, 64), strconv.FormatInt(e.Timestamp, 10)})
	}
	table.Render()
	return nil
}
