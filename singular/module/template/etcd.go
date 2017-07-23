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

var EtcdCATemplate = map[string]string{
	"etcd-3.2.2": `
{
	"CN": "etcd",
	"hosts": [
		"127.0.0.1",
		"{{.IP}}"
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
}`,
}

var EtcdSystemdTemplate = map[string]string{
	"etcd-3.2.2": `
[Unit]
Description=Etcd Server
After=network.target
After=network-online.target
Wants=network-online.target
Documentation=https://github.com/coreos

[Service]
Type=notify
WorkingDirectory=/var/lib/etcd/
ExecStart=/usr/local/bin/etcd \
  --name={{.Name}} \
  --cert-file=/etc/etcd/ssl/etcd.pem \
  --key-file=/etc/etcd/ssl/etcd-key.pem \
  --peer-cert-file=/etc/etcd/ssl/etcd.pem \
  --peer-key-file=/etc/etcd/ssl/etcd-key.pem \
  --trusted-ca-file=/etc/kubernetes/ssl/ca.pem \
  --peer-trusted-ca-file=/etc/kubernetes/ssl/ca.pem \
  --initial-advertise-peer-urls=https://{{.IP}}:2380 \
  --listen-peer-urls=https://{{.IP}}:2380 \
  --listen-client-urls=https://{{.IP}}:2379,http://127.0.0.1:2379 \
  --advertise-client-urls=https://{{.IP}}:2379 \
  --initial-cluster-token=etcd-cluster-0 \
  --initial-cluster={{.Nodes}} \
  --initial-cluster-state=new \
  --data-dir=/var/lib/etcd
Restart=on-failure
RestartSec=5
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target`,
}
