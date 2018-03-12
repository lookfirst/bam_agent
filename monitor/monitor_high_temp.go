package monitor

import (
	"log"
	"time"

	"github.com/blockassets/bam_agent/service"
)

type HighTempConfig struct {
	Enabled  bool          `json:"enabled"`
	Period   time.Duration `json:"period"`
	HighTemp float64       `json:"highTemp"`
}

// Implements the Monitor interface
type HighTempMonitor struct {
	*Context
	config     *HighTempConfig
	miner      service.Miner
	OnHighTemp func()
}

func newHighTempMonitor(context *Context, config *HighTempConfig, miner service.Miner, OnHighTemp func()) Monitor {
	return &HighTempMonitor{
		Context:    context,
		config:     config,
		miner:      miner,
		OnHighTemp: OnHighTemp,
	}
}

func (mon *HighTempMonitor) Start() error {
	if mon.config.Enabled {
		log.Printf("HighTempMonitor: Checking for temp every %v\n", mon.config.Period)
		go mon.makeTickerFunc(func() {
			overTemp, err := mon.checkHighTemp()
			if err == nil && overTemp {
				mon.OnHighTemp()
			}
		}, mon.config.Period)()
	} else {
		log.Println("HighTempMonitor: Not enabled")
	}
	return nil
}

func (mon *HighTempMonitor) checkHighTemp() (bool, error) {
	temp, err := service.GetTemp(mon.miner)
	if err != nil {
		return false, err
	}
	if temp < mon.config.HighTemp {
		return false, nil
	}
	return true, nil
}
