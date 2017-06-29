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

package vm

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/Huawei/containerops/singular/init_config"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

//	"github.com/digitalocean/godo/context"

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type Client struct {
	// HTTP client used to communicate with the DO API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// Rate contains the current rate limit for the client as determined by the most recent
	// API call.
	// // Rate Rate

	// // // Services used for communicating with the API
	// // Account           AccountService
	// // Actions           ActionsService
	// // Domains           DomainsService
	// // Droplets          DropletsService
	// // DropletActions    DropletActionsService
	// // Images            ImagesService
	// // ImageActions      ImageActionsService
	// // Keys              KeysService
	// // Regions           RegionsService
	// Sizes             SizesService
	// FloatingIPs       FloatingIPsService
	// FloatingIPActions FloatingIPActionsService
	// Snapshots         SnapshotsService
	// Storage           StorageService
	// StorageActions    StorageActionsService
	// Tags              TagsService
	// LoadBalancers     LoadBalancersService
	// Certificates      CertificatesService

	// // Optional function called after every successful request made to the DO APIs
	// onRequestCompleted RequestCompletionCallback
}

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *Client

	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	//client = NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}
func CreateNewVM(vmname string) {

	tokenSource := &TokenSource{
		AccessToken: init_config.TSpet,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	dropletName := vmname

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: init_config.Region, //"sfo2",
		Size:   init_config.MSize,  //"512mb",
		Image: godo.DropletCreateImage{
			Slug: init_config.Slug, //17.04 x64
		},
		SSHKeys: []godo.DropletCreateSSHKey{
			{Fingerprint: init_config.Fingerprint},
		},
		PrivateNetworking: true,
	}

	ctx := context.TODO()
	//newDroplet have  sync issue
	newDroplet, newResponse, err := client.Droplets.Create(ctx, createRequest)
	fmt.Printf("%s\n\n", err, newDroplet, newResponse)

	//dropletIP := newDroplet.Networks.V4
	//newDroplet.PublicIPv4()
	//dropletIP1, _ := newDroplet.PrivateIPv4()
	//newDroplet.Networks.V4String()
	//fmt.Printf("%s\n\n", err, dropletIP, newResponse, dropletIP1)
}

func CreateNode_VM(vmname string) {

	tokenSource := &TokenSource{
		AccessToken: init_config.TSpet,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	dropletName := vmname

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: init_config.Region, //"sfo2",
		Size:   init_config.MSize,  //"512mb",
		Image: godo.DropletCreateImage{
			Slug: init_config.Slug, //17.04 x64
		},
		SSHKeys: []godo.DropletCreateSSHKey{
			{Fingerprint: init_config.Fingerprint},
		},
		PrivateNetworking: true,
	}

	ctx := context.TODO()
	//newDroplet have  sync issue
	newDroplet, newResponse, err := client.Droplets.Create(ctx, createRequest)
	fmt.Printf("%s\n\n", err, newDroplet, newResponse)

	//dropletIP := newDroplet.Networks.V4
	//newDroplet.PublicIPv4()
	//dropletIP1, _ := newDroplet.PrivateIPv4()
	//newDroplet.Networks.V4String()
	//fmt.Printf("%s\n\n", err, dropletIP, newResponse, dropletIP1)
}

func GetListDropletsByTag(tag string) {

	// mux.HandleFunc("/v2/droplets", func(w http.ResponseWriter, r *http.Request) {
	// 	if r.URL.Query().Get("tag_name") != "testing-1" {
	// 		t.Errorf("Droplets.ListByTag did not request with a tag parameter")
	// 	}

	// 	testMethod(t, r, "GET")
	// 	fmt.Fprint(w, `{"droplets": [{"id":1},{"id":2}]}`)
	// })

	//droplets, _, err := client.Droplets.ListByTag(ctx, "testing-1", nil)
	// if err != nil {
	// 	//t.Errorf("Droplets.ListByTag returned error: %v", err)
	// }

	// expected := []Droplet{{ID: 1}, {ID: 2}}
	// if !reflect.DeepEqual(droplets, expected) {
	// 	t.Errorf("Droplets.ListByTag returned %+v, expected %+v", droplets, expected)
	// }
}
