package api

import (
	"encoding/json"
	"net/http"

	cfg "github.com/AVENTER-UG/mesos-m3s/types"
	"github.com/sirupsen/logrus"
)

// V0GetCountK3SAgent will write out the current count of Kubernetes Agents and the expected one.
// example:
// curl -X GET http://user:password@127.0.0.1:10000/v0/agent/scale'
func (e *API) V0GetCountK3SAgent(w http.ResponseWriter, r *http.Request) {
	if !e.CheckAuth(r, w) {
		return
	}

	var count cfg.Count

	count.Scale = e.Config.K3SAgentMax
	count.Running = e.CountRedisKey(e.Framework.FrameworkName + ":agent:*")

	d, err := json.Marshal(count)

	if err != nil {
		logrus.Error("V0GetCountK3SAgent Error 1: ", err)
		return
	}

	logrus.Debug("HTTP GET V0GetCountK3SAgent")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	w.Write(d)
}
