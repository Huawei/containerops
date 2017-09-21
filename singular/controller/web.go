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
	"regexp"
	"strings"

	"github.com/Huawei/containerops/singular/model"
)

type HtmlSingular struct {
	Namespace   string
	Repository  string
	Name        string
	Result      bool
	Deployments []HtmlDeployment
}

type HtmlDeployment struct {
	SingularNamespace  string
	SingularRepository string
	SingularName       string
	CreatedTime        string
	ID                 int64
	Infras             []HtmlInfra
	Data               template.HTML // The YAML file content
	CA                 template.HTML // The YAML file content
	Log                template.HTML // The log

	Version int64
	Tag     string

	InfraName   string
	InfraLogo   string
	InfraLog    template.HTML
	StatusIcon  string
	StatusColor string
}

type HtmlInfra struct {
	ID         int64
	Name       string
	Version    string
	Logo       string
	Log        string
	Components []HtmlComponent
}

type HtmlComponent struct {
	// Log           template.HTML
	Name          string
	Log           string
	ImageSrc, Alt string
	Width, Height int
}

type HtmlInfraTitle struct {
	Name string
	Logo string
}

func GetHtmlDeploymentList() ([]HtmlDeployment, error) {
	var deployments []model.DeploymentV1
	err := model.DB.Order("create_at desc").Find(&deployments).Error
	if err != nil {
		return nil, err
	}

	htmlDeployments := []HtmlDeployment{}
	for i := 0; i < len(deployments); i++ {
		deployment := deployments[i]

		statusIcon, statusColor := "clear", "red"
		if deployment.Result == true {
			statusIcon, statusColor = "check", "green"
		}
		htmlDeployment := HtmlDeployment{
			// Name: deployment.Name,
			CreatedTime: deployment.CreatedAt.Format("2006-01-02 15:04:05"),
			ID:          deployment.ID,
			StatusIcon:  statusIcon,
			StatusColor: statusColor,
		}
		htmlDeployments = append(htmlDeployments, htmlDeployment)
	}
	// Find the corresponding infrastructure & components info
	sig := make(chan error)
	for i := 0; i < len(deployments); i++ {
		go func(htmlDeployment *HtmlDeployment) {
			var infras []model.InfraV1
			err := model.DB.Where("deployment_v1=?", htmlDeployment.ID).Find(&infras).Error
			if err != nil {
				sig <- err
				return
			}

			htmlInfras := []HtmlInfra{}
			for i := 0; i < len(InfraOrder); i++ {
				found := false
				var htmlInfra HtmlInfra
				for j := 0; j < len(infras); j++ {
					infra := infras[j]
					if infra.Name == InfraOrder[i] {
						found = true
						htmlInfra = HtmlInfra{
							// ID:      infra.ID,
							Logo:    getInfraLogo(infra.Name),
							Name:    infra.Name,
							Version: getInfraSemver(infra.Name, infra.Version),
						}
						break
					}
				}

				if !found {
					htmlInfra = HtmlInfra{
						Name:    InfraOrder[i],
						Version: "N/A",
					}
				}
				htmlInfras = append(htmlInfras, htmlInfra)
			}
			htmlDeployment.Infras = htmlInfras
			sig <- nil

			// htmlDeployments = append(htmlDeployments, htmlDeployment)
		}(&htmlDeployments[i])
	}

	for i := 0; i < len(deployments); i++ {
		if err := <-sig; err != nil {
			return nil, err
		}
	}
	return htmlDeployments, nil
}

func GetHtmlDeploymentDetail(deploymentID int) *HtmlDeployment {
	var singular model.SingularV1
	var deployment model.DeploymentV1
	var infras []model.InfraV1
	htmlInfras := []HtmlInfra{}

	// Get the infra and components
	err := model.DB.Where("id=?", deploymentID).First(&deployment).Error
	if err != nil && err.Error() == "record not found" {
		return nil
	}

	model.DB.Where("id=?", deployment.SingularV1).First(&singular)

	model.DB.Where("deployment_v1=?", deployment.ID).Find(&infras)
	for i := 0; i < len(infras); i++ {
		infra := infras[i]
		var components []model.ComponentV1
		model.DB.Where("infra_v1=?", infra.ID).Find(&components)

		htmlComponents := []HtmlComponent{}
		for j := 0; j < len(components); j++ {
			component := components[j]
			htmlComponents = append(htmlComponents, HtmlComponent{
				Name: component.Name,
				Log:  component.Log,
			})
		}

		htmlInfra := convertInfra(&infra)
		htmlInfra.Components = htmlComponents
		htmlInfras = append(htmlInfras, htmlInfra)
	}

	htmlDeployment := HtmlDeployment{
		Version:            deployment.Version,
		Tag:                deployment.Tag,
		Log:                convertToBr(deployment.Log),
		Data:               convertToBr(deployment.Data),
		CA:                 convertToBr(deployment.CA),
		Infras:             htmlInfras,
		SingularName:       singular.Name,
		SingularNamespace:  singular.Namespace,
		SingularRepository: singular.Repository,
	}

	return &htmlDeployment
}

func StringifyComponentsNames(args ...interface{}) (string, error) {
	v := reflect.ValueOf(args)
	numArgs := v.Len()
	if numArgs != 1 {
		fmt.Fprintf(os.Stderr, "component_names function expect 1 argument, but got %d", numArgs)
		return "", fmt.Errorf("expect 1 argument")
	}

	if components, ok := v.Index(0).Interface().([]HtmlComponent); !ok {
		fmt.Fprintln(os.Stderr, "function component_names receives an argument which is not []Component")
		return "", fmt.Errorf("argument is not []Component")
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

func convertInfra(input *model.InfraV1) HtmlInfra {
	infra := HtmlInfra{
		Name: input.Name,
		// Log:  convertToBr(input.Log),
		Log: input.Log,
	}
	// Got the size and icon
	switch infra.Name {
	case "kubernetes":
		// infra.Width = 40
		// infra.Height = 40
		infra.Logo = "./public/icons/kubernetes.svg"
		break
	case "docker":
		// infra.Width = 40
		// infra.Height = 40
		infra.Logo = "./public/icons/docker.svg"
		break
	case "flannel":
		// infra.Width = 40
		// infra.Height = 40
		infra.Logo = "./public/icons/flannel.svg"
		break
	case "etcd":
		// infra.Width = 40
		// infra.Height = 40
		infra.Logo = "./public/icons/etcd.svg"
		break
	default:
		break
	}

	return infra
}

func getInfraSemver(name, version string) string {
	re := regexp.MustCompile(fmt.Sprintf("%s-", name))
	locations := re.FindStringIndex(version)
	if locations == nil {
		return "N/A"
	}
	return version[locations[1]:]
}

var infraNames map[string]string = map[string]string{
	"kubernetes":  "Kubernetes",
	"etcd":        "etcd",
	"flannel":     "Flannel",
	"docker":      "Docker",
	"prometheus":  "Prometheus",
	"opentracing": "OpenTracking",
	"fluentd":     "Fluentd",
	"linerd":      "linkerd",
	"grpc":        "gRPC",
	"coredns":     "CoreDNS",
	"containerd":  "containerd",
	"rkt":         "rkt",
	"cni":         "CNI",
	"envoy":       "Envoy",
	"jaeger":      "Jaeger",
}

var infraLogos map[string]string = map[string]string{
	"kubernetes":  "./public/icons/kubernetes.png",
	"etcd":        "./public/icons/etcd.svg",
	"flannel":     "./public/icons/flannel.svg",
	"docker":      "./public/icons/docker.svg",
	"prometheus":  "./public/icons/prometheus.png",
	"opentracing": "./public/icons/opentracing.png",
	"fluentd":     "./public/icons/fluentd.png",
	"linkerd":     "./public/icons/linkerd.png",
	"grpc":        "./public/icons/grpc.png",
	"coredns":     "./public/icons/coredns.png",
	"containerd":  "./public/icons/containerd.png",
	"rkt":         "./public/icons/rkt.png",
	"cni":         "./public/icons/cni.png",
	"envoy":       "./public/icons/envoy.png",
	"jaeger":      "./public/icons/jaeger.png",
}

var InfraOrder []string = []string{"kubernetes", "etcd", "flannel", "docker", "prometheus", "opentracing", "fluentd", "linerd", "grpc", "coredns", "containerd", "rkt", "cni", "envoy", "jaeger"}
var InfraTitles []HtmlInfraTitle

func init() {
	for i := 0; i < len(InfraOrder); i++ {
		infra_name := InfraOrder[i]
		InfraTitles = append(InfraTitles, HtmlInfraTitle{
			Name: getInfraName(infra_name),
			Logo: getInfraLogo(infra_name),
		})
	}
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
