---
title: V1 Specification
keywords: component
tags: [API]
sidebar: home_sidebar
permalink: v1-specification.html
summary: V1 Specification
---

## API V1 Operations

### Create Workflow

| HTTP Method |  Request Address |
| -------- | -------- |
| POST  | /pipeline/v1/:namespace/:repository |

#### Body:

```
{
  "name": "pythoncheck",
  "version": "1.0"
}
```

#### response json

```
{
  "message": "create new pipeline success"
}
```

### getPipelineList

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/:namespace/:repository |

#### response json

```
{
  "list": [
    {
      "id": 68,  //pipelineID
      "name": "pythoncheck", //pipelineName
      "version": [
        {
          "id": 68,
          "version": "1.0",
          "versionCode": 1
        }
      ]
    }
  ]
}
```


### getPipelineInfo

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/:namespace/:repository/:pipelineName/json?id=:pipelineID |

#### response json

```
{
  "lineList": [
    {
      "endData": {
        "component": {
          "name": "pythoncheck",
          "versionid": 137
        },
        "env": [
          {
            "key": "gitUrl",
            "value": "123456"
          }
        ],
        "height": 38,
        "id": "pipeline-action-224ec5e0-a3fd-11e6-8e43-dbedb3b31745",
        "inputJson": {
          "gitUrl": "https://github.com/Huawei/containerops.git"
        },
        "outputJson": {
          "status": true
        },
        "setupData": {
          "action": {
            "apiserver": "http://192.168.10.131:8080",
            "datafrom": "{}",
            "image": {
              "name": "xiechuan/pythoncheck",
              "tag": "1.0"
            },
            "ip": "192.168.10.131",
            "name": "pythoncheck",
            "timeout": "30000",
            "type": "Kubernetes",
            "useAdvanced": false
          },
          "pod": {
            "spec": {
              "containers": [
                {
                  "resources": {
                    "limits": {
                      "cpu": "0.1",
                      "memory": "128Mi"
                    },
                    "requests": {
                      "cpu": "0.1",
                      "memory": "64Mi"
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
        "translateY": 224,
        "type": "pipeline-action",
        "width": 38
      },
      "endPoint": {
        "x": 253.5,
        "y": 224
      },
      "id": "start-stage-pipeline-action-224ec5e0-a3fd-11e6-8e43-dbedb3b31745",
      "pipelineLineViewId": "pipeline-line-view",
      "relation": [
        {
          "finalPath": "start-stage.gitUrl",
          "from": ".gitUrl",
          "to": ".gitUrl"
        }
      ],
      "startData": {
        "height": 52,
        "id": "start-stage",
        "outputJson": {
          "gitUrl": ""
        },
        "setupData": {
          "event": "PullRequest",
          "type": "customize"
        },
        "translateX": 50,
        "translateY": 107,
        "type": "pipeline-start",
        "width": 45
      },
      "startPoint": {
        "x": 50,
        "y": 107
      }
    }
  ],
  "stageList": [
    {
      "height": 52,
      "id": "start-stage",
      "outputJson": {
        "gitUrl": ""
      },
      "setupData": {
        "event": "PullRequest",
        "type": "customize"
      },
      "translateX": 50,
      "translateY": 107,
      "type": "pipeline-start",
      "width": 45
    },
    {
      "actions": [
        {
          "component": {
            "name": "pythoncheck",
            "versionid": 137
          },
          "env": [
            {
              "key": "TEST",
              "value": "123456"
            }
          ],
          "height": 38,
          "id": "pipeline-action-224ec5e0-a3fd-11e6-8e43-dbedb3b31745",
          "inputJson": {
            "gitUrl": "https://github.com/Huawei/containerops.git"
          },
          "outputJson": {
            "status": true
          },
          "setupData": {
            "action": {
              "apiserver": "http://192.168.10.131:8080",
              "datafrom": "{}",
              "image": {
                "name": "xiechuan/pythoncheck",
                "tag": "1.0"
              },
              "ip": "192.168.10.131",
              "name": "pythoncheck",
              "timeout": "30000",
              "type": "Kubernetes",
              "useAdvanced": false
            },
            "pod": {
              "spec": {
                "containers": [
                  {
                    "resources": {
                      "limits": {
                        "cpu": "0.1",
                        "memory": "128Mi"
                      },
                      "requests": {
                        "cpu": "0.1",
                        "memory": "64Mi"
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
          "translateY": 224,
          "type": "pipeline-action",
          "width": 38
        }
      ],
      "class": "pipeline-stage",
      "drawX": 0,
      "drawY": 0,
      "height": 52,
      "id": "pipeline-stage-1d4ce1d0-a3fd-11e6-8e43-dbedb3b31745",
      "setupData": {
        "name": "pythoncheck",
        "timeout": "3000"
      },
      "translateX": 250,
      "translateY": 107,
      "type": "pipeline-stage",
      "width": 45
    },
    {
      "height": 52,
      "id": "add-stage",
      "translateX": 450,
      "translateY": 107,
      "type": "pipeline-add-stage",
      "width": 45
    },
    {
      "height": 52,
      "id": "end-stage",
      "setupData": {},
      "translateX": 650,
      "translateY": 107,
      "type": "pipeline-end",
      "width": 45
    }
  ],
  "status": false
}
```

###  savePipelineInfo

Save workflow as the new version use same API.

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  | /pipeline/v1/:namespace/:repository/:pipelineName|

#### body

```
{
  "id": 68,
  "version": "1.0",
  "define": {
    "lineList": [
      {
        "pipelineLineViewId": "pipeline-line-view",
        "startData": {
          "id": "start-stage",
          "setupData": {
            "type": "customize",
            "event": "PullRequest"
          },
          "type": "pipeline-start",
          "width": 45,
          "height": 52,
          "translateX": 50,
          "translateY": 107,
          "outputJson": {
            "gitUrl": ""
          }
        },
        "endData": {
          "id": "pipeline-action-224ec5e0-a3fd-11e6-8e43-dbedb3b31745",
          "type": "pipeline-action",
          "setupData": {
            "action": {
              "apiserver": "http://192.168.10.131:8080",
              "datafrom": "{}",
              "image": {
                "name": "xiechuan/pythoncheck",
                "tag": "1.0"
              },
              "ip": "192.168.10.131",
              "name": "pythoncheck",
              "timeout": "30000",
              "type": "Kubernetes",
              "useAdvanced": false
            },
            "pod": {
              "spec": {
                "containers": [
                  {
                    "resources": {
                      "limits": {
                        "cpu": "0.1",
                        "memory": "128Mi"
                      },
                      "requests": {
                        "cpu": "0.1",
                        "memory": "64Mi"
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
          "translateY": 224,
          "width": 38,
          "height": 38,
          "translateX": 253.5,
          "inputJson": {
            "gitUrl": "https://github.com/Huawei/containerops.git"
          },
          "outputJson": {
            "status": true
          },
          "env": [
            {
              "key": "TEST",
              "value": "123456"
            }
          ],
          "component": {
            "name": "pythoncheck",
            "versionid": 137
          }
        },
        "startPoint": {
          "x": 50,
          "y": 107
        },
        "endPoint": {
          "x": 253.5,
          "y": 224
        },
        "id": "start-stage-pipeline-action-224ec5e0-a3fd-11e6-8e43-dbedb3b31745",
        "relation": [
          {
            "to": ".gitUrl",
            "from": ".gitUrl",
            "finalPath": "start-stage.gitUrl"
          }
        ]
      }
    ],
    "stageList": [
      {
        "id": "start-stage",
        "setupData": {
          "type": "customize",
          "event": "PullRequest"
        },
        "type": "pipeline-start",
        "width": 45,
        "height": 52,
        "translateX": 50,
        "translateY": 107,
        "outputJson": {
          "gitUrl": ""
        }
      },
      {
        "id": "pipeline-stage-1d4ce1d0-a3fd-11e6-8e43-dbedb3b31745",
        "type": "pipeline-stage",
        "class": "pipeline-stage",
        "drawX": 0,
        "drawY": 0,
        "width": 45,
        "height": 52,
        "translateX": 250,
        "translateY": 107,
        "actions": [
          {
            "id": "pipeline-action-224ec5e0-a3fd-11e6-8e43-dbedb3b31745",
            "type": "pipeline-action",
            "setupData": {
              "action": {
                "apiserver": "http://192.168.10.131:8080",
                "datafrom": "{}",
                "image": {
                  "name": "xiechuan/pythoncheck",
                  "tag": "1.0"
                },
                "ip": "192.168.10.131",
                "name": "pythoncheck",
                "timeout": "30000",
                "type": "Kubernetes",
                "useAdvanced": false
              },
              "pod": {
                "spec": {
                  "containers": [
                    {
                      "resources": {
                        "limits": {
                          "cpu": "0.1",
                          "memory": "128Mi"
                        },
                        "requests": {
                          "cpu": "0.1",
                          "memory": "64Mi"
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
            "translateY": 224,
            "width": 38,
            "height": 38,
            "translateX": 253.5,
            "inputJson": {
              "gitUrl": "https://github.com/Huawei/containerops.git"
            },
            "outputJson": {
              "status": true
            },
            "env": [
              {
                "key": "TEST",
                "value": "123456"
              }
            ],
            "component": {
              "name": "pythoncheck",
              "versionid": 137
            }
          }
        ],
        "setupData": {
          "name": "pythoncheck",
          "timeout": "3000"
        }
      },
      {
        "id": "add-stage",
        "type": "pipeline-add-stage",
        "width": 45,
        "height": 52,
        "translateX": 450,
        "translateY": 107
      },
      {
        "id": "end-stage",
        "setupData": {},
        "type": "pipeline-end",
        "width": 45,
        "height": 52,
        "translateX": 650,
        "translateY": 107
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

### set pipeline env

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  |/pipeline/v1/:namespace/:repository/:pipelineName/env|

#### body

```
{
  "id": 68,
  "env": {
    "GITURL": "https://github.com/Huawei/containerops.git"
  }
}
```

#### response json

```
{
  "message": "success"
}
```

### GET pipeline env

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  |/pipeline/v1/:namespace/:repository/:pipelineName/env?id=:pipelineID|

#### response json

```
{
  "env": {
    "GITURL": "https://github.com/Huawei/containerops.git"
  }
}
```

### git event json

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/eventJson/github/:eventName|

#### response json

```
{
  "output": {
    "description": "",
    "master_branch": "master",
    "pusher_type": "user",
    "ref": "0.0.1",
    "ref_type": "tag",
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

### change pipeline state

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  |/pipeline/v1/:namespace/:repository/:pipelineName/state|

#### body

```
{
  "id": 68,
  "state": 1 #0 OFF, 1 ON
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
| POST  | /pipeline/v1/:namespace/:repository/:pipelineName/exec?version=:version|

#### body
```
{
  "gitUrl": "https://github.com/xiechuanj/python-sonar-runner.git"
}
```

#### response json
```
{
  "result": "pipeline start ..."
}
```

### createComponent

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  |/pipeline/v1/demo/component|

#### body

```
{
  "name": "pythoncheck",
  "version": "1.0"
}
```

#### response json

```
{
  "message": "create new component success"
}
```

### get component list

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  |/pipeline/v1/:namespace/component|

#### response json

```
{
  "list": [
    {
      "id": 248,   //componentID
      "name": "pythoncheck",   //componentName
      "version": [
        {
          "id": 248,
          "version": "1.0",
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
| GET  |/pipeline/v1/:namespace/component/:componentName?id=:componentID|

#### response json

```
{
  "env": [
    {
      "key": "CO_DATA",
      "value": "{'contents':'sonar.projectKey:python-sonar-runner sonar.projectName=python-sonar-runner sonar.projectVersion=1.0 sonar.sources=src sonar.language=py sonar.sourceEncoding=UTF-8','filename':'sonar-project.properties'}"
    }
  ],
  "inputJson": {
    "gitUrl": "https://github.com/xiechuanj/python-sonar-runner.git"
  },
  "outputJson": {
    "status": true
  },
  "setupData": {
    "action": {
      "apiserver": "",
      "datafrom": "{}",
      "image": {
        "name": "xiechuan/pythoncheck",
        "tag": "1.0"
      },
      "ip": "",
      "name": "pythoncheck",
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
                "cpu": "0.1",
                "memory": "128Mi"
              },
              "requests": {
                "cpu": "0.1",
                "memory": "64Mi"
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
| PUT  |/pipeline/v1/:namespace/component/:componentName|

#### body

```
{
  "id": 248,
  "version": "1.0",
  "define": {
    "inputJson": {
      "gitUrl": "https://github.com/xiechuanj/python-sonar-runner.git"
    },
    "outputJson": {
      "status": true,
      "result":""
    },
    "setupData": {
      "action": {
        "type": "Kubernetes",
        "name": "pythoncheck",
        "timeout": "3000",
        "ip": "",
        "apiserver": "",
        "image": {
          "name": "xiechuan/pythoncheck",
          "tag": "1.0"
        },
        "useAdvanced": false,
        "datafrom": "{}"
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
                  "cpu": "0.1",
                  "memory": "128Mi"
                },
                "requests": {
                  "cpu": "0.1",
                  "memory": "64Mi"
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
        "value": "{'contents':'sonar.projectKey:python-sonar-runner sonar.projectName=python-sonar-runner sonar.projectVersion=1.0 sonar.sources=src sonar.language=py sonar.sourceEncoding=UTF-8','filename':'sonar-project.properties'}"
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



### get pipelien token and url

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/:namespace/:repository/:pipelineName/token?id=:pipelineID|

#### response json

```
{
  "token": "621c09e7bf910401b9a514c075def56d",
  "url": "http://192.168.10.131:10000/demo/demo/pythoncheck"
}
```

### getPipelineHistories

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/:namespace/:repository/histories|

#### response json

```
{
  "pipelineList": [
    {
      "id": 68,
      "name": "pythoncheck",
      "versionList": [
        {
          "id": 68,
          "info": "Success :0 Total :1",
          "name": "1.0",
          "sequenceList": [
            {
              "pipelineSequenceID": 26,   //pipelineSequenceID
              "sequence": 1,
              "status": false,
              "time": "2016-11-08 14:46:09"
            }
          ]
        }
      ]
    }
  ]
}
```

### getPipelineHistory

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/:namespace/:repository/:pipelineName/:version/define?sequenceId=:pipelineSequenceID|

#### response json

```
{
  "define": {
    "lineList": [
      {
        "endData": {
          "id": "a-196",   //actionLogID
          "setupData": {
            "action": {
              "name": "pythoncheck"
            }
          },
          "type": "pipeline-action"
        },
        "id": "s-400-a-196",
        "pipelineLineViewId": "pipeline-line-view",
        "startData": {
          "id": "s-400",
          "type": "pipeline-start"
        }
      }
    ],
    "stageList": [
      {
        "id": "s-400",  //stageLogID
        "setupData": {
          "name": "pythoncheck-start-stage"
        },
        "status": true,
        "type": "pipeline-start"
      },
      {
        "actions": [
          {
            "id": "a-196",
            "setupData": {
              "name": "pythoncheck"
            },
            "status": false,
            "type": "pipeline-action"
          }
        ],
        "id": "s-401",
        "setupData": {
          "name": "pythoncheck"
        },
        "status": false,
        "type": "pipeline-stage"
      },
      {
        "id": "s-402",
        "setupData": {
          "name": "pythoncheck-end-stage"
        },
        "status": false,
        "type": "pipeline-end"
      }
    ],
    "status": false
  }
}
```

### getStageRunHistory

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/:namespace/:repository/:pipelineName/stage/:stageName/history?stageLogId=:stageLogID|


### getActionRunHistory

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/:namespace/:repository/:pipelineName/:version/:sequence/stage/:stageName/action/:actionName/define?actionLogId=:actionLogID|

#### response json

```
{
  "result": {
    "data": {
      "input": {
        "gitUrl": "https://github.com/xiechuanj/python-sonar-runner.git"
      },
      "output": {}
    },
    "logList": [
      "2016-11-08 19:35:09 -> {\"EVENTID\": 1207, \"EVENT\": \"COMPONENT_START\", \"RUN_ID\": \"30,419,202,1,137\"}",
      "2016-11-08 19:35:09 -> {\"data\":\"{\\\"gitUrl\\\":\\\"https://github.com/xiechuanj/python-sonar-runner.git\\\"}\",\"resp\":\"{\\\"gitUrl\\\":\\\"https://github.com/xiechuanj/python-sonar-runner.git\\\"}\\r\\n\"}",
      "2016-11-08 19:35:09 -> {\"EVENTID\": 1209, \"EVENT\": \"TASK_START\", \"RUN_ID\": \"30,419,202,1,137\"}",
      "2016-11-08 19:35:09 -> {\"EVENTID\": 1211, \"INFO\": {\"TASK_STATUS\": \"RUNNING\"}, \"EVENT\": \"TASK_STATUS\", \"RUN_ID\": \"30,419,202,1,137\"}",
      "2016-11-08 19:35:09 -> {\"EVENTID\": 1211, \"INFO\": {\"TASK_STATUS\": \"GET RESULT\"}, \"EVENT\": \"TASK_STATUS\", \"RUN_ID\": \"30,419,202,1,137\"}",
      "2016-11-08 19:35:09 -> {\"EVENTID\": 1210, \"INFO\": {\"status\": true, \"result\": {}}, \"EVENT\": \"TASK_RESULT\", \"RUN_ID\": \"30,419,202,1,137\"}",
      "2016-11-08 19:35:10 -> {\"EVENTID\": 1208, \"EVENT\": \"COMPONENT_STOP\", \"RUN_ID\": \"30,419,202,1,137\"}"
    ]
  }
}
```

### getLineDataInfo

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/:namespace/:repository/:pipelineName/:version/:lineid|

#### response json

```
{
  "define": {
    "input": {
      "gitUrl": "https://github.com/xiechuanj/python-sonar-runner.git"
    },
    "output": {
      "gitUrl": "https://github.com/xiechuanj/python-sonar-runner.git"
    }
  }
}
```

### component event

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  | /pipeline/v1/:namespace/:repository/:pipelineName/event|

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
| PUT  | /pipeline/v1/:namespace/:repository/:pipelineName/register|

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
