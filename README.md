# Varnishbeat

Current status: **beta release**

Varnishbeat is an [Elastic Beat](https://github.com/elastic/beats) for
varnishlog and varnishstat. It uses [vago](https://github.com/phenomenes/vago)
to read logs and stats from a Varnish Shared Memory file.

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

Harvest Varnish logs

```
varnishbeat -log -c varnishlogbeat.yml
```

Harvest Varnish stats

```
varnishbeat -stats -c varnishstatbeat.yml
```
