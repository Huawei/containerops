## Java Gradle Jar Component

### What's the Component?

This image is java runtime image, used for compile your project to a jar file, and upload it to the target

### Learn how to build it?

Use the `docker build` command build the image, and your project must build with gradle
<br>
<br> gradle jar

```
docker build -t containerops/compile/java_gradle_jar ./
```
### Component Usage
```
docker run \
    --rm \
    --env CO_DATA=" \
    version=gradle3 \
    git-url=https://github.com/wangkirin/bmi.git \
    build-path=calculator
    target=https://hub.opshub.sh/binary/v1/containerops/component/binary/2.2.4/demo.jar" \
    containerops/compile/java_gradle_jar
```

### Parameters 
- `version` gradle version available value: gradle3,gradle4
- `git-url` where your code is located
- `build-path` the module which you want to build,default value is project root path  
- `target`  where your package file to upload, curl -i -X PUT -T file target
### Versions 1.0.0



