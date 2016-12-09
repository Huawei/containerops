/* 
Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
s
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
 */
import { resizeWidget } from "../theme/widget";

let setting;
export function initWorkflowSetting(workflowSettingData){
	if(!_.isUndefined(workflowSettingData) && !_.isEmpty(workflowSettingData)){
      setting = workflowSettingData;
    }else{
      setting = $.extend(true,{},metadata);
      workflowSettingData = setting;
    }

    showWorkflowSetting();
}

function showWorkflowSetting(){
	$.ajax({
		url: "../../templates/workflow/workflowSetting.html",
		type: "GET",
		cache: false,
		success: function(data) {
			$("#workflow-info-edit").html($(data));
			resizeWidget();
		}
	})
}

var metadata = {
	"runningInstances" : {
		"available" : true,
		"number" : 10
	},
	"timedTasks" : {
		"available" : true,
		"tasks" : []
	}
}