# API spec of assembling


### POST  /assembling/build

Receive the Dockerfile or archive file to build the image in an isolated environment.

#### Request

- **Syntax**
```http
HTTP/1.1
POST  /assembling/build
```

- **Parameters**

|Parameter|type|From|Description|
--|--|--
|registry|string|Query| The registry where to push the built image |
|namespace|string|Query| The namespace |
|image|string|Query| The image name|
|tag|string|Query| The tag of the image|
|buildargs|json string|Query| Arguments that will be sent to docker daemon when build, equivalent to the `--build-arg `option of `docker build`|
|insecure_registry|array of string|Query| Specify insecure registries when starting the docker daemon in dind image, user can specify multiple insecure registries|
|autstr|string|Query| A base64 encoded auth string, more details could be found [here](https://docs.docker.com/engine/api/v1.30/#section/Authentication) |
|Dockerfile & archive file|binary|Body| The body, see below |

> The body should be a single Dockerfile or, if some extra files are included, an archive file, which can be in the format of tar, gzip, bzip2 or xz, according to the [Docker Engine API](https://docs.docker.com/engine/api/v1.31/#operation/ImageBuild), section `REQUEST BODY`.



- **Example**

```http
POST /assembling/build?registry=hub.opshub.sh&namespace=containerops&image=ubuntu&tag=14.04&buildargs={"RELEASE":"1.0.1", "DEV_TAG":"FOXHOUND"}&insecure_registry=hub1.myhost&insecure_registry=hub2.myhost&authstr=0bdaf9cf38....

Body(Binary file)

```

#### Response On Success

- **Syntax:**
```
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
  "endpoint" : "hub.opshub.sh/containerops/ubuntu:14.04"
}
```
