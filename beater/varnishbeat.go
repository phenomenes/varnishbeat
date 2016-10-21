package beater

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/phenomenes/vago"
	"github.com/phenomenes/varnishbeat/config"
)

type Varnishbeat struct {
	done    chan struct{}
	config  config.Config
	client  publisher.Client
	varnish *vago.Varnish
}

var logFlag, statsFlag bool
var varnishDirectoryFlag string

func init() {
	flag.BoolVar(&logFlag, "log", false, "Read data from varnishlog")
	flag.BoolVar(&statsFlag, "stats", false, "Read data from varnishstat")
	flag.StringVar(&varnishDirectoryFlag, "directory", "", "Directory including the name if given to where varnish writes vsm and vcls.")
}

// New creates a beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	vb := &Varnishbeat{
		done:   make(chan struct{}),
		config: config,
	}
	return vb, nil
}

func (vb *Varnishbeat) Run(b *beat.Beat) error {
	var err error

	vb.varnish, err = vago.Open(varnishDirectoryFlag)
	if err != nil {
		return err
	}

	vb.client = b.Publisher.Connect()

	logp.Info("varnishbeat is running! Hit CTRL-C to stop it.")
	if logFlag {
		err := vb.harvestLog()
		if err != nil {
			logp.Err("%s", err)
		}

	} else {
		ticker := time.NewTicker(vb.config.Period)
		for {
			select {
			case <-vb.done:
				return nil
			case <-ticker.C:
			}
			event, err := vb.harvestStats()
			if err != nil {
				logp.Err("%s", err)
				break
			}
			logp.Info("Event sent")
			vb.client.PublishEvent(event)
		}
	}
	return err
}

func (vb *Varnishbeat) harvestStats() (common.MapStr, error) {
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

func (vb *Varnishbeat) harvestLog() error {
	tx := make(common.MapStr)
	vb.varnish.Log("", vago.REQ, func(vxid uint32, tag, _type, data string) int {
		switch _type {
		case "c":
			_type = "client"
		case "b":
			_type = "backend"
		default:
			return 0
		}
		switch tag {
		case "BereqHeader", "BerespHeader", "ObjHeader", "ReqHeader", "RespHeader":
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

func (vb *Varnishbeat) Stop() {
	vb.varnish.Stop()
	vb.varnish.Close()
	close(vb.done)
}
