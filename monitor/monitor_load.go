package monitor

import (
	"github.com/blockassets/bam_agent/service"
	"log"
	"time"
)

func monitorLoad(sr statRetriever, interval time.Duration) {
	for {
		high, err := checkLoadAvg(sr)
		if (err == nil) && high {
			cmds := service.Command{}
			cmds.Reboot()
		}
		time.Sleep(interval)
	}
}

func checkLoadAvg(sr statRetriever) (bool, error) {
	loads, err := sr.getLoad()
	high := false
	if (err == nil) && (loads.fiveMinAvg > 5.0) {
		high = true
	} else {
		log.Println("Monitor load error:", err)
	}
	return high, err
}
