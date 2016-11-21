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

import { initDesigner } from "./initDesigner";
import { initPipeline } from "./initPipeline";
import { initAction } from "./initAction";
import * as pipelineDataService from "./pipelineData";
import { notify, confirm } from "../common/notify";
import { loading } from "../common/loading";
import { setLinePathAry, linePathAry, setCurrentSelectedItem } from "../common/constant";
import { pipelineCheck } from "../common/check";
import { initButton } from "./initButton";
import {getSequenceDetail} from "../history/main";


export let allPipelines;

export let pipelineData;
let pipelineDataOriginalCopy,linePathAryOriginalCopy;
let pipelineName, pipelineVersion, pipelineVersionID,pipelineHasHistory;
let pipelineEnvs;

let splitStartY;

export function initPipelinePage() {
    var promise = pipelineDataService.getAllPipelines();
    promise.done(function(data) {
        loading.hide();
        allPipelines = data.list;
        if (allPipelines.length > 0) {
            showPipelineList();
        } else {
            showNoPipeline();
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

function showPipelineList() {
    $.ajax({
        url: "../../templates/pipeline/pipelineList.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").html($(data));
            $("#pipelinelist").show("slow");

            $(".newpipeline").on('click', function() {
                showNewPipeline();
            })

            $(".pipelinelist_body").empty();
            _.each(allPipelines, function(item) {
                var pprow = `<tr class="pp-row">
                                <td class="pptd">
                                    <span class="glyphicon glyphicon-menu-down treeclose treecontroller" data-name=` 
                                    + item.name + `></span><span style="margin-left:10px">` 
                                    + item.name + `</span></td><td></td><td></td><td></td></tr>`;
                $(".pipelinelist_body").append(pprow);

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

                    $(".pipelinelist_body").append(vrow);
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
                pipelineName = target.parent().parent().data("pname");
                pipelineVersion = target.parent().parent().data("version");
                pipelineVersionID = target.parent().parent().data("versionid");
                pipelineHasHistory = target.parent().data("hashistory");
                getPipelineData();
            });
        }
    });
}

function getPipelineData() {
    setCurrentSelectedItem(null);
    var promise = pipelineDataService.getPipeline(pipelineName, pipelineVersionID);
    promise.done(function(data) {
        // pipelineDataOriginalCopy = _.map(data.stageList,function(item){
        //     return $.extend(true,{},item);
        // });
        // linePathAryOriginalCopy = _.map(data.lineList,function(item){
        //     return $.extend(true,{},item);
        // });
        loading.hide();
        pipelineData = data.stageList;
        setLinePathAry(data.lineList); 
        showPipelineDesigner(data.status);
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

function showNoPipeline() {
    $.ajax({
        url: "../../templates/pipeline/noPipeline.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").html($(data));
            $("#nopipeline").show("slow");
            $(".newpipeline").on('click', function() {
                showNewPipeline();
            })
        }
    });
}

function beforeShowNewPipeline() {
    if(_.isEqual(pipelineDataOriginalCopy,pipelineData) && _.isEqual(linePathAryOriginalCopy,linePathAry)){
        showNewPipeline();
    }else{
        var actions = [{
            "name": "save",
            "label": "Yes",
            "action": function() {
                savePipelineData(showNewPipeline);
            }
        }, {
            "name": "show",
            "label": "No",
            "action": function() {
                showNewPipeline();
            }
        }]
        confirm("The pipeline design has been modified, would you like to save the changes at first.", "info", actions);
    }
}

function showNewPipeline() {
    $.ajax({
        url: "../../templates/pipeline/newPipeline.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").children().hide();
            $("#main").append($(data));
            $("#newpipeline").show("slow");
            $("#newppBtn").on('click', function() {
                var promise = pipelineDataService.addPipeline();
                if (promise) {
                    promise.done(function(data) {
                        loading.hide();
                        notify(data.message, "success");
                        initPipelinePage();
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

function showPipelineDesigner(state) {
    $.ajax({
        url: "../../templates/pipeline/pipelineDesign.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").html($(data));
            $("#pipelinedesign").show("slow");

            $("#selected_pipeline").text(pipelineName + " / " + pipelineVersion);

            if(state){
                $(".pipeline-state").addClass("pipeline-on");
            }else{
                $(".pipeline-state").addClass("pipeline-off");
            }

            initDesigner();
            drawPipeline();

            $(".backtolist").on('click', function() {
                beforeBackToList();
            });

            $(".pipeline-state").on('click',function(event){
                if($(event.currentTarget).hasClass("pipeline-off")){
                    if(!pipelineCheck(pipelineData)){
                        notify("This pipeline does not pass the availability check, please make it available before ", "error");
                    }else{
                        beforeRunPipeline();
                    }
                }else if($(event.currentTarget).hasClass("pipeline-on")){
                    beforeStopPipeline();
                }
            });

            $(".checkpipeline").on('click', function() {
                pipelineCheck(pipelineData);
            });

            $(".savepipeline").on('click', function() {
                savePipelineData();
            });

            $(".newpipelineversion").on('click', function() {
                showNewPipelineVersion();
            });

            $(".loghistory").on('click',function(){
                if(pipelineHasHistory == "no"){
                    notify("The pipeline has never been run, there's no logs.", "info");
                }else if(pipelineHasHistory == "yes"){
                    beforeShowLog();
                }
            });

            $(".newpipelineindesigner").on('click', function() {
                beforeShowNewPipeline();
            });

            $(".envsetting").on("click", function(event) {
                showPipelineEnv();
            });

            $(".designer-split").on("dragstart",function(event){
                splitStartY = event.originalEvent.y;
            })

            $(".designer-split").on("dragend",function(event){
                var svgDiv = $("#div-d3-main-svg");
                svgDiv.height(svgDiv.height() + event.originalEvent.y - splitStartY);
            })

            pipelineDataOriginalCopy = _.map(pipelineData,function(item){
                return $.extend(true,{},item);
            });
            linePathAryOriginalCopy = _.map(linePathAry,function(item){
                return $.extend(true,{},item);
            });
            
        }
    });
}

function drawPipeline() {
    $("#pipeline-info-edit").empty();
    initPipeline();
    initButton();
}

export function savePipelineData(next) {
    var promise = pipelineDataService.savePipeline(pipelineName, pipelineVersion, pipelineVersionID, pipelineData, linePathAry);
    promise.done(function(data) {
        pipelineDataOriginalCopy = _.map(pipelineData,function(item){
            return $.extend(true,{},item);
        });
        linePathAryOriginalCopy = _.map(linePathAry,function(item){
            return $.extend(true,{},item);
        });
        loading.hide();
        if (!next) {
            notify(data.message, "success");
        } else {
            next();
        }
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        if (!next) {
            if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
                notify(xhr.responseJSON.errMsg, "error");
            } else if(xhr.statusText != "abort") {
                notify("Server is unreachable", "error");
            }
        } else {
            next();
        }
    });
}

function showNewPipelineVersion() {
    $.ajax({
        url: "../../templates/pipeline/newPipelineVersion.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").children().hide();
            $("#main").append($(data));
            $("#newpipelineversion").show("slow");

            $("#pp-name-newversion").val(pipelineName);

            $("#newppVersionBtn").on('click', function() {
                var promise = pipelineDataService.addPipelineVersion(pipelineName, pipelineVersionID, pipelineData, linePathAry);
                if (promise) {
                    promise.done(function(data) {
                        loading.hide();
                        notify(data.message, "success");
                        initPipelinePage();
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
    $("#nopipeline").hide();
    $("#newpipeline").hide();
    $("#newpipelineversion").show("slow");
}

function cancelNewPPPage() {
    $("#newpipeline").remove();
    $("#main").children().show("slow");
}

function cancelNewPPVersionPage() {
    $("#newpipelineversion").remove();
    $("#main").children().show("slow");
}


function showPipelineEnv() {
    if ($("#env-setting").hasClass("env-setting-closed")) {
        $("#env-setting").removeClass("env-setting-closed");
        $("#env-setting").addClass("env-setting-opened");
        $("#close_pp_env").removeClass("pipeline-open-env");
        $("#close_pp_env").addClass("pipeline-close-env");

        $.ajax({
            url: "../../templates/pipeline/envSetting.html",
            type: "GET",
            cache: false,
            success: function(data) {
                $("#env-setting").html($(data));

                $(".add-env").on('click', function() {
                    pipelineEnvs.push(["", ""]);
                    showEnvKVs();
                });

                $(".pipeline-close-env").on('click', function() {
                    hidePipelineEnv();
                });

                $(".save-env").on('click', function() {
                    savePipelineEnvs();
                });

                getEnvList();
            }
        });

    } else {
        $("#env-setting").removeClass("env-setting-opened");
        $("#env-setting").addClass("env-setting-closed");
        $("#close_pp_env").removeClass("pipeline-close-env");
        $("#close_pp_env").addClass("pipeline-open-env");
    }
}

function hidePipelineEnv() {
    $("#env-setting").removeClass("env-setting-opened");
    $("#env-setting").addClass("env-setting-closed");
    $("#close_pp_env").removeClass("pipeline-close-env");
    $("#close_pp_env").addClass("pipeline-open-env");
}

function getEnvList() {
    var promise = pipelineDataService.getEnvs(pipelineName, pipelineVersionID);
    promise.done(function(data) {
        loading.hide();
        pipelineEnvs = _.pairs(data.env);
        showEnvKVs();
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

function showEnvKVs() {
    $("#envs").empty();
    _.each(pipelineEnvs,function(item,index){
         var row = '<div class="env-row"><div class="env-key-div">'
                        +'<div>'
                            +'<label for="normal-field" class="col-sm-3 control-label" style="margin-top:5px">'
                                +'KEY'
                            +'</label>'
                            +'<div class="col-sm-9" data-index="' + index + '">'
                                +'<input type="text" value="' + item[0] + '" class="form-control pp-env-input pp-env-key" required>'
                            +'</div>'
                        +'</div>'
                    +'</div>'
                    +'<div class="env-value-div" style="margin-left:0">'
                        +'<div>'
                            +'<label for="normal-field" class="col-sm-3 control-label" style="margin-top:5px">'
                                +'VALUE'
                            +'</label>'
                            +'<div class="col-sm-9" data-index="' + index + '">' 
                                +'<input type="text" class="form-control pp-env-input pp-env-value" required>'
                            +'</div>'
                        +'</div>'
                    +'</div>'
                    +'<div class="env-remove-div pp-rm-kv" data-index="' + index + '">'
                        +'<span class="glyphicon glyphicon-remove"></span>'
                    +'</div></div>';
        $("#envs").append(row);
        $("#envs").find("div[data-index="+index+"]").find(".pp-env-value").val(item[1]);
    });

    $(".pp-env-key").on('input',function(event){
        var key = $(event.currentTarget).val();
        $(event.currentTarget).val(key.toUpperCase());
    });

    $(".pp-env-key").on('blur',function(event){
        var index = $(event.currentTarget).parent().data("index");
        pipelineEnvs[index][0] = $(event.currentTarget).val();
    });

    $(".pp-env-value").on('blur',function(event){
        var index = $(event.currentTarget).parent().data("index");
        pipelineEnvs[index][1] = $(event.currentTarget).val();
    });

    $(".pp-rm-kv").on('click',function(event){
        var index = $(event.currentTarget).data("index");
        pipelineEnvs.splice(index, 1);
        showEnvKVs();
    }); 
}

function savePipelineEnvs() {
    var promise = pipelineDataService.setEnvs(pipelineName, pipelineVersionID, pipelineEnvs);
    if (promise) {
        promise.done(function(data) {
            loading.hide();
            notify(data.message, "success");
            hidePipelineEnv();
        });
        promise.fail(function(xhr, status, error) {
            loading.hide();
            if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
                notify(xhr.responseJSON.errMsg, "error");
            } else if(xhr.statusText != "abort") {
                notify("Server is unreachable", "error");
            }
            hidePipelineEnv();
        });
    }
}

// run pipeline
function beforeRunPipeline() {
    if(_.isEqual(pipelineDataOriginalCopy,pipelineData) && _.isEqual(linePathAryOriginalCopy,linePathAry)){
        runPipeline();
    }else{
        var actions = [{
            "name": "saveAndRun",
            "label": "Yes, save it first.",
            "action": function() {
                savePipelineData(runPipeline);
            }
        }, {
            "name": "run",
            "label": "No, just run it.",
            "action": function() {
                runPipeline();
            }
        }]
        confirm("The pipeline design has been modified, would you like to save the changes before run it.", "info", actions);
    }
}

function runPipeline() {
    var promise = pipelineDataService.changeState(pipelineName, pipelineVersionID, 1);
    promise.done(function(data) {
        loading.hide();
        notify(data.message, "success");
        $(".pipeline-state").removeClass("pipeline-off").addClass("pipeline-on");
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

//stop pipeline
function beforeStopPipeline() {
    if(_.isEqual(pipelineDataOriginalCopy,pipelineData) && _.isEqual(linePathAryOriginalCopy,linePathAry)){
        stopPipeline();
    }else{
        var actions = [{
            "name": "saveAndStop",
            "label": "Yes, save it first.",
            "action": function() {
                savePipelineData(stopPipeline);
            }
        }, {
            "name": "stop",
            "label": "No, just stop it.",
            "action": function() {
                stopPipeline();
            }
        }]
        confirm("The pipeline design has been modified, would you like to save the changes before stop it.", "info", actions);
    }
}

function stopPipeline() {
    var promise = pipelineDataService.changeState(pipelineName, pipelineVersionID, 0);
    promise.done(function(data) {
        loading.hide();
        notify(data.message, "success");
        $(".pipeline-state").removeClass("pipeline-on").addClass("pipeline-off");
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
    if(_.isEqual(pipelineDataOriginalCopy,pipelineData) && _.isEqual(linePathAryOriginalCopy,linePathAry)){
        initPipelinePage();
    }else{
        var actions = [{
            "name": "save",
            "label": "Yes",
            "action": function() {
                savePipelineData(initPipelinePage);
            }
        }, {
            "name": "back",
            "label": "No",
            "action": function() {
                initPipelinePage();
            }
        }]
        confirm("The pipeline design has been modified, would you like to save the changes before go back to list.", "info", actions);
    }
}

export function getPipelineToken(){
    return pipelineDataService.getToken(pipelineName,pipelineVersionID);
}

function beforeShowLog() {
    if(_.isEqual(pipelineDataOriginalCopy,pipelineData) && _.isEqual(linePathAryOriginalCopy,linePathAry)){
        showLogHistory();
    }else{
        var actions = [{
            "name": "saveAndLog",
            "label": "Yes",
            "action": function() {
                savePipelineData(showLogHistory);
            }
        }, {
            "name": "log",
            "label": "No",
            "action": function() {
                showLogHistory();
            }
        }]
        confirm("The pipeline design has been modified, would you like to save the changes before show the log history.", "info", actions);
    }
}

function showLogHistory(){
    $(".menu-history").parent().addClass("active");
    $(".menu-pipeline").parent().removeClass("active");
                    
    var pipelineInfo = {
        "pipelineName" : pipelineName,
        "pipelineVersionID" : pipelineVersionID,
        "pipelineVersion" : pipelineVersion,
        "sequenceID" : ""
    };
    getSequenceDetail(pipelineInfo);
}
// $("#pipeline-select").on('change',function(){
//     showVersionList();
// })
// $("#version-select").on('change',function(){
//     showPipeline();
// })

// function showPipelineList(){
//     $("#pipeline-select").empty();
//     d3.select("#pipeline-select")
//         .selectAll("option")
//         .data(allPipelines)
//         .enter()
//         .append("option")
//         .attr("value",function(d,i){
//             return d.name;
//         })
//         .text(function(d,i){
//             return d.name;
//         }); 
//      $("#pipeline-select").select2({
//        minimumResultsForSearch: Infinity
//      });   
//     showVersionList();
// }

// function showVersionList(){
//     var pipeline = $("#pipeline-select").val();
//     var versions = _.find(allPipelines,function(item){
//         return item.name == pipeline;
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

//     showPipeline(); 
// }
