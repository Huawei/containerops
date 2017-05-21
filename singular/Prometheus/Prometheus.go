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

package Prometheus

//Installing Prometheus Server

// mkdir ~/Downloads
// cd ~/Downloads
// wget "https://github.com/prometheus/prometheus/releases/download/0.15.1/prometheus-0.15.1.linux-amd64.tar.gz"
// mkdir -p ~/Prometheus/server
// cd ~/Prometheus/server
// tar -xvzf ~/Downloads/prometheus-0.15.1.linux-amd64.tar.gz
// ./prometheus -version

//Installing Node Exporter
// mkdir -p ~/Prometheus/node_exporter
// cd ~/Prometheus/node_exporter
// wget https://github.com/prometheus/node_exporter/releases/download/0.11.0/node_exporter-0.11.0.linux-amd64.tar.gz -O ~/Downloads/node_exporter-0.11.0.linux-amd64.tar.gz
// tar -xvzf ~/Downloads/node_exporter-0.11.0.linux-amd64.tar.gz

//Running Node Exporter as a Service

//  ln -s ~/Prometheus/node_exporter/node_exporter /usr/bin

//  nano /etc/init/node_exporter.conf
// # Run node_exporter

// start on startup

// script
//    /usr/bin/node_exporter
// end script
// service node_exporter start

//return http://your_server_ip:9100/metrics

// Starting Prometheus Server
//cd ~/Prometheus/server
//nano ~/Prometheus/server/prometheus.yml
// scrape_configs:
//   - job_name: "node"
//     scrape_interval: "15s"
//     target_groups:
//     - targets: ['localhost:9100']
// nohup ./prometheus > prometheus.log 2>&1 &
// tail ~/Prometheus/server/prometheus.log
//return  http://your_server_ip:9090
