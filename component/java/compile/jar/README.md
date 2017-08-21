## Java Gradle Jar Component

### What's the Component?

This image is java runtime image, used for compile your project to a jar file, and upload it to the target

### Learn how to build it?

Use the `docker build` command build the image, and your project must build with gradle

```
docker build -t containerops/compile/java_gradle_jar ./
```
### Component Usage
```
docker run \
    --rm \
    --env CO_DATA=" \
    git-url=https://github.com/vanniuner/gradle-demo.git \
    target=https://hub.opshub.sh/binary/v1/containerops/component/binary/2.2.4/demo.jar" \
    huawei/compile/java_gradle_jar
```

### Parameters 
- `git-url` where your code is located
- `target`  where your package file to upload, curl -i -X PUT -T file target
### Versions 1.0.0



