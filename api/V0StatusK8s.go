package api

import (
	"crypto/tls"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// V0StatusK8s gives out the current status of the K8s services
// example:
// curl -X GET 127.0.0.1:10000/v0/status/k8s
func (e *API) V0StatusK8s(w http.ResponseWriter, r *http.Request) {
	logrus.Debug("HTTP GET V0StatusK8s ")
	vars := mux.Vars(r)

	if vars == nil || !e.CheckAuth(r, w) {
		return
	}

	client := &http.Client{}
	// #nosec G402
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: e.Config.SkipSSL},
	}
	req, _ := http.NewRequest("GET", e.BootstrapProtocol+"://"+e.Config.K3SServerHostname+":"+strconv.Itoa(e.Config.K3SServerContainerPort)+"/api/m3s/bootstrap/v0/status?verbose", nil)
	req.SetBasicAuth(e.Config.BootstrapCredentials.Username, e.Config.BootstrapCredentials.Password)
	req.Close = true
	res, err := client.Do(req)

	if err != nil {
		logrus.Error("StatusK8s: Error 1: ", err, res)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		logrus.Error("StatusK8s: Error Status is not 200")
		return
	}

	content, err := io.ReadAll(res.Body)

	if err != nil {
		logrus.Error("StatusK8s: Error 2: ", err, res)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Api-Service", "v0")
	w.Write(content)
}
