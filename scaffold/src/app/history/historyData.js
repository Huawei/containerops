/*
Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.

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

import {historyApi} from "../common/api";

export function getHistoryData(params){
  return historyApi.getHistoryData(params);
}

export function getWorkflowHistories( ){
    return historyApi.workflowHistories( );
}

export function getWorkflowHistory(workflowName,versionName,workflowRunSequence){
    return historyApi.workflowHistory(workflowName,versionName,workflowRunSequence);
}

export function getActionRunHistory(workflowName,versionName,workflowRunSequence,stageName,actionName){
    return historyApi.action(workflowName,versionName,workflowRunSequence,stageName,actionName);
}

export function getLineDataInfo(workflowName,versionName,workflowRunSequence,sequenceLineId){
    return historyApi.relation(workflowName,versionName,workflowRunSequence,sequenceLineId);
}

// export function getContainerLog(){
// 	return historyApi.containerLog();
// }
// export function getScheduleLog(){
// 	return historyApi.scheduleLog();
// }

export function getWorkflows(page,workflowNum,isInitPages){
  return historyApi.getWorkflows(page,workflowNum,isInitPages);
}

export function getVersions(workflowName,workflowId){
  return historyApi.getVersions(workflowName,workflowId);
}

export function getSequences(workflowName,workflowId,versionName,versionId,sequenceNum){
  return historyApi.getSequences(workflowName,workflowId,versionName,versionId,sequenceNum);
}

export function getStartedWorkflows(workflowName,workflowId,version,sequence,sequenceId,stageName,actionId,actionName){
  return historyApi.getStartedWorkflows(workflowName,workflowId,version,sequence,sequenceId,stageName,actionId,actionName);
}




// export function sequenceData(workflowName,versionID,workflowRunSequenceID){
//     return historyApi.sequenceData(workflowName,versionID,workflowRunSequenceID);
// }

// export function sequenceList( ){
//     return historyApi.sequenceList( );
//}
