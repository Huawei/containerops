/*
Copyright 2014 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

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

package deploy

func Sflanneld() {

	// cat > flanneld-csr.json <<EOF
	// {
	//   "CN": "flanneld",
	//   "hosts": [],
	//   "key": {
	//     "algo": "rsa",
	//     "size": 2048
	//   },
	//   "names": [
	//     {
	//       "C": "CN",
	//       "ST": "BeiJing",
	//       "L": "BeiJing",
	//       "O": "k8s",
	//       "OU": "System"
	//     }
	//   ]
	// }
	// EOF

	//    cfssl gencert -ca=/etc/kubernetes/ssl/ca.pem \
	//   -ca-key=/etc/kubernetes/ssl/ca-key.pem \
	//   -config=/etc/kubernetes/ssl/ca-config.json \
	//   -profile=kubernetes flanneld-csr.json | cfssljson -bare flanneld
	//  ls flanneld*
	// flanneld.csr  flanneld-csr.json  flanneld-key.pem flanneld.pem
	//  sudo mkdir -p /etc/flanneld/ssl
	//  sudo mv flanneld*.pem /etc/flanneld/ssl
	//  rm flanneld.csr  flanneld-csr.json
}
func Setcd() {

	//  cat > etcd-csr.json <<EOF
	// {
	//   "CN": "etcd",
	//   "hosts": [
	//     "127.0.0.1",
	//     "${NODE_IP}"
	//   ],
	//   "key": {
	//     "algo": "rsa",
	//     "size": 2048
	//   },
	//   "names": [
	//     {
	//       "C": "CN",
	//       "ST": "BeiJing",
	//       "L": "BeiJing",
	//       "O": "k8s",
	//       "OU": "System"
	//     }
	//   ]
	// }
	// EOF

	//   $ cfssl gencert -ca=/etc/kubernetes/ssl/ca.pem \
	//   -ca-key=/etc/kubernetes/ssl/ca-key.pem \
	//   -config=/etc/kubernetes/ssl/ca-config.json \
	//   -profile=kubernetes etcd-csr.json | cfssljson -bare etcd
	// $ ls etcd*
	// etcd.csr  etcd-csr.json  etcd-key.pem etcd.pem
	// $ sudo mkdir -p /etc/etcd/ssl
	// $ sudo mv etcd*.pem /etc/etcd/ssl
	// $ rm etcd.csr  etcd-csr.json
}

func Skubernetes() {
	//   cat > kubernetes-csr.json <<EOF
	// {
	//   "CN": "kubernetes",
	//   "hosts": [
	//     "127.0.0.1",
	//     "${MASTER_IP}",
	//     "${CLUSTER_KUBERNETES_SVC_IP}",
	//     "kubernetes",
	//     "kubernetes.default",
	//     "kubernetes.default.svc",
	//     "kubernetes.default.svc.cluster",
	//     "kubernetes.default.svc.cluster.local"
	//   ],
	//   "key": {
	//     "algo": "rsa",
	//     "size": 2048
	//   },
	//   "names": [
	//     {
	//       "C": "CN",
	//       "ST": "BeiJing",
	//       "L": "BeiJing",
	//       "O": "k8s",
	//       "OU": "System"
	//     }
	//   ]
	// }
	// EOF

	//    cfssl gencert -ca=/etc/kubernetes/ssl/ca.pem \
	//   -ca-key=/etc/kubernetes/ssl/ca-key.pem \
	//   -config=/etc/kubernetes/ssl/ca-config.json \
	//   -profile=kubernetes kubernetes-csr.json | cfssljson -bare kubernetes
	// $ ls kubernetes*
	// kubernetes.csr  kubernetes-csr.json  kubernetes-key.pem  kubernetes.pem
	// $ sudo mkdir -p /etc/kubernetes/ssl/
	// $ sudo mv kubernetes*.pem /etc/kubernetes/ssl/
	// $ rm kubernetes.csr  kubernetes-csr.json
}

func Skube_apiserver() {

	//  cat > token.csv <<EOF
	// ${BOOTSTRAP_TOKEN},kubelet-bootstrap,10001,"system:kubelet-bootstrap"
	// EOF
	// $ mv token.csv /etc/kubernetes/
}
