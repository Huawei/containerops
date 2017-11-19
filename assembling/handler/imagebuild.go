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
	"net/http"
	"net/url"
	"os"
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
	queries := getQueryParameters(mctx.Req.Request.URL)
	registry := queries["registry"]
	namespace := queries["namespace"]
	image := queries["image"]
	tag := queries["tag"]
	buildArgsJSON := queries["buildargs"]
	authstr := queries["authstr"]

	var buildArgs map[string]*string
	if buildArgsJSON == "" {
		buildArgs = map[string]*string{}
	} else if err := json.Unmarshal([]byte(buildArgsJSON), &buildArgs); err != nil {
		log.Errorf("Failed to parse buildargs: %s", err.Error())
		return http.StatusBadRequest, []byte("{}")
	}

	isBodyDockerArchive, buf, err := isDockerArchive(mctx.Req.Request.Body)
	if err != nil {
		log.Errorf("Failed to check gzip format: %s", err.Error())
		return http.StatusInternalServerError, []byte("{}")
	}
	if buf.Len() == 0 {
		log.Errorf("Empty file")
		return http.StatusBadRequest, []byte("{}")
	}

	var tarfile io.Reader
	if !isBodyDockerArchive {
		tarfile, err = createTarFile(buf)
	} else {
		tarfile = buf
	}

	log.Infof("Init k8s resources")
	podClient, serviceClient, err := initK8SResourceInterfaces(common.Assembling.KubeConfig)
	if err != nil {
		log.Errorf("Failed to init k8s pod client: %s", err.Error())
		return http.StatusInternalServerError, []byte(err.Error())
	}

	buildId := uuid.NewV4().String()
	podName := fmt.Sprintf("containerops-build-pod-%s", buildId)
	serviceName := fmt.Sprintf("containerops-build-svc-%s", buildId)

	log.Infof("Create pod %s for build %s", podName, buildId)
	_, err = createPod(podClient, podName, buildId)
	if err != nil {
		log.Errorf("Failed to create pod: %s", err.Error())
		return http.StatusInternalServerError, []byte("{}")
	}
	defer deletePod(podClient, podName)

	log.Infof("Create load balancer for build %s", buildId)
	loadBalancer, err := createLoadBalancer(serviceClient, serviceName, buildId)
	if err != nil {
		log.Errorf("Failed to create load balancer: %s", err.Error())
		return http.StatusInternalServerError, []byte("{}")
	}

	servicePort := 2375
	defer func() {
		if err := deleteLoadBalancer(serviceClient, serviceName); err != nil {
			log.Errorf("Failed to delete load balancer: %s", err.Error())
		}
	}()

	if len(loadBalancer.Status.LoadBalancer.Ingress) == 0 {
		log.Errorf("Load balancer: no ingress created")
		return http.StatusInternalServerError, []byte("{}")
	}
	serviceIP := loadBalancer.Status.LoadBalancer.Ingress[0].IP
	dockerDaemonHost := fmt.Sprintf("%s:%d", serviceIP, servicePort)
	ctx, dockerClient := initDockerCli(dockerDaemonHost)

	log.Infof("Build image, id: %s", buildId)
	if err := buildImage(ctx, dockerClient, registry, namespace, image, tag, buildArgs, tarfile); err != nil {
		log.Errorf("Failed to build image: %s", err.Error())
		return http.StatusInternalServerError, []byte("{}")
	}

	// TODO Support pushing to registries that need authorization
	// authStr, _ := generateAuthStr("", "")
	authStr := authstr
	if authStr == "" {
		authStr, _ = generateAuthStr("", "")
	}

	log.Infof("Push image, id: %s", buildId)
	if err := pushImage(ctx, dockerClient, registry, namespace, image, tag, authStr); err != nil {
		log.Errorf("Failed to push image: %s", err.Error())
		return http.StatusInternalServerError, []byte("{}")
	}

	builtImage := fmt.Sprintf("%s/%s/%s:%s", registry, namespace, image, tag)
	log.Infof("Image pushed: %s", builtImage)
	return http.StatusOK, []byte(fmt.Sprintf("{\"endpoint\":\"%s\"}", builtImage))
}

// Take the first value of the query
func getQueryParameters(u *url.URL) map[string]string {
	ret := map[string]string{}
	for key, ary := range u.Query() {
		if len(ary) == 0 {
			ret[key] = ""
		} else {
			ret[key] = ary[0]
		}
	}
	return ret
}

func isDockerArchive(src io.Reader) (bool, *bytes.Buffer, error) {
	var buf bytes.Buffer

	// We should not make constraint on the number of bytes read, and should skip the io.EOF error
	// Since the file might not be tar file and the length might be shorter than 265
	_, err := io.CopyN(&buf, src, 265)
	if err != nil && err != io.EOF {
		return false, nil, err
	} /* else if n != 265 {
		return false, nil, fmt.Errorf("Failed to read first 265 bytes")
	} */

	bs := buf.Bytes()
	is_docker_tar_file := isBZip(bs) || isGZip(bs) || isXZ(bs) || isTar(bs)
	_, err = io.Copy(&buf, src)

	return is_docker_tar_file, &buf, err
}

func isBZip(header []byte) bool {
	return len(header) >= 3 &&
		header[0] == 0x42 &&
		header[1] == 0x5a &&
		header[2] == 0x68
}

func isGZip(header []byte) bool {
	return len(header) >= 2 &&
		header[0] == 0x1f &&
		header[1] == 0x8b
}

func isXZ(header []byte) bool {
	return len(header) >= 6 &&
		header[0] == 0xfd &&
		header[1] == 0x37 &&
		header[2] == 0x7a &&
		header[3] == 0x58 &&
		header[4] == 0x5a &&
		header[5] == 0x00
}

func isTar(header []byte) bool {
	if len(header) < 264 {
		return false
	}

	magic := header[257:265]
	return isPosixTar(magic) || isGnuTar(magic)
}

func isPosixTar(magic []byte) bool {
	return len(magic) >= 8 &&
		magic[0] == 0x75 &&
		magic[1] == 0x73 &&
		magic[2] == 0x74 &&
		magic[3] == 0x61 &&
		magic[4] == 0x72 &&
		magic[5] == 0x00 &&
		magic[6] == 0x30 &&
		magic[7] == 0x30
}
func isGnuTar(magic []byte) bool {
	return len(magic) >= 8 &&
		magic[0] == 0x75 &&
		magic[1] == 0x73 &&
		magic[2] == 0x74 &&
		magic[3] == 0x61 &&
		magic[4] == 0x72 &&
		magic[5] == 0x20 &&
		magic[6] == 0x20 &&
		magic[7] == 0x00
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
					Image: common.Assembling.DockerDaemonImage,
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
	cli, err := client.NewClient(targetUrl, "v1.23", httpClient, buildClientHeaders)
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

func buildImage(ctx context.Context, cli *client.Client, host, namespace, imageName, tag string, buildArgs map[string]*string, tarFileReader io.Reader) error {
	targetTag := fmt.Sprintf("%s/%s/%s:%s", host, namespace, imageName, tag)
	buildOptions := types.ImageBuildOptions{
		Tags:      []string{targetTag},
		BuildArgs: buildArgs,
	}

	out, err := cli.ImageBuild(ctx, tarFileReader, buildOptions)
	if err != nil {
		return err
	}

	defer out.Body.Close()
	// io.Copy(ioutil.Discard, out.Body)
	io.Copy(os.Stdout, out.Body)
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
	// io.Copy(ioutil.Discard, pushResult)
	io.Copy(os.Stdout, pushResult)
	return nil
}

func createLoadBalancer(serviceClient v1.ServiceInterface, serviceName, buildId string) (*corev1.Service, error) {
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

	// Create LoadBalancer
	_, err := serviceClient.Create(svc)
	if err != nil {
		return nil, err
	}

	var loadBalancer *corev1.Service
	var e error
	start := time.Now()

	for {
		loadBalancer, e = serviceClient.Get(serviceName, metav1.GetOptions{})
		if e != nil || len(loadBalancer.Status.LoadBalancer.Ingress) != 0 {
			break
		}
		time.Sleep(time.Second * 3)
		if time.Since(start).Seconds() > 180 {
			e = fmt.Errorf("NoadBalancer creation timeout")
			// If the error is not nil, the deletion will most likely to be ignored ouside the function.
			if err := deleteLoadBalancer(serviceClient, serviceName); err != nil {
				log.Errorf("Failed to delete load balancer: %s", err.Error())
			}
			break
		}
	}

	return loadBalancer, e
}

func deleteLoadBalancer(serviceClient v1.ServiceInterface, serviceName string) error {
	deletePolicy := metav1.DeletePropagationForeground
	if err := serviceClient.Delete(serviceName, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	return nil
}
