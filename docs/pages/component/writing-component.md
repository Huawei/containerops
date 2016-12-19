---
title: How Write A DevOps Component 
keywords: component
tags: [component]
sidebar: home_sidebar
permalink: writing-component.html
summary: How Write A DevOps Component
folder: component
---

## Write A Execute Task

What's a task?

## How To Got And Export A Pipeline Event Data

## Get A Event Data From REST API

## Get A Event Data From Environment Variable

| Variable |  Value |
| -------- | ------ |
| CO_DATA  | The component data| 

## Export Event Data

## How To Definition The Event Data With Editor

## System Callback Environment Variable

When a new DevOps component container created, the DevOps workflow engine will set some environment variables automatically. All the variables is REST API URL, the DevOps component should call it for passing the component information or status to the DevOps workflow execute engine.

|  Variable       |  Value |
| --------------- | --------- |
| CO_COMPONENT_START | When the container of DevOps component start completely include all the dependencies started, the component should call the REST API of *COMPONENT_START* passing the start status status to the engine. Then the workflow execute engine will monitor the container status via the orchestration tools like Kubernetes. And the execute engine will call the component passing the *event* data. |
| CO_TASK_START      | When the container of DevOps component get all datas via *event* data or volume data, call the REST API of *TASK_START* passing the task start execute status to the engine. |
| CO_TASK_STATUS     | When the container of DevOps component execute the task, it should call the REST API of *TASK_STATUS* passing the interim outputs to the execute engine repeatly. |
| CO_TASK_RESULT     | When the container of DevOps component execute successfully or failure, it should call the REST API of *TASK_RESULT* passing the result and final output to the execute engine. |
| CO_COMPONENT_STOP  | When the program of task in the container of DevOps component stop completely, it should call REST API of *COMPONENT_STOP* passing the stop status to the execute engine. The engine will notify the orchestration tools destory the container and release the resource. |