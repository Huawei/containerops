# Dockyard - Container And Artifact Repository 

![Dockyard](docs/images/dockyard.jpg "Dockyard - Container And Artifact Repository")

## What is Dockyard ?
Dockyard is a container and artifact repository storing and distributing container image, software artifact and virtual images of KVM or XEN. It's key features and goals include:
- Multi supported distribute protocols include Docker Registry V1 & V2 and App Container Image Discovery.
- Multi supported software artifact format include jar, war, tar and so on.
- Multi supported virtual images of KVM, XEN, VirtualBox and so on.
- Container image, software artifact and virtual images encryption, verification and vulnerability analytsis.
- Custome distribute protocol by framework base HTTPS and peer to peer. 
- Authentication in distributing process and authorization for public and private repository.
- Supporting mainstream object storage service like Amazon S3, Google Cloud Storage. 
- Built-in object storage service for deployment convenience.
- Web UI portal for all functions above.

## Why it matters ?

## The Dockyard's Story :)

## Runtime configuration

```
runmode = "dev"

listenmode = "https"
httpscertfile = "cert/containerops/containerops.crt"
httpskeyfile = "cert/containerops/containerops.key"

[site]
domain = "containerops.me"

[log]
filepath = "log/backend.log"
level = "info"

[database]
driver = "mysql"
uri = "containerops:containerops@/containerops?charset=utf8&parseTime=True&loc=Asia%2FShanghai"

[deployment]
domains = "containerops.me"

[dockerv1]
standalone = "true"
version = "0.9"
storage = "/tmp/data/dockerv1"

[dockerv2]
distribution = "registry/2.0"
storage = "/tmp/data/dockerv2"

[appc]
storage = "/tmp/data/appc"

[updateserver]
keymanager = "/tmp/containerops_keymanager_cache"
storage = "/tmp/containerops_storage_cache"

```

#### Nginx configuration
It's a Nginx config example. You can change **client_max_body_size** what limited upload file size. You should copy `containerops.me` keys from `cert/containerops.me` to `/etc/nginx`, then run **Dockyard** with `http` mode and listen on `127.0.0.1:9911`.

```nginx
upstream dockyard_upstream {
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
    proxy_pass         http://dockyard_upstream;
  }
}
```

### Database Configuration

#### Database SQL

```
INSERT INTO mysql.user(Host,User,Password) VALUES ('localhost', 'containerops', password('containerops'));
CREATE DATABASE `containerops` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci;
GRANT ALL PRIVILEGES ON containerops.* TO containerops@localhost IDENTIFIED BY 'containerops';
FLUSH PRIVILEGES;
```

#### Initlization Tables

```
./dockyard database migrate
```

### Start dockyard service
- Run directly:

```bash
./dockyard daemon run --address 0.0.0.0 --port 443
```

- Run with Nginx:

```bash
./dockyard daemon run --address 127.0.0.1 --port 9911 &
```

## How to build

We are using [glide](https://glide.sh/) as package manager/

* retrieve dependencies
```
glide install
```

* build (with go 1.6+)
```
go build
```

## How to involve
If any issues are encountered while using the dockyard project, several avenues are available for support:
<table>
<tr>
	<th align="left">
	Issue Tracker
	</th>
	<td>
	https://github.com/Huawei/dockyard/issues
	</td>
</tr>
<tr>
	<th align="left">
	Google Groups
	</th>
	<td>
	https://groups.google.com/forum/#!forum/dockyard-dev
	</td>
</tr>
</table>

### Pull Requests

If you want to contribute to the template, you can create pull requests. All pull requests must be done to the `develop` branch. We are working on build an automated tests with ourself means *containerops*, now we just add *Travis CI* instead.

## Who should join
- Ones who want to choose a container image hub instead of docker hub.
- Ones who want to ease the burden of container image management.

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

    Signed-off-by: Meaglith Ma <genedna@gmail.com>

Use your real name (sorry, no pseudonyms or anonymous contributions.)

If you set your `user.name` and `user.email` git configs, you can sign your
commit automatically with `git commit -s`.
