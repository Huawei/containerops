/* Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

import { initDesigner } from "./initDesigner";
import { initPipeline } from "./initPipeline";
import { initAction } from "./initAction";
import { getAllPipelines, getPipeline, addPipeline, savePipeline, addPipelineVersion, getEnvs, setEnvs, changeState } from "./pipelineData";
import { notify, confirm } from "../common/notify";
import { loading } from "../common/loading";
import { setLinePathAry, linePathAry } from "../common/constant";
import { pipelineCheck } from "../common/check";
import { initButton } from "./initButton";


export let allPipelines;

export let pipelineData;
let pipelineName, pipelineVersion, pipelineVersionID;
let pipelineEnvs;

export function initPipelinePage() {
    loading.show();
    var promise = getAllPipelines();
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
        if (xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else {
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
                                    + item.name + `></span>&nbsp;&nbsp;&nbsp;&nbsp;` 
                                    + item.name + `</td><td></td><td></td></tr>`;
                $(".pipelinelist_body").append(pprow);

                _.each(item.version, function(version) {
                    var vrow = `<tr data-pname=` + item.name + ` data-version=` + version.version + ` data-versionid=`
                                + version.id + ` class="ppversion-row">
                                    <td></td>
                                    <td class="pptd">` + version.version + `</td>
                                    <td>
                                        <button type="button" class="btn btn-success ppview">
                                            <i class="glyphicon glyphicon-eye-open" style="font-size:16px"></i>&nbsp;&nbsp;View
                                        </button>
                                    </td>
                                </tr>`;
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
                getPipelineData();
            });
        }
    });
}

function getPipelineData() {
    loading.show();
    var promise = getPipeline(pipelineName, pipelineVersionID);
    promise.done(function(data) {
        loading.hide();
        pipelineData = data.stageList;
        setLinePathAry(data.lineList);
        showPipelineDesigner();
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        if (xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else {
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
    confirm("The pipeline design may be modified, would you like to save the pipeline at first.", "info", actions);
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
                var promise = addPipeline();
                if (promise) {
                    loading.show();
                    promise.done(function(data) {
                        loading.hide();
                        notify(data.message, "success");
                        initPipelinePage();
                    });
                    promise.fail(function(xhr, status, error) {
                        loading.hide();
                        if (xhr.responseJSON.errMsg) {
                            notify(xhr.responseJSON.errMsg, "error");
                        } else {
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

function showPipelineDesigner() {
    $.ajax({
        url: "../../templates/pipeline/pipelineDesign.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").html($(data));
            $("#pipelinedesign").show("slow");

            $("#selected_pipeline").text(pipelineName + " / " + pipelineVersion);

            initDesigner();
            drawPipeline();

            $(".backtolist").on('click', function() {
                beforeBackToList();
            });

            $(".runpipeline").on('click', function() {
                beforeRunPipeline();
            });

            $(".stoppipeline").on('click', function() {
                beforeStopPipeline();
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

            $(".newpipelineindesigner").on('click', function() {
                beforeShowNewPipeline();
            });

            $(".envsetting").on("click", function(event) {
                showPipelineEnv();
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
    loading.show();
    var promise = savePipeline(pipelineName, pipelineVersion, pipelineVersionID, pipelineData, linePathAry);
    promise.done(function(data) {
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
            if (xhr.responseJSON.errMsg) {
                notify(xhr.responseJSON.errMsg, "error");
            } else {
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
                var promise = addPipelineVersion(pipelineName, pipelineVersionID, pipelineData, linePathAry);
                if (promise) {
                    loading.show();
                    promise.done(function(data) {
                        loading.hide();
                        notify(data.message, "success");
                        initPipelinePage();
                    });
                    promise.fail(function(xhr, status, error) {
                        loading.hide();
                        if (xhr.responseJSON.errMsg) {
                            notify(xhr.responseJSON.errMsg, "error");
                        } else {
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

        $.ajax({
            url: "../../templates/pipeline/envSetting.html",
            type: "GET",
            cache: false,
            success: function(data) {
                $("#env-setting").html($(data));

                $(".new-kv").on('click', function() {
                    pipelineEnvs.push(["", ""]);
                    showEnvKVs();
                });

                $(".close-env").on('click', function() {
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
    }
}

function hidePipelineEnv() {
    $("#env-setting").removeClass("env-setting-opened");
    $("#env-setting").addClass("env-setting-closed");
}

function getEnvList() {
    loading.show();
    var promise = getEnvs(pipelineName, pipelineVersionID);
    promise.done(function(data) {
        loading.hide();
        pipelineEnvs = _.pairs(data.env);
        showEnvKVs();
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        if (xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else {
            notify("Server is unreachable", "error");
        }
    });
}

function showEnvKVs() {
    $("#envs").empty();
    _.each(pipelineEnvs, function(item, index) {
        var row = '<tr data-index="' + index + '"><td>' + '<input type="text" class="form-control col-md-5 env-key" value="' + item[0] + '" required>' + '</td><td>' + '<input type="text" class="form-control col-md-5 env-value" required value=' + item[1] + '>' + '</td><td>' + '<span class="glyphicon glyphicon-minus rm-kv"></span>' + '</td></tr>';
        $("#envs").append(row);
    });

    $(".env-key").on('input', function(event) {
        var key = $(event.currentTarget).val();
        $(event.currentTarget).val(key.toUpperCase());
    });

    $(".env-key").on('blur', function(event) {
        var index = $(event.currentTarget).parent().parent().data("index");
        pipelineEnvs[index][0] = $(event.currentTarget).val();
    });

    $(".env-value").on('blur', function(event) {
        var index = $(event.currentTarget).parent().parent().data("index");
        pipelineEnvs[index][1] = $(event.currentTarget).val();
    });

    $(".rm-kv").on('click', function(event) {
        var index = $(event.currentTarget).parent().parent().data("index");
        pipelineEnvs.splice(index, 1);
        showEnvKVs();
    });
}

function savePipelineEnvs() {
    var promise = setEnvs(pipelineName, pipelineVersionID, pipelineEnvs);
    if (promise) {
        loading.show();
        promise.done(function(data) {
            loading.hide();
            notify(data.message, "success");
            hidePipelineEnv();
        });
        promise.fail(function(xhr, status, error) {
            loading.hide();
            if (xhr.responseJSON.errMsg) {
                notify(xhr.responseJSON.errMsg, "error");
            } else {
                notify("Server is unreachable", "error");
            }
            hidePipelineEnv();
        });
    }
}

// run pipeline
function beforeRunPipeline() {
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
    confirm("The pipeline design may be modified, would you like to save the pipeline before run it.", "info", actions);
}

function runPipeline() {
    loading.show();
    var promise = changeState(pipelineName, pipelineVersionID, 1);
    promise.done(function(data) {
        loading.hide();
        notify(data.message, "success");
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        if (xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else {
            notify("Server is unreachable", "error");
        }
    });
}

//stop pipeline
function beforeStopPipeline() {
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
    confirm("The pipeline design may be modified, would you like to save the pipeline before stop it.", "info", actions);
}

function stopPipeline() {
    loading.show();
    var promise = changeState(pipelineName, pipelineVersionID, 0);
    promise.done(function(data) {
        loading.hide();
        notify(data.message, "success");
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        if (xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else {
            notify("Server is unreachable", "error");
        }
    });
}

function beforeBackToList() {
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
    confirm("The pipeline design may be modified, would you like to save the pipeline before go back to list.", "info", actions);
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
