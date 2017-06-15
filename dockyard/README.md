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
Dockyard reads the configs from a `.toml` file. You can specify the config file path with option `-c` or `--config`:
``` bash
./dockyard daemon start -c PATH_TO_CONFIG
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
