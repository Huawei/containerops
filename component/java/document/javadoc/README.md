## Java Gradle Javadoc Component

### What's the Component?

This image is java runtime image, used for generate a javadoc compressed file

### Learn how to build it?

Use the `docker build` command build the image, and your project must build with gradle

```
docker build -t containerops/document/java_gradle_javadoc ./
```
### Component Usage
```
docker run \
    --rm \
    --env CO_DATA=" \
    git-url=https://github.com/vanniuner/gradle-demo.git \
    target=https://hub.opshub.sh/binary/v1/lidian/test/binary/1.1.0/javadoc.tar" \
    containerops/document/java_gradle_javadoc
```

### Parameters 
- `git-url` where your code is located
- `target`  where your package file to upload, curl -i -X PUT -T file target
### Versions 1.0.0



