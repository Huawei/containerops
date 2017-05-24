## Bazel Docker Image

### What's the Bazel?

> Bazel is a build tool which coordinates builds and run tests. It works with source files written in any language, with native support for Java, C, C++ and Python. Bazel produces builds and runs tests for multiple platforms.

[Bazel](https://bazel.io) is Google's build tool, it has build-in support for building both client and server software. It also provides an extensible framework that you can use to develop your own build rules.

* _Speed_
* _Scalability_
* _Flexibility_
* _Correctness_
* _Reliability_
* _Repeatability_

### How to build the image?

Use the `docker build` command to build the image.

```bash
docker build -t containerops/bazel .
```

### How to use the image?

The image just is a base image, we don't set `WORKSPACE` location. You should use it as the base image, like this:
 
```dockerfile
FROM containerops/bazel:latesd

RUN git clone https://github.com/kubernetes/kubernetes.git /var/opt/src/kubernetes
WORKDIR /var/opt/src/kubernetes

RUN make bazel-test
```


