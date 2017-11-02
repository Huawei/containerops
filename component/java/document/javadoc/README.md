## Java Gradle Javadoc Component

### What's the Component?

This image is java runtime image, used for generate a javadoc compressed file, and upload it to the target
<br>
<br> gradle javadoc

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
    version=gradle3 \ 
    git-url=https://github.com/wangkirin/demo-bmi.git \
    target=https://hub.opshub.sh/binary/v1/lidian/test/binary/1.1.0/javadoc.tar \
    build-path=./calculator" \
    containerops/document/java_gradle_javadoc
```

### Parameters 
- `version` gradle version available value: gradle3,gradle4
- `git-url` where your code is located
- `target`  where you want the file to upload, curl -i -X PUT -T file target
- `build-path` not required which item package you want 
### Versions 1.0.0



