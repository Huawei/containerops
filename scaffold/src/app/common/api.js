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
 import {loading} from "./loading";

let apiUrlConf = {
	"host" : "https://test-1.containerops.sh",
	"pipeline" : {
		"list" : "/pipeline/v1/demo/demo",
		"data" : "/pipeline/v1/demo/demo/{pipelineName}/json?id={pipelineID}",
		"add" : "/pipeline/v1/demo/demo",
		"save" : "/pipeline/v1/demo/demo/{pipelineName}",
		"eventOutput" : "/pipeline/v1/eventJson/github/{eventName}",
		"getEnv" : "/pipeline/v1/demo/demo/{pipelineName}/env?id={pipelineID}",
		"setEnv" : "/pipeline/v1/demo/demo/{pipelineName}/env",
		"changeState" : "/pipeline/v1/demo/demo/{pipelineName}/state",
		"getToken" : "/pipeline/v1/demo/demo/{pipelineName}/token?id={pipelineID}"
	},
	"component" : {
		"list" : "/pipeline/v1/demo/component",
		"data" : "/pipeline/v1/demo/component/{componentName}?id={componentID}",
		"add" : "/pipeline/v1/demo/component",
		"save" : "/pipeline/v1/demo/component/{componentName}"
	},
	"history" : {
		"sequenceList" : "/pipeline/v1/demo/demo/histories",
		"sequenceData" : "/pipeline/v1/demo/demo/{pipelineName}/historyDefine?versionId={versionID}&sequenceId={pipelineSequenceID}",
		"action" : "/pipeline/v1/demo/demo/{pipelineName}/stage/{stageName}/{actionName}/history?actionLogId={actionLogID}",
		"relation" : "/pipeline/v1/demo/demo/{pipelineName}/{pipelineSequenceID}/lineHistory?startActionId={startActionId}&endActionId={endActionId}"
	}
}

let pendingPromise;

// abort
function abortPendingPromise(){
	if(pendingPromise){
		pendingPromise.abort();
	}
	loading.show();
}

// pipeline
export let pipelineApi = {
	"list" : function(){
		abortPendingPromise();
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.list,
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return pendingPromise;
	},
	"data" : function(name,id){
		abortPendingPromise();
		pendingPromise= $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.data.replace(/{pipelineName}/g, name).replace(/{pipelineID}/g, id),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return pendingPromise;
	},
	"add" : function(name,version){
		abortPendingPromise();
		var data = JSON.stringify({
				"name":name,
				"version":version
			});
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.add,
	        "type": "POST",
	        "dataType": "json",
	        "data": data
	    });
	    return pendingPromise;
	},
	"save" : function(name,reqbody){
		abortPendingPromise();
		var data = JSON.stringify(reqbody);
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.save.replace(/{pipelineName}/g, name),
	        "type": "PUT",
	        "dataType": "json",
	        "data": data
	    });
	    return pendingPromise;
	},
	"eventOutput" : function(name){
		abortPendingPromise();
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.eventOutput.replace(/{eventName}/g, name),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return pendingPromise;
	},
	"getEnv" : function(name,id){
		abortPendingPromise();
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.getEnv.replace(/{pipelineName}/g, name).replace(/{pipelineID}/g, id),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return pendingPromise;
	},
	"setEnv" : function(name,reqbody){
		abortPendingPromise();
		var data = JSON.stringify(reqbody);
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.setEnv.replace(/{pipelineName}/g, name),
	        "type": "PUT",
	        "dataType": "json",
	        "data": data
	    });
	    return pendingPromise;
	},
	"changeState" : function(name,reqbody){
		abortPendingPromise();
		var data = JSON.stringify(reqbody);
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.changeState.replace(/{pipelineName}/g, name),
	        "type": "PUT",
	        "dataType": "json",
	        "data": data
	    });
	    return pendingPromise;
	},
	"getToken" : function(name,id){
		abortPendingPromise();
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.getToken.replace(/{pipelineName}/g, name).replace(/{pipelineID}/g, id),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return pendingPromise;
	}
}

// component
export let componentApi = {
	"list" : function(){
		abortPendingPromise();
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.component.list,
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return pendingPromise;
	},
	"data" : function(name,id){
		abortPendingPromise();
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.component.data.replace(/{componentName}/g, name).replace(/{componentID}/g, id),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return pendingPromise;
	},
	"add" : function(name,version){
		abortPendingPromise();
		var data = JSON.stringify({
				"name":name,
				"version":version
			});
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.component.add,
	        "type": "POST",
	        "dataType": "json",
	        "data": data
	    });
	    return pendingPromise;
	},
	"save" : function(name,reqbody){
		abortPendingPromise();
		var data = JSON.stringify(reqbody);
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.component.save.replace(/{componentName}/g, name),
	        "type": "PUT",
	        "dataType": "json",
	        "data": data
	    });
	    return pendingPromise;
	}
}

// history
export let historyApi = {
	"sequenceData" : function(pipelineName,versionID,pipelineRunSequenceID){
		abortPendingPromise();
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.history.sequenceData.replace(/{pipelineName}/g, pipelineName).replace(/{versionID}/g, versionID).replace(/{pipelineSequenceID}/g, pipelineRunSequenceID),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return pendingPromise;
	},
	"sequenceList" : function () {
		abortPendingPromise();
		pendingPromise = $.ajax({
			"url" : apiUrlConf.host + apiUrlConf.history.sequenceList,
			"type" : "GET",
			"dataType" : "json",
			"cache": false
		});
		return pendingPromise;
	},
	"action" : function(pipelineName,stageName,actionName,actionLogID){
		abortPendingPromise();
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.history.action.replace(/{pipelineName}/g, pipelineName).replace(/{stageName}/g, stageName).replace(/{actionName}/g, actionName).replace(/{actionLogID}/g, actionLogID),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return pendingPromise;
	},
	"relation" : function(pipelineName,pipelineSequenceID,startActionId,endActionId){
		abortPendingPromise();
		pendingPromise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.history.relation.replace(/{pipelineName}/g, pipelineName).replace(/{pipelineSequenceID}/g, pipelineSequenceID).replace(/{startActionId}/g, startActionId).replace(/{endActionId}/g, endActionId),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return pendingPromise;
	}
}
