package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/phenomenes/varnishbeat/beater"
)

func main() {
	if err := beat.Run("varnishbeat", "", beater.New); err != nil {
		os.Exit(1)
	}
}
