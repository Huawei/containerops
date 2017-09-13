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

package controller

import (
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strings"

	"github.com/Huawei/containerops/singular/model"
)

type HtmlDeployment struct {
	ID          int64
	InfraName   string
	InfraLogo   string
	InfraLog    template.HTML
	Components  []HtmlComponent
	StatusIcon  string
	StatusColor string
}

type HtmlComponent struct {
	Name          string
	Log           template.HTML
	ImageSrc, Alt string
	Width, Height int
}

func GetHtmlDeploymentList() ([]HtmlDeployment, error) {
	var deployments []model.DeploymentV1
	err := model.DB.Find(&deployments).Error
	if err != nil {
		return nil, err
	}

	var retDeployments []HtmlDeployment
	// Find the corresponding infrastructure & components info
	sig := make(chan error)
	for i := 0; i < len(deployments); i++ {
		go func(deployment *model.DeploymentV1) {
			var infra model.InfraV1
			model.DB.Where("deployment_v1=?", deployment.ID).First(&infra)
			var components []model.ComponentV1
			model.DB.Where("infra_v1=?", infra.ID).Find(&components)

			statusIcon, statusColor := "clear", "red"
			if deployment.Result == true {
				statusIcon, statusColor = "check", "green"
			}
			retDeployment := HtmlDeployment{
				// Name: deployment.Name,
				ID:          deployment.ID,
				InfraName:   infra.Name,
				StatusIcon:  statusIcon,
				StatusColor: statusColor,
			}
			retComponents := []HtmlComponent{}

			for j := 0; j < len(components); j++ {
				retComponents = append(retComponents, convertComponent(&components[j]))
			}

			retDeployment.Components = retComponents

			retDeployments = append(retDeployments, retDeployment)

			sig <- nil
		}(&deployments[i])
	}

	for i := 0; i < len(deployments); i++ {
		if err := <-sig; err != nil {
			return nil, err
		}
	}
	return retDeployments, nil
}

func GetHtmlDeploymentDetail(deploymentID int) *HtmlDeployment {
	var dep model.DeploymentV1
	var infra model.InfraV1
	var comps []model.ComponentV1

	// Get the infra and components
	err := model.DB.Where("id=?", deploymentID).First(&dep).Error
	if err != nil && err.Error() == "record not found" {
		return nil
	}
	model.DB.Where("deployment_v1=?", dep.ID).First(&infra)
	model.DB.Where("infra_v1=?", infra.ID).Find(&comps)

	components := []HtmlComponent{}
	for i := 0; i < len(comps); i++ {
		components = append(components, convertComponent(&comps[i]))
	}
	deployment := HtmlDeployment{
		InfraName:  getInfraName(infra.Name),
		InfraLogo:  getInfraLogo(infra.Name),
		InfraLog:   convertToBr(infra.Log),
		Components: components,
	}

	return &deployment
}

func StringifyComponentsNames(args ...interface{}) (string, error) {
	v := reflect.ValueOf(args)
	numArgs := v.Len()
	if numArgs != 1 {
		fmt.Fprintf(os.Stderr, "component_names function expect 1 argument, but got %d", numArgs)
		return "", fmt.Errorf("Expect 1 argument")
	}

	if components, ok := v.Index(0).Interface().([]HtmlComponent); !ok {
		fmt.Fprintln(os.Stderr, "function component_names receives an argument which is not []Component")
		return "", fmt.Errorf("Argument is not []Component!")
	} else {
		s := ""
		for i := 0; i < len(components); i++ {
			c := components[i]
			s += c.Name + ", "
		}
		return s[:len(s)-1], nil
	}
}

// Convert \n to <br />
func convertToBr(src string) template.HTML {
	replaced := strings.Replace(src, "\n", "\n<br />", -1)
	return template.HTML(replaced)
	// return strings.Join(strings.Split(src, "\n"), "<br />")
}

func convertComponent(input *model.ComponentV1) HtmlComponent {
	c := HtmlComponent{
		Name: input.Binary,
		Log:  convertToBr(input.Log),
	}
	// Got the size and icon
	switch c.Name {
	case "kubernetes":
		c.Width = 40
		c.Height = 40
		c.ImageSrc = "./public/icons/kubernetes.svg"
		break
	case "docker":
		c.Width = 40
		c.Height = 40
		c.ImageSrc = "./public/icons/docker.svg"
		break
	case "flannel":
		c.Width = 40
		c.Height = 40
		c.ImageSrc = "./public/icons/flannel.svg"
		break
	case "etcd":
		c.Width = 40
		c.Height = 40
		c.ImageSrc = "./public/icons/etcd.svg"
		break
	default:
		break
	}

	return c
}

var infraNames map[string]string = map[string]string{
	"digital_ocean": "Digital Ocean",
	"gke":           "Google Container Engine",
	"aws":           "AWS EC2",
	"azure":         "Microsoft Azure",
}
var infraLogos map[string]string = map[string]string{
	"digital_ocean": "./public/icons/digital-ocean.svg",
	"gke":           "./public/icons/google-cloud.svg",
	"aws":           "./public/icons/aws-ec2.svg",
	"azure":         "./public/icons/azure.svg",
}

func getInfraName(key string) string {
	name := infraNames[key]
	if name == "" {
		name = "Unsupported Cloud Service Provider"
	}
	return name
}

func getInfraLogo(name string) string {
	logo := infraLogos[name]
	if logo == "" {
		logo = "./public/icons/google-cloud.svg"
	}
	return logo
}
