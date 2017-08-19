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
docker build -t huawei/compile/java_gradle_war ./
```

## DETAIL RUN
```bash
docker run --rm --env CO_DATA="git-url=https://github.com/rominirani/GradleWebAppSample.git target=https://hub.opshub.sh/binary/v1/lidian/test/binary/2.2.4/web.war" huawei/compile/java_gradle_war
```