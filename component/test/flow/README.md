## ContainerOps Components Test Flow Component

### What's the Component?
Push YML File of ContainerOps Components to The orchestration engine service for auto-testing.

### Learn how to build it?

```
docker build -t containerops/component/component-ctest-flow ./

```
### Component Usage
```
docker run \
    --rm \
    --env CO_DATA=" \
    git-url=https://github.com/Huawei/containerops.git " \
    containerops/component/component-ctest-flow \ 
```

### Parameters 

- `git-url` where your code is located

### Versions 1.0.0
