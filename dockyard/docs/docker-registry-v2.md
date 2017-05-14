## Docker Registry V2 Protocol

Offical Docker Registry V2 Doc is [here](https://github.com/docker/distribution/blob/master/docs/spec/api.md).

### Docker Registry V2 Ping

1.1 (Docker Client -> Docker Registry) `GET /v2` 
> https://github.com/docker/distribution/blob/master/docs/spec/api.md#api-version-check

### Docker Registry V2 Push

2.1 Push the all images to the registry.

(Docker Client -> Docker Registry) `HEAD /:namespace/:repository/blobs/:digest`
(Docker Client -> Docker Registry) `PUT /:namespace/:repository/blobs/uploads`
(Docker Client -> Docker Registry) `PATCH /:namespace/:repository/blobs/uploads/:uuid`
(Docker Client -> Docker Registry) `PUT /:namspace/:repository/blobs/uploads/:uuid`

2.2 Push the tag information.

(Docker Client -> Docker Registry) `PUT /:namespace/:repsitory/manifests/:tag`

### Docker Registry V2 Pull 

3.1 Get the tag list and information.

(Docker Client -> Docker Registry) `GET /:namespace/:repository/tags/list`
(Docker Client -> Docker Registry) `GET /:namespace/:repository/manifests/:tag`

3.2 Get the images.

(Docker Client -> Docker Registry) `GET /:namespace/:repository/blobs/:digest`

### Delete Methods

(Docker Client -> Docker Registry) `DELETE /:namespace/:repository/blobs/:digest`
(Docker Client -> Docker Registry) `DELETE /:namespace/:repository/manifests/:reference`
(Docker Client -> Docker Registry) `DELETE /:namespace/:repository/blobs/:uuid`
