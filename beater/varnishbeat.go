package beater

import (
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	"github.com/phenomenes/vago"
)

type Varnishbeat struct {
	alive   bool
	client  publisher.Client
	varnish *vago.Varnish
}

func New() *Varnishbeat {
	return &Varnishbeat{}
}

func (vb *Varnishbeat) Config(b *beat.Beat) error {
	return nil
}

func (vb *Varnishbeat) Setup(b *beat.Beat) error {
	var err error
	vb.varnish, err = vago.Open("")
	if err != nil {
		logp.Err("%s", err)
		return err
	}
	vb.client = b.Publisher.Connect()
	return nil
}

func (vb *Varnishbeat) Run(b *beat.Beat) error {
	err := vb.exportLog()
	if err != nil {
		logp.Err("%s", err)
	}
	return err
}

func (vb *Varnishbeat) exportLog() error {
	vb.alive = true
	vb.varnish.Log("", vago.RAW, func(vxid uint32, tag, _type, data string) int {
		if vb.alive == false {
			return -1
		}
		switch _type {
		default:
			_type = ""
		case "c":
			_type = "client"
		case "b":
			_type = "backend"
		}
		event := common.MapStr{
			"@timestamp": common.Time(time.Now()),
			"type":       _type,
			"count":      1,
			"vxid":       vxid,
			"tag":        tag,
			"data":       data,
		}
		vb.client.PublishEvent(event)
		return 0
	})
	return nil
}

func (vb *Varnishbeat) Cleanup(b *beat.Beat) error {
	vb.varnish.Close()
	return nil
}

func (vb *Varnishbeat) Stop() {
	vb.alive = false
}
