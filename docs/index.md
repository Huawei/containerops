---
title: Getting Started With ContainerOps
keywords: homepage
tags: [getting_start]
sidebar: home_sidebar
permalink: index.html
summary:
---

## ContainerOps - DevOps Orchestration

### Why DevOps orchestration?

There are many tools, projects, plugins, services adopted in the DevOps workflow. However, no one can cover all DevOps tasks. When developers move from one tool to another, they are facing the huge risk of reinvestment. How should we promote the DevOps process and make iteration more speedy? The principle of DevOps orchestration is to keep your original process working without any changes and just assemble tools, projects, plugins, services into orchestration engine. Improvements to the DevOps process can be done by gradually adding or replacing tools, projects, plugins or services for a smooth migration, and overall, DevOps needs to be promoted step by step, invasive changes are dangerous.

{% include image.html file="CloudNativeLandscape_v0.9.2.jpg" url="https://github.com/cncf/landscape" alt="Cloud Native Landscape" caption="Cloud Native Landscape" %}

### ContainerOps is a DevOps orchestration platform with the container technology.

ContainerOps is a DevOps orchestration platform built with containers. It has an orchestrating engine to drive components or services, and runs on Kubernetes. ContainerOps provides tools to encapsulate plugins or any programs into a container, and a set of environment variables are used for interaction with the engine. We call this container encapsulated DevOps task a component. All components run in Kubernetes, and the lifecycle of the tasks are managed by the engine. At the same time, the engine can be integrated with DevOps services like Github or Travis CI through REST API. The ContainerOps is designed for cloud native app development and runs within container cluster.