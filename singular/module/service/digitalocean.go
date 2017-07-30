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
	"io/ioutil"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// DigitalOcean struct use for manage create/delete DigitalOcean droplets.
type DigitalOcean struct {
	Token    string         `json:"token" yaml:"token"`
	Region   string         `json:"region" yaml:"region"`
	Size     string         `json:"size" yaml:"size"`
	Image    string         `json:"image" yaml:"image"`
	Droplets map[string]int `json:"droplets,omitempty" yaml:"droplets,omitempty"`

	// Runtime Properties
	client *godo.Client
}

type TokenSource struct {
	AccessToken string
}

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

func (do *DigitalOcean) UpdateSSHKey(publicFile string) error {

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

func (do *DigitalOcean) CreateDroplet(nodes int, fingerprint string) error {
	names := []string{}

	for i := 0; i < nodes; i++ {
		droplet := fmt.Sprintf("singular-node-%d", i+1)
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
		Tags:              []string{"singular", "containerops"},
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
