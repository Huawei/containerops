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
docker build -t containerops/test/java_gradle_junit ./
```

## DETAIL RUN
```bash
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/test/java_gradle_junit
```