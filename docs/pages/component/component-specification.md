---
title: Component Sepecification
keywords: component
tags: [component]
sidebar: home_sidebar
permalink: component-specification.html
summary: Component Sepecification
folder: component
---

### 1. Following all [] value can be found the docker container env variable

### 2. Each component of the implementation process
- __After the component started to [COMPONENT_START] url send EVENT type for the COMPONENT_START EVENT notifications__

- __After the task started to [TASK_START] url send EVENT type for the TASK_START EVENT notifications__

- __When the task execing to [TASK_RESULT] url send EVENT type for the TASK_RESULT EVENT notifications__

- __After the task finished to [TASK_STATUS] url send EVENT type for the TASK_STATUS EVENT notifications__

- __After the task finished to [COMPONENT_END] url send EVENT type for the COMPONENT_END EVENT notifications__

### 3.  send EVENT notifications

| HTTP Method |  Request Address |
| -------- | ------ |
| POST  |[COMPONENT_START]/[TASK_START]/[TASK_RESULT]/[TASK_STATUS]/[COMPONENT_END]|

#### body

```
{
  "EVENT": "TASK_RESULT", #COMPONENT_START/TASK_START/TASK_RESULT/TASK_STATUS/COMPONENT_END
  "EVENTID": 5260,
  "RUN_ID": "288,1197,877,36,84"
  "INFO": {
    "output": {
      "status":true,
      "result":"",
      "output":""
    },
  }
}
```

#### response json

```
{
  "message": "ok"
}
```

### 4.  Env variable
All env variable will start with CO_
- __EVENT types__

| KEY| VALUE |
| -------- | ------ |
|CO_COMPONENT_START  |[CO)COMPONENT_START]|
|CO_TASK_START|[CO_TASK_START]|
|CO_TASK_RESULT|[CO_TASK_RESULT]|
|CO_TASK_STATUS|[CO_TASK_STATUS]|
|CO_COMPONENT_END|[CO_COMPONENT_END]|

- __SERVER ADDRESS__

| KEY| VALUE |
| -------- | ------ |
| CO_SERVICE_ADDR  |[Cluster IP]:[Cluster PORT]:[ContainerListenPort]|

- __Component type__

| KEY| VALUE |
| -------- | ------ |
| CO_RUN_ID  |[CO_RUN_ID]|


- __Event List__

| KEY| VALUE |
| -------- | ------ |
| CO_EVENT_LIST  |CO_COMPONENT_START,1;CO_COMPONENT_STOP,2;CO_TASK_START,8;CO_TASK_RESULT,9;CO_TASK_STATUS,10;CO_REGISTER_URL,11|

- __Change Env__

| KEY| VALUE |
| -------- | ------ |
| CO_SET_GLOBAL_VAR_URL  |[CO_SET_GLOBAL_VAR_URL]|

#### body

```
{
  "RUN_ID"：[CO_RUN_ID],
  "varMap": {
      "KEY":[KEY],
      "VALUE":[VALUE],
  }
}
```

#### response json

```
{
  "message": "ok"
}
```

- __Set Gobal Env__

| KEY| VALUE |
| -------- | ------ |
| CO_SET_GLOBAL_VAR_URL  |[CO_SET_GLOBAL_VAR_URL]|

#### body

```
{
  "RUN_ID"：[CO_RUN_ID],
  "varMap": {
      "KEY":[KEY],
      "VALUE":[VALUE],
  }
}
```

#### response json

```
{
  "message": "ok"
}
```

- __Call other workflow__

| KEY| VALUE |
| -------- | ------ |
| CO_LINKSTART_URL  |[CO_LINKSTART_URL]|

#### body

```
{
  "RUN_ID"：[CO_RUN_ID],
  "linkInfoMap": {
      "token":[token],
      "workflowName":[workflowName],
      "workflowVersion":[workflowVersion],
      "startJson":[startJson],
  }
}
```

#### response json

```
{
  "message": "ok"
}
```
