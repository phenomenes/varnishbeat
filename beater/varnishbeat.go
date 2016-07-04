package beater

import (
	"flag"
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
	done    chan uint
	period  time.Duration
	varnish *vago.Varnish
}

var logFlag, statsFlag bool

func init() {
	flag.BoolVar(&logFlag, "log", false, "Read data from varnishlog")
	flag.BoolVar(&statsFlag, "stats", false, "Read data from varnishstat")
}

func New() *Varnishbeat {
	return &Varnishbeat{}
}

func (vb *Varnishbeat) HandleFlags(b *beat.Beat) error {
	return nil
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
	vb.alive = true
	var err error
	if logFlag {
		err := vb.exportLog()
		if err != nil {
			logp.Err("%s", err)
		}
	} else {
		vb.period = 10 * time.Second
		ticker := time.NewTicker(vb.period)
		defer ticker.Stop()

		for vb.alive {
			select {
			case <-vb.done:
				return nil
			case <-ticker.C:
			}
			timerStart := time.Now()
			event, err := vb.exportStats()
			if err != nil {
				logp.Err("%s", err)
				break
			}
			vb.client.PublishEvent(event)
			timerEnd := time.Now()
			duration := timerEnd.Sub(timerStart)
			if duration.Nanoseconds() > vb.period.Nanoseconds() {
				logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
			}
		}
	}
	return err
}

func (vb *Varnishbeat) exportStats() (common.MapStr, error) {
	stats := make(common.MapStr)
	for k, v := range vb.varnish.Stats() {
		k1 := strings.Replace(k, ".", "_", -1)
		stats[k1] = v
	}
	event := common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "stats",
		"stats":      stats,
	}
	return event, nil
}

func (vb *Varnishbeat) exportLog() error {
	tx := make(common.MapStr)
	vb.varnish.Log("", vago.RAW, func(vxid uint32, tag, _type, data string) int {
		if vb.alive == false {
			return -1
		}
		switch _type {
		case "c":
			_type = "client"
		case "b":
			_type = "backend"
		default:
			_type = "ping"
		}
		switch tag {
		case "ReqHeader", "BereqHeader", "BerespHeader", "ObjHeader":
			header := strings.SplitN(data, ": ", 2)
			k := header[0]
			v := header[1]
			if _, ok := tx[tag]; ok {
				tx[tag].(common.MapStr)[k] = v
			} else {
				tx[tag] = common.MapStr{k: v}
			}
		case "End":
			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"count":      1,
				"type":       _type,
				"vxid":       vxid,
				"tx":         tx,
			}
			vb.client.PublishEvent(event)
			// destroy and re-create the map
			tx = nil
			tx = make(common.MapStr)
		default:
			tx[tag] = data
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
