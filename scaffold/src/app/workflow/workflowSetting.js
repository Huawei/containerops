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
import { cron } from "../common/cron";

let setting;
export function initWorkflowSetting(workflowSettingData){
	if(!_.isUndefined(workflowSettingData.data) && !_.isEmpty(workflowSettingData.data)){
      setting = workflowSettingData.data;
    }else{
      setting = $.extend(true,{},metadata);
      workflowSettingData["data"] = setting;
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

			initNumofRunningInstances();

			initTimedTasks();

			resizeWidget();
		}
	})
}

function initNumofRunningInstances(){
	if(setting.runningInstances.available){
		$(".switch-numof-running-instances").removeClass("switch-off");
		$(".switch-numof-running-instances").addClass("switch-on");
		$(".switch-numof-running-instances").prop("title","Click to turn off");
		$("#numof-running-instances-div").show();
	}else{
		$(".switch-numof-running-instances").removeClass("switch-on");
		$(".switch-numof-running-instances").addClass("switch-off");
		$(".switch-numof-running-instances").prop("title","Click to turn on");
		$("#numof-running-instances-div").hide();
	}	
	$(".switch-numof-running-instances").on("click",function(){
		if($(".switch-numof-running-instances").hasClass("switch-on")){
			$(".switch-numof-running-instances").removeClass("switch-on");
			$(".switch-numof-running-instances").addClass("switch-off");
			$(".switch-numof-running-instances").prop("title","Click to turn on");
			$("#numof-running-instances-div").hide();
			setting.runningInstances.available = false;
		}else if($(".switch-numof-running-instances").hasClass("switch-off")){
			$(".switch-numof-running-instances").removeClass("switch-off");
			$(".switch-numof-running-instances").addClass("switch-on");
			$(".switch-numof-running-instances").prop("title","Click to turn off");
			$("#numof-running-instances-div").show();
			setting.runningInstances.available = true;
		}
	});

	$("#numof-running-instances").val(setting.runningInstances.number);
	$("#numof-running-instances").on("blur",function(){
		setting.runningInstances.number = $("#numof-running-instances").val();
	});
}

function initTimedTasks(){
	if(setting.timedTasks.available){
		$(".switch-timed-task").removeClass("switch-off");
		$(".switch-timed-task").addClass("switch-on");
		$(".switch-timed-task").prop("title","Click to turn off");
		$("#timed-tasks-div").show();
		$(".new-timed-task").show();
	}else{
		$(".switch-timed-task").removeClass("switch-on");
		$(".switch-timed-task").addClass("switch-off");
		$(".switch-timed-task").prop("title","Click to turn on");
		$("#timed-tasks-div").hide();
		$(".new-timed-task").hide();
	}	
	$(".switch-timed-task").on("click",function(){
		if($(".switch-timed-task").hasClass("switch-on")){
			$(".switch-timed-task").removeClass("switch-on");
			$(".switch-timed-task").addClass("switch-off");
			$(".switch-timed-task").prop("title","Click to turn on");
			$("#timed-tasks-div").hide();
			$(".new-timed-task").hide();
			setting.timedTasks.available = false;
		}else if($(".switch-timed-task").hasClass("switch-off")){
			$(".switch-timed-task").removeClass("switch-off");
			$(".switch-timed-task").addClass("switch-on");
			$(".switch-timed-task").prop("title","Click to turn off");
			$("#timed-tasks-div").show();
			$(".new-timed-task").show();
			setting.timedTasks.available = true;
		}
	});

	showTimedTasks();
	$(".new-timed-task").on("click",function(){
		addTimedTask();
	});
}

function showTimedTasks(){
	$("#timed-tasks-div").empty();
    _.each(setting.timedTasks.tasks,function(task,index){
        var row = `<div class="timed-task-row" data-index="`+ index +`">
        				<div class="task-design-div col-md-10"></div>
                    	<div class="task-action-div">`;

        if(task.byDesigner){
        	row += `<div><span class="task-editor" title="Use Task Editor"></span></div>`;
        	row += `<div><span class="task-delete" title="Delete Task"></span></div>`;
        }else{
        	row += `<div><span class="task-designer" title="Use Task Designer"></span></div>`;
        	row += `<div><span class="task-delete" title="Delete Task"></span></div>`;
        }
                    	
        row += `</div></div>`
        $("#timed-tasks-div").append(row);
        showTask(task,index);
    });

    $(".task-delete").on('click',function(event){
        deleteTask(event);
        showTimedTasks();
    });

    $(".task-editor").on('click',function(event){
    	var index = $(event.currentTarget).parent().parent().parent().data("index");
    	setting.timedTasks.tasks[index].byDesigner = false;
    	setting.timedTasks.tasks[index].cronEntry = "";
    	showTimedTasks();
    });

    $(".task-designer").on('click',function(event){
    	var index = $(event.currentTarget).parent().parent().parent().data("index");
    	setting.timedTasks.tasks[index].byDesigner = true;
    	setting.timedTasks.tasks[index].cronEntry = "* * * * *";
    	showTimedTasks();
    });

    $(".cron-editor").on('blur',function(event){
    	var index = $(event.currentTarget).parent().parent().parent().parent().data("index");
		setting.timedTasks.tasks[index].cronEntry = $(event.currentTarget).val();
		$(event.currentTarget).parent().parent().find(".cron-val").text(setting.timedTasks.tasks[index].cronEntry);
    });
}

function addTimedTask(){
	var task = $.extend(true,{},metatask);
	setting.timedTasks.tasks.push(task);
	showTimedTasks();
}

function deleteTask(event){
	var index = $(event.currentTarget).parent().parent().parent().data("index");
	setting.timedTasks.tasks.splice(index,1);
}

function showTask(task,index){
	if(task.byDesigner){
		var cronInstance = $.extend(true,{},cron);
		cronInstance.initCronEntry($(".timed-task-row[data-index="+index+"]").find(".task-design-div"),task);
	}else{
		var editor = `<div class="row">
					    <div class="cron-designer">
					    	<input class="cron-editor" type="text" placeholder="Please enter your cron task expression" value="`+ task.cronEntry +`">
						</div>
					    <div class="cron-result">
					    	<span>Generated cron entry: </span>
					    	<span class="cron-val">`+ task.cronEntry +`</span>
					    </div>
					</div>`;
		$(".timed-task-row[data-index="+index+"]").find(".task-design-div").empty();
		$(".timed-task-row[data-index="+index+"]").find(".task-design-div").append(editor);	
	}
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

var metatask = {
	"byDesigner" : true,
	"cronEntry" : "* * * * *"
}