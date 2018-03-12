package service

import (
	"log"

	"github.com/blockassets/cgminer_client"
)

//
//	Interface to decouple cgminer_client implementation. Helps in testing.
//
type Miner interface {
	Devs() (*[]cgminer_client.Dev, error)
	Quit() error
}

// Helper functions to derive facts from the miner
// They are unit tested as part of monitor_accepted_test and monitor_high_temp_test
func GetAccepted(miner Miner) (int64, error) {
	devs, err := miner.Devs()
	if err != nil {
		log.Printf("Error getting accepted shares: %v", err)
		return 0, err
	}
	accepted := int64(0)
	for _, dev := range *devs {
		accepted += dev.Accepted
	}
	return accepted, nil
}

func GetTemp(miner Miner) (float64, error) {
	devs, err := miner.Devs()
	if err != nil {
		log.Printf("Error getting Temp: %v", err)
		return 0, err
	}
	// find top temp...
	highestTemp := float64(0)
	for _, dev := range *devs {
		if dev.Temperature > highestTemp {
			highestTemp = dev.Temperature
		}
	}
	return highestTemp, nil
}
