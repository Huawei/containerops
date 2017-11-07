## ContainerOps Components Build Flow Component

### What's the Component?
Build  ContainerOps Components and push them into repository of hub.opshub.sh. At same time, the component 

### Learn how to build it?

```
docker build -t containerops/component/component-ctest-build ./build

```
### Component Usage
```
docker run \
    --rm \
    --env CO_DATA=" \
    git-url=https://github.com/Huawei/containerops.git " \
    containerops/component/component-ctest-build \ 
```

### Parameters 

- `git-url` where your code is located

### Versions 1.0.0
