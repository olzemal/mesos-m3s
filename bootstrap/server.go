package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os/exec"

	"github.com/AVENTER-UG/util"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var MinVersion string
var DashboardInstalled bool

// Commands is the main function of this package
func Commands() *mux.Router {
	// Connect with database

	rtr := mux.NewRouter()
	rtr.HandleFunc("/versions", APIVersions).Methods("GET")
	rtr.HandleFunc("/status", APIHealth).Methods("GET")
	rtr.HandleFunc("/api/k3s/v0/config", APIGetKubeConfig).Methods("GET")

	return rtr
}

// APIVersions give out a list of Versions
func APIVersions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Api-Service", "-")
	w.Write([]byte("/api/k3s/v0"))
}

// APIGetKubeConfig
func APIGetKubeConfig(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("/mnt/mesos/sandbox/kubeconfig.yaml")
	if err != nil {
		logrus.Error("Error reading file:", err)
		w.Write([]byte("Error reading kubeconfig.yaml"))
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Api-Service", "v0")

		w.Write(content)
	}
}

// APIHealth give out the status of the kubernetes server
func APIHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Api-Service", "v0")

	logrus.Debug("Health Check")

	// check if the kubernetes server is working
	stdout, err := exec.Command("kubectl", "get", "--raw=/livez/ping").Output()

	if err != nil {
		logrus.Error("Health to Kubernetes Server: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if string(stdout) == "ok" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))

		// if kubernetes server is running and the dashboard is not installed, then do it
		if !DashboardInstalled {
			deployDashboard()
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// deployDashboard will deploy the kubernetes dashboard
// if the server is in the running state
func deployDashboard() {
	err := exec.Command("kubectl", "apply", "-f", "/mnt/mesos/sandbox/dashboard_auth.yaml").Run()
	logrus.Info("Install Kubernetes Dashboard")

	if err != nil {
		logrus.Error("Install Kubernetes Dashboard Auth: ", err)
		return
	}

	err = exec.Command("kubectl", "apply", "-f", "/mnt/mesos/sandbox/dashboard.yaml").Run()

	if err != nil {
		logrus.Error("Install Kubernetes Dashboard: ", err)
		return
	}

	logrus.Info("Install Kubernetes Dashboard: Done")
	DashboardInstalled = true
}

func main() {
	util.SetLogging("INFO", false, "GO-K3S-API")

	bind := flag.String("bind", "0.0.0.0", "The IP address to bind")
	port := flag.String("port", "10422", "The port to listen")

	logrus.Println("GO-K3S-API build"+MinVersion, *bind, *port)

	DashboardInstalled = false

	http.Handle("/", Commands())

	if err := http.ListenAndServe(*bind+":"+*port, nil); err != nil {
		logrus.Fatalln("ListenAndServe: ", err)
	}

}
