/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const (
	DORootUser = "root"
)

// DigitalOcean struct use for manage create/delete DigitalOcean droplets.
type DigitalOcean struct {
	Token    string         `json:"token" yaml:"token"`
	Region   string         `json:"region" yaml:"region"`
	Size     string         `json:"size" yaml:"size"`
	Image    string         `json:"image" yaml:"image"`
	Droplets map[string]int `json:"droplets,omitempty" yaml:"droplets,omitempty"`
	Logs     []string       `json:"logs,omitempty" yaml:"logs,omitempty"`

	// Runtime Properties
	client *godo.Client
}

//WriteLog implement Logger interface.
func (d *DigitalOcean) WriteLog(log string, writer io.Writer, output bool) error {
	d.Logs = append(d.Logs, log)

	if output == true {
		if _, err := io.WriteString(writer, fmt.Sprintf("%s\n", log)); err != nil {
			return err
		}
	}

	return nil
}

// TokenSource is access token of DigitalOcean
type TokenSource struct {
	AccessToken string
}

//
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func (do *DigitalOcean) InitClient() error {
	tokenSource := &TokenSource{
		AccessToken: do.Token,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	do.client = godo.NewClient(oauthClient)

	return nil
}

// TODO Customize SSH key name.
func (do *DigitalOcean) UploadSSHKey(publicFile string) error {
	if public, err := ioutil.ReadFile(publicFile); err != nil {
		return err
	} else {
		createRequest := &godo.KeyCreateRequest{
			Name:      "singular",
			PublicKey: string(public),
		}

		ctx := context.TODO()
		if _, _, err := do.client.Keys.Create(ctx, createRequest); err != nil {
			return err
		}
	}

	return nil
}

func (do *DigitalOcean) CreateDroplet(nodes int, fingerprint, name string, tags []string) error {
	names := []string{}

	for i := 0; i < nodes; i++ {
		droplet := fmt.Sprintf("%s-node-%d", name, i+1)
		names = append(names, droplet)
	}

	sshFingerprint := godo.DropletCreateSSHKey{
		Fingerprint: fingerprint,
	}

	createRequest := &godo.DropletMultiCreateRequest{
		Names:  names,
		Region: do.Region,
		Size:   do.Size,
		Image: godo.DropletCreateImage{
			Slug: do.Image,
		},
		SSHKeys:           []godo.DropletCreateSSHKey{sshFingerprint},
		Backups:           false,
		IPv6:              false,
		PrivateNetworking: false,
		Monitoring:        true,
		Tags:              tags,
		UserData:          "",
	}

	ctx := context.TODO()

	droplets, _, err := do.client.Droplets.CreateMultiple(ctx, createRequest)

	if err != nil {
		fmt.Printf("Something bad happened: %s\n\n", err)
		return err
	}

	time.Sleep(10 * time.Second)

	do.Droplets = map[string]int{}
	for {
		for _, value := range droplets {
			ctx := context.TODO()
			droplet, _, err := do.client.Droplets.Get(ctx, value.ID)

			if err != nil {
				fmt.Printf("Something bad happened: %s\n\n", err)
				return err
			}

			if len(droplet.Networks.V4) > 0 {
				v4 := droplet.Networks.V4[0]
				do.Droplets[v4.IPAddress] = droplet.ID
			}
		}

		if len(do.Droplets) == nodes {
			break
		}

		time.Sleep(5 * time.Second)
	}

	return nil
}
