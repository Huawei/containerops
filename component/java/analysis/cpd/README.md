## Java Gradle Checkstyle Component

### What's the Component?

This image is java runtime image, used for find theDuplicate code 

### Learn how to build it?

Use the `docker build` command build the image, and your project must build with gradle

```
docker build -t containerops/analysis/java_gradle_cpd ./
```
### Component Usage
```
docker run \
    --rm \
    --env CO_DATA=" \
    git-url=https://github.com/vanniuner/gradle-demo.git \
    out-put-type=xml" 
    containerops/analysis/java_gradle_cpd \
```

### Parameters 
- `git-url` where your code is located
- `out-put-type`  xml,yaml,json
- `report-path`   not required,if you defined reports path
### Versions 1.0.0



