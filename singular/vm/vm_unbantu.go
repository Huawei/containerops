package vm

import (
	"context"
	"fmt"

	"github.com/Huawei/containerops/singular/init_config"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
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
	dropletIP := newDroplet.Networks.V4
	//newDroplet.PublicIPv4()
	dropletIP1, _ := newDroplet.PrivateIPv4()
	//newDroplet.Networks.V4String()
	fmt.Printf("%s\n\n", err, dropletIP, newResponse, dropletIP1)
}
