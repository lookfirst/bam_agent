package service

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func GetUptime() (time.Duration, error) {
	data, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return time.Duration(0), err
	}
	return parseUptime(string(data))
}

func parseUptime(data string) (time.Duration, error) {
	parts := strings.Fields(data)
	if len(parts) < 2 {
		return time.Duration(0), fmt.Errorf("unexpected content in /proc/uptime: %s", data)
	}

	uptimeInSeconds := strings.Split(parts[0], ".")
	if len(uptimeInSeconds) < 2 {
		return time.Duration(0), fmt.Errorf("no period found %s", parts[0])
	}

	flt, err := strconv.ParseInt(uptimeInSeconds[0], 10, 64)
	if err != nil {
		return time.Duration(0), err
	}

	return time.Duration(flt), nil
}