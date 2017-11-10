## ContainerOps Components Test Flow

### What's the Flow?
Build And Test ContainerOps Components with orchestration engine

### Component Usage
```
$ /.ci-component-run.sh

or
  
pilotage cli run cncf-demo.yaml --verbose true --timestamp true

```
### Flow Configuration File
```
uri: containerops/auto-test-for-components/ci-component-ctest
title: CI build and test  for components
version: 1
tag: latest
timeout: 0
receivers:
  -
    type: mail
    address: lidian@containerops.sh
  -
    type: mail
    address: hubopsnotifier@gmail.com
stages:
  -
    type: start
    name: start
    title: Start
  -
    type: normal
    name: build-component
    title: Component auto test action
    sequencing: sequence
    actions:
      -
        name: build-component
        title: build-component-action
        jobs:
          -
            type: component
            kubectl:
            endpoint: hub.opshub.sh/containerops/component-ctest-build:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: ""
      -
        name: flow-component
        title: flow-component-action
        jobs:
          -
            type: component
            kubectl:
            endpoint: hub.opshub.sh/containerops/component-ctest-flow:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: ""
      -
  -
    type: end
    name: end
    title: End
```


### Versions 1.0.0
