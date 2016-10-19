---
title: What's Customized DevOps Component 
keywords: component
tags: [component]
sidebar: home_sidebar
permalink: concept-component.html
summary: What's Customized DevOps Component
folder: component
---

## What's Customized DevOps Component

The Customized DevOps Component is the core concept in the ContainerOps, it's the part **DevOps Base Container**. Usually developer write DevOps plugins with *Bash*, *Python* and other script languages, but it has *runtime environment consistency* problem. When we got a script, we must setup a correct runtime environment to execute, it's a trivial and easy makes mistakes job. Now the container just fitly solved the *environment consistency* problem. Developer would use *Golang*, *Java*, *Rust* or any program languages or tools for the DevOps tasks, and share with runtime environment in contaienr image format. When we got the a container image, we don't care the environment at all and just care the function. 

The customized DevOps component is a container image include DevOps programs or tools and runtime environment, it has specify inputs with environment parameters or REST API interface and specify outputs with callback interface.Now we use [Docker Registry Image specification](https://github.com/docker/distribution/blob/master/docs/spec/manifest-v2-2.md), and will support [ACI](https://github.com/appc/spec/blob/master/spec/aci.md) and [OCI Image specification](https://github.com/opencontainers/image-spec) soon.

We don't follow the [pods concept of Kubernetes](http://kubernetes.io/docs/user-guide/pods), but we execute the component in a pod of Kubernetes. We wanna to keep the component simply for sharing, it's very important for the community. If a DevOps task is complex which need more than one component, you should use multiple **Stages** and **Action**, and define the workflow in the Pilotage.

The component definition is open, we are glad to hear you feedback. Please open an [issue](https://github.com/containerops/pilotage/issues) to discuss.  
