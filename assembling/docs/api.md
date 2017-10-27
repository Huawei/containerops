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

|Parameter|type|From|
--|--|--
|registry|string|Query|
|namespace|string|Query|
|image|string|Query|
|tag|string|Query|
|buildargs|json string|Query|
|Dockerfile & archive file|binary|Body|

> The body should be a single Dockerfile or, if some extra files are included, an archive file, which can be in the format of tar, gzip, bzip2 or xz, according to the [Docker Engine API](https://docs.docker.com/engine/api/v1.31/#operation/ImageBuild), section `REQUEST BODY`.



- **Example**

```http
POST /assembling/build?registry=hub.opshub.sh&namespace=containerops&image=ubuntu&tag=14.04&buildargs={"RELEASE":"1.0.1", "DEV_TAG":"FOXHOUND"}

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
