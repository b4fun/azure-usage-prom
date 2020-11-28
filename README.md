# azure-usage-prom

Convert Azure resource usage to prometheus metrics.

| Resource | Link |
|:----|:----|
| Reference | [![API](https://godoc.org/github.com/b4fun/azure-usage-prom?status.svg)](https://pkg.go.dev/github.com/b4fun/azure-usage-prom?tab=overview) |
| Docker Image | [![Docker Build Status](https://img.shields.io/docker/build/b4fun/azure-usage-prom)](https://hub.docker.com/r/b4fun/azure-usage-prom) |

## Usage

```
$ ./azure-usage-prom \
    -query-targets "microsoft.compute|0000000-000-0000-0000-0000000000|eastus,microsoft.network|0000000-000-0000-0000-0000000000|eastus"
I1128 14:00:12.316210   23445 main.go:98] azure-usage-prom listening at :8080
$ curl -v http://localhost:8080/metrics
```

## LICENSE

MIT
