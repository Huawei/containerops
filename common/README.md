## Common - Public modules of ContainerOps

### What is Common?

`Common` is the public modules of `Containerops` . For now, it includes flowing parts :
 
 * **utils** 
    General utils such as dir/file operations, parameter validation etc.
 * **configuration**
    Configuration module of containerops
 
 
#### How to set configurations of ContainerOps

The format of `ContainerOps` configuration file is [`.toml`](https://github.com/toml-lang/toml). To make configurations file work, user should specific configuration file path by `--config` option or put default configuration file `containerops.toml` in one of `/etc/containerops/config` or `$HOME/.containerops/config` or binary execution path.  
 

##### 1. Configurations of database.
```toml
[database]
driver = "mysql"
host = "127.0.0.1"
port = 3306
user = "root"
password = "containerops_database"
db = "containerops_password"
```


#####  2. Configurations for HTTPS or Unix Socket
######    2.1 If multi modules deploy in one node, there should have a proxy like Caddy or Nginx.
```toml
[web]
mode = "unix"
address = "/var/run/${module}.socket"
```
######    2.2 If module deploys in one node alone, it only supports HTTPS model and must have the SSL certification files.
```toml
[web]
domain = "opshub.sh"
mode = "https"
address = "127.0.0.1"
port = 443
cert = "PATH_TO_CERT_FILE"
key = "PATH_TO_KEY_FILE"
```

##### 3. Configurations for storage path of Dockyard module.
######    3.1 TODO Using the Object Storage Service in the Dockyard module.
```toml
[storage]
dockerv2 = "/tmp/dockerv2" # path for image files of Docker Distribution V2 Protocol
binaryv1 = "/tmp/binaryv1" # path for binary files of Dockyard Binary V1 Protocol

```
#####  4. Configurations for Warship of Dockyard client.
```toml
[warship]
domain = "hub.opshub.sh"
```
#####  5. Configurations for Singular modules.
```toml
[singular]
```
#####  6. Configurations for Mail Notifier (Optional).
```toml
[mail]
smtp_address = "smtp.gmail.com"
smtp_port = "587"
user = "notify@containerops.sh"
password = "password"
```
