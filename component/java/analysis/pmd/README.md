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
docker build -t containerops/analysis/java_gradle_pmd -f analysis/pmd/Dockerfile
```

## DETAIL RUN
```bash
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git out-put-type=json" containerops/analysis/java_gradle_pmd
```