package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"encoding/json"

	"github.com/Sirupsen/logrus"
	dockerapi "github.com/docker/docker/api"
	dockerclient "github.com/docker/engine-api/client"
	"github.com/docker/go-plugins-helpers/authorization"
)

func newPlugin(dockerHost, certPath string, tlsVerify bool) (*authPlugin, error) {
	var httpClient *http.Client
	if certPath != "" {
		tlsc := &tls.Config{}

		cert, err := tls.LoadX509KeyPair(filepath.Join(certPath, "cert.pem"), filepath.Join(certPath, "key.pem"))
		if err != nil {
			return nil, fmt.Errorf("Error loading x509 key pair: %s", err)
		}

		tlsc.Certificates = append(tlsc.Certificates, cert)
		tlsc.InsecureSkipVerify = !tlsVerify
		transport := &http.Transport{
			TLSClientConfig: tlsc,
		}
		httpClient = &http.Client{Transport: transport}

	}

	client, err := dockerclient.NewClient(dockerHost, dockerapi.DefaultVersion, httpClient, nil)
	if err != nil {
		return nil, err
	}
	return &authPlugin{client: client}, nil
}

type authPlugin struct {
	client *dockerclient.Client
}

func (p *authPlugin) AuthZReq(req authorization.Request) authorization.Response {
	if !strings.Contains(req.RequestURI, "/version") {
		logrus.WithFields(logrus.Fields{"REQ": req}).Debug("rEQ")
		if req.RequestBody != nil {
			var dat map[string]interface{}
			if err := json.Unmarshal(req.RequestBody, &dat); err != nil {
				panic(err)
			}
			fmt.Println("\n", dat, "\n")
		}
	}
	if !strings.Contains(req.RequestURI, "/containers/create") {
		return authorization.Response{Allow: true}
	}

	logrus.WithFields(logrus.Fields{"DOCKER_CREATE_BLOCK": isBlocked()}).Debug("Container Create requested")
	if isBlocked() {
		return authorization.Response{Msg: "Create New Container Blocked"}
	}

	return authorization.Response{Allow: true}
}

func (p *authPlugin) AuthZRes(req authorization.Request) authorization.Response {
//	if !strings.Contains(req.RequestURI, "/version") {
//		logrus.WithFields(logrus.Fields{"RES": req}).Debug("rES")
//	}
	if strings.Contains(req.RequestURI, "exec") {
		if strings.Contains(req.RequestURI, "json") {
			logrus.WithFields(logrus.Fields{"RES": req}).Debug("rES")
			if req.ResponseStatusCode == 200 {
				var dat map[string]interface{}
				if err := json.Unmarshal(req.ResponseBody, &dat); err != nil {
					panic(err)
				}
				fmt.Println("\n", dat, "\n")
			}
		}
	}
	return authorization.Response{Allow: true}
}

