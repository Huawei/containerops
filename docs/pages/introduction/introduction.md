---
title: ContainerOps - DevOps Orchestration  
keywords: introduction
tags: [introduction]
sidebar: home_sidebar
permalink: introduction.html
summary: ContainerOps is a DevOps orchestration for cloud native. Defining a DevOps component base container like Docker or rkt. Drawing the DevOps workflow with a WYSIWYG editor in the browser mixed DevOps components and the exist DevOps services like Github, Jenkins, Travis CI and so on. Running the components with container orchestration like Kubernetes. ContainerOps introduction can help you set up ContainerOps, learn about the system, and get your DevOps workflow running on Kubernetes.  
---

## ContainerOps - DevOps Orchestration

There is a lot of projects, plugins and services implemeted build, test and deployment in the DevOps workflow. But no one cover all tasks, we could only choose different projects and services together assembling into DevOps workflow. How to determine the order of tools and pass the data between them?  

ContainerOps is DevOps orchestration. It build a DevOps workflow orchestrating tools or services, and passing the trigger data or customize data from the workflow start trigger to the every stage. So don’t need abandon the original DevOps solution, just add ContainerOps as the top orchestration, scheduling the existing build system, continuous integration system or continuous deployment system, adding any features used and reused container DevOps component.

{% raw %}
{% include image.html file="docs/images/workflow-running.gif" url="https://containerops.sh" alt="ContainerOps" caption="" %}
{% endraw %}