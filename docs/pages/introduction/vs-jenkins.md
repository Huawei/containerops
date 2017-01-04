---
title: ContainerOps VS Jenkins
keywords: introduction
tags: [introduction]
sidebar: home_sidebar
permalink: vs-jenkins.html
summary: Introduction Of ContainerOps
folder: introduction  
---

## ContainerOps VS Jenkins

1. Jenkins servers become snowflakes. Maintaining a DevOps workflow needs pasting shell script into the textboxes of UI, uploading and installing multiple plugins and dependencies in the slaves. The ContainerOps could clearly define the DevOps workflow by drawing and dragging through a GUI IDE. It uses the component endpoint URLs and doesn't need to download and install. All dependencies are already self-contained in the component.


2. Jenkins has no first class support for pipelines, and Jenkins 2.0 tries to address this by introducing a Pipeline plugins. However, it misses the point of DevOps. There are many tools and services throughout the whole Jenkins process, and each has its own pipeline. With the cloud native app becomes more and more sophisticated, this increasing complexity makes tools to be isolated in the chain and finding a good way to pass context in the whole workflow becomes critical. These are the problems ContainerOps wants to resolve.


3. Plugins accomplish everything you care about in Jenkins, and the plugin has dependencies on the environment. The ContainerOps use container-encapsulated components instead of plugins, and all dependencies have been self-contained in the component. The resources of running a component are handled by the ContainerOps engine through Kubernetes.