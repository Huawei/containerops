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

import {initDesigner} from "./initDesigner";
import {initPipeline} from "./initPipeline";
import {initAction} from "./initAction";
import {getAllPipelines,getPipeline,addPipeline,savePipeline,addPipelineVersion,getEnvs} from "./pipelineData";
import {notify} from "../common/notify";
import {loading} from "../common/loading";
import {setLinePathAry,linePathAry} from "../common/constant";

export let allPipelines;

export let pipelineData;
let pipelineName, pipelineVersion, pipelineVersionID;
let pipelineEnvs;

export function initPipelinePage(){
    loading.show();
    var promise = getAllPipelines();
    promise.done(function(data){
        loading.hide();
        allPipelines = data.list;
        if(allPipelines.length>0){
            showPipelineList();
        }else{
            showNoPipeline();
        }
    });
    promise.fail(function(xhr,status,error){
        loading.hide();
        if(xhr.responseJSON.errMsg){
            notify(xhr.responseJSON.errMsg,"error");
        }else{
            notify("Server is unreachable","error");
        }
    });
}

function showPipelineList(){
    $.ajax({
        url: "../../templates/pipeline/pipelineList.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").html($(data));    
            $("#pipelinelist").show("slow");

            $(".newpipeline").on('click',function(){
                showNewPipeline();
            }) 

            $(".pipelinelist_body").empty();
            _.each(allPipelines,function(item){
                var pprow = '<tr style="height:50px"><td class="pptd">'
                        +'<span class="glyphicon glyphicon-menu-down treeclose" data-name="'+item.name+'"></span>&nbsp;'
                        +'<span class="glyphicon glyphicon-menu-right treeopen" data-name="'+item.name+'"></span>&nbsp;' 
                        + item.name + '</td><td></td><td></td></tr>';
                $(".pipelinelist_body").append(pprow);
                _.each(item.version,function(version){
                    var vrow = '<tr data-pname="' + item.name + '" data-version="' + version.version + '" data-versionid="' + version.id + '" style="height:50px">'
                            +'<td></td><td class="pptd">' + version.version + '</td>'
                            +'<td><button type="button" class="btn btn-primary ppview">View</button></td></tr>';
                    $(".pipelinelist_body").append(vrow);
                })
            }) ;

            $(".treeclose").on("click",function(event){
                var target = $(event.currentTarget);
                target.hide();
                target.next().show();
                var name = target.data("name");
                $('*[data-pname='+name+']').hide();
            });

            $(".treeopen").on("click",function(event){
                var target = $(event.currentTarget);
                target.hide();
                target.prev().show();
                var name = target.data("name");
                $('*[data-pname='+name+']').show();
            });

            $(".ppview").on("click",function(event){
                var target = $(event.currentTarget);
                pipelineName = target.parent().parent().data("pname");
                pipelineVersion = target.parent().parent().data("version");
                pipelineVersionID = target.parent().parent().data("versionid");
                getPipelineData();
            });
        }
    });
}

function getPipelineData(){
    loading.show();
    var promise = getPipeline(pipelineName,pipelineVersionID);
    promise.done(function(data){
        loading.hide();
        pipelineData = data.stageList;
        setLinePathAry(data.lineList);
        showPipelineDesigner();
    });
    promise.fail(function(xhr,status,error){
        loading.hide();
        if(xhr.responseJSON.errMsg){
            notify(xhr.responseJSON.errMsg,"error");
        }else{
            notify("Server is unreachable","error");
        }
    });
}

function showNoPipeline(){
    $.ajax({
        url: "../../templates/pipeline/noPipeline.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").html($(data));    
            $("#nopipeline").show("slow");
            $(".newpipeline").on('click',function(){
                showNewPipeline();
            })  
        }
    });
}

function showNewPipeline(){
    $.ajax({
        url: "../../templates/pipeline/newPipeline.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").children().hide();
            $("#main").append($(data));    
            $("#newpipeline").show("slow");
            $("#newppBtn").on('click',function(){
                var promise = addPipeline();
                if(promise){
                    loading.show();
                    promise.done(function(data){
                        loading.hide();
                        notify(data.message,"success");
                        initPipelinePage();
                    });
                    promise.fail(function(xhr,status,error){
                        loading.hide();
                        if(xhr.responseJSON.errMsg){
                            notify(xhr.responseJSON.errMsg,"error");
                        }else{
                            notify("Server is unreachable","error");
                        }
                    });
                }
            })
            $("#cancelNewppBtn").on('click',function(){
                cancelNewPPPage();
            })
        }
    });
}

function showPipelineDesigner(){ 
    $.ajax({
        url: "../../templates/pipeline/pipelineDesign.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").html($(data));    
            $("#pipelinedesign").show("slow"); 

            $("#selected_pipeline").text(pipelineName + " / " + pipelineVersion); 

            initDesigner();
            drawPipeline();

            $(".backtolist").on('click',function(){
                initPipelinePage();
            });

            $(".savepipeline").on('click',function(){
                savePipelineData();
            });

            $(".newpipelineversion").on('click',function(){
                showNewPipelineVersion();
            });

            $(".newpipeline").on('click',function(){
                showNewPipeline();
            });

            $(".envsetting").on("click",function(event){
                showPipelineEnv();
            });
        }
    }); 
}

function drawPipeline(){
    $("#pipeline-info-edit").empty();
    
    initPipeline();
    // initAction();
}

export function savePipelineData(silence){
    if(!silence){
        loading.show();
    }
    
    var promise = savePipeline(pipelineName,pipelineVersion,pipelineVersionID,pipelineData,linePathAry);
    promise.done(function(data){
        if(!silence){
            loading.hide();
            notify(data.message,"success");
        }
        
    });
    promise.fail(function(xhr,status,error){
        if(!silence){
            loading.hide();
            if(xhr.responseJSON.errMsg){
                notify(xhr.responseJSON.errMsg,"error");
            }else{
                notify("Server is unreachable","error");
            }
        } 
    });
}

function showNewPipelineVersion(){
    $.ajax({
        url: "../../templates/pipeline/newPipelineVersion.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").children().hide();
            $("#main").append($(data));    
            $("#newpipelineversion").show("slow"); 

            $("#pp-name-newversion").val(pipelineName);

            $("#newppVersionBtn").on('click',function(){
                // addPipelineVersion(pipelineVersion);

                // to be removed below
                if(addPipelineVersion(pipelineVersion)){
                    initPipelinePage();
                } 
            })
            $("#cancelNewppVersionBtn").on('click',function(){
                cancelNewPPVersionPage();
            })      
        }
    }); 
    
    $("#content").hide();
    $("#nopipeline").hide();
    $("#newpipeline").hide();
    $("#newpipelineversion").show("slow");
}

function cancelNewPPPage(){
    $("#newpipeline").remove();
    $("#main").children().show("slow");
}

function cancelNewPPVersionPage(){
    $("#newpipelineversion").remove();
    $("#main").children().show("slow");
}


function showPipelineEnv(){
    if($("#env-setting").hasClass("env-setting-closed")){
        $("#env-setting").removeClass("env-setting-closed");
        $("#env-setting").addClass("env-setting-opened");

        $.ajax({
            url: "../../templates/pipeline/envSetting.html",
            type: "GET",
            cache: false,
            success: function (data) {
                $("#env-setting").html($(data));
              
                $(".new-kv").on('click',function(){
                    pipelineEnvs.push({
                        "key" : "",
                        "value" : ""
                    });
                    showEnvKVs();
                });

                $(".close-env").on('click',function(){
                    hidePipelineEnv();
                });

                $(".save-env").on('click',function(){
                    savePipelineEnvs();
                });   
                
                getEnvList();
            }   
        }); 
        
    }else{
        $("#env-setting").removeClass("env-setting-opened");
        $("#env-setting").addClass("env-setting-closed");
    }
}

function hidePipelineEnv(){
    $("#env-setting").removeClass("env-setting-opened");
    $("#env-setting").addClass("env-setting-closed");
}

function getEnvList(){
    pipelineEnvs = getEnvs();
    showEnvKVs();   
}

function showEnvKVs(){
    $("#envs").empty();
    _.each(pipelineEnvs,function(item,index){
        var row = '<tr data-index="'+index+'"><td>'
                    +'<input type="text" class="form-control col-md-5 env-key" value="'+item.key+'">'
                    + '</td><td>'
                    +'<input type="text" class="form-control col-md-5 env-value" value="'+item.value+'">'
                    + '</td><td>'
                    +'<span class="glyphicon glyphicon-minus rm-kv"></span>'
                    +'</td></tr>';
        $("#envs").append(row);
    });

    $(".env-key").on('input',function(event){
        var key = $(event.currentTarget).val();
        $(event.currentTarget).val(key.toUpperCase());
    });

    $(".env-key").on('blur',function(event){
        var index = $(event.currentTarget).parent().parent().data("index");
        pipelineEnvs[index].key = $(event.currentTarget).val();
    });

    $(".env-value").on('blur',function(event){
        var index = $(event.currentTarget).parent().parent().data("index");
        pipelineEnvs[index].value = $(event.currentTarget).val();
    });

    $(".rm-kv").on('click',function(event){
        var index = $(event.currentTarget).parent().parent().data("index");
        pipelineEnvs.splice(index,1);
        showEnvKVs();
    });
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
