## Java Gradle Testng Component

### What's the Component?

This image is java runtime image, used for generate a testng report
<br>
<br> gradle test

### Learn how to build it?

Use the `docker build` command build the image, and your project must build with gradle

```
docker build -t containerops/test/java_gradle_testng ./
```
### Component Usage
```
docker run \
    --rm \
    --env CO_DATA=" \
    version=gradle3 \
    git-url=https://github.com/vanniuner/gradle-demo.git \
    out-put-type=json" \
    containerops/test/java_gradle_testng \
```

### Parameters 
- `version` gradle version available value: gradle3,gradle4
- `git-url` where your code is located
- `out-put-type`  xml,yaml,json
- `report-path`   not required,if you defined the testng reports path
### Versions 1.0.0



