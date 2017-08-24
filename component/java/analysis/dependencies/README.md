## Java Gradle Dependencies Component

### What's the Component?

This image is java runtime image, used for print dependencies 
<br>
<br> grep dependencyReport

### Learn how to build it?

Use the `docker build` command build the image, and your project must build with gradle

```
docker build -t containerops/analysis/java_gradle_dependencies ./
```
### Component Usage
```
docker run \
    --rm \
    --env CO_DATA=" \
    version=gradle3 \
    git-url=https://github.com/vanniuner/gradle-demo.git" \
    containerops/analysis/java_gradle_dependencies \
```

### Parameters 
- `version` gradle version available value: gradle3,gradle4
- `git-url` where your code is located
### Versions 1.0.0



