# ContainerOps - DevOps Orchestration

### Why DevOps orchestration?

There are countless projects, plugins, services in the DevOps workflow. However, no one covers all DevOps tasks. The developers are facing huge risk and investment when transitioning from any tools to another. How should improve the DevOps process and increase iteration speed? The principle of DevOps orchestration is keeping the process and tools, assemble them into orchestration engine and deliver the context between them. Then improve the DevOps by add or replace through orchestration engine without interrupt the process. DevOps is a step by step process, and invasive change is dangerous.

### ContainerOps is a DevOps orchestration platform with the container technology.

ContainerOps is DevOps orchestration platform with the container. It has an engine orchestrating tools, components or services, and running with Kubernetes. ContainerOps provides tools encapsulate plugins or any programs into a container, and a set of environment variables used for interaction with the engine. We call this container included DevOps task is component. All components running in the Kubernetes, and the lifestyle managed by the engine. At the same time, the engine integrated with DevOps services like Github or Travis CI, and interacted with them through REST API. The ContainerOps designed for cloud native App development and running within container cluster.

### Concept - Component 

A DevOps function developed in any programming language like Bash, Golang, Python or C++ could encapsulate as a containerized DevOps component. This guarantee that the DevOps task always run the same, regardless of its environment. Moreover, the developers do not care the resource of running components, and the orchestration engine use Kubernetes run them. It is more easily use than scripts and share. We are working on a component registry at https://opshub.sh.

The component has environment variables. Some variables value is JSON data format, set by orchestration engine when running. Some variables are REST API URL, post data to interact with the engine.

### Concept - Event

![Event Linking](docs/images/event-link.jpg)

![Event Conflict](docs/images/event-conflict.jpg)

### Concept - Workflow

The ContainerOps has a DevOps workflow WYSIWYG editor in the browser. The user could define the DevOps workflow by drag or drop the line or set the data flow from workflow trigger to the end by the same way.

![Workflow Running](docs/images/workflow-running.gif)

### ContainerOps VS. Jenkins

### ContainerOps VS. Travis CI And Other Hosted Offering 

### ContainerOps VS. Concourse 

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

