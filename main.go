package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/phenomenes/varnishbeat/beater"
)

func main() {
	err := beat.Run("varnishbeat", "", beater.New())
	if err != nil {
		os.Exit(1)
	}
}
