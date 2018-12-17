mackerel-plugin-mirakurun
=======================

mirakurun custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
mackerel-plugin-mirakurun -host=<hostname or ip> -port=<port> [-metric-key-prefix=<metric-key-prefix> [-tempfile=<tempfile>]
```

## Example of mackerel-agent.conf

```
[plugin.metrics.mirakurun]
command = "/path/to/mackerel-plugin-mirakurun -host=localhost -port=40772"
```
