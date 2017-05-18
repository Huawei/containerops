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
