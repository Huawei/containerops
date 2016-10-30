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
		// "sequenceData" : "/pipeline/v1/demo/demo/{pipelineName}/historyDefine?sequenceId={pipelineSequenceID}",
		// https://test-1.containerops.sh/pipeline/v1/demo/demo/histories
		"sequenceData" : ":8080/pipelineInfo?id={pipelineRunSequenceID}",
		"host" : "http://localhost"
		// "sequenceEventShow" : ""
	}
}

// pipeline
export let pipelineApi = {
	"list" : function(){
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.list,
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return promise;
	},
	"data" : function(name,id){
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.data.replace(/{pipelineName}/g, name).replace(/{pipelineID}/g, id),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return promise;
	},
	"add" : function(name,version){
		var data = JSON.stringify({
				"name":name,
				"version":version
			});
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.add,
	        "type": "POST",
	        "dataType": "json",
	        "data": data
	    });
	    return promise;
	},
	"save" : function(name,reqbody){
		var data = JSON.stringify(reqbody);
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.save.replace(/{pipelineName}/g, name),
	        "type": "PUT",
	        "dataType": "json",
	        "data": data
	    });
	    return promise;
	},
	"eventOutput" : function(name){
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.eventOutput.replace(/{eventName}/g, name),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return promise;
	},
	"getEnv" : function(name,id){
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.getEnv.replace(/{pipelineName}/g, name).replace(/{pipelineID}/g, id),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return promise;
	},
	"setEnv" : function(name,reqbody){
		var data = JSON.stringify(reqbody);
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.setEnv.replace(/{pipelineName}/g, name),
	        "type": "PUT",
	        "dataType": "json",
	        "data": data
	    });
	    return promise;
	},
	"changeState" : function(name,reqbody){
		var data = JSON.stringify(reqbody);
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.changeState.replace(/{pipelineName}/g, name),
	        "type": "PUT",
	        "dataType": "json",
	        "data": data
	    });
	    return promise;
	},
	"getToken" : function(name,id){
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.getToken.replace(/{pipelineName}/g, name).replace(/{pipelineID}/g, id),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return promise;
	}
}

// component
export let componentApi = {
	"list" : function(){
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.component.list,
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return promise;
	},
	"data" : function(name,id){
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.component.data.replace(/{componentName}/g, name).replace(/{componentID}/g, id),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return promise;
	},
	"add" : function(name,version){
		var data = JSON.stringify({
				"name":name,
				"version":version
			});
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.component.add,
	        "type": "POST",
	        "dataType": "json",
	        "data": data
	    });
	    return promise;
	},
	"save" : function(name,reqbody){
		var data = JSON.stringify(reqbody);
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.component.save.replace(/{componentName}/g, name),
	        "type": "PUT",
	        "dataType": "json",
	        "data": data
	    });
	    return promise;
	}
}

// history
export let historyApi = {
	"sequenceData" : function(pipelineRunSequenceID){
		var promise = $.ajax({
	        "url": apiUrlConf.history.host + apiUrlConf.history.sequenceData.replace(/{pipelineRunSequenceID}/g, pipelineRunSequenceID) ,
	        // "url": apiUrlConf.host + apiUrlConf.history.sequenceData.replace(/{pipelineName}/g, pipelineName).replace(/{pipelineSequenceID}/g, pipelineRunSequenceID),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    });
	    return promise;
	},
	"sequenceList" : function () {
		var promise = $.ajax({
			"url" : apiUrlConf.host + apiUrlConf.history.sequenceList,
			"type" : "GET",
			"dataType" : "json",
			"cache": false
		});
		return promise;
	}
}
