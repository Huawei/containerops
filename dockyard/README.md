## Dockyard - Container And Artifact Repository 

![Dockyard](../docs/images/dockyard.jpg "Dockyard - Container And Artifact Repository")

### The Dockyard's Story :)

The Dockyard is my startup project before joining Huawei in April 2015, after that I donated it to Huawei company. Now the git repository is in the Huawei organization on Github. The Dockyard is a combo storage service including container registry and artifact repository. The storage of docker images and binary files is an important part of any DevOps process, it's even more important in ContainerOps, which is a DevOps Orchestration engine powered by containers.

The version of Dockyard in the ContainerOps repository is shipped with the minimal functions, and we will develop it continuously as the next release. The repository in the Huawei organization is frozen for a while, and this version will be migrated to the repository finally in the feature.


### Getting started

#### Notice! Dockyard only support Golang 1.8 or above!
Since we introduced [graceful shutdown](https://beta.golang.org/doc/go1.8#http_shutdown), the new feature of Golang 1.8, the source code SHOUL be built with Go 1.8 or higher. If you really need to compile it in a former Go version, just find the incompatible code and delete them in `cmd/daemon.go` :)

#### Initialize the database
Dockyard starts with a database(currently we only support MySQL):
``` bash
CREATE DATABASE containerops_dockyard DEFAULT CHARACTER SET utf8 COLLATE utf8_bin;
```
And then call the `database` subcommand with action `migrate`:
``` bash
./dockyard database migrate
```

#### Start dockyard daemon
Now all the tables are created, you can start the daemon by:
``` bash
./dockyard daemon start
```

#### The config file
Dockyard reads the configs from a `.toml` file. You can specify the config file path with option or `--config`:
``` bash
./dockyard daemon start --config PATH_TO_CONFIG
```

If the path is not given, dockyard will take the default path `./conf/runtime.toml`

A config file is like this(you can find it under `./conf/runtime.toml.example`):

```toml
[database]
driver = "mysql"
host = "127.0.0.1"
port = 3306
user = "root"
password = "root"
db = 'containerops_dockyard'

[web]
#dockyard support only https and unix
mode = 'https'
address = "127.0.0.1"
port = 8990
cert = "./cert/nginx.crt"
key = "./cert/nginx.key"
# For unix socket listening:
# mode = 'unix'
# address = "/var/run/containerops/dockyard.socket"

```

You can also override the address and port by passing command line arguments:
``` bash
./dockyard daemon start -p 8990 -a 127.0.0.1
```

For more details of the command line arguments, just type the sub command without any action:
```
./dockyard daemon
```

### Put dockyard behind a proxy server
You might put dockyard behind a proxy server(like Nginx, Caddy etc.), because of the design of docker registry API, you'll have to take care of the header forwarding, you should pass the `scheme` and `host` headers to dockyard, or dockyard might not work as expected. 

What's more, since most of the proxy servers would have a default request body size limit, it's better to make it larger.

The example of a nginx config:
```
    server {
      listen    443;
      server_name    containerops.io;
      ssl on;
      ssl_certificate ssl/containerops.crt;
      ssl_certificate_key ssl/containerops.key;

      client_max_body_size 200M;  # Set it as you wish
      location / {
        proxy_pass    http://unix:/var/run/dockyard.socket;
        proxy_set_header X-Forwarded-Proto $scheme; # This is essential for nginx!
        proxy_set_header Host $http_host;  # This is essential for nginx!
      }
    }
```
