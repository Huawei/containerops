# Pilotage - DevOps Workflow Engine 

## What's ContainerOps?

> DevOps With Container, DevOps For Container!

## What is Pilotage ?
Pilotage is a DevOps workflow engine with customizable DevOps component repository base container, and it's core project of ContainerOps. It's key features and goals include:
- Customizable DevOps component with container running in [Kuberenetes](https://kubernetes.io)..
- Customizable DevOps workflow with components and services.
- Customizable event type and event in the DevOps workflow.
- Web UI portal for all functions above.

## What's differents between component and service?

|Item|Component|Service|
|------|----|------|
|duration|Destroy after DevOps task done or failure.|Long-time running waiting for DevOps taks.|
|Provider|User|User or third service.|
|Format|Container Image|Functions could be called|

## Why it matters ?

## The Pilotage's story :)

## Runtime configuration

```
runmode = dev

listenmode = https
httpscertfile = cert/containerops/containerops.crt
httpskeyfile = cert/containerops/containerops.key

[log]
filepath = log/backend.log
level = info

[database]
driver = mysql
uri = containerops:containerops@/containerops?charset=utf8&parseTime=True&loc=Asia%2FShanghai
```

#### Nginx configuration
It's a Nginx config example. You can change **client_max_body_size** what limited upload file size. You should copy `containerops.me` keys from `cert/containerops.me` to `/etc/nginx`, then run **pilotage** with `http` mode and listen on `127.0.0.1:9911`.

```nginx
upstream pilostage_upstream {
  server 127.0.0.1:9911;
}

server {
  listen 80;
  server_name containerops.me;
  rewrite  ^/(.*)$  https://containerops.me/$1  permanent;
}

server {
  listen 443;

  server_name containerops.me;

  access_log /var/log/nginx/containerops-me.log;
  error_log /var/log/nginx/containerops-me-errror.log;

  ssl on;
  ssl_certificate /etc/nginx/containerops.me.crt;
  ssl_certificate_key /etc/nginx/containerops.me.key;

  client_max_body_size 1024m;
  chunked_transfer_encoding on;

  proxy_redirect     off;
  proxy_set_header   X-Real-IP $remote_addr;
  proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
  proxy_set_header   X-Forwarded-Proto $scheme;
  proxy_set_header   Host $http_host;
  proxy_set_header   X-NginX-Proxy true;
  proxy_set_header   Connection "";
  proxy_http_version 1.1;

  location / {
    proxy_pass         http://pilostage_upstream;
  }
}
```

### Test pilotage service
- Run directly:

```bash
./pilotage web --address 0.0.0.0
```

- Run with Nginx:

```bash
./pilotage web --address 127.0.0.1 --port 9911 &
```

## Update The Libraries Dependencies

```
go get -u -v github.com/urfave/cli
go get -u -v gopkg.in/macaron.v1
go get -u -v github.com/jinzhu/gorm
```

## Certificate of Origin
By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. 

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
660 York Street, Suite 102,
San Francisco, CA 94110 USA

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.

Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

## Format of the Commit Message

You just add a line to every git commit message, like this:

    Signed-off-by: Meaglith Ma <maquanyi@huawei.com>

Use your real name (sorry, no pseudonyms or anonymous contributions.)

If you set your `user.name` and `user.email` git configs, you can sign your
commit automatically with `git commit -s`.

