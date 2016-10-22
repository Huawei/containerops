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

import * as constant from "../common/constant";
import {initPipeline} from "../pipeline/initPipeline";
import {initAction} from "../pipeline/initAction";
import {initLine} from "../pipeline/initLine";
import {pipelineData,savePipelineData} from "../pipeline/main";
import {resizeWidget} from "../theme/widget";
import {initActionIO} from "./actionIO";
import {initActionSetup} from "./actionSetup";
import {getAllComponents,getComponent} from "../component/componentData";
import {showNewComponent} from "../component/main";
import {notify} from "../common/notify";
import {loading} from "../common/loading";

export function clickAction(sd, si) {
    $.ajax({
        url: "../../templates/action/actionMain.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#pipeline-info-edit").html($(data));
            if(sd.component){
                showActionEditor(sd);
            }else{
                $(".usecomponent").on('click',function(){
                    getComponents(sd);
                });
            }  
            resizeWidget(); 
        }
    });
}

function showActionEditor(action){
    $.ajax({
        url: "../../templates/action/actionEdit.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#actionMain").html($(data));

            initActionSetup(action);

            initActionIO(action);

            $("#uuid").attr("value", action.id);

            // view select init
            $("#action-component-select").select2({
                minimumResultsForSearch: Infinity
            });
            // $("#k8s-service-protocol").select2({
            //     minimumResultsForSearch: Infinity
            // });     
        }
    });
}

let allComponents;
function getComponents(action){
    loading.show();
    var promise = getAllComponents();
    promise.done(function(data){
        loading.hide();
        allComponents = data.list;
        showComponentList(action);
        if(allComponents.length == 0){
            notify("You have no components to reuse, please go to 'Component' to create one.","info");     
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

function showComponentList(action){
    $.ajax({
        url: "../../templates/action/actionComponentList.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#actionMain").html($(data));    

            $(".newcomponent").on('click',function(){
                $(".menu-component").parent().addClass("active");
                $(".menu-pipeline").parent().removeClass("active");
                notify("Saving current pipeline automatically.","info");
                savePipelineData();
                showNewComponent(true);
            })

            $(".componentlist_body").empty();
            _.each(allComponents,function(item){
                var pprow = '<tr style="height:50px"><td class="pptd">'
                +'<span class="glyphicon glyphicon-menu-down treeclose" data-name="'+item.name+'"></span>&nbsp;'
                +'<span class="glyphicon glyphicon-menu-right treeopen" data-name="'+item.name+'"></span>&nbsp;' 
                + item.name + '</td><td></td><td></td></tr>';
                $(".componentlist_body").append(pprow);
                _.each(item.version,function(version){
                    var vrow = '<tr data-pname="' + item.name + '" data-version="' + version.version + '" data-versionid="' + version.id + '" style="height:50px">'
                    +'<td></td><td class="pptd">' + version.version + '</td>'
                    +'<td><button type="button" class="btn btn-primary cload">Load</button></td></tr>';
                    $(".componentlist_body").append(vrow);
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

            $(".cload").on("click",function(event){
                var target = $(event.currentTarget);
                var componentName = target.parent().parent().data("pname");
                var componentVersionID = target.parent().parent().data("versionid");
                LoadComponentToAction(componentName,componentVersionID,action);
            })
        }
    });
}

function LoadComponentToAction(componentName,componentVersionID,action){
    loading.show();
    var promise = getComponent(componentName,componentVersionID);
    promise.done(function(data){
        loading.hide();
        if(_.isEmpty(data.setupData)){
            notify("Selected component lack base config, can not be loaded.","error");
        }else if(_.isEmpty(data.inputJson)){
            notify("Selected component lack input json, can not be loaded.","error");
        }else if(_.isEmpty(data.outputJson)){
            notify("Selected component lack output json, can not be loaded.","error");
        }else{
            action.setupData = $.extend(true,{},data.setupData);
            action.inputJson = $.extend(true,{},data.inputJson);
            action.outputJson = $.extend(true,{},data.outputJson);
            action.component = {
                "name" : componentName,
                "versionid" : componentVersionID
            }
            showActionEditor(action);
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

function jsonChanged(root,json){
    root.val(JSON.stringify(json));
}