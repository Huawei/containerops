## Java Gradle Pmd Component

### What's the Component?

This image is java runtime image, used for source code analyzer
PMD can finds common programming flaws like unused variables, empty catch blocks, unnecessary object creation, and so forth

### Learn how to build it?

Use the `docker build` command build the image, and your project must build with gradle

```
docker build -t containerops/analysis/java_gradle_pmd ./
```
### Component Usage
```
docker run \
    --rm \
    --env CO_DATA=" \
    git-url=https://github.com/vanniuner/gradle-demo.git \
    out-put-type=json" 
    containerops/analysis/java_gradle_pmd \
```

### Parameters 
- `git-url` where your code is located
- `out-put-type`  xml,yaml,json
- `report-path`   not required,if you defined reports path
### Versions 1.0.0



