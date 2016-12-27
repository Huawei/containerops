---
title: Install Local
keywords: install
tags: [component]
sidebar: home_sidebar
permalink: install-local.html
summary: Install Local
---

## Install Local

- __Check mysql remote login__

 ```
 mysql -h xxx.xxx.xxx.xxx -uxxx -pxxx
 ```

- __Create mysql database__

  ```
  mysql> create database containerops;
  ```

- __Git containerops__

  ```
  git clone https://github.com/Huawei/containerops.git
  ```

- __Create conf folder__

  ```
  cd containerops/pilotage
  mkdir conf
  vim conf/containerops.toml
  ```

- __containerops.toml contents, projectaddr, uri replaced with actual values__

  ```
  appname = "pilotage"
  usage = "DevOps Workflow Engine"
  version = "0.0.1"
  author = "Meaglith Ma"
  email = "genedna@gmail.com"

  # include runtime.conf
  runmode = "dev"

  listenmode = "http"
  projectaddr = "192.168.10.180:10000"
  httpscertfile = "cert/containerops/containerops.crt"
  httpskeyfile = "cert/containerops/containerops.key"

  [log]
  filepath = "log/backend.log"
  level = "info"

  [database]
  driver = "mysql"
  uri = "root:123456@tcp(192.168.10.180:3306)/containerops?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
  ```

- __get the dependency packages__

  ```
  go get -u -v github.com/urfave/cli
  go get -u -v gopkg.in/macaron.v1
  go get -u -v github.com/jinzhu/gorm
  ```

- __generate database structs__

  ```
  go  run main.go database migrate
  mysql> show tables;
  ```

- __run daemon__

  ```
  go run main.go daemon start --port 10000
  ```

- __vim api.js, set host__

  ```
  vim src/app/common/api.js

  "host" : "http://192.168.10.180:10000"
  ```

- __set scaffold environment__

  ```
  cd ../scaffold/
  apt-get install nodejs npm -y
  npm install -g gulp bower
  ```

- __install scaffold dependencies__

  ```
  npm install
  bower install
  ```

- __Compile all items and start web server__

  ```
  gulp
  ```
