## ContainerOps - DevOps Orchestration

### Why DevOps orchestration?

There are many tools, projects, plugins, services adopted in the DevOps workflow. However, no one can cover all DevOps tasks. When developers move from one tool to another, they are facing the huge risk of reinvestment. How should we promote the DevOps process and make iteration more speedy? The principle of DevOps orchestration is to keep your original process working without any changes and just assemble tools, projects, plugins, services into orchestration engine. Improvements to the DevOps process can be done by gradually adding or replacing tools, projects, plugins or services for a smooth migration, and overall, DevOps needs to be promoted step by step, invasive changes are dangerous.

### How DevOps Orchestration?

Combo the different DevOps services, tools and plugins to implementing DevOps orchestration is very complex, and it should resolve many challenges like deliver data between jobs, or resolve environment consistency for tools or plugins.

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

### Contribute

If you interest the ContainerOps and want to get involved in developing. Getting start with this reading: 

* The [contributor guidelines](CONTRIBUTING.md)

### Community

### Community

