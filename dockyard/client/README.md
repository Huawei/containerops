## How To Use Warship

### Set configuration file

In your `$HOME/.containerops/config` folder create a file named **containerops.toml**. Add **warship** section like this:

```toml
[warship]
domain = "hub.opshub.sh"
```

### Compile the warship client

Run `go build warship.go` in `containerops/dockyard/client` to compile warship binary file, then run `cp warship /usr/local/bin` command. Also you could add the `$GOPATH/src/github.com/Huawei/containerops/dockyard/client` to your _$PATH_ .

### Create binary repository or docker image repository

```bash
warship create --type binary cncf/demo
```

* The **URI** patten is `<namespace>/<repository>`
* the **type** has _binary_ or _docker_ options.

### Upload binary file to hub

```bash
warship binary upload /home/meaglith/Downloads/code_amd64.deb containerops/cncf-demo/latest
```

* The **URI** patten is `<namespace>/<repository>/<tag>`

### Download binary file

```bash
warship binary download containerops/cncf-demo/latest/code_amd64.deb /home/meaglith
```

* The **URI** patten is `<namespace>/<repository>/<tag>/<filename>`