# ContainerOps - DevOps Orchestration

### Why DevOps orchestration?

There are many tools, projects, plugins, services adopted in the DevOps workflow. However, no one can cover all DevOps tasks. When developers move from one tool to another, they are facing the huge risk of investment. How should we promote the DevOps process and make iteration more speedy? The principle of DevOps orchestration is to keep your original process working without any changes and just assemble tools, projects, plugins, services into orchestration engine.  Improvements to the DevOps process through by adding or replacing tools, projects, plugins, services via orchestration engine without interrupting the process. DevOps needs a step by step promotion, and invasive change is dangerous. 

### ContainerOps is a DevOps orchestration platform with the container technology.

ContainerOps is DevOps orchestration platform with the container. It has an engine orchestrating tools, components or services, and running with Kubernetes. ContainerOps provides tools encapsulate plugins or any programs into a container, and a set of environment variables used for interaction with the engine. We call this container included DevOps task is component. All components running in the Kubernetes, and the lifestyle managed by the engine. At the same time, the engine integrated with DevOps services like Github or Travis CI, and interacted with them through REST API. The ContainerOps designed for cloud native App development and running within container cluster.

### Concept - Component 

A DevOps function developed in any programming languages like Bash, Golang, Python or C++ could be encapsulated as a containerized DevOps Component. The Components guarantees that the DevOps task always runs the same way, regardless of its environment. Moreover, the developers do not care the resource of running Components, and the Orchestration Engine uses Kubernetes to run them. It is easier to use and share than scripts. We are working on a Component Registry to invigorate the community. The component registry will release in Q1, 2017 at https://opshub.sh. 

The Component has a set of environment variables. Some variables value is JSON data, and they are the input values which are the dependency of the running component. Some variables value is REST API URL, and the program posts data to interact with Orchestration Engine. 

### Concept - Event

There is a special environment variable named CO_DATA, the value is JSON data. When the developer designs a component, customize the JSON what's the component running requirements. The JSON data is the component input event. The CO_DATA is not the only environment variable which provides input data, and the developer set any global environment variables for all elements of workflow could read or write.

There is another special environment variable named CO_TASK_RESULT, the value is REST API URL. When the developer designs the component, customize the output JSON data post to the URL in the HTTP body. The JSON data is the component output event.

Orchestration engine collects the input data and output data, and transfer from one or more output event to one input event. The developer designs the DevOps workflow, at the same time designs the component event flow by drag and draw. The engine map the one or more output JSON data to input JSON data of component automatically, the developer could change the mapping relation in the workflow designer.

![Event Linking](docs/images/event-link.jpg)

However, the engine could not merge the conflict events mapping. It shows all conflicts, only remove this could be saved to run.

![Event Conflict](docs/images/event-conflict.jpg)

### Concept - Workflow

Usually, we define different stages in the DevOps orchestration. The stage is the phase of DevOps workflow, and the developer could define any number of the stage according to organizational structure, project property, programming language, other factors. How define stage is a philosophy, you should define follow the best practices in DevOps.

The ContainerOps define any task on stage named action. The action links the component or a service. All tasks of stage execute at the same time, and the events deliver to actions of another stage.

![Workflow Running](docs/images/workflow-running.gif)

### ContainerOps VS. Jenkins

1. Jenkins servers become snowflakes. Maintain a DevOps workflow need past shell script into the textboxes of UI, upload and install multiple plugins, install dependencies in the slaves.The ContainerOps cloud clearly defines the DevOps workflow drawing and dragging. It uses the component endpoint URL and doesn't need to download and install. All dependencies already include in the component, doesn't need install.

2. Jenkins has no first class support for pipelines, and Jenkins 2.0 tries to address this by promoting a Pipeline plugins. However, it misses the point of DevOps. There are many tools and services throughout the whole DevOps process and each with the pipeline of their own. A pipeline does not resolve the split DevOps tools chain and difficult to pass the context in the whole workflow. These are the problems the ContainerOps want to resolve.

3. Plugins accomplish everything you care about in Jenkins. The ContainerOps aims never to have plugins instead of with component. The component already has dependencies. The resources of running a component handle by ContainerOps engine through the Kubernetes.

## Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. 

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
660 York Street, Suite 102,
San Francisco, CA 94110 USA

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.

Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

## Format of the Commit Message

You just add a line to every git commit message, like this:

    Signed-off-by: Meaglith Ma <maquanyi@huawei.com>

Use your real name (sorry, no pseudonyms or anonymous contributions.)

If you set your `user.name` and `user.email` git configs, you can sign your
commit automatically with `git commit -s`.

