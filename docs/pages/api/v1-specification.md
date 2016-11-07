---
title: V1 Specification
keywords: component
tags: [component]
sidebar: home_sidebar
permalink: v1-specification.html
summary: V1 Specification
---

## API V1 Operations

### createPipeline

| HTTP Method |  Request Address |
| -------- | ------ |
| POST  | /pipeline/v1/:namespace/:repository |

#### Body:

```
{
  "name": "pythonSonarCheck",
  "version": "V1.0"
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
  "list" : [
    {
      "id" : 15,
      "name" : "pythonSonarCheck",
      "version" : [
        {
          "id" : 15,
          "version" : "V1.0",
          "versionCode" : 1
        }
      ]
    },
    {
      "id" : 1,
      "name" : "python",
      "version" : [
        {
          "id" : 1,
          "version" : "1.0",
          "versionCode" : 1
        }
      ]
    }
  ]
}
```

| Response Field | Field Type | Descripiton |
| -------- | ------ |--------|
| id          | integer | id of the pipeLine   |
| name        | string  | name of the pipeLine |
| id          | integer | -------- |
| version     | string  | -------- |
| versionCode | integer | -------- |

### getPipelineInfo

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/:namespace/:repository/:pipelineName/json?id=:pipelineID |

#### response json

```
{
  "lineList" : [ ],
  "stageList" : [
    {
      "id" : "start-stage",
      "setupData" : { },
      "type" : "pipeline-start"
    },
    {
      "id" : "add-stage",
      "type" : "pipeline-add-stage"
    },
    {
      "id" : "end-stage",
      "setupData" : { },
      "type" : "pipeline-end"
    }
  ],
  "status" : false
}
```
###  savePipelineInfo/savePipelineAsNewVersion

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  | /pipeline/v1/:namespace/:repository/:pipelineName|


#### body

```
{
  "id": 15,
  "version": "V1.0",
  "define": {
    "lineList": [],
    "stageList": [
      {
        "id": "start-stage",
        "setupData": {},
        "type": "pipeline-start",
        "width": 45,
        "height": 52,
        "translateX": 50,
        "translateY": 96.2
      },
      {
        "id": "pipeline-stage-8b3dc560-a2f4-11e6-baa2-615f3fa81dab",
        "type": "pipeline-stage",
        "class": "pipeline-stage",
        "drawX": 0,
        "drawY": 0,
        "width": 45,
        "height": 52,
        "translateX": 250,
        "translateY": 96.2,
        "actions": [
          {
            "id": "pipeline-action-8d0cb900-a2f4-11e6-baa2-615f3fa81dab",
            "type": "pipeline-action",
            "setupData": {},
            "translateY": 213.2,
            "width": 38,
            "height": 38,
            "translateX": 253.5
          }
        ],
        "setupData": {
          "name": "",
          "timeout": ""
        }
      },
      {
        "id": "add-stage",
        "type": "pipeline-add-stage",
        "width": 45,
        "height": 52,
        "translateX": 450,
        "translateY": 96.2
      },
      {
        "id": "end-stage",
        "setupData": {},
        "type": "pipeline-end",
        "width": 45,
        "height": 52,
        "translateX": 650,
        "translateY": 96.2
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
  "id": 15,
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
| GET  |/pipeline/v1/:namespace/:repository/:pipelineName/env?id=:pipelineId|

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
| PUT  | /pipeline/v1/eventJson/github/:eventName|

### change pipeline state

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  |/pipeline/v1/:namespace/:repository/:pipelineName/state|

#### body

```
{
  "id": 15,
  "state": 1 #0 OFF, 1 ON
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
| get  |/pipeline/v1/:namespace/component|

#### response json

```
{
  "list" : [
    {
      "id" : 2,
      "name" : "pythonrun",
      "version" : [
        {
          "id" : 2,
          "version" : "v1.0",
          "versionCode" : 1
        }
      ]
    },
    {
      "id" : 1,
      "name" : "pythoncheck",
      "version" : [
        {
          "id" : 1,
          "version" : "v1.0",
          "versionCode" : 1
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
  "env" : [
    {
      "key" : "CO_DATA",
      "value" : "{     \"serverities\": \"MAJOR\",     \"contents\": \"sonar.projectKey=python-sonar-runner\\nsonar.projectName=python-sonar-runner\\nsonar.projectVersion=1.0\\nsonar.sources=src\\nsonar.language=py\\nsonar.sourceEncoding=UTF-8\",     \"filename\": \"sonar-project.properties\" }"
    }
  ],
  "inputJson" : {
    "gitUrl" : "https://github.com/xiechuanj/python-sonar-runner.git"
  },
  "outputJson" : {
    "EVENT" : "TASK_RESULT",
    "EVENTID" : "",
    "INFO" : {
      "result" : "",
      "status" : ""
    },
    "RUN_ID" : ""
  },
  "setupData" : {
    "action" : {
      "apiserver" : "",
      "datafrom" : "{}",
      "image" : {
        "name" : "xiechuan/pythoncheck",
        "tag" : "1.0"
      },
      "ip" : "",
      "name" : "pythoncheck",
      "timeout" : "9000",
      "type" : "Kubernetes",
      "useAdvanced" : false
    },
    "pod" : {
      "spec" : {
        "containers" : [
          {
            "resources" : {
              "limits" : {
                "cpu" : 2,
                "memory" : "1024Mi"
              },
              "requests" : {
                "cpu" : 1,
                "memory" : "128Mi"
              }
            }
          }
        ]
      }
    },
    "pod_advanced" : { },
    "service" : {
      "spec" : {
        "ports" : [
          {
            "nodePort" : 32007,
            "port" : 8000,
            "targetPort" : 8000
          }
        ],
        "type" : "NodePort"
      }
    },
    "service_advanced" : { }
  }
}
```

### saveComponentInfo/saveComponentAsNewVersion

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  |/pipeline/v1/:namespace/component/:componentName|

#### body

```
{
  "id": 1,
  "version": "v1.0",
  "define": {
    "env": [
      {
        "key": "CO_DATA",
        "value": "{     \"serverities\": \"MAJOR\",     \"contents\": \"sonar.projectKey:python-sonar-runner\\nsonar.projectName=python-sonar-runner\\nsonar.projectVersion=1.0\\nsonar.sources=src\\nsonar.language=py\\nsonar.sourceEncoding=UTF-8\",     \"filename\": \"sonar-project.properties\" }"
      }
    ],
    "inputJson": {
      "gitUrl": "https://github.com/xiechuanj/python-sonar-runner.git"
    },
    "outputJson": {
      "EVENT": "TASK_RESULT",
      "EVENTID": "",
      "INFO": {
        "result": "",
        "status": ""
      },
      "RUN_ID": ""
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
        "timeout": "9000",
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
                  "memory": "1024Mi"
                },
                "requests": {
                  "cpu": "0.1",
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
              "nodePort": 32007,
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

### createComponent

| HTTP Method |  Request Address |
| -------- | ------ |
| PUT  |/pipeline/v1/demo/component|

#### body

```
{
  "name": "javaCheck",
  "version": "v3.0"
}
```

#### response json

```
{
  "message": "create new component success"
}
```

### get pipelien token and url

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/:namespace/:repository/:pipelineName/token?id=:pipelineId|

#### response json

```
{
  "token": "ed97b3cba1426429423fa13eeb97c1b2",
  "url": "http://192.168.137.1/demo/demo/go-codecheck"
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
      "id": 25,
      "name": "python-demo",
      "versionList": [
        {
          "id": 25,
          "info": "Success :0 Total :0",
          "name": "v1",
          "sequenceList": []
        }
      ]
    },
    {
      "id": 26,
      "name": "python",
      "versionList": [
        {
          "id": 26,
          "info": "Success :0 Total :0",
          "name": "v2",
          "sequenceList": []
        }
      ]
    },
    {
      "id": 24,
      "name": "python",
      "versionList": [
        {
          "id": 24,
          "info": "Success :0 Total :8",
          "name": "v1",
          "sequenceList": [
            {
              "pipelineSequenceID": 141,
              "sequence": 1,
              "status": false,
              "time": "2016-10-27 11:32:46"
            },
            {
              "pipelineSequenceID": 142,
              "sequence": 2,
              "status": false,
              "time": "2016-10-27 11:34:54"
            },
            {
              "pipelineSequenceID": 143,
              "sequence": 3,
              "status": false,
              "time": "2016-10-27 11:36:51"
            },
            {
              "pipelineSequenceID": 144,
              "sequence": 4,
              "status": false,
              "time": "2016-10-27 11:39:13"
            },
            {
              "pipelineSequenceID": 145,
              "sequence": 5,
              "status": false,
              "time": "2016-10-27 11:42:46"
            },
            {
              "pipelineSequenceID": 146,
              "sequence": 6,
              "status": false,
              "time": "2016-10-27 11:46:26"
            },
            {
              "pipelineSequenceID": 147,
              "sequence": 7,
              "status": false,
              "time": "2016-10-27 11:48:56"
            },
            {
              "pipelineSequenceID": 148,
              "sequence": 8,
              "status": false,
              "time": "2016-10-27 11:52:11"
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
| GET  | /pipeline/v1/:namespace/:repository/:pipelineName/historyDefine?versionId={versionId}sequenceId={pipelineSequenceID}|

### getStageRunHistory

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/:namespace/:repository/:pipelineName/stage/:stageName/history?stageLogId={stageLogID}|

### getActionRunHistory

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/:namespace/:repository/:pipelineName/stage/:stageName/:actionName/history?actionLogId={actionLogID}|

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
      "2016-10-27 11:52:13 -> {\"EVENTID\": 1048, \"EVENT\": \"COMPONENT_START\", \"RUN_ID\": \"148,497,210,8,9\"}",
      "2016-10-27 11:52:23 -> {\"data\":\"{\"gitUrl\":\"https://github.com/xiechuanj/python-sonar-runner.git\"}\",\"resp\":\"{\"gitUrl\":\"https://github.com/xiechuanj/python-sonar-runner.git\"}\\r\\n\"}",
      "2016-10-27 11:52:23 -> {\"EVENTID\": 1050, \"EVENT\": \"TASK_START\", \"RUN_ID\": \"148,497,210,8,9\"}",
      "2016-10-27 11:53:04 -> {\"EVENTID\": 1052, \"INFO\": {\"TASK_STATUS\": \"RUNNING\"}, \"EVENT\": \"TASK_STATUS\", \"RUN_ID\": \"148,497,210,8,9\"}",
      "2016-10-27 11:53:12 -> {\"EVENTID\": 1052, \"INFO\": {\"TASK_STATUS\": \"GET RESULT\"}, \"EVENT\": \"TASK_STATUS\", \"RUN_ID\": \"148,497,210,8,9\"}",
      "2016-10-27 11:53:22 -> {\"EVENTID\": 1051, \"INFO\": {\"status\": false, \"result\": \"{\"total\":27,\"p\":1,\"ps\":100,\"paging\":{\"pageIndex\":1,\"pageSize\":100,\"total\":27},\"issues\":[{\"key\":\"AVf_WUAXuMltxfy4_Gy5\",\"rule\":\"python:S1110\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/badfortune.py\",\"componentId\":4,\"project\":\"python-sonar-runner\",\"line\":30,\"textRange\":{\"startLine\":30,\"endLine\":30,\"startOffset\":12,\"endOffset\":28},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Remove those useless parentheses\",\"effort\":\"1min\",\"debt\":\"1min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"confusing\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUAZuMltxfy4_Gy6\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/badfortune.py\",\"componentId\":4,\"project\":\"python-sonar-runner\",\"line\":90,\"textRange\":{\"startLine\":90,\"endLine\":90,\"startOffset\":8,\"endOffset\":13},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUAauMltxfy4_Gy7\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/badfortune.py\",\"componentId\":4,\"project\":\"python-sonar-runner\",\"line\":92,\"textRange\":{\"startLine\":92,\"endLine\":92,\"startOffset\":4,\"endOffset\":9},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUB0uMltxfy4_GzZ\",\"rule\":\"python:S1110\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/strfile.py\",\"componentId\":14,\"project\":\"python-sonar-runner\",\"line\":60,\"textRange\":{\"startLine\":60,\"endLine\":60,\"startOffset\":6,\"endOffset\":9},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Remove those useless parentheses\",\"effort\":\"1min\",\"debt\":\"1min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"confusing\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUB1uMltxfy4_Gza\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/strfile.py\",\"componentId\":14,\"project\":\"python-sonar-runner\",\"line\":28,\"textRange\":{\"startLine\":28,\"endLine\":28,\"startOffset\":4,\"endOffset\":9},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUB1uMltxfy4_Gzb\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/strfile.py\",\"componentId\":14,\"project\":\"python-sonar-runner\",\"line\":97,\"textRange\":{\"startLine\":97,\"endLine\":97,\"startOffset\":0,\"endOffset\":5},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBDuMltxfy4_Gy9\",\"rule\":\"python:S1110\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/fortune.py\",\"componentId\":12,\"project\":\"python-sonar-runner\",\"line\":30,\"textRange\":{\"startLine\":30,\"endLine\":30,\"startOffset\":12,\"endOffset\":28},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Remove those useless parentheses\",\"effort\":\"1min\",\"debt\":\"1min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"confusing\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBEuMltxfy4_Gy-\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/fortune.py\",\"componentId\":12,\"project\":\"python-sonar-runner\",\"line\":90,\"textRange\":{\"startLine\":90,\"endLine\":90,\"startOffset\":8,\"endOffset\":13},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBFuMltxfy4_Gy_\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/fortune.py\",\"componentId\":12,\"project\":\"python-sonar-runner\",\"line\":92,\"textRange\":{\"startLine\":92,\"endLine\":92,\"startOffset\":4,\"endOffset\":9},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBTuMltxfy4_GzD\",\"rule\":\"python:S1110\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":84,\"textRange\":{\"startLine\":84,\"endLine\":84,\"startOffset\":10,\"endOffset\":13},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Remove those useless parentheses\",\"effort\":\"1min\",\"debt\":\"1min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"confusing\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBUuMltxfy4_GzE\",\"rule\":\"python:S1110\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":149,\"textRange\":{\"startLine\":149,\"endLine\":149,\"startOffset\":15,\"endOffset\":44},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Remove those useless parentheses\",\"effort\":\"1min\",\"debt\":\"1min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"confusing\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBVuMltxfy4_GzF\",\"rule\":\"python:S1110\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":151,\"textRange\":{\"startLine\":151,\"endLine\":151,\"startOffset\":17,\"endOffset\":64},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Remove those useless parentheses\",\"effort\":\"1min\",\"debt\":\"1min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"confusing\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBWuMltxfy4_GzG\",\"rule\":\"python:S125\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":73,\"textRange\":{\"startLine\":73,\"endLine\":73,\"startOffset\":4,\"endOffset\":73},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Remove this commented out code.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"misra\",\"unused\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBZuMltxfy4_GzH\",\"rule\":\"python:S125\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":86,\"textRange\":{\"startLine\":86,\"endLine\":86,\"startOffset\":12,\"endOffset\":93},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Remove this commented out code.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"misra\",\"unused\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBcuMltxfy4_GzI\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":99,\"textRange\":{\"startLine\":99,\"endLine\":99,\"startOffset\":12,\"endOffset\":17},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBduMltxfy4_GzJ\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":110,\"textRange\":{\"startLine\":110,\"endLine\":110,\"startOffset\":16,\"endOffset\":21},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBeuMltxfy4_GzK\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":115,\"textRange\":{\"startLine\":115,\"endLine\":115,\"startOffset\":16,\"endOffset\":21},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBfuMltxfy4_GzL\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":116,\"textRange\":{\"startLine\":116,\"endLine\":116,\"startOffset\":16,\"endOffset\":21},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBguMltxfy4_GzM\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":120,\"textRange\":{\"startLine\":120,\"endLine\":120,\"startOffset\":16,\"endOffset\":21},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBhuMltxfy4_GzN\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":145,\"textRange\":{\"startLine\":145,\"endLine\":145,\"startOffset\":16,\"endOffset\":21},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBhuMltxfy4_GzO\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":146,\"textRange\":{\"startLine\":146,\"endLine\":146,\"startOffset\":16,\"endOffset\":21},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBiuMltxfy4_GzP\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":148,\"textRange\":{\"startLine\":148,\"endLine\":148,\"startOffset\":12,\"endOffset\":17},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBjuMltxfy4_GzQ\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":150,\"textRange\":{\"startLine\":150,\"endLine\":150,\"startOffset\":16,\"endOffset\":21},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBjuMltxfy4_GzR\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":152,\"textRange\":{\"startLine\":152,\"endLine\":152,\"startOffset\":16,\"endOffset\":21},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBkuMltxfy4_GzS\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":154,\"textRange\":{\"startLine\":154,\"endLine\":154,\"startOffset\":12,\"endOffset\":17},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBluMltxfy4_GzT\",\"rule\":\"python:PrintStatementUsage\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/letters.py\",\"componentId\":13,\"project\":\"python-sonar-runner\",\"line\":157,\"textRange\":{\"startLine\":157,\"endLine\":157,\"startOffset\":0,\"endOffset\":5},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Replace print statement by built-in function.\",\"effort\":\"5min\",\"debt\":\"5min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"obsolete\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"},{\"key\":\"AVf_WUBzuMltxfy4_GzY\",\"rule\":\"python:S1110\",\"severity\":\"MAJOR\",\"component\":\"python-sonar-runner:src/samples/strfile.py\",\"componentId\":14,\"project\":\"python-sonar-runner\",\"line\":34,\"textRange\":{\"startLine\":34,\"endLine\":34,\"startOffset\":12,\"endOffset\":28},\"flows\":[],\"status\":\"OPEN\",\"message\":\"Remove those useless parentheses\",\"effort\":\"1min\",\"debt\":\"1min\",\"author\":\"xiechuanj@gmail.com\",\"tags\":[\"confusing\"],\"creationDate\":\"2016-10-26T04:56:56+0000\",\"updateDate\":\"2016-10-26T04:56:56+0000\",\"type\":\"CODE_SMELL\"}],\"components\":[{\"id\":12,\"key\":\"python-sonar-runner:src/samples/fortune.py\",\"uuid\":\"AVf_WTteuMltxfy4_Gy1\",\"enabled\":true,\"qualifier\":\"FIL\",\"name\":\"fortune.py\",\"longName\":\"src/samples/fortune.py\",\"path\":\"src/samples/fortune.py\",\"projectId\":1,\"subProjectId\":1},{\"id\":13,\"key\":\"python-sonar-runner:src/samples/letters.py\",\"uuid\":\"AVf_WTteuMltxfy4_Gy2\",\"enabled\":true,\"qualifier\":\"FIL\",\"name\":\"letters.py\",\"longName\":\"src/samples/letters.py\",\"path\":\"src/samples/letters.py\",\"projectId\":1,\"subProjectId\":1},{\"id\":14,\"key\":\"python-sonar-runner:src/samples/strfile.py\",\"uuid\":\"AVf_WTteuMltxfy4_Gy3\",\"enabled\":true,\"qualifier\":\"FIL\",\"name\":\"strfile.py\",\"longName\":\"src/samples/strfile.py\",\"path\":\"src/samples/strfile.py\",\"projectId\":1,\"subProjectId\":1},{\"id\":1,\"key\":\"python-sonar-runner\",\"uuid\":\"AVf_WTIKl7uABRzX3_me\",\"enabled\":true,\"qualifier\":\"TRK\",\"name\":\"python-sonar-runner\",\"longName\":\"python-sonar-runner\"},{\"id\":4,\"key\":\"python-sonar-runner:src/badfortune.py\",\"uuid\":\"AVf_WTtcuMltxfy4_Gyt\",\"enabled\":true,\"qualifier\":\"FIL\",\"name\":\"badfortune.py\",\"longName\":\"src/badfortune.py\",\"path\":\"src/badfortune.py\",\"projectId\":1,\"subProjectId\":1}]}\"}, \"EVENT\": \"TASK_RESULT\", \"RUN_ID\": \"148,497,210,8,9\"}",
      "2016-10-27 11:53:22 -> \t\t\t{\"EVENTID\": 1049, \"EVENT\": \"COMPONENT_STOP\", \"RUN_ID\": \"148,497,210,8,9\"}"
    ]
  }
}
```

### getLineDataInfo

| HTTP Method |  Request Address |
| -------- | ------ |
| GET  | /pipeline/v1/:namespace/:repository/:pipelineName/:pipelineSequenceID/lineHistory?startActionId={startActionId}&endActionId={endActionId}|

#### response json

```
{
  "define": {
    "input": {
      "data":"XXXXXXXXXX"
    },
    "output": {
      "gocyclo": "{\"url\":\"https://github.com/baxterthehacker/public-repo.git\"}"
    }
  }
}
```
