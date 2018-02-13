package controller

import (
	"net/http"

	"github.com/blockassets/bam_agent/service"
)

// command to update the ip address for
// pools
// POST
//Just POST this file:
//{
//"pool1":"",
//"pool2":"",
//"pool3":""
//}
// eg { "pool1":"111.2.3.4", "pool2":"112.3.4.5", "pool3":"113.4.5.6"}
// and we update the conf.default file on the miner

const defConfigPath = "usr/app/conf.default"

// Implements Controller interface
type PutPoolsCtrl struct {
}

func (c PutPoolsCtrl) build() *Controller {
	return &Controller{
		Methods: []string{http.MethodPut},
		Path:    "/config/pools",
		Handler: c.makeHandler(),
	}
}

func (c PutPoolsCtrl) makeHandler() http.HandlerFunc {
	return makeJsonHandler(
		func(w http.ResponseWriter, r *http.Request) {

			bamStat := BAMStatus{"OK", nil}
			httpStat := http.StatusOK

			cmds := service.Command{}
			err := cmds.UpdatePools(r.Body, defConfigPath)
			if err != nil {
				httpStat = http.StatusBadGateway
				bamStat = BAMStatus{"Error", err}
			}
			w.WriteHeader(httpStat)
			resp, _ := json.Marshal(bamStat)
			w.Write(resp)

		})
}
