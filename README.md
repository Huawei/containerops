# ContainerOps - DevOps Orchestration

### Why DevOps orchestration?
There are countless tools, projects, plugins, services adopted in the DevOps workflow. However, no one can cover all DevOps tasks. When developers move from one tool to another ,they are facing huge risk of investment. How should we promote the DevOps process and make iteration more speedy? The principle of DevOps orchestration is to keep your original process working without any changes and just assemble tools, projects, plugins, services into orchestration engine.  Improvements of the DevOps process could be done by adding or replacing tools, projects, plugins, services via orchestration engine without interrupting the process. DevOps needs a step by step promotion, invasive change is dangerous.

### ContainerOps is a DevOps orchestration platform with the container technology.

ContainerOps is a container-based DevOps orchestration platform. It comprises an Orchestration Engine, couple of Components or Services, which are all running in the Kubernetes cluster. ContainerOps provides a serials of tools encapsulating plugins or any programing languages into a container,  a sort of mechanism for interaction with Orchestration Engine by environment variables. We call this specific Container including DevOps task as Component. The lifecycle of all Components  is managed by the Orchestration Engine. At the same time, the Engine can be integrated with  the third-party DevOps services like Github or Travis CI, and interacts with them via REST API. The ContainerOps is designed for Cloud-Native App development and deployment within Kubernetes cluster.

### Concept - Component 

A DevOps function developed in any programming language like Bash, Golang, Python or C++ could encapsulate as a containerized DevOps component. This guarantee that the DevOps task always run the same, regardless of its environment. Moreover, the developers do not care the resource of running components, and the orchestration engine use Kubernetes run them. It is more easily use than scripts and share. We are working on a component registry to invigorate the community. The registry will release in Q1, 2017 at https://opshub.sh. 

The component has a set of environment variables. Some variables value is JSON data, and they are the input values which are the dependency of the component running. Some variables value is REST API URL, the program post data to interact with orchestration engine.

### Concept - Event

There is a special environment variable named CO_DATA, the value is JSON data. When the developer designs a Component, customize the JSON what's the component running requirements. The JSON data is the Component input event. The CO_DATA is not only the environment variable which provides input data, and the developer set any global environment variables for all elements of workflow could read or write.

There is another special environment variable named CO_TASK_RESULT, the value is REST API URL. When the developer designs the Component, customize the output JSON data post to the URL in the HTTP body. The JSON data is the Component output event.

Orchestration Engine collects the input data and output data, and transfer from one or more outputs event to the input event. The developer designs the DevOps workflow, at the same time designs the Component event flow by drag and draw. The Engine maps  one or more output JSON data to input JSON data of Component automatically, the developer could change the mapping relation in the workflow designer.

![Event Linking](docs/images/event-link.jpg)

However, the Engine could not merge the conflict events mapping. It shows all conflicts, only removing this could be saved to run.

![Event Conflict](docs/images/event-conflict.jpg)

### Concept - Workflow

Usually, we define different stages in the DevOps orchestration. The stage is the phase of DevOps workflow, and the developer could define any number of the stage according to organizational structure, project property, programming language, other factors. How define stage is a philosophy, you should define follow the best practices in DevOps.

The ContainerOps define any task on stage named action. The action links the component or a service. All tasks of stage execute at the same time, and the events deliver to actions of another stage.

![Workflow Running](docs/images/workflow-running.gif)

### ContainerOps VS. Jenkins

Jenkins need shell scripts and dependencies, and the developer should maintain the pipeline of jobs.  The component of ContainerOps is easily maintaining.

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

