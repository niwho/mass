package utils

import (
	"github.com/robfig/cron"
)

var TickerInstance *Ticker

type Ticker struct {
	cronIn *cron.Cron
}

func InitTicker() {
	TickerInstance = &Ticker{
		cronIn: cron.New(),
	}
	TickerInstance.cronIn.Start()
}

func (tk *Ticker) AddFunc(spec string, cmd func()) error {
	return tk.cronIn.AddFunc(spec, cmd)
}
