package main

import (
	"fmt"

	"github.com/Huawei/containerops/singular/nodes"
	//"github.com/Huawei/containerops/singular/nodes"
	//"github.com/Huawei/containerops_he3io/singular/cmd"
)

type SSHCommander struct {
	User string
	IP   string
}

//var nodes = [2][2]string{{"192.168.60.141", "centos-master"}, {"192.168.60.150", "centos-minion"}}
var m = make(map[string]string)

func main() {
	// SSHCommander.IP
	m["centos-master"] = "192.168.60.158"
	m["centos-minion"] = "192.168.60.157"
	for k, ip := range m {
		fmt.Printf("k=%v, v=%v\n", k, ip)
		if k == "centos-master" {
			nodes.Deploymaster(m, ip)
		}
		if k == "centos-minion" {
			nodes.Deploynode(m, ip)
		}
	}
}
