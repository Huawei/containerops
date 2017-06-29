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

import {loading} from "./loading";
import {notify} from "./notify";

let apiUrlConf = {
	host : "",

	workflow : {
		rootUrl : "/v2/{namespace}/{repository}/workflow/v1/define",

		list : {
			"url" : "/list",
			"type" : "GET",
		},

		detail : {
			"url" :	"/{workflowName}?id={workflowID}",
			"type" : "GET"
		},

		add : {
			"url" :	"",
			"type" : "POST"
		},

		save : {
			"url" :	"/{workflowName}",
			"type" : "PUT"
		},

		eventOutput : {
			"url" :	"/event/{site}/{eventName}",
			"type" : "GET",
			"skipAbort" : true
		},

		getEnv : {
			"url" :	"/{workflowName}/env?id={workflowID}",
			"type" : "GET"
		},

		setEnv : {
			"url" :	"/{workflowName}/env",
			"type" : "PUT"
		},

		getVar : {
			"url" :	"/{workflowName}/var?id={workflowID}",
			"type" : "GET"
		},

		setVar : {
			"url" :	"/{workflowName}/var",
			"type" : "PUT"
		},

		changeState : {
			"url" :	"/{workflowName}/state",
			"type" : "PUT"
		},

		getToken : {
			"url" :	"/{workflowName}/token?id={workflowID}",
			"type" : "GET"
		}
	},

	component : {
		rootUrl : "/v2/{namespace}/component",

		list : {
			"url" :	"/list",
			"type" : "GET"
		},

		detail : {
			"url" :	"/{componentName}?id={componentID}",
			"type" : "GET"
		},

		add : {
			"url" :	"",
			"type" : "POST"
		},

		save : {
			"url" :	"/{componentName}",
			"type" : "PUT"
		}
	},

	history : {
		rootUrl : "/v2/{namespace}/{repository}/workflow/v1/history",

		workflowHistory : {
			"url" :	"/{workflowName}/{version}?sequence={sequence}",
			"type" : "GET"
		},

		action : {
			"url" :	"/{workflowName}/{version}/{sequence}/stage/{stageName}/action/{actionName}",
			"type" : "GET"
		},

		relation : {
			"url" :	"/{workflowName}/{version}/{sequence}/{lineId}",
			"type" : "GET"
		},

		containerLog: {
			"url" :	"/{workflowName}/{version}/{sequence}/stage/{stageName}/action/{actionName}/console/log?key={key}&size=10",
			"type" : "GET",
			"skipAbort" : true
		},

		list: {
			"url" :	"/workflow/list?page={page}&prePageCount={workflowNum}&filter={keywords}&filtertype={filterType}",
			"type" : "GET"
		},

		version : {
			"url" :	"/workflow/{workflowName}/version/list?id={workflowID}",
			"type" : "GET"
		},

		sequence : {
			"url" :	"/workflow/{workflowName}/version/{versionName}/list?id={versionID}&sequenceNum={sequenceNum}",
			"type" : "GET"
		},

		startedWorkflow : {
			"url" :	"/workflow/{workflowName}/version/{versionName}/sequence/{sequence}/action/{actionName}/linkstart/list?workflowId={sequenceID}&actionId={actionID}",
			"type" : "GET"
		}
	},

	setting : {
		rootUrl : "/v2/{namespace}/{repository}/system/v1/setting",

		list : {
			"url" :	"",
			"type" : "GET"
		},

		save : {
			"url" :	"",
			"type" : "PUT"
		}
	}
}

let userInfo, pendingPromise;

export function initApi(namespace,repository){
	userInfo = {
		"namespace" : namespace,
		"repository" : repository
	}
}

// abort
function beforeApiInvocation(skipAbort){
	if(_.isEmpty(apiUrlConf.host)){
		$.ajax({
	        "url": "/host.json",
	        "async" : false,
	        "type": "GET",
	        "dataType": "json",
	        "cache": false,
	        "success" : function(obj) {
			    apiUrlConf.host = obj.host;
			},
			"error" : function(error){
				notify("Can not find API host configuration file.","error");
			}
	    });
	}

	if(!skipAbort){
		_.each(pendingPromise,function(promise){
			promise.abort();
		});			
		pendingPromise = [];
	}
	loading.show();
}

export function ajaxCall(target, params, reqbody){
	var urls = target.split(".");
	beforeApiInvocation(apiUrlConf[urls[0]][urls[1]].skipAbort);
	
	var urlroot = getUrlRoot(urls[0]);
	var urlext = getUrlExt(urls[0],urls[1],params);
	var type = apiUrlConf[urls[0]][urls[1]].type;

	var options;
	if(type == "GET"){
		options = {
	        "url": apiUrlConf.host + urlroot + urlext,
	        "type": type,
	        "dataType": "json",
	        "cache": false
	    }
	}else if(type == "POST" || type == "PUT"){
		var data = JSON.stringify(reqbody);
		options = {
	        "url": apiUrlConf.host + urlroot + urlext,
	        "type": type,
	        "dataType": "json",
	        "data": data
	    }
	}
	var promise = $.ajax(options);
	pendingPromise.push(promise);
	return promise;
}

function getUrlRoot(type){
	var rootUrl = apiUrlConf[type].rootUrl;
	return rootUrl.replace(/{namespace}/g, userInfo.namespace).replace(/{repository}/g, userInfo.repository);
}

function getUrlExt(type,ext,params){
	var extensionUrl = apiUrlConf[type][ext].url;
	var paramKeys = _.keys(params);
	_.each(paramKeys,function(key){
		if(extensionUrl.indexOf("{"+key+"}") >= 0){
			var regexp = new RegExp("{"+key+"}","g");
			extensionUrl = extensionUrl.replace(regexp, params[key]);
		}
	});
	return extensionUrl;
}
