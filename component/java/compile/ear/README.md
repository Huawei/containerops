## FAST BUILD

```bash
./build.sh
```

## FAST RUN

```bash
./run.sh
```

## DETAIL BUILD
```bash
docker build -t huawei/compile/java_gradle_ear ./
```

## DETAIL RUN
```bash
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git target=https://hub.opshub.sh/binary/v1/containerops/component/binary/2.2.4/demo.ear" huawei/compile/java_gradle_ear
```