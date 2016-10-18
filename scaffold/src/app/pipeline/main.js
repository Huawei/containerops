
import {initDesigner} from "./initDesigner";
import {initPipeline} from "./initPipeline";
import {initAction} from "./initAction";
import {getAllPipelines,getPipeline,addPipeline,addPipelineVersion,getEnvs} from "./pipelineData";
import {notify} from "../common/notify";

export let allPipelines;

export let pipelineData;
let pipelineName, pipelineVersion;
let pipelineEnvs;

export function initPipelinePage(){
    var promise = getAllPipelines();
    promise.done(function(data){
        allPipelines = data.list;
        if(allPipelines.length>0){
            showPipelineList();
        }else{
            showNoPipeline();
        }
    });
    promise.fail(function(xhr,status,error){
        notify(xhr.responseJSON.errMsg,"error");
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
                    var vrow = '<tr data-pname="' + item.name + '" data-version="' + version.id + '" style="height:50px">'
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
                showPipelineDesigner();
            });
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
                    promise.done(function(data){
                        notify(data.msg,"success");
                        initPipelinePage();
                    });
                    promise.fail(function(xhr,status,error){
                        notify(xhr.responseJSON.errMsg,"error");
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

            // var selectedpp = _.find(allPipelines,function(pp){
            //     return pp.name == pipelineName;
            // });
            // var selectedversion = _.find(selectedpp.versions,function(version){
            //     return version.version == pipelineVersion;
            // });
            pipelineData = getPipeline();

            $("#selected_pipeline").text(pipelineName + " / " + pipelineVersion); 

            initDesigner();
            drawPipeline();

            $(".backtolist").on('click',function(){
                initPipelinePage();
            });

            $(".newpipelineversion").on('click',function(){
                showNewPipelineVersion();
            })

            $(".newpipeline").on('click',function(){
                showNewPipeline();
            })

            $(".envsetting").on("click",function(event){
                showPipelineEnv();
            });
        }
    }); 
}

function drawPipeline(){
    $("#pipeline-info-edit").empty();
    
    initPipeline();
    initAction();
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
