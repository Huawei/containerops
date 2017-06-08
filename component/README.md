## DevOps Component

### What's the DevOps Component?

The DevOps component is a container image which encapsulates one or more programs completed DevOps task.

### Why choose the Phusion BaseImage as the base image?

### How to transfer data to the component inside with environment variables?

### How to collecting the data from stdout/stderr?

#### Exit Code

The ContainerOps engine checks the exit code of the process determined the result.

* _0_ - Successful termination
* _1_ - Catchall for general errors
* _2_ - Misuse of shell builtins (according to Bash documentation)
* _64_ - Command line usage error
* _65_ - Data format error
* _66_ - Cannot open input   
* _67_ - Addressee unknown
* _68_ - Host name unknown
* _69_ - Service unavailable
* _70_ - Internal software error
* _71_ - System error (e.g., can't fork)
* _72_ - Critical OS file missing
* _73_ - Can't create (user) output file
* _74_ - Input/Output error
* _75 - Temp failure; user is invited to retry
* _76_ - Remote error in protocol
* _77_ - Permission denied
* _78_ - Configuration error
* _126_ - Command invoked cannot execute
* _127_ - Command not found
* _128_ - Invalid argument to exit
* _128+n_ - Fatal error signal "n."
* _255_ - Exit status out of range (exit takes only integer args in the range 0 - 255)

If the component has one or more NodeJS programs, there is an exit code [list](https://github.com/nodejs/node-v0.x-archive/blob/master/doc/api/process.markdown#exit-codes). 

The developers of Component could define exit code by themselves also and should paste it in README.

### What's the official DevOps component?
