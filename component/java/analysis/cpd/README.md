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
docker build -t containerops/analysis/java_gradle_cpd -f analysis/cpd/Dockerfile
```

## DETAIL RUN
```bash
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=xml" containerops/analysis/java_gradle_cpd
```