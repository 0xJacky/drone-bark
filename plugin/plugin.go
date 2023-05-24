// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	// Plugin Config
	ServerUrl  string `envconfig:"PLUGIN_SERVER_URL"`
	BarkDevice string `envconfig:"PLUGIN_BARK_DEVICE"`
	Icon       string `envconfig:"PLUGIN_ICON"`
	Group      string `envconfig:"PLUGIN_BARK_GROUP"`
	BarkLevel  string `envconfig:"PLUGIN_BARK_LEVEL"`
	Sound      string `envconfig:"PLUGIN_BARK_SOUND"`
}

type BarkRequestBody struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	Level     string `json:"level"`
	Icon      string `json:"icon"`
	Group     string `json:"group"`
	Url       string `json:"url"`
	Sound     string `json:"sound"`
	DeviceKey string `json:"device_key"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) (err error) {
	// write code here
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// set default url
	if args.ServerUrl == "" {
		args.ServerUrl = "https://api.day.app/"
	}

	u, err := url.JoinPath(args.ServerUrl, "/push")

	if err != nil {
		return
	}

	shortSHA := ""

	if len(args.Commit.Rev) > 8 {
		shortSHA = args.Commit.Rev[0:8]
	} else {
		shortSHA = args.Commit.Rev
	}

	title := ""

	// https://docs.drone.io/pipeline/environment/reference/drone-build-status/
	if args.Build.Status == "success" {
		title = "Drone CI Run Succeeded"
	} else {
		title = "Drone CI Run Failed"
	}

	body := fmt.Sprintf("Project: %s/%s\nBranch:%s\nCommit: %s",
		args.Repo.Namespace, args.Repo.Name, args.Build.Branch, shortSHA)

	reqBody := BarkRequestBody{
		Title:     title,
		Body:      body,
		Icon:      args.Icon,
		Group:     args.Group,
		Url:       args.Pipeline.Build.Link,
		Level:     args.BarkLevel,
		Sound:     args.Sound,
		DeviceKey: args.BarkDevice,
	}

	reqBodyBytes, err := json.Marshal(reqBody)

	if err != nil {
		return
	}

	logrus.Debugf("%s\n", reqBodyBytes)

	req, err := http.NewRequest("POST", u, bytes.NewReader(reqBodyBytes))

	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)

	logrus.Infof("%s\n", content)

	return
}
