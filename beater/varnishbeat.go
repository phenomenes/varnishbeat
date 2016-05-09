package beater

import (
	"strings"
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
	tx := make(map[string]interface{})
	vb.varnish.Log("", vago.RAW, func(vxid uint32, tag, _type, data string) int {
		if vb.alive == false {
			return -1
		}
		switch _type {
		default:
			_type = "ping"
		case "c":
			_type = "client"
		case "b":
			_type = "backend"
		}
		if tag == "ReqHeader" || tag == "BereqHeader" || tag == "BerespHeader" || tag == "ObjHeader" {
			header := strings.SplitN(data, ":", 2)
			logp.Info(tag, data)
			k := header[0]
			v := header[1]
			if _, ok := tx[tag]; ok {
				tx[tag].(map[string]interface{})[k] = v
			} else {
				tx[tag] = map[string]interface{}{k: v}
			}
		} else {
			tx[tag] = data
		}
		if tag == "End" {
			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"count":      1,
				"type":       _type,
				"vxid":       vxid,
				"tx":         tx,
			}
			vb.client.PublishEvent(event)
			tx = nil
			tx = make(map[string]interface{})
		}
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
