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

package template

var BooststrapToken string = "e3304987e085aad4c62038b2f2e05c94"

var KubernetesCATemplate = map[string]string{
	"kubernetes-1.6.7": `
{
  "CN": "kubernetes",
  "hosts": [
    "127.0.0.1",
    "{{.MasterIP}}",
    "10.254.0.1",
    "kubernetes",
    "kubernetes.default",
    "kubernetes.default.svc",
    "kubernetes.default.svc.cluster",
    "kubernetes.default.svc.cluster.local"
  ],
  "key": {
    "algo": "rsa",
    "size": 4096
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "k8s",
      "OU": "System"
    }
  ]
}
	`,
}

var KubernetesAPIServerSystemdTemplate = map[string]string{
	"kubernetes-1.6.7": `
[Unit]
Description=Kubernetes API Server
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
After=network.target

[Service]
ExecStart=/usr/local/bin/kube-apiserver \
  --admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota \
  --advertise-address={{.MasterIP}} \
  --bind-address={{.MasterIP}} \
  --insecure-bind-address={{.MasterIP}} \
  --authorization-mode=RBAC \
  --runtime-config=rbac.authorization.k8s.io/v1alpha1 \
  --kubelet-https=true \
  --experimental-bootstrap-token-auth \
  --token-auth-file=/etc/kubernetes/token.csv \
  --service-cluster-ip-range=10.254.0.0/16 \
  --service-node-port-range=8400-9000 \
  --tls-cert-file=/etc/kubernetes/ssl/kubernetes.pem \
  --tls-private-key-file=/etc/kubernetes/ssl/kubernetes-key.pem \
  --client-ca-file=/etc/kubernetes/ssl/ca.pem \
  --service-account-key-file=/etc/kubernetes/ssl/ca-key.pem \
  --etcd-cafile=/etc/kubernetes/ssl/ca.pem \
  --etcd-certfile=/etc/kubernetes/ssl/kubernetes.pem \
  --etcd-keyfile=/etc/kubernetes/ssl/kubernetes-key.pem \
  --etcd-servers={{.Nodes}} \
  --enable-swagger-ui=true \
  --allow-privileged=true \
  --apiserver-count=3 \
  --audit-log-maxage=30 \
  --audit-log-maxbackup=3 \
  --audit-log-maxsize=100 \
  --audit-log-path=/var/lib/audit.log \
  --event-ttl=1h \
  --v=2
Restart=on-failure
RestartSec=5
Type=notify
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
`,
}

var KubernetesControllerManagerSystemdTemplate = map[string]string{
	"kubernetes-1.6.7": `
[Unit]
Description=Kubernetes Controller Manager
Documentation=https://github.com/GoogleCloudPlatform/kubernetes

[Service]
ExecStart=/usr/local/bin/kube-controller-manager \
  --address=127.0.0.1 \
  --master=http://{{.MasterIP}}:8080 \
  --allocate-node-cidrs=true \
  --service-cluster-ip-range=10.254.0.0/16 \
  --cluster-cidr=172.30.0.0/16 \
  --cluster-name=kubernetes \
  --cluster-signing-cert-file=/etc/kubernetes/ssl/ca.pem \
  --cluster-signing-key-file=/etc/kubernetes/ssl/ca-key.pem \
  --service-account-private-key-file=/etc/kubernetes/ssl/ca-key.pem \
  --root-ca-file=/etc/kubernetes/ssl/ca.pem \
  --leader-elect=true \
  --v=2
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
`,
}

var KubernetesSchedulerSystemdTemplate = map[string]string{
	"kubernetes-1.6.7": `
[Unit]
Description=Kubernetes Scheduler
Documentation=https://github.com/GoogleCloudPlatform/kubernetes

[Service]
ExecStart=/usr/local/bin/kube-scheduler \
  --address=127.0.0.1 \
  --master=http://{{.MasterIP}}:8080 \
  --leader-elect=true \
  --v=2
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
`,
}

var KubeletSystemdTemplate = map[string]string{
	"kubernetes-1.6.7": `
[Unit]
Description=Kubernetes Kubelet
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
After=docker.service
Requires=docker.service

[Service]
WorkingDirectory=/var/lib/kubelet
ExecStart=/usr/local/bin/kubelet \
  --address={{.IP}} \
  --hostname-override={{.IP}} \
  --pod-infra-container-image=registry.access.redhat.com/rhel7/pod-infrastructure:latest \
  --experimental-bootstrap-kubeconfig=/etc/kubernetes/bootstrap.kubeconfig \
  --kubeconfig=/etc/kubernetes/kubelet.kubeconfig \
  --require-kubeconfig \
  --cert-dir=/etc/kubernetes/ssl \
  --cluster-dns=8.8.8.8 \
  --cluster-domain=cluster.local. \
  --hairpin-mode promiscuous-bridge \
  --allow-privileged=true \
  --serialize-image-pulls=false \
  --logtostderr=true \
  --v=2
ExecStopPost=iptables -A INPUT -s 10.0.0.0/8 -p tcp --dport 4194 -j ACCEPT
ExecStopPost=iptables -A INPUT -s 172.16.0.0/12 -p tcp --dport 4194 -j ACCEPT
ExecStopPost=iptables -A INPUT -s 192.168.0.0/16 -p tcp --dport 4194 -j ACCEPT
ExecStopPost=iptables -A INPUT -p tcp --dport 4194 -j DROP
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
`,
}

var KubeProxyCATemplate = map[string]string{
	"kubernetes-1.6.7": `
{
  "CN": "system:kube-proxy",
  "hosts": [],
  "key": {
    "algo": "rsa",
    "size": 4096
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "k8s",
      "OU": "System"
    }
  ]
}
`,
}

var KubeProxySystemdTemplate = map[string]string{
	"kubernetes-1.6.7": `
[Unit]
Description=Kubernetes Kube-Proxy Server
Documentation=https://github.com/GoogleCloudPlatform/kubernetes
After=network.target

[Service]
WorkingDirectory=/var/lib/kube-proxy
ExecStart=/usr/local/bin/kube-proxy \
  --bind-address={{.IP}} \
  --hostname-override={{.IP}} \
  --cluster-cidr=10.254.0.0/16 \
  --kubeconfig=/etc/kubernetes/kube-proxy.kubeconfig \
  --logtostderr=true \
  --v=2
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
`,
}
