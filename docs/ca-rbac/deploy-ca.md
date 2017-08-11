CA Design
============
# config ca
-> root-ca </br>
--> etcd-root-ca </br>
----> etcd-server-ca </br>
----> etcd-peer-ca </br>
----> etcd-client-kube-ca </br>
--> network-root-ca </br>
----> network-flannel-ca </br>
--> kubernetes-root-ca </br>
----> kubernetes-client-ca </br>
------> kubernetes-client-ca </br>
------> kubernetes-client-ca </br>
----> kubernetes-server-ca </br>
# root ca
# etcd
# flannel ca
# kubernetes ca



Signed-off-by: Fanliang Meng <mengfanliang@huawei.com>

curl  \
--cert ./shell/kubernetes-ca/kubernetes-admin.pem \
--key ./shell/kubernetes-ca/kubernetes-admin-key.pem \
--cacert ./shell/root-ca/root-ca.pem \
https://10.138.48.164:6443/api/v1/endpoints

./bin/kubernetes/client/kubectl \
--certificate-authority="./shell/root-ca/root-ca.pem" \
--client-certificate="./shell/kubernetes-ca/kubernetes-admin.pem" \
--client-key="./shell/kubernetes-ca/kubernetes-admin-key.pem" \
--server="https://10.138.48.164:6443" \
get svc
