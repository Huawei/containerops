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

import {getAllComponents,getComponent,addComponent,addComponentVersion,saveComponent,validateComponent} from "./componentData";
import {initComponentIO} from "./componentIO";
import {initComponentSetup} from "./componentSetup";
import {initComponentEnv} from "./componentEnv";
import {notify,confirm} from "../common/notify";
import {loading} from "../common/loading";

export let allComponents;

export let componentData;
let componentDataOriginalCopy;
let componentName, componentVersion,componentVersionID;

export function initComponentPage(){
    var promise = getAllComponents();
    promise.done(function(data){
        loading.hide();
        allComponents = data.list;
        if(allComponents.length>0){
            showComponentList();
        }else{
            showNoComponent();
        }
    });
    promise.fail(function(xhr,status,error){
        loading.hide();
        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg,"error");
        }else if(xhr.statusText != "abort"){
            notify("Server is unreachable","error");
        }
    });    
}

function showComponentList(){
    $.ajax({
        url: "../../templates/component/componentList.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").html($(data));    
            $("#componentlist").show("slow");

            $(".newcomponent").on('click',function(){
                showNewComponent();
            }) 

            $(".componentlist_body").empty();
            _.each(allComponents,function(item){
                var pprow = `<tr class="pp-row">
                                <td class="pptd">
                                    <span class="glyphicon glyphicon-menu-down treeclose treecontroller" data-name=` 
                                    + item.name +`></span><span style="margin-left:10px">`
                                    + item.name 
                                +`</span></td><td></td><td></td></tr>`;
                $(".componentlist_body").append(pprow);

                _.each(item.version,function(version){
                    var vrow = `<tr data-pname=` + item.name + ` data-version=` + version.version + ` data-versionid=` 
                                + version.id + ` class="ppversion-row">
                                    <td></td>
                                    <td class="pptd">` + version.version + `</td>
                                    <td>
                                        <button type="button" class="btn btn-success ppview">
                                            <i class="glyphicon glyphicon-eye-open" style="font-size:16px"></i>
                                            <span style="margin-left:5px">View</span>
                                        </button>
                                    </td>
                                </tr>`;
                
                $(".componentlist_body").append(vrow);
                })
            }) ;

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

            $(".ppview").on("click",function(event){
                var target = $(event.currentTarget);
                componentName = target.parent().parent().data("pname");
                componentVersion = target.parent().parent().data("version");
                componentVersionID = target.parent().parent().data("versionid");
                getComponentData();
            })
        }
    });
}

function getComponentData(){
    var promise = getComponent(componentName,componentVersionID);
    promise.done(function(data){
        loading.hide();
        componentData = data;
        componentDataOriginalCopy = $.extend(true,{},data);
        showComponentDesigner();
    });
    promise.fail(function(xhr,status,error){
        loading.hide();
        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg,"error");
        }else if(xhr.statusText != "abort"){
            notify("Server is unreachable","error");
        }
    });
}

function showNoComponent(){
    $.ajax({
        url: "../../templates/component/noComponent.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").html($(data));    
            $("#nocomponent").show("slow");
            $(".newcomponent").on('click',function(){
                showNewComponent();
            })  
        }
    });
}

export function showNewComponent(fromWorkflow){
    $.ajax({
        url: "../../templates/component/newComponent.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").children().hide();
            $("#main").append($(data));    
            $("#newcomponent").show("slow");
            $("#newComponentBtn").on('click',function(){
                var promise = addComponent();
                if(promise){
                    promise.done(function(data){
                        loading.hide();
                        notify(data.message,"success");
                        initComponentPage();
                    });
                    promise.fail(function(xhr,status,error){
                        loading.hide();
                        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
                            notify(xhr.responseJSON.errMsg,"error");
                        }else if(xhr.statusText != "abort"){
                            notify("Server is unreachable","error");
                        }
                    });
                }
            })
            $("#cancelNewComponentBtn").on('click',function(){
                if(fromWorkflow){
                    $(".menu-component").parent().removeClass("active");
                    $(".menu-workflow").parent().addClass("active");
                }
                cancelNewComponentPage();
            })
        }
    });
}

function showComponentDesigner(){  
    $.ajax({
        url: "../../templates/component/componentDesign.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").html($(data));    
            $("#componentdesign").show("slow"); 

            $("#selected_component").text(componentName + " / " + componentVersion); 

            $(".backtolist").on('click',function(){
                beforeBackToList();
            });

            $(".savecomponent").on('click',function(){
                saveComponentData();
            });

            $(".newcomponentversion").on('click',function(){
                if(validateComponent(componentData)){
                    showNewComponentVersion();
                }
            });

            $(".newcomponentindesigner").on('click',function(){
                beforeShowNewComponent();
            });

            initComponentEdit();
        }
    }); 
}

function initComponentEdit(){
    $.ajax({
        url: "../../templates/component/componentEdit.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#componentDesigner").html($(data));

            initComponentSetup(componentData);

            initComponentIO(componentData);

            initComponentEnv(componentData);

            // view select init
            $("#action-component-select").select2({
               minimumResultsForSearch: Infinity
            });
            
            $("#service-type-select").select2({
               minimumResultsForSearch: Infinity
            });
        }
    });
}

function showNewComponentVersion(){
    $.ajax({
        url: "../../templates/component/newComponentVersion.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#main").children().hide();
            $("#main").append($(data));    
            $("#newcomponentversion").show("slow"); 

            $("#c-name-newversion").val(componentName);

            $("#newComponentVersionBtn").on('click',function(){
                var promise = addComponentVersion(componentName,componentVersionID,componentData);
                if(promise){
                    promise.done(function(data){
                        loading.hide();
                        notify(data.message,"success");
                        initComponentPage();
                    });
                    promise.fail(function(xhr,status,error){
                        loading.hide();
                        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
                            notify(xhr.responseJSON.errMsg,"error");
                        }else if(xhr.statusText != "abort"){
                            notify("Server is unreachable","error");
                        }
                    });
                }
            })
            $("#cancelNewComponentVersionBtn").on('click',function(){
                cancelNewComponentVersionPage();
            })      
        }
    }); 
    
    $("#content").hide();
    $("#nocomponent").hide();
    $("#newcomponent").hide();
    $("#newcomponentversion").show("slow");
}

function cancelNewComponentPage(){
    $("#newcomponent").remove();
    $("#main").children().show("slow");
}

function cancelNewComponentVersionPage(){
    $("#newcomponentversion").remove();
    $("#main").children().show("slow");
}

function beforeBackToList(){
    if(_.isEqual(componentDataOriginalCopy,componentData)){
        initComponentPage();
    }else{
        var actions = [
            {
                "name":"save",
                "label":"Yes",
                "action":function(){
                    saveComponentData(initComponentPage);
                }
            },{
                "name":"back",
                "label":"No",
                "action":function(){
                    initComponentPage();
                }
            }
        ]
        confirm("The component design has been modified, would you like to save the changes before go back to list.","info",actions);
    }
}

function saveComponentData(next){
    if(validateComponent(componentData)){
        var promise = saveComponent(componentName, componentVersion, componentVersionID, componentData);
        promise.done(function(data){
            componentDataOriginalCopy = $.extend(true,{},componentData);
            loading.hide();
            if(!next){
                notify(data.message,"success");
            }else{
                next();
            }
        });
        promise.fail(function(xhr,status,error){
            loading.hide();
            if(!next){
                if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
                    notify(xhr.responseJSON.errMsg,"error");
                }else if(xhr.statusText != "abort"){
                    notify("Server is unreachable","error");
                }
            }else{
                next();
            } 
        });
    } 
}

function beforeShowNewComponent(){
    if(_.isEqual(componentDataOriginalCopy,componentData)){
        showNewComponent();
    }else{
        var actions = [
            {
                "name":"save",
                "label":"Yes",
                "action":function(){
                    saveComponentData(showNewComponent);
                }
            },{
                "name":"show",
                "label":"No",
                "action":function(){
                    showNewComponent();
                }
            }
        ]
        confirm("The component design has been modified, would you like to save the changes at first.","info",actions);
    }
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
