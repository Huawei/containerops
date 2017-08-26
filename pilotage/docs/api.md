#API spec of pilotage


### POST  /flow/v1/:namespace/:repository/:flow/:tag/:type

receive the definition file of a `flow` and execute   

#### Request

- **Syntax:**
```http
POST  /flow/v1/:namespace/:repository/:flow/:tag/:type HTTP/1.1
```

```
flow definition file content
```

#### Response On Success

- **Syntax:**
```
HTTP/1.1 201 Created
Content-Type: application/json
```

```json
{
  "id": "abcd-123",
  "namespace": "cncf",
  "repository": "kubernetes",
  "name": "kubernetes-flow",
  "tag": "v1",
  "title": "Demo For pilotage",
  "version": "4",
  "status": "Running"
}
```