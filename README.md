# azure-usage-prom

Convert Azure resource usage to prometheus metrics.

## Usage

```
$ ./azure-usage-prom \
    -query-targets "microsoft.compute|0000000-000-0000-0000-0000000000|eastus,microsoft.network|0000000-000-0000-0000-0000000000|eastus"
```

## LICENSE

MIT
