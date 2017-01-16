package main

import (
	"flag"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/authorization"
)

const (
	defaultDockerHost = "unix:///var/run/docker.sock"
	pluginSocket      = "/run/docker/plugins/docker-auth-plugin.sock"
)

var (
	flDockerHost = flag.String("host", defaultDockerHost, "Specifies the host where to contact the docker daemon")
	flCertPath   = flag.String("cert-path", "", "Certificates path to connect to Docker (cert.pem, key.pem)")
	flTLSVerify  = flag.Bool("tls-verify", false, "Whether to verify certificates or not")
	flDebug      = flag.Bool("debug", false, "Print debug message")
	flAddr       = flag.String("addr", "0.0.0.0:8080", "[HOSTNAME:PORT] for running web server")
)

func main() {
	flag.Parse()

	if *flDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.Debug("Docker auth plugin started")

	http.HandleFunc("/block", blockHandler)
	http.HandleFunc("/unblock", unblockHandler)
	http.HandleFunc("/status", statusHandler)
	go func() {
		logrus.Debugf("Control server Running on %s", *flAddr)
		if err := http.ListenAndServe(*flAddr, nil); err != nil {
			logrus.Fatal(err)
		}
	}()

	plugin, err := newPlugin(*flDockerHost, *flCertPath, *flTLSVerify)
	if err != nil {
		logrus.Fatal(err)
	}

	h := authorization.NewHandler(plugin)

	logrus.Debugf("Plugin  Running on %s", pluginSocket)
	if err := h.ServeUnix(pluginSocket, 0); err != nil {
		logrus.Fatal(err)
	}
}
