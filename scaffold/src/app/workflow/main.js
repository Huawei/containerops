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

import { initDesigner } from "./initDesigner";
import { initWorkflow } from "./initWorkflow";
import { initAction } from "./initAction";
import * as workflowDataService from "./workflowData";
import { notify, confirm } from "../common/notify";
import { loading } from "../common/loading";
import { setLinePathAry, linePathAry, setCurrentSelectedItem } from "../common/constant";
import { workflowCheck } from "../common/check";
import { initButton } from "./initButton";
import {initHistoryPage} from "../history/main";
import {initWorkflowEnv,showWorkflowEnv} from "./workflowEnv";
import {initWorkflowVar,showWorkflowVar} from "./workflowVar";
import {initWorkflowSetting} from "./workflowSetting";


export let allWorkflows;

export let workflowData,workflowSettingData;
let workflowDataOriginalCopy,linePathAryOriginalCopy,workflowSettingDataOriginalCopy;
let workflowName, workflowVersion, workflowVersionID,workflowHasHistory;

let splitStartY;

export function initWorkflowPage() {
    var promise = workflowDataService.getAllWorkflows();
    promise.done(function(data) {
        loading.hide();
        allWorkflows = data.list;
        if (allWorkflows.length > 0) {
            showWorkflowList();
        } else {
            showNoWorkflow();
        }
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else if(xhr.statusText != "abort") {
            notify("Server is unreachable", "error");
        }
    });
}

function showWorkflowList() {
    $.ajax({
        url: "../../templates/workflow/workflowList.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").html($(data));
            $("#workflowlist").show("slow");

            $(".newworkflow").on('click', function() {
                showNewWorkflow();
            })

            $(".workflowlist_body").empty();
            _.each(allWorkflows, function(item) {
                var pprow = `<tr class="pp-row">
                                <td class="pptd">
                                    <span class="glyphicon glyphicon-menu-down treeclose treecontroller" data-name=` 
                                    + item.name + `></span><span style="margin-left:10px">` 
                                    + item.name + `</span></td><td></td><td></td><td></td></tr>`;
                $(".workflowlist_body").append(pprow);

                _.each(item.version, function(version) {
                    var vrow = `<tr data-pname=` + item.name + ` data-version=` + version.version + ` data-versionid=`
                                + version.id + ` class="ppversion-row">
                                    <td></td>
                                    <td class="pptd">` + version.version + `</td><td>`;

                    if(_.isUndefined(version.status)){
                        vrow += `<div class="state-list">
                                    <div class="state-icon-list state-norun"></div>
                                </div></td><td data-hashistory="no">`;
                    }else if(version.status.status){
                        vrow += `<div class="state-list">
                                    <div class="state-icon-list state-success"></div>
                                    <span class="state-label-list">` + version.status.time + `</span>
                                </div></td><td data-hashistory="yes">`;
                    }else{
                        vrow += `<div class="state-list">
                                    <div class="state-icon-list state-fail"></div>
                                    <span class="state-label-list">` + version.status.time + `</span>
                                </div></td><td data-hashistory="yes">`
                    }

                    vrow += `<button type="button" class="btn btn-success ppview">
                                    <i class="glyphicon glyphicon-eye-open" style="font-size:16px"></i>
                                    <span style="margin-left:5px">View</span>
                                </button>
                            </td></tr>`;

                    $(".workflowlist_body").append(vrow);
                })
            });
            
            $(".treecontroller").on("click",function(event){
                var target = $(event.currentTarget);
                if(target.hasClass("treeclose")){
                    target.removeClass("glyphicon-menu-down treeclose");
                    target.addClass("glyphicon-menu-right treeopen");

                    var name = target.data("name");
                    $('*[data-pname="'+name+'"]').hide();
                }else{
                    target.addClass("glyphicon-menu-down treeclose");
                    target.removeClass("glyphicon-menu-right treeopen");

                    var name = target.data("name");
                    $('*[data-pname="'+name+'"]').show();
                }  
            });

            $(".ppview").on("click", function(event) {
                var target = $(event.currentTarget);
                workflowName = target.parent().parent().data("pname");
                workflowVersion = target.parent().parent().data("version");
                workflowVersionID = target.parent().parent().data("versionid");
                workflowHasHistory = target.parent().data("hashistory");
                getWorkflowData();
            });
        }
    });
}

function getWorkflowData() {
    setCurrentSelectedItem(null);
    var promise = workflowDataService.getWorkflow(workflowName, workflowVersionID);
    promise.done(function(data) {
        loading.hide();
        workflowData = data.stageList;
        setLinePathAry(data.lineList); 
        workflowSettingData = data.setting;
        showWorkflowDesigner(data.status);
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else if(xhr.statusText != "abort") {
            notify("Server is unreachable", "error");
        }
    });
}

function showNoWorkflow() {
    $.ajax({
        url: "../../templates/workflow/noWorkflow.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").html($(data));
            $("#noworkflow").show("slow");
            $(".newworkflow").on('click', function() {
                showNewWorkflow();
            })
        }
    });
}

function beforeShowNewWorkflow() {
    if(isWorkflowChanged()){
        showNewWorkflow();
    }else{
        var actions = [{
            "name": "save",
            "label": "Yes",
            "action": function() {
                saveWorkflowData(showNewWorkflow);
            }
        }, {
            "name": "show",
            "label": "No",
            "action": function() {
                showNewWorkflow();
            }
        }]
        confirm("The workflow design has been modified, would you like to save the changes at first.", "info", actions);
    }
}

function showNewWorkflow() {
    $.ajax({
        url: "../../templates/workflow/newWorkflow.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").children().hide();
            $("#main").append($(data));
            $("#newworkflow").show("slow");
            $("#newppBtn").on('click', function() {
                var promise = workflowDataService.addWorkflow();
                if (promise) {
                    promise.done(function(data) {
                        loading.hide();
                        notify(data.message, "success");
                        initWorkflowPage();
                    });
                    promise.fail(function(xhr, status, error) {
                        loading.hide();
                        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
                            notify(xhr.responseJSON.errMsg, "error");
                        } else if(xhr.statusText != "abort") {
                            notify("Server is unreachable", "error");
                        }
                    });
                }
            })
            $("#cancelNewppBtn").on('click', function() {
                cancelNewPPPage();
            })
        }
    });
}

function showWorkflowDesigner(state) {
    $.ajax({
        url: "../../templates/workflow/workflowDesign.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").html($(data));
            $("#workflowdesign").show("slow");

            $("#selected_workflow").text(workflowName + " / " + workflowVersion);
            $("#selected_workflow").prop("title",workflowName + " / " + workflowVersion);

            if(state){
                $(".workflow-state").addClass("workflow-on");
            }else{
                $(".workflow-state").addClass("workflow-off");
            }

            initDesigner();
            drawWorkflow();

            initWorkflowEnv(workflowName,workflowVersionID);
            initWorkflowVar(workflowName,workflowVersionID);

            $(".backtolist").on('click', function() {
                beforeBackToList();
            });

            $(".workflow-state").on('click',function(event){
                if($(event.currentTarget).hasClass("workflow-off")){
                    if(!workflowCheck(workflowData,workflowSettingData)){
                        notify("This workflow does not pass the availability check, please make it available before ", "error");
                    }else{
                        beforeRunWorkflow();
                    }
                }else if($(event.currentTarget).hasClass("workflow-on")){
                    beforeStopWorkflow();
                }
            });

            $(".checkworkflow").on('click', function() {
                workflowCheck(workflowData,workflowSettingData);
            });

            $(".saveworkflow").on('click', function() {
                saveWorkflowData(getWorkflowData);
            });

            $(".newworkflowversion").on('click', function() {
                showNewWorkflowVersion();
            });

            $(".loghistory").on('click',function(){
                if(workflowHasHistory == "no"){
                    notify("The workflow has never been run, there's no logs.", "info");
                }else if(workflowHasHistory == "yes"){
                    beforeShowLog();
                }
            });

            $(".newworkflowindesigner").on('click', function() {
                beforeShowNewWorkflow();
            });

            $(".envsetting").on("click", function(event) {
                showWorkflowEnv();
            });

            $(".varsetting").on("click", function(event) {
                showWorkflowVar();
            });

            $(".designer-split").on("dragstart",function(event){
                splitStartY = event.originalEvent.y;
            })

            $(".designer-split").on("dragend",function(event){
                var svgDiv = $("#div-d3-main-svg");
                svgDiv.height(svgDiv.height() + event.originalEvent.y - splitStartY);
            })

            workflowDataOriginalCopy = _.map(workflowData,function(item){
                return $.extend(true,{},item);
            });
            linePathAryOriginalCopy = _.map(linePathAry,function(item){
                return $.extend(true,{},item);
            });
            workflowSettingDataOriginalCopy = $.extend(true,{},workflowSettingData);

            initWorkflowSetting(workflowSettingData);
        }
    });
}

function drawWorkflow() {
    $("#workflow-info-edit").empty();
    initWorkflow();
    initButton();
}

export function saveWorkflowData(next) {
    var promise = workflowDataService.saveWorkflow(workflowName, workflowVersion, workflowVersionID, workflowData, linePathAry, workflowSettingData);
    promise.done(function(data) {
        workflowDataOriginalCopy = _.map(workflowData,function(item){
            return $.extend(true,{},item);
        });
        linePathAryOriginalCopy = _.map(linePathAry,function(item){
            return $.extend(true,{},item);
        });
        workflowSettingDataOriginalCopy = $.extend(true,{},workflowSettingData);

        loading.hide();
        if (!next) {
            notify(data.message, "success");
        } else {
            next();
        }
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        // if (!next) {
        //     if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
        //         notify(xhr.responseJSON.errMsg, "error");
        //     } else if(xhr.statusText != "abort") {
        //         notify("Server is unreachable", "error");
        //     }
        // } else {
        //     next();
        // }
        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else if(xhr.statusText != "abort") {
            notify("Server is unreachable", "error");
        }
    });
}

function showNewWorkflowVersion() {
    $.ajax({
        url: "../../templates/workflow/newWorkflowVersion.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").children().hide();
            $("#main").append($(data));
            $("#newworkflowversion").show("slow");

            $("#pp-name-newversion").val(workflowName);

            $("#newppVersionBtn").on('click', function() {
                var promise = workflowDataService.addWorkflowVersion(workflowName, workflowVersionID, workflowData, linePathAry, workflowSettingData);
                if (promise) {
                    promise.done(function(data) {
                        loading.hide();
                        notify(data.message, "success");
                        initWorkflowPage();
                    });
                    promise.fail(function(xhr, status, error) {
                        loading.hide();
                        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
                            notify(xhr.responseJSON.errMsg, "error");
                        } else if(xhr.statusText != "abort") {
                            notify("Server is unreachable", "error");
                        }
                    });
                }
            })
            $("#cancelNewppVersionBtn").on('click', function() {
                cancelNewPPVersionPage();
            })
        }
    });

    $("#content").hide();
    $("#noworkflow").hide();
    $("#newworkflow").hide();
    $("#newworkflowversion").show("slow");
}

function cancelNewPPPage() {
    $("#newworkflow").remove();
    $("#main").children().show("slow");
}

function cancelNewPPVersionPage() {
    $("#newworkflowversion").remove();
    $("#main").children().show("slow");
}

// run workflow
function beforeRunWorkflow() {
    if(isWorkflowChanged()){
        runWorkflow();
    }else{
        var actions = [{
            "name": "saveAndRun",
            "label": "Yes, save it first.",
            "action": function() {
                saveWorkflowData(runWorkflow);
            }
        }, {
            "name": "run",
            "label": "No, just run it.",
            "action": function() {
                runWorkflow();
            }
        }]
        confirm("The workflow design has been modified, would you like to save the changes before run it.", "info", actions);
    }
}

function runWorkflow() {
    var promise = workflowDataService.changeState(workflowName, workflowVersionID, 1);
    promise.done(function(data) {
        loading.hide();
        notify(data.message, "success");
        $(".workflow-state").removeClass("workflow-off").addClass("workflow-on");
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else if(xhr.statusText != "abort") {
            notify("Server is unreachable", "error");
        }
    });
}

//stop workflow
function beforeStopWorkflow() {
    if(isWorkflowChanged()){
        stopWorkflow();
    }else{
        var actions = [{
            "name": "saveAndStop",
            "label": "Yes, save it first.",
            "action": function() {
                saveWorkflowData(stopWorkflow);
            }
        }, {
            "name": "stop",
            "label": "No, just stop it.",
            "action": function() {
                stopWorkflow();
            }
        }]
        confirm("The workflow design has been modified, would you like to save the changes before stop it.", "info", actions);
    }
}

function stopWorkflow() {
    var promise = workflowDataService.changeState(workflowName, workflowVersionID, 0);
    promise.done(function(data) {
        loading.hide();
        notify(data.message, "success");
        $(".workflow-state").removeClass("workflow-on").addClass("workflow-off");
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else if(xhr.statusText != "abort") {
            notify("Server is unreachable", "error");
        }
    });
}

function beforeBackToList() {
    if(isWorkflowChanged()){
        initWorkflowPage();
    }else{
        var actions = [{
            "name": "save",
            "label": "Yes",
            "action": function() {
                saveWorkflowData(initWorkflowPage);
            }
        }, {
            "name": "back",
            "label": "No",
            "action": function() {
                initWorkflowPage();
            }
        }]
        confirm("The workflow design has been modified, would you like to save the changes before go back to list.", "info", actions);
    }
}

export function getWorkflowToken(){
    return workflowDataService.getToken(workflowName,workflowVersionID);
}

function beforeShowLog() {
    if(isWorkflowChanged()){
        showLogHistory();
    }else{
        var actions = [{
            "name": "saveAndLog",
            "label": "Yes",
            "action": function() {
                saveWorkflowData(showLogHistory);
            }
        }, {
            "name": "log",
            "label": "No",
            "action": function() {
                showLogHistory();
            }
        }]
        confirm("The workflow design has been modified, would you like to save the changes before show the log history.", "info", actions);
    }
}

function showLogHistory(){
    $(".menu-history").parent().addClass("active");
    $(".menu-workflow").parent().removeClass("active");
                    
    initHistoryPage(workflowName);
}

function isWorkflowChanged(){
    return _.isEqual(workflowDataOriginalCopy,workflowData) && _.isEqual(linePathAryOriginalCopy,linePathAry) &&  _.isEqual(workflowSettingDataOriginalCopy,workflowSettingData);
}

// $("#workflow-select").on('change',function(){
//     showVersionList();
// })
// $("#version-select").on('change',function(){
//     showWorkflow();
// })

// function showWorkflowList(){
//     $("#workflow-select").empty();
//     d3.select("#workflow-select")
//         .selectAll("option")
//         .data(allWorkflows)
//         .enter()
//         .append("option")
//         .attr("value",function(d,i){
//             return d.name;
//         })
//         .text(function(d,i){
//             return d.name;
//         }); 
//      $("#workflow-select").select2({
//        minimumResultsForSearch: Infinity
//      });   
//     showVersionList();
// }

// function showVersionList(){
//     var workflow = $("#workflow-select").val();
//     var versions = _.find(allWorkflows,function(item){
//         return item.name == workflow;
//     }).versions;

//     $("#version-select").empty();
//     d3.select("#version-select")
//         .selectAll("option")
//         .data(versions)
//         .enter()
//         .append("option")
//         .attr("value",function(d,i){
//             return d.version;
//         })
//         .text(function(d,i){
//             return d.version;
//         }); 
//     $("#version-select").select2({
//        minimumResultsForSearch: Infinity
//      });

//     versions_shown = versions;

//     showWorkflow(); 
// }
