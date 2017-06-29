/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import {ajaxCall} from "../common/api";

export function getWorkflowHistory(workflowName,versionName,workflowRunSequence){
    var params = {
      "workflowName" : workflowName,
      "version" : versionName,
      "sequence" : workflowRunSequence
    }

    return ajaxCall("history.workflowHistory",params);
}

export function getActionRunHistory(workflowName,versionName,workflowRunSequence,stageName,actionName){
    var params = {
      "workflowName" : workflowName,
      "version" : versionName,
      "sequence" : workflowRunSequence,
      "stageName" : stageName,
      "actionName" : actionName
    }

    return ajaxCall("history.action",params);
}

export function getLineDataInfo(workflowName,versionName,workflowRunSequence,sequenceLineId){
    var params = {
      "workflowName" : workflowName,
      "version" : versionName,
      "sequence" : workflowRunSequence,
      "lineId" : sequenceLineId
    }

    return ajaxCall("history.relation",params);
}

export function getContainerLogsData(workflowName,versionName,workflowRunSequence,stageName,actionName,key){
    var params = {
      "workflowName" : workflowName,
      "version" : versionName,
      "sequence" : workflowRunSequence,
      "stageName" : stageName,
      "actionName" : actionName,
      "key" : key
    }

	  return ajaxCall("history.containerLog",params);
}

export function getWorkflows(page,workflowNum,keywords,filterType){
    var params = {
      "page" : page,
      "workflowNum" : workflowNum,
      "keywords" : keywords,
      "filterType" : filterType
    }

    return ajaxCall("history.list",params);
}

export function getVersions(workflowName,workflowId){
    var params = {
      "workflowName" : workflowName,
      "workflowID" : workflowId
    }

    return ajaxCall("history.version",params);
}

export function getSequences(workflowName,workflowId,versionName,versionId,sequenceNum){
    var params = {
      "workflowName" : workflowName,
      "versionName" : versionName,
      "versionID" : versionId,
      "sequenceNum" : sequenceNum
    }

    return ajaxCall("history.sequence",params);
}

export function getStartedWorkflows(workflowName,workflowId,version,sequence,sequenceId,stageName,actionId,actionName){
    var params = {
      "workflowName" : workflowName,
      "versionName" : version.versionName,
      "versionId" : version.versionId,
      "sequence" : sequence,
      "sequenceID" : sequenceId,
      "stageName" : stageName,
      "actionName" : actionName,
      "workflowID" : workflowId,
      "actionID" : actionId
    }

    return ajaxCall("history.startedWorkflow",params);
}
