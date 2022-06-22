# gcp-prom-query - Query GCP Prometheus API

## Usage

```shell-session
$ ./gcp-prom-query -h
usage: gcp-prom-query [<flags>] <command> [<args> ...]

Runs query on the gcp prometheus api

Flags:
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
  -v, --version                  Show version
  -d, --debug                    Debug level logging
  -t, --timeout=10               Timeout in seconds for the query
  -u, --prom-api="localhost:9090"  
                                 URL to API server. Used when gcp-project is not provided.
  -p, --gcp-project=GCP-PROJECT  Name of the GCP project. If not given, uses 'prom-api'
  -a, --gcp-token=GCP-TOKEN      Name of the GCP project. Required if project is given. Can also be provided via env var GCP_ACCESS_TOKEN

Commands:
  help [<command>...]
    Show help.

  instant [<flags>] <query>
    Instant query

$ ./gcp-prom-query instant -h
usage: gcp-prom-query instant [<flags>] <query>

Instant query

Flags:
...
      --now=NOW                  Time to run instant query as Unix epoch time

Args:
  <query>  Promql query
```

## Sample output

```shell-session
$ ./gcp-prom-query -p some-gcp-project instant 'sum by(job)(scrape_samples_scraped)'
+--------------------+-------+------------+
|       METRIC       | VALUE | TIMESTAMP  |
+--------------------+-------+------------+
| {job="prometheus"} | 446.0 | 1655917668 |
+--------------------+-------+------------+
```
