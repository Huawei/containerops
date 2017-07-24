package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func main() {
	data := os.Getenv("CO_DATA")
	if len(data) == 0 {
		fmt.Fprintf(os.Stderr, "[COUT] %s\n", "The CO_DATA value is null.")
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = %s\n", "false")
		os.Exit(1)
	}

	action, token, kubeApiServerUrl, kubeControllerManagerUrl, kubeSchedulerUrl, kubectlUrl, kubeleteUrl, kubeProxyUrl, err := parseEnv(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Parse the CO_DATA error: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
	}

	if action != "release" {
		fmt.Fprintf(os.Stderr, "[COUT] %s\n", "Unknown action, the component only support build, test and release action.")
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
	}

	bs, err := ioutil.ReadFile("./singular.template.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Failed to read the template file: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
	}
	singularTemplate := string(bs)

	singular := map[string]interface{}{
		"Token":                    token,
		"KubeApiServerUrl":         kubeApiServerUrl,
		"KubeControllerManagerUrl": kubeControllerManagerUrl,
		"KubeSchedulerUrl":         kubeSchedulerUrl,
		"KubectlUrl":               kubectlUrl,
		"KubeletUrl":               kubeleteUrl,
		"KubeProxyUrl":             kubeProxyUrl,
	}
	t := template.Must(template.New("Singular").Parse(singularTemplate))
	var buf bytes.Buffer
	err = t.Execute(&buf, singular)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate yaml from template: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
	}

	// Write the yaml to file
	if err := ioutil.WriteFile("./singular.yaml", buf.Bytes(), os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Failed to write yaml into file: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
	}

	// Call singular binary to deploy the k8s cluster
	// Since the Provider and Token is provied in the yaml, the config file is not necessary here, create a mock
	mockConfigContent := `
	[singular]
	provider="digitalocean"
	token="helloworld"
	`
	if err := ioutil.WriteFile("./runtime.toml", []byte(mockConfigContent), os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "[COUT] Failed to generate mock config: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
	}
	cmd := exec.Command("./singular", "deploy", "template", "./singular.yaml", "--config", "./runtime.toml", "--verbose", "true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to deploy k8s cluster with singular: %s\n", err.Error())
		fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = false\n")
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "[COUT] CO_RESULT = true\n")
}

// Parse CO_DATA value, and return token, and the urls of kubernetes components
func parseEnv(env string) (action, token, kubeApiServerUrl, kubeControllerManagerUrl, kubeSchedulerUrl, kubectlUrl, kubeleteUrl, kubeProxyUrl string, err error) {
	files := strings.Fields(env)
	if len(files) == 0 {
		err = fmt.Errorf("CO_DATA value is null\n")
		return
	}

	for _, v := range files {
		s := strings.Split(v, "=")
		key, value := s[0], s[1]

		switch key {
		case "action":
			action = value
		case "token":
			token = value
		case "kube_apiserver_url":
			kubeApiServerUrl = value
		case "kube_controllermanager_url":
			kubeControllerManagerUrl = value
		case "kube_scheduler_url":
			kubeSchedulerUrl = value
		case "kubectl_url":
			kubectlUrl = value
		case "kubelete_url":
			kubeleteUrl = value
		case "kube_proxy_url":
			kubeProxyUrl = value
		default:
			fmt.Fprintf(os.Stdout, "[COUT] Unknown Parameter: [%s]\n", s)
		}
	}
	return
}
