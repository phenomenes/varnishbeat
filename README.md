# Varnishbeat

Current status: **beta release**

Varnishbeat is an [Elastic Beat](https://github.com/elastic/beats) for
varnishlog and varnishstat. It uses [vago](https://github.com/phenomenes/vago)
to read logs and stats from a Varnish Shared Memory file.

## Build Varnishbeat

```
go get -u github.com/phenomenes/varnishbeat
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
