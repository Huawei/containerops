## Java Gradle Checkstyle Component

### What's the Component?

This image is java runtime image, used for analysis your java coding style. 
Checkstyle is a development tool to help programmers write Java code that adheres to a coding standard. It automates the process of checking Java code to spare humans of this boring (but important) task. This makes it ideal for projects that want to enforce a coding standard.
<br>
<br>gradle checkstyleMain
<br>gradle checkstyleTest

### Learn how to build it?

Use the `docker build` command build the image, and your project must build with gradle

```
docker build -t containerops/analysis/java_gradle_checkstyle ./
```
### Component Usage
```
docker run \
    --rm \
    --env CO_DATA=" \
    version=gradle3 \
    git-url=https://github.com/vanniuner/gradle-demo.git \
    out-put-type=json" \
    containerops/analysis/java_gradle_checkstyle \ 
```

### Parameters 
- `version` gradle version available value: gradle3,gradle4
- `git-url` where your code is located
- `out-put-type`  available value: xml,yaml,json
- `report-path`   not required,if you defined reports path,
### Versions 1.0.0



