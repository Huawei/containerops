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
docker build -t huawei/document/java_gradle_javadoc ./
```

## DETAIL RUN
```bash
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git target=https://hub.opshub.sh/binary/v1/lidian/test/binary/1.1.0/javadoc.tar" huawei/document/java_gradle_javadoc
```