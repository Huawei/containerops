## ContainerOps - DevOps Orchestration

### Why DevOps orchestration?

There are many tools, projects, plugins, services adopted in the DevOps workflow. However, no one can cover all DevOps tasks. When developers move from one tool to another, they are facing the huge risk of reinvestment. How should we promote the DevOps process and make iteration more speedy? The principle of DevOps orchestration is to keep your original process working without any changes and just assemble tools, projects, plugins, services into orchestration engine. Improvements to the DevOps process can be done by gradually adding or replacing tools, projects, plugins or services for a smooth migration, and overall, DevOps needs to be promoted step by step, invasive changes are dangerous.

### ContainerOps is a DevOps orchestration platform with the container technology.

ContainerOps is a DevOps orchestration platform built with containers. It has an orchestrating engine to drive components or services, and runs on Kubernetes. ContainerOps provides tools to encapsulate plugins or any programs into a container, and a set of environment variables are used for interaction with the engine. We call this container encapsulated DevOps task a component. All components run in Kubernetes, and the lifecycle of the tasks are managed by the engine. At the same time, the engine can be integrated with DevOps services like Github or Travis CI through REST API. The ContainerOps is designed for cloud native app development and runs within container cluster.

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

1. Jenkins servers become snowflakes. Maintaining a DevOps workflow needs pasting shell script into the textboxes of UI, uploading and installing multiple plugins and dependencies in the slaves. The ContainerOps could clearly define the DevOps workflow by drawing and dragging through a GUI IDE. It uses the component endpoint URLs and doesn't need to download and install. All dependencies are already self-contained in the component.

2. Jenkins has no first class support for pipelines, and Jenkins 2.0 tries to address this by introducing a Pipeline plugins. However, it misses the point of DevOps. There are many tools and services throughout the whole Jenkins process, and each has its own pipeline. With the cloud native app becomes more and more sophisticated, this increasing complexity makes tools to be isolated in the chain and finding a good way to pass context in the whole workflow becomes critical. These are the problems ContainerOps wants to resolve.

3. Plugins accomplish everything you care about in Jenkins, and the plugin has dependencies on the environment. The ContainerOps use container-encapsulated components instead of plugins, and all dependencies have been self-contained in the component. The resources of running a component are handled by the ContainerOps engine through Kubernetes.

### Architecture of project

The ContainerOps is a DevOps orchestration platform, and its architecture is micro services. All the codes in one repository, each service has an own folder and maintainer.

* **component** - Some components maintained by ContainerOps team.
* **pilotage** - The orchestration engine service.
* **scaffold** - The orchestration engine UI.
* **joints** - The orchestration engine UI of AngularJS version 1. 
* **tenant** - The UI of containerops.sh.
* **crew** - It's a RABC service.
* **dockyard** - It's a artifact & container repository.
* **assembling** - It's a aritfact & container build service.
* **scaler** - The third service manage service.
* **nucleus** - The component manage service.
* **singular** - The deployment and operations tools.
* **logarithm** - The log service collect from platform.
* **dashboard** - The dashboard UI.

### How to try?

### Contribute

If you interest the ContainerOps and want to get involved in developing. Getting start with this reading: 

* The [contributor guidelines](CONTRIBUTING.md)

### Community