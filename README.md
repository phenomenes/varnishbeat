# Varnishbeat

Current status: **beta release**

Varnishbeat is an [Elastic Beat](https://github.com/elastic/beats) for
varnishlog and varnishstat. It uses [vago](https://github.com/phenomenes/vago)
to read logs and stats from Varnish Shared Memory.

## Requirements

To build this package you will need:
- pkg-config
- libvarnishapi-dev >= 4.0.0

You will also need to set PKG_CONFIG_PATH to the directory where varnishapi.pc
is located before running `go get`. For example:
```
export PKG_CONFIG_PATH=/usr/local/lib/pkgconfig
```

## Build Varnishbeat

```
go get github.com/phenomenes/varnishbeat
```

## Run Varnishbeat

Currently `varnishbeat` operates in two modes: `stats` or `logs`.
Each mode needs its own ES index.

You can run `varnishbeat` to collect `stats` and `logs` but you'll need
to execute the binary with different configuration files. For example:


* To collect Varnish logs, add these lines to the configuration:

```
# varnishlogbeat.yml

varnishbeat:
  log: true

output:
  elasticsearch:
    index: "varnishlogbeat"
```

* To collect Varnish stats: 
 
```
# varnishstatsbeat.yml

varnishbeat:
  stats: true

output:
  elasticsearch:
    index: "varnishstatsbeat"
```

Run `varnishbeat`

```
$ varnishbeat -c varnishlogbeat.yml
$ varnishbeat -c varnishstatsbeat.yml
```
