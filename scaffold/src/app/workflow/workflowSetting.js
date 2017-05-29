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
	var x = window.scrollX;
	var y = window.scrollY;

	$("#timed-tasks-div").empty();
    _.each(setting.timedTasks.tasks,function(task,index){
        var row = `<div class="timed-task-row" data-index="`+ index +`">`;

        if(task.collapse){
        	row +=	`<div class="task-design-div" style="display:none">`;
        }else{
        	row +=	`<div class="task-design-div">`;
        }
        
        row += `<div class="task-header">Expression</div>
	        	<div class="task-cron-div"></div>
	        	<div class="task-header">Start Json</div>
	        	<div class="task-input-div"></div>
	        	</div>`;

	    if(task.collapse){
        	row += `<div class="task-simplify-div">Task : `  + task.cronEntry 
								+ `  |  Event Type: ` + task.eventType + ` |  Event Name: ` 
								+ task.eventName + `</div>`;
        }else{
        	row += `<div class="task-simplify-div" style="display:none"></div>`;
        }

        row += `<div class="task-action-div">`;

        if(task.byDesigner){
        	row += `<span class="fa fa-edit action task-editor" title="Use Task Editor"></span>`;
        	row += `<span class="glyphicon glyphicon-remove action task-delete" title="Delete Task"></span>`;
        }else{
        	row += `<span class="fa fa-list-ul action task-designer" title="Use Task Designer"></span>`;
        	row += `<span class="glyphicon glyphicon-remove action task-delete" title="Delete Task"></span>`;
        }

        if(task.collapse){
        	row += `<span class="glyphicon glyphicon-chevron-up action task-collapse"></span>`;
        }else{
        	row += `<span class="glyphicon glyphicon-chevron-down action task-collapse"></span>`;
        }
                    	
        row += `</div></div>`
        $("#timed-tasks-div").append(row);
        showTask(task,index);
    });
	
	$(".task-collapse").on("click",function(event){
		var target = $(event.currentTarget);
		var index = target.parent().parent().data("index");
		if(target.hasClass("glyphicon-chevron-down")){
			target.removeClass("glyphicon-chevron-down").addClass("glyphicon-chevron-up");
			target.parent().parent().find(".task-design-div").hide();
			target.parent().parent().find(".task-simplify-div").show();
			var simpiliedTask = "Task : " + setting.timedTasks.tasks[index].cronEntry 
								+ "  |  Event Type: " + setting.timedTasks.tasks[index].eventType + " |  Event Name: " 
								+ setting.timedTasks.tasks[index].eventName; 
			target.parent().parent().find(".task-simplify-div").text(simpiliedTask);

			setting.timedTasks.tasks[index].collapse = true;
		}else if(target.hasClass("glyphicon-chevron-up")){
			target.removeClass("glyphicon-chevron-up").addClass("glyphicon-chevron-down");
			target.parent().parent().find(".task-design-div").show();
			target.parent().parent().find(".task-simplify-div").hide();
			setting.timedTasks.tasks[index].collapse = false;
		}
	})

    $(".task-delete").on('click',function(event){
        deleteTask(event);
        showTimedTasks();
    });

    $(".task-editor").on('click',function(event){
    	var index = $(event.currentTarget).parent().parent().data("index");
    	setting.timedTasks.tasks[index].byDesigner = false;
    	setting.timedTasks.tasks[index].cronEntry = "";
    	setting.timedTasks.tasks[index].collapse = false;
    	showTimedTasks();
    });

    $(".task-designer").on('click',function(event){
    	var index = $(event.currentTarget).parent().parent().data("index");
    	setting.timedTasks.tasks[index].byDesigner = true;
    	setting.timedTasks.tasks[index].cronEntry = "* * * * *";
    	setting.timedTasks.tasks[index].collapse = false;
    	showTimedTasks();
    });

    $(".cron-editor").on('blur',function(event){
    	var index = $(event.currentTarget).parent().parent().parent().parent().parent().data("index");
		setting.timedTasks.tasks[index].cronEntry = $(event.currentTarget).val();
		$(event.currentTarget).parent().parent().find(".cron-val").text(setting.timedTasks.tasks[index].cronEntry);
    });

    window.scrollTo(x,y);
}

function addTimedTask(){
	var task = $.extend(true,{},metatask);
	setting.timedTasks.tasks.push(task);
	showTimedTasks();
}

function deleteTask(event){
	var index = $(event.currentTarget).parent().parent().data("index");
	setting.timedTasks.tasks.splice(index,1);
}

function showTask(task,index){
	if(task.byDesigner){
		var cronInstance = $.extend(true,{},cron);
		cronInstance.initCronEntry($(".timed-task-row[data-index="+index+"]").find(".task-cron-div"),task);
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
		$(".timed-task-row[data-index="+index+"]").find(".task-cron-div").empty();
		$(".timed-task-row[data-index="+index+"]").find(".task-cron-div").append(editor);	
	}

	showInputJson(task,index);
}

function showInputJson(task,index){
	var inputJsonDom = `<div class="row">
							<div class="form-group col-md-6">
								<div for="hint-field" class="col-sm-4 control-label">
									Event Type
								</div>
								<div class="col-sm-7 input-group">
									<input type="text" class="form-control eventType"> 
								</div>
							</div>
							<div class="form-group col-md-6">
								<div for="normal-field" class="col-sm-4 control-label">
									Event Name
								</div>
								<div class="col-sm-7">
									<input type="text" class="form-control eventName">
								</div>
							</div>
						</div>
						<div class="row">
							<div class="form-group col-md-12">
								<div for="normal-field" class="col-sm-2 control-label">
									Start Json
								</div>
								<div class="col-sm-10">
									<div class="cron-input-json"></div>
								</div>
							</div>
						</div>`;

	$(".timed-task-row[data-index="+index+"]").find(".task-input-div").append(inputJsonDom);

	var codeOptions = {
        "mode": "code",
        "indentation": 2,
        "onChange" : function(){
        	try{
        		task.startJson = codeEditor.get();
        	}catch(e){
        		console.log("Start json of timed task " + (index+1) + " is invalid.")
        	}
        }
    };

    var parentDom = $(".timed-task-row[data-index="+index+"]").find(".task-input-div");
	var codeContainer = parentDom.find(".cron-input-json")[0];
	var codeEditor = new JSONEditor(codeContainer, codeOptions);
	codeEditor.set(task.startJson);

	parentDom.find(".eventType").val(task.eventType);
	parentDom.find(".eventType").on("blur",function(){
		task.eventType = parentDom.find(".eventType").val();
	})

	parentDom.find(".eventName").val(task.eventName);
	parentDom.find(".eventName").on("blur",function(){
		task.eventName = parentDom.find(".eventName").val();
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

var metatask = {
	"collapse" : false,
	"byDesigner" : true,
	"cronEntry" : "* * * * *",
	"startJson" : {},
	"eventType" : "",
	"eventName" : ""
}