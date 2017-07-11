## Build, Release golint component

```bash
docker build -t containerops/golint .
```


```bash
docker run --env CO_DATA="coderepo=https://github.com/haijunTan/gohello.git" containerops/golint:latest
```
