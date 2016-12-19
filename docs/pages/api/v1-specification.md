---
title: V1 Specification
keywords: component
tags: [API]
sidebar: home_sidebar
permalink: v1-specification.html
summary: V1 Specification
---

## API V1 Operations

### createComponent

| HTTP Method |  Request Address |
| -------- | ------ |
| POST  |/v2/:namespace/component|

#### body

```
{
  "name": "pythoncheck",
  "version": "v1.0.1"
}
```

#### response json

```
{
  "message": "create new component success"
}
```
### deleteComponent

| HTTP Method |  Request Address |
| -------- | ------ |
| Delete  |/v2/:namespace/component/:component|

#### response json

```
{
  "message": "delete component success"
}
```

### changeComponent

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  |/v2/:namespace/component/:component|

#### body

```
{
  "id": 2,
  "version": "v1.0.1",
  "define": {
    "env": [
      {
        "key": "CO_DATA",
        "value": "{\"test\":true}"
      }
    ],
    "inputJson": {
      "gitUrl": "www"
    },
    "outputJson": {
      "status": true
    },
    "setupData": {
      "action": {
        "apiserver": "",
        "datafrom": ":3000/",
        "image": {
          "name": "pythoncheck",
          "tag": "1.0.1"
        },
        "ip": "",
        "name": "",
        "timeout": "2000",
        "type": "Kubernetes",
        "useAdvanced": false
      },
      "pod": {
        "spec": {
          "containers": [
            {
              "resources": {
                "limits": {
                  "cpu": 0.2,
                  "memory": "256Mi"
                },
                "requests": {
                  "cpu": 0.1,
                  "memory": "128Mi"
                }
              }
            }
          ]
        }
      },
      "pod_advanced": {},
      "service": {
        "spec": {
          "ports": [
            {
              "nodePort": 32001,
              "port": 8000,
              "targetPort": 8000
            }
          ],
          "type": "NodePort"
        }
      },
      "service_advanced": {}
    }
  }
}
```

#### response json

```
{
  "message": "success"
}
```

### get component list

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  |/v2/:namespace/component/list|

#### response json

```
{
  "list": [
    {
      "id": 2,
      "name": "pythoncheck",
      "version": [
        {
          "id": 2,
          "version": "v1.0.1",
          "versionCode": 1
        }
      ]
    },
    {
      "id": 1,
      "name": "busybox",
      "version": [
        {
          "id": 1,
          "version": "1.0.1",
          "versionCode": 1
        }
      ]
    }
  ]
}
```

### get component info

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  |/v2/:namespace/component/:componentName?id=:componentID|

#### response json

```
{
  "env": [
    {
      "key": "CO_DATA",
      "value": "{\"test\":true}"
    }
  ],
  "inputJson": {
    "gitUrl": "www"
  },
  "outputJson": {
    "status": true
  },
  "setupData": {
    "action": {
      "apiserver": "",
      "datafrom": ":3000/",
      "image": {
        "name": "pythoncheck",
        "tag": "1.0.1"
      },
      "ip": "",
      "name": "",
      "timeout": "3000",
      "type": "Kubernetes",
      "useAdvanced": false
    },
    "pod": {
      "spec": {
        "containers": [
          {
            "resources": {
              "limits": {
                "cpu": 0.2,
                "memory": "256Mi"
              },
              "requests": {
                "cpu": 0.1,
                "memory": "128Mi"
              }
            }
          }
        ]
      }
    },
    "pod_advanced": {},
    "service": {
      "spec": {
        "ports": [
          {
            "nodePort": 32001,
            "port": 8000,
            "targetPort": 8000
          }
        ],
        "type": "NodePort"
      }
    },
    "service_advanced": {}
  }
}
```

### saveComponentInfo

Save component as new version used the same api.

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  |/v2/:namespace/component/:componentName|

#### body

```
{
  "id": 2,
  "version": "v1.0.1",
  "define": {
    "inputJson": {
      "gitUrl": "www"
    },
    "outputJson": {
      "status": true
    },
    "setupData": {
      "action": {
        "type": "Kubernetes",
        "name": "",
        "timeout": "3000",
        "ip": "",
        "apiserver": "",
        "image": {
          "name": "pythoncheck",
          "tag": "1.0.1"
        },
        "useAdvanced": false,
        "datafrom": ":3000/"
      },
      "service": {
        "spec": {
          "type": "NodePort",
          "ports": [
            {
              "port": 8000,
              "targetPort": 8000,
              "nodePort": 32001
            }
          ]
        }
      },
      "pod": {
        "spec": {
          "containers": [
            {
              "resources": {
                "limits": {
                  "cpu": 0.2,
                  "memory": "256Mi"
                },
                "requests": {
                  "cpu": 0.1,
                  "memory": "128Mi"
                }
              }
            }
          ]
        }
      },
      "service_advanced": {},
      "pod_advanced": {}
    },
    "env": [
      {
        "key": "CO_DATA",
        "value": "{\"test\":true}"
      }
    ]
  }
}
```

#### response json

```
{
  "message": "success"
}
```

### Create Workflow

| HTTP Method |  Request Address |
| -------- | -------- |
| POST  | /v2/:namespace/:repository/workflow/v1/define |

#### Body:

```
{
  "name": "pythoncheck",
  "version": "v1.0.1"
}
```

#### response json

```
{
  "message": "create new workflow success"
}
```

### getWorkflowList

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /v2/:namespace/:repository/workflow/v1/define/list |

#### response json

```
{
  "list": [
    {
      "id": 2,
      "name": "pythoncheck",
      "version": [
        {
          "id": 2,
          "version": "v1.0.1",
          "versionCode": 1
        }
      ]
    },
    {
      "id": 1,
      "name": "busybox",
      "version": [
        {
          "id": 1,
          "status": {
            "status": true,
            "time": "2016-12-05 22:41:38"
          },
          "version": "1.0.1",
          "versionCode": 1
        }
      ]
    }
  ]
}
```


### getWorkflowInfo

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /v2/:namespace/:repository/workflow/v1/define/:workflowName?id=:workflowID |

#### response json

```
{
  "lineList": [
    {
      "endData": {
        "component": {
          "name": "pythoncheck",
          "versionid": 2,
          "versionname": "v1.0.1"
        },
        "env": [
          {
            "key": "CO_DATA",
            "value": "{\"test\":true}"
          }
        ],
        "height": 38,
        "id": "workflow-action-ac85dff0-bb63-11e6-be23-4d9fe51a1df7",
        "inputJson": {
          "gitUrl": "www"
        },
        "outputJson": {
          "status": true
        },
        "setupData": {
          "action": {
            "apiserver": "192.168.10.131:8080",
            "datafrom": ":3000/",
            "image": {
              "name": "pythoncheck",
              "tag": "1.0.1"
            },
            "ip": "192.168.10.131",
            "name": "pythoncheck ",
            "timeout": "2000",
            "type": "Kubernetes",
            "useAdvanced": false
          },
          "pod": {
            "spec": {
              "containers": [
                {
                  "resources": {
                    "limits": {
                      "cpu": 0.2,
                      "memory": "256Mi"
                    },
                    "requests": {
                      "cpu": 0.1,
                      "memory": "128Mi"
                    }
                  }
                }
              ]
            }
          },
          "pod_advanced": {},
          "service": {
            "spec": {
              "ports": [
                {
                  "nodePort": 32001,
                  "port": 8000,
                  "targetPort": 8000
                }
              ],
              "type": "NodePort"
            }
          },
          "service_advanced": {}
        },
        "translateX": 253.5,
        "translateY": 201,
        "type": "workflow-action",
        "width": 38
      },
      "endPoint": {
        "x": 253.5,
        "y": 201
      },
      "id": "start-stage-workflow-action-ac85dff0-bb63-11e6-be23-4d9fe51a1df7",
      "relation": {
        "pull_customize": [
          {
            "from": ".gitUrl",
            "to": ".gitUrl"
          }
        ]
      },
      "startData": {
        "height": 52,
        "id": "start-stage",
        "outputJson": [
          {
            "event": "pull",
            "json": {
              "gitUrl": ""
            },
            "type": "customize"
          }
        ],
        "setupData": {},
        "translateX": 50,
        "translateY": 84,
        "type": "workflow-start",
        "width": 45
      },
      "startPoint": {
        "x": 50,
        "y": 84
      },
      "workflowLineViewId": "workflow-line-view"
    }
  ],
  "stageList": [
    {
      "height": 52,
      "id": "start-stage",
      "outputJson": [
        {
          "event": "pull",
          "json": {
            "gitUrl": ""
          },
          "type": "customize"
        }
      ],
      "setupData": {},
      "translateX": 50,
      "translateY": 84,
      "type": "workflow-start",
      "width": 45
    },
    {
      "actions": [
        {
          "component": {
            "name": "pythoncheck",
            "versionid": 2,
            "versionname": "v1.0.1"
          },
          "env": [
            {
              "key": "CO_DATA",
              "value": "{\"test\":true}"
            }
          ],
          "height": 38,
          "id": "workflow-action-ac85dff0-bb63-11e6-be23-4d9fe51a1df7",
          "inputJson": {
            "gitUrl": "www"
          },
          "outputJson": {
            "status": true
          },
          "setupData": {
            "action": {
              "apiserver": "192.168.10.131:8080",
              "datafrom": ":3000/",
              "image": {
                "name": "pythoncheck",
                "tag": "1.0.1"
              },
              "ip": "192.168.10.131",
              "name": "pythoncheck ",
              "timeout": "2000",
              "type": "Kubernetes",
              "useAdvanced": false
            },
            "pod": {
              "spec": {
                "containers": [
                  {
                    "resources": {
                      "limits": {
                        "cpu": 0.2,
                        "memory": "256Mi"
                      },
                      "requests": {
                        "cpu": 0.1,
                        "memory": "128Mi"
                      }
                    }
                  }
                ]
              }
            },
            "pod_advanced": {},
            "service": {
              "spec": {
                "ports": [
                  {
                    "nodePort": 32001,
                    "port": 8000,
                    "targetPort": 8000
                  }
                ],
                "type": "NodePort"
              }
            },
            "service_advanced": {}
          },
          "translateX": 253.5,
          "translateY": 201,
          "type": "workflow-action",
          "width": 38
        }
      ],
      "class": "workflow-stage",
      "drawX": 0,
      "drawY": 0,
      "height": 52,
      "id": "workflow-stage-a400c890-bb63-11e6-be23-4d9fe51a1df7",
      "setupData": {
        "name": "pythoncheck",
        "timeout": "3000"
      },
      "translateX": 250,
      "translateY": 84,
      "type": "workflow-stage",
      "width": 45
    },
    {
      "height": 52,
      "id": "add-stage",
      "translateX": 450,
      "translateY": 84,
      "type": "workflow-add-stage",
      "width": 45
    },
    {
      "height": 52,
      "id": "end-stage",
      "setupData": {},
      "translateX": 650,
      "translateY": 84,
      "type": "workflow-end",
      "width": 45
    }
  ],
  "status": false
}
```

###  saveWorkflowInfo

Save workflow as the new version use same API.

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  | /v2/:namespace/:repository/workflow/v1/define/:workflowName|

#### body

```
{
  "id": 2,
  "version": "v1.0.1",
  "define": {
    "lineList": [
      {
        "endData": {
          "component": {
            "name": "pythoncheck",
            "versionid": 2,
            "versionname": "v1.0.1"
          },
          "env": [
            {
              "key": "CO_DATA",
              "value": "{\"test\":true}"
            }
          ],
          "height": 38,
          "id": "workflow-action-ac85dff0-bb63-11e6-be23-4d9fe51a1df7",
          "inputJson": {
            "gitUrl": "www"
          },
          "outputJson": {
            "status": true
          },
          "setupData": {
            "action": {
              "apiserver": "192.168.10.131:8080",
              "datafrom": ":3000/",
              "image": {
                "name": "pythoncheck",
                "tag": "1.0.1"
              },
              "ip": "192.168.10.131",
              "name": "pythoncheck ",
              "timeout": "2000",
              "type": "Kubernetes",
              "useAdvanced": false
            },
            "pod": {
              "spec": {
                "containers": [
                  {
                    "resources": {
                      "limits": {
                        "cpu": 0.2,
                        "memory": "256Mi"
                      },
                      "requests": {
                        "cpu": 0.1,
                        "memory": "128Mi"
                      }
                    }
                  }
                ]
              }
            },
            "pod_advanced": {},
            "service": {
              "spec": {
                "ports": [
                  {
                    "nodePort": 32001,
                    "port": 8000,
                    "targetPort": 8000
                  }
                ],
                "type": "NodePort"
              }
            },
            "service_advanced": {}
          },
          "translateX": 253.5,
          "translateY": 201,
          "type": "workflow-action",
          "width": 38
        },
        "endPoint": {
          "x": 253.5,
          "y": 201
        },
        "id": "start-stage-workflow-action-ac85dff0-bb63-11e6-be23-4d9fe51a1df7",
        "relation": {
          "pull_customize": [
            {
              "from": ".gitUrl",
              "to": ".gitUrl"
            }
          ]
        },
        "startData": {
          "height": 52,
          "id": "start-stage",
          "outputJson": [
            {
              "event": "pull",
              "json": {
                "gitUrl": ""
              },
              "type": "customize"
            }
          ],
          "setupData": {},
          "translateX": 50,
          "translateY": 84,
          "type": "workflow-start",
          "width": 45
        },
        "startPoint": {
          "x": 50,
          "y": 84
        },
        "workflowLineViewId": "workflow-line-view"
      }
    ],
    "stageList": [
      {
        "height": 52,
        "id": "start-stage",
        "outputJson": [
          {
            "event": "pull",
            "json": {
              "gitUrl": ""
            },
            "type": "customize"
          }
        ],
        "setupData": {},
        "translateX": 50,
        "translateY": 100.4,
        "type": "workflow-start",
        "width": 45
      },
      {
        "actions": [
          {
            "component": {
              "name": "pythoncheck",
              "versionid": 2,
              "versionname": "v1.0.1"
            },
            "env": [
              {
                "key": "CO_DATA",
                "value": "{\"test\":true}"
              }
            ],
            "height": 38,
            "id": "workflow-action-ac85dff0-bb63-11e6-be23-4d9fe51a1df7",
            "inputJson": {
              "gitUrl": "www"
            },
            "outputJson": {
              "status": true
            },
            "setupData": {
              "action": {
                "apiserver": "192.168.10.131:8080",
                "datafrom": ":3000/",
                "image": {
                  "name": "pythoncheck",
                  "tag": "1.0.1"
                },
                "ip": "192.168.10.131",
                "name": "pythoncheck ",
                "timeout": "2000",
                "type": "Kubernetes",
                "useAdvanced": false
              },
              "pod": {
                "spec": {
                  "containers": [
                    {
                      "resources": {
                        "limits": {
                          "cpu": 0.2,
                          "memory": "256Mi"
                        },
                        "requests": {
                          "cpu": 0.1,
                          "memory": "128Mi"
                        }
                      }
                    }
                  ]
                }
              },
              "pod_advanced": {},
              "service": {
                "spec": {
                  "ports": [
                    {
                      "nodePort": 32001,
                      "port": 8000,
                      "targetPort": 8000
                    }
                  ],
                  "type": "NodePort"
                }
              },
              "service_advanced": {}
            },
            "translateX": 253.5,
            "translateY": 217.4,
            "type": "workflow-action",
            "width": 38
          }
        ],
        "class": "workflow-stage",
        "drawX": 0,
        "drawY": 0,
        "height": 52,
        "id": "workflow-stage-a400c890-bb63-11e6-be23-4d9fe51a1df7",
        "setupData": {
          "name": "pythoncheck",
          "timeout": "3000"
        },
        "translateX": 250,
        "translateY": 100.4,
        "type": "workflow-stage",
        "width": 45
      },
      {
        "height": 52,
        "id": "add-stage",
        "translateX": 450,
        "translateY": 100.4,
        "type": "workflow-add-stage",
        "width": 45
      },
      {
        "height": 52,
        "id": "end-stage",
        "setupData": {},
        "translateX": 650,
        "translateY": 100.4,
        "type": "workflow-end",
        "width": 45
      }
    ]
  }
}
```


#### response json

```
{
  "message": "success"
}
```

###  deleteWorkflowInfo


| HTTP Method |  Request Address |
| -------- | ------ |
| DELETE  | /v2/:namespace/:repository/workflow/v1/define/:workflowName|

#### response json

```
{
  "message": "success"
}
```

### set workflow env

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  |/v2/:namespace/:repository/workflow/v1/define/:workflowName/env|

#### body

```
{
  "id": 2,
  "env": {
    "TEST": "123456"
  }
}
```

#### response json

```
{
  "message": "success"
}
```

### GET workflow env

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  |/v2/:namespace/:repository/workflow/v1/define/:workflowName/env?id=:workflowID|

#### response json

```
{
  "env": {
    "TEST": "123456"
  }
}
```

### git event json

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /v2/:namespace/:repository/workflow/v1/define/event/github/:eventName|

#### response json

```
{
  "output": {
    "deployment": {
      "created_at": "2015-05-05T23:40:38Z",
      "creator": {
        "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
        "events_url": "https://api.github.com/users/baxterthehacker/events{/privacy}",
        "followers_url": "https://api.github.com/users/baxterthehacker/followers",
        "following_url": "https://api.github.com/users/baxterthehacker/following{/other_user}",
        "gists_url": "https://api.github.com/users/baxterthehacker/gists{/gist_id}",
        "gravatar_id": "",
        "html_url": "https://github.com/baxterthehacker",
        "id": 6752317,
        "login": "baxterthehacker",
        "organizations_url": "https://api.github.com/users/baxterthehacker/orgs",
        "received_events_url": "https://api.github.com/users/baxterthehacker/received_events",
        "repos_url": "https://api.github.com/users/baxterthehacker/repos",
        "site_admin": false,
        "starred_url": "https://api.github.com/users/baxterthehacker/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/baxterthehacker/subscriptions",
        "type": "User",
        "url": "https://api.github.com/users/baxterthehacker"
      },
      "description": null,
      "environment": "production",
      "id": 710692,
      "payload": {},
      "ref": "master",
      "repository_url": "https://api.github.com/repos/baxterthehacker/public-repo",
      "sha": "9049f1265b7d61be4a8904a9a27120d2064dab3b",
      "statuses_url": "https://api.github.com/repos/baxterthehacker/public-repo/deployments/710692/statuses",
      "task": "deploy",
      "updated_at": "2015-05-05T23:40:38Z",
      "url": "https://api.github.com/repos/baxterthehacker/public-repo/deployments/710692"
    },
    "repository": {
      "archive_url": "https://api.github.com/repos/baxterthehacker/public-repo/{archive_format}{/ref}",
      "assignees_url": "https://api.github.com/repos/baxterthehacker/public-repo/assignees{/user}",
      "blobs_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/blobs{/sha}",
      "branches_url": "https://api.github.com/repos/baxterthehacker/public-repo/branches{/branch}",
      "clone_url": "https://github.com/baxterthehacker/public-repo.git",
      "collaborators_url": "https://api.github.com/repos/baxterthehacker/public-repo/collaborators{/collaborator}",
      "comments_url": "https://api.github.com/repos/baxterthehacker/public-repo/comments{/number}",
      "commits_url": "https://api.github.com/repos/baxterthehacker/public-repo/commits{/sha}",
      "compare_url": "https://api.github.com/repos/baxterthehacker/public-repo/compare/{base}...{head}",
      "contents_url": "https://api.github.com/repos/baxterthehacker/public-repo/contents/{+path}",
      "contributors_url": "https://api.github.com/repos/baxterthehacker/public-repo/contributors",
      "created_at": "2015-05-05T23:40:12Z",
      "default_branch": "master",
      "description": "",
      "downloads_url": "https://api.github.com/repos/baxterthehacker/public-repo/downloads",
      "events_url": "https://api.github.com/repos/baxterthehacker/public-repo/events",
      "fork": false,
      "forks": 0,
      "forks_count": 0,
      "forks_url": "https://api.github.com/repos/baxterthehacker/public-repo/forks",
      "full_name": "baxterthehacker/public-repo",
      "git_commits_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/commits{/sha}",
      "git_refs_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/refs{/sha}",
      "git_tags_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/tags{/sha}",
      "git_url": "git://github.com/baxterthehacker/public-repo.git",
      "has_downloads": true,
      "has_issues": true,
      "has_pages": true,
      "has_wiki": true,
      "homepage": null,
      "hooks_url": "https://api.github.com/repos/baxterthehacker/public-repo/hooks",
      "html_url": "https://github.com/baxterthehacker/public-repo",
      "id": 35129377,
      "issue_comment_url": "https://api.github.com/repos/baxterthehacker/public-repo/issues/comments{/number}",
      "issue_events_url": "https://api.github.com/repos/baxterthehacker/public-repo/issues/events{/number}",
      "issues_url": "https://api.github.com/repos/baxterthehacker/public-repo/issues{/number}",
      "keys_url": "https://api.github.com/repos/baxterthehacker/public-repo/keys{/key_id}",
      "labels_url": "https://api.github.com/repos/baxterthehacker/public-repo/labels{/name}",
      "language": null,
      "languages_url": "https://api.github.com/repos/baxterthehacker/public-repo/languages",
      "merges_url": "https://api.github.com/repos/baxterthehacker/public-repo/merges",
      "milestones_url": "https://api.github.com/repos/baxterthehacker/public-repo/milestones{/number}",
      "mirror_url": null,
      "name": "public-repo",
      "notifications_url": "https://api.github.com/repos/baxterthehacker/public-repo/notifications{?since,all,participating}",
      "open_issues": 2,
      "open_issues_count": 2,
      "owner": {
        "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
        "events_url": "https://api.github.com/users/baxterthehacker/events{/privacy}",
        "followers_url": "https://api.github.com/users/baxterthehacker/followers",
        "following_url": "https://api.github.com/users/baxterthehacker/following{/other_user}",
        "gists_url": "https://api.github.com/users/baxterthehacker/gists{/gist_id}",
        "gravatar_id": "",
        "html_url": "https://github.com/baxterthehacker",
        "id": 6752317,
        "login": "baxterthehacker",
        "organizations_url": "https://api.github.com/users/baxterthehacker/orgs",
        "received_events_url": "https://api.github.com/users/baxterthehacker/received_events",
        "repos_url": "https://api.github.com/users/baxterthehacker/repos",
        "site_admin": false,
        "starred_url": "https://api.github.com/users/baxterthehacker/starred{/owner}{/repo}",
        "subscriptions_url": "https://api.github.com/users/baxterthehacker/subscriptions",
        "type": "User",
        "url": "https://api.github.com/users/baxterthehacker"
      },
      "private": false,
      "pulls_url": "https://api.github.com/repos/baxterthehacker/public-repo/pulls{/number}",
      "pushed_at": "2015-05-05T23:40:38Z",
      "releases_url": "https://api.github.com/repos/baxterthehacker/public-repo/releases{/id}",
      "size": 0,
      "ssh_url": "git@github.com:baxterthehacker/public-repo.git",
      "stargazers_count": 0,
      "stargazers_url": "https://api.github.com/repos/baxterthehacker/public-repo/stargazers",
      "statuses_url": "https://api.github.com/repos/baxterthehacker/public-repo/statuses/{sha}",
      "subscribers_url": "https://api.github.com/repos/baxterthehacker/public-repo/subscribers",
      "subscription_url": "https://api.github.com/repos/baxterthehacker/public-repo/subscription",
      "svn_url": "https://github.com/baxterthehacker/public-repo",
      "tags_url": "https://api.github.com/repos/baxterthehacker/public-repo/tags",
      "teams_url": "https://api.github.com/repos/baxterthehacker/public-repo/teams",
      "trees_url": "https://api.github.com/repos/baxterthehacker/public-repo/git/trees{/sha}",
      "updated_at": "2015-05-05T23:40:30Z",
      "url": "https://api.github.com/repos/baxterthehacker/public-repo",
      "watchers": 0,
      "watchers_count": 0
    },
    "sender": {
      "avatar_url": "https://avatars.githubusercontent.com/u/6752317?v=3",
      "events_url": "https://api.github.com/users/baxterthehacker/events{/privacy}",
      "followers_url": "https://api.github.com/users/baxterthehacker/followers",
      "following_url": "https://api.github.com/users/baxterthehacker/following{/other_user}",
      "gists_url": "https://api.github.com/users/baxterthehacker/gists{/gist_id}",
      "gravatar_id": "",
      "html_url": "https://github.com/baxterthehacker",
      "id": 6752317,
      "login": "baxterthehacker",
      "organizations_url": "https://api.github.com/users/baxterthehacker/orgs",
      "received_events_url": "https://api.github.com/users/baxterthehacker/received_events",
      "repos_url": "https://api.github.com/users/baxterthehacker/repos",
      "site_admin": false,
      "starred_url": "https://api.github.com/users/baxterthehacker/starred{/owner}{/repo}",
      "subscriptions_url": "https://api.github.com/users/baxterthehacker/subscriptions",
      "type": "User",
      "url": "https://api.github.com/users/baxterthehacker"
    }
  }
}
```

#### eventName
|Create | Delete | Deployment | DeploymentStatus | Fork | Gollum | IssueComment | Issues | Member | PageBuild | Public | PullRequestReviewComment | PullRequestReview | PullRequest | Push | Repository|Release|Status|TeamAdd|Watch|

### change workflow state

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  |/v2/:namespace/:repository/workflow/v1/define/:workflowName/state|

#### body

```
{
  "id": 2,
  "state": 1   #1 ON, 0 OFF
}
```

#### response json

```
{
  "message": "success"
}
```

### run workflow

| HTTP Method |  Request Address |
| -------- | ------ |
| POST  | /v2/:namespace/:repository/workflow/v1/exec/:workflowName|

#### body
```
{
  "gitUrl": "https://github.com/xiechuanj/python-sonar-runner.git"
}
```

#### response json
```
{
  "result": "workflow start ..."
}
```



### get workflow token and url

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /v2/:namespace/:repository/workflow/v1/define/:workflow/token?id=:workflowID|

#### response json

```
{
  "token": "621c09e7bf910401b9a514c075def56d",
  "url": "http://192.168.10.131:10000/demo/demo/pythoncheck"
}
```

### getWorkflowHistories

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  |/v2/:namespace/:repository/workflow/v1/log/list|

#### response json

```
{
  "workflowList": [
    {
      "id": 1,
      "name": "busybox",
      "versionList": [
        {
          "id": 1,
          "info": "Success: 0 Total: 1",
          "name": "1.0.1",
          "sequenceList": [
            {
              "sequence": 1,
              "status": 1,
              "time": "2016-12-05 22:41:38",
              "workflowSequenceID": 1
            }
          ],
          "success": 0,
          "total": 1
        }
      ]
    }
  ]
}
```

### getWorkflowHistory

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /v2/:namespace/:repository/workflow/v1/log/:workflowName/:version?sequence=:workflowSequenceID|

#### response json

```
{
  "define": {
    "lineList": [
      {
        "endData": {
          "id": "a-1",
          "setupData": {
            "action": {
              "name": "busybox"
            }
          }
        },
        "id": "s-1-a-1",
        "startData": {
          "id": "s-1",
          "type": "workflow-start"
        },
        "workflowLineViewId": "workflow-line-view"
      }
    ],
    "sequence": 1,
    "stageList": [
      {
        "id": "s-1",
        "name": "busybox-start-stage",
        "runTime": "2016-12-05 22:41:38 - 2016-12-05 22:41:40",
        "setupData": {
          "name": "busybox-start-stage"
        },
        "status": 3,
        "type": "workflow-start"
      },
      {
        "actions": [
          {
            "id": "a-1",
            "setupData": {
              "name": "busybox"
            },
            "status": 4,
            "type": "workflow-action"
          }
        ],
        "id": "s-2",
        "name": "ffsdafd",
        "runTime": "2016-12-05 22:41:38 - ",
        "setupData": {
          "name": "ffsdafd"
        },
        "status": 1,
        "type": "workflow-stage"
      },
      {
        "id": "s-3",
        "name": "busybox-end-stage",
        "runTime": "2016-12-05 22:41:38 - ",
        "setupData": {
          "name": "busybox-end-stage"
        },
        "status": 0,
        "type": "workflow-end"
      }
    ],
    "status": 1,
    "version": "1.0.1",
    "workflow": "busybox"
  }
}
```

### getStageRunHistory

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /v2/:namespace/:repository/workflow/v1/log/:workflowName/:version/:sequenceID/stage/:stageName|


### getActionRunHistory

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /v2/:namespace/:repository/workflow/v1/log/:workflowName/:version/:sequenceID/stage/:stageName/action/:actionName|

#### response json

```
{
  "result": {
    "data": {
      "input": {},
      "output": {}
    },
    "logList": []
  }
}
```

### getLineDataInfo

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  |/v2/:namespace/:repository/workflow/v1/log/:workflowName/:version/:sequenceID/:relation|

#### response json

```
{
  "define": {
    "input": {
      "gitUrl": "https://github.com/xiechuanj/python-sonar-runner.git"
    },
    "output": {}
  }
}
```

### component event

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  | /workflow/v1/:namespace/:repository/:workflowName/event|

#### body

```
{
  "EVENT": "TASK_RESULT",
  "EVENTID": 5260,
  "INFO": {
    "output": {
      "binaryFileUrl": "aaa",
      "data": "task start"
    },
    "result": "",
    "status": true
  },
  "RUN_ID": "288,1197,877,36,84"
}
```

#### response json

```
{
  "message": "ok"
}
```

#### EVENT
|COMPONENT_START |TASK_START |TASK_STATUS|TASK_RESULT|COMPONENT_STOP|

### component register

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  | /workflow/v1/:namespace/:repository/:workflowName/register|

#### body

```
{
  "RUN_ID": "4400,12678,10113,1,137",
  "POD_NAME": "pod-4400-12678-10113-1-137",
  "RECEIVE_URL": "32001"
}
```

#### response json

```
{
  "message": "ok"
}
```
