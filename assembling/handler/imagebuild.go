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

package handler

import (
	"archive/tar"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Huawei/containerops/common"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
	macaron "gopkg.in/macaron.v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func BuildImageHandler(mctx *macaron.Context) (int, []byte) {
	// TODO image, namespace, registry, tag pattern validation with regex
	registry := mctx.Req.Request.FormValue("registry")
	namespace := mctx.Req.Request.FormValue("namespace")
	image := mctx.Req.Request.FormValue("image")
	tag := mctx.Req.Request.FormValue("tag")

	podClient, serviceClient, err := initK8SResourceInterfaces(common.Assembling.KubeConfig)
	if err != nil {
		log.Errorf("Failed to init k8s pod client: %s", err.Error())
		return http.StatusInternalServerError, []byte(err.Error())
	}

	buildId := uuid.NewV4().String()
	podName := fmt.Sprintf("containerops-build-pod-%s", buildId)
	serviceName := fmt.Sprintf("containerops-build-svc-%s", buildId)

	_, err = createPod(podClient, podName, buildId)
	if err != nil {
		log.Errorf("Failed to create pod: %s", err.Error())
		return http.StatusInternalServerError, []byte("{}")
	}
	defer deletePod(podClient, podName)

	nodeBalancer, err := createNodeBalancer(serviceClient, serviceName, buildId)
	if err != nil {
		log.Errorf("Failed to create pod: %s", err.Error())
		return http.StatusInternalServerError, []byte("{}")
	}

	servicePort := 2375
	defer deleteNodeBalancer(serviceClient, serviceName)

	serviceIP := nodeBalancer.Status.LoadBalancer.Ingress[0].IP
	dockerDaemonHost := fmt.Sprintf("%s:%d", serviceIP, servicePort)
	ctx, dockerClient := initDockerCli(dockerDaemonHost)

	tarfile, err := createTarFile(mctx.Req.Request.Body)
	if err := buildImage(ctx, dockerClient, registry, namespace, image, tag, tarfile); err != nil {
		log.Errorf("Failed to build image: %s", err.Error())
		return http.StatusInternalServerError, []byte("{}")
	}

	// TODO Support pushing to registries that need authorization
	authStr, _ := generateAuthStr("", "")

	if err := pushImage(ctx, dockerClient, registry, namespace, image, tag, authStr); err != nil {
		log.Errorf("Failed to push image: %s", err.Error())
		return http.StatusInternalServerError, []byte("{}")
	}

	builtImage := fmt.Sprintf("%s/%s/%s:%s", registry, namespace, image, tag)
	return http.StatusOK, []byte(fmt.Sprintf("{\"endpoint\":\"%s\"}", builtImage))
}

func generateAuthStr(username, password string) (string, error) {
	authConfig := types.AuthConfig{
		Username: username,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		log.Error(err)
		return "", err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	return authStr, nil

}

func initK8SResourceInterfaces(kubeconfig string) (v1.PodInterface, v1.ServiceInterface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, err
	}

	podClient := clientset.CoreV1().Pods(corev1.NamespaceDefault)
	serviceClient := clientset.CoreV1().Services(corev1.NamespaceDefault)
	return podClient, serviceClient, nil
}

func createPod(podClient v1.PodInterface, podName, buildId string) (*corev1.Pod, error) {
	isPrivileged := true
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
			Labels: map[string]string{
				"build-id": buildId,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "docker-dind",
					Image: common.Assembling.DockerDamonImage,
					// Reservation for Args
					Args: []string{},
					SecurityContext: &corev1.SecurityContext{
						Privileged: &isPrivileged,
					},
				},
			},
			// Reservation for NodeSelector
			NodeSelector: map[string]string{},
		},
	}

	// Create pod
	_, err := podClient.Create(pod)
	if err != nil {
		return nil, err
	}

	// Monit pod creation status in a ticker, since the monitoring API in the k8s client is too complicated and lack of docs
	var buildPod *corev1.Pod
	var e error
	start := time.Now()
	for {
		buildPod, e = podClient.Get(podName, metav1.GetOptions{})
		if e != nil || buildPod.Status.Phase == "Running" {
			break
		}
		time.Sleep(time.Second)
		if time.Since(start).Seconds() > 30 {
			buildPod, e = nil, fmt.Errorf("Pod creation timeout")
			break
		}
	}

	return buildPod, e
}

func deletePod(podClient v1.PodInterface, podName string) error {
	deletePolicy := metav1.DeletePropagationForeground
	if err := podClient.Delete(podName, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	return nil
}

func initDockerCli(registryHost string) (context.Context, *client.Client) {
	ctx := context.Background()
	var httpClient *http.Client
	buildClientHeaders := map[string]string{"Content-Type": "application/tar"}

	targetUrl := fmt.Sprintf("http://%s", registryHost)
	cli, err := client.NewClient(targetUrl, "v1.27", httpClient, buildClientHeaders)
	if err != nil {
		panic(err)
	}

	return ctx, cli
}

func createTarFile(dockerfile io.Reader) (io.Reader, error) {
	// Create a new tar archive.
	tarBuf := new(bytes.Buffer)
	tw := tar.NewWriter(tarBuf)

	// Add dockerfile to the archive.
	contentBuf := new(bytes.Buffer)
	contentBuf.ReadFrom(dockerfile)
	contentBytes := contentBuf.Bytes()

	hdr := &tar.Header{
		Name: "Dockerfile",
		Mode: 0600,
		Size: int64(len(contentBytes)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		log.Fatalln(err)
		return nil, err
	}
	_, err := tw.Write(contentBytes)
	if err != nil {
		return nil, err
	}
	tw.Close()

	return bytes.NewReader(tarBuf.Bytes()), nil
}

func buildImage(ctx context.Context, cli *client.Client, host, namespace, imageName, tag string, tarFileReader io.Reader) error {
	targetTag := fmt.Sprintf("%s/%s/%s:%s", host, namespace, imageName, tag)
	buildOptions := types.ImageBuildOptions{
		Tags: []string{targetTag},
	}

	out, err := cli.ImageBuild(ctx, tarFileReader, buildOptions)
	if err != nil {
		return err
	}

	defer out.Body.Close()
	io.Copy(ioutil.Discard, out.Body)
	// io.Copy(os.Stdout, out.Body)
	return nil
}

func pushImage(ctx context.Context, cli *client.Client, host, namespace, imageName, tag, authStr string) error {
	imagePushOptions := types.ImagePushOptions{
		RegistryAuth: authStr,
	}
	targetTag := fmt.Sprintf("%s/%s/%s:%s", host, namespace, imageName, tag)

	pushResult, err := cli.ImagePush(ctx, targetTag, imagePushOptions)
	if err != nil {
		return err
	}

	defer pushResult.Close()
	io.Copy(ioutil.Discard, pushResult)
	// io.Copy(os.Stdout, pushResult)
	return nil
}

func createNodeBalancer(serviceClient v1.ServiceInterface, serviceName, buildId string) (*corev1.Service, error) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeLoadBalancer,
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Port: 2375,
				},
			},
			Selector: map[string]string{
				"build-id": buildId,
			},
		},
	}

	// Create NodeBalancer
	_, err := serviceClient.Create(svc)
	if err != nil {
		return nil, err
	}

	var nodeBalancer *corev1.Service
	var e error
	start := time.Now()

	for {
		nodeBalancer, e = serviceClient.Get(serviceName, metav1.GetOptions{})
		if e != nil || len(nodeBalancer.Status.LoadBalancer.Ingress) != 0 {
			break
		}
		time.Sleep(time.Second * 3)
		if time.Since(start).Seconds() > 180 {
			nodeBalancer, e = nil, fmt.Errorf("NodeBalancer creation timeout")
			break
		}
	}

	return nodeBalancer, nil
}

func deleteNodeBalancer(serviceClient v1.ServiceInterface, serviceName string) error {
	deletePolicy := metav1.DeletePropagationForeground
	if err := serviceClient.Delete(serviceName, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	return nil
}
