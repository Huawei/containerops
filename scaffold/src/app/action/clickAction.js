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

import * as constant from "../common/constant";
import { initWorkflow } from "../workflow/initWorkflow";
import { initAction } from "../workflow/initAction";
import { initLine } from "../workflow/initLine";
import { workflowData, saveWorkflowData } from "../workflow/main";
import { resizeWidget } from "../theme/widget";
import { initActionIO } from "./actionIO";
import { initActionSetup } from "./actionSetup";
import { initActionEnv } from "./actionEnv";
import { getAllComponents, getComponent } from "../component/componentData";
import { showNewComponent } from "../component/main";
import { notify } from "../common/notify";
import { loading } from "../common/loading";
import { getConflict, svgTree } from "./actionConflict";
import {workflowVars} from "../workflow/workflowVar";

let filter = "";

export function clickAction(sd, si) {
    filter = "";
    
    if (sd.component) {
        showActionEditor(sd);
    } else {
        $.ajax({
            url: "../../templates/action/actionMain.html",
            type: "GET",
            cache: false,
            success: function(data) {
                $("#workflow-info-edit").html($(data));
                $(".usecomponent").on('click', function() {
                    getComponents(sd);
                });
                resizeWidget();
            }
        })
    }
}

function showActionEditor(action) {
    $.ajax({
        url: "../../templates/action/actionEdit.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#workflow-info-edit").html($(data));

            initActionSetup(action);

            initActionIO(action);

            initActionEnv(action);

            getConflict(action.id);

            $("#uuid").attr("value", action.id);

            // view select init
            $("#action-component-select").select2({
                minimumResultsForSearch: Infinity
            });

            $("#service-type-select").select2({
               minimumResultsForSearch: Infinity
            });
            
            // use global vars
            var globalvars = _.map(workflowVars,function(item){
                                return "@"+item[0]+"@";
                            });
            $(".allowFromVar").autocomplete({
                source:[globalvars],
                limit: 100,
                visibleLimit: 5
            }); 

            resizeWidget();
        }
    });
}

let allComponents;

function getComponents(action) {
    var promise = getAllComponents();
    promise.done(function(data) {
        loading.hide();
        allComponents = data.list;
        showComponentList(action);
        if (allComponents.length == 0) {
            notify("You have no components to reuse, please go to 'Component' to create one.", "info");
        }
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else if(xhr.statusText != "abort"){
            notify("Server is unreachable", "error");
        }
    });
}

function showComponentList(action) {
    $.ajax({
        url: "../../templates/action/actionComponentList.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#actionMain").html($(data));

            $(".component-filter-input").val(filter);

            $(".newcomponent").on('click', function() {
                $(".menu-component").parent().addClass("active");
                $(".menu-workflow").parent().removeClass("active");
                notify("Saving current workflow automatically.", "info");
                saveWorkflowData();
                showNewComponent(true);
            })

            $("#searchComponent").on('click', function() {
                filter = $(".component-filter-input").val();
                showComponentList(action);
            })

            var components = doFilter(filter);
            
            $(".componentlist_body").empty();
            _.each(components, function(item) {
                var pprow = `<tr class="pp-row">
                                <td class="pptd">
                                    <span class="glyphicon glyphicon-menu-down treeclose treecontroller" data-name=` 
                                    + item.name +`></span><span style="margin-left:10px">`
                                    + item.name 
                                +`</span></td><td></td><td></td></tr>`;
                $(".componentlist_body").append(pprow);
                _.each(item.version, function(version) {
                    var vrow = `<tr data-pname=` + item.name + ` data-version=` + version.version + ` data-versionid=` 
                                + version.id + ` class="ppversion-row">
                                    <td></td>
                                    <td class="pptd">` + version.version + `</td>
                                    <td>
                                        <button type="button" class="btn btn-success ppview cload">
                                            <i class="fa fa-copy" style="font-size:16px"></i>
                                            <span style="margin-left:5px">Load</span>
                                        </button>
                                    </td>
                                </tr>`;
                    $(".componentlist_body").append(vrow);
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

            $(".cload").on("click", function(event) {
                var target = $(event.currentTarget);
                var componentName = target.parent().parent().data("pname");
                var componentVersionName = target.parent().parent().data("version");
                var componentVersionID = target.parent().parent().data("versionid");
                LoadComponentToAction(componentName, componentVersionName, componentVersionID, action);
            })
        }
    });
}

function LoadComponentToAction(componentName, componentVersionName, componentVersionID, action) {
    var promise = getComponent(componentName, componentVersionID);
    promise.done(function(data) {
        loading.hide();
        if (_.isEmpty(data.setupData)) {
            notify("Selected component lack base config, can not be loaded.", "error");
        } else if (_.isEmpty(data.inputJson)) {
            notify("Selected component lack input json, can not be loaded.", "error");
        } else if (_.isEmpty(data.outputJson)) {
            notify("Selected component lack output json, can not be loaded.", "error");
        } else {
            action.setupData = $.extend(true, {}, data.setupData);
            action.inputJson = $.extend(true, {}, data.inputJson);
            action.outputJson = $.extend(true, {}, data.outputJson);
            action.env = [].concat(data.env);
            action.component = {
                "name": componentName,
                "versionid": componentVersionID,
                "versionname" : componentVersionName
            }
            showActionEditor(action);
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

function jsonChanged(root, json) {
    root.val(JSON.stringify(json));
}

function doFilter(filter){
    var tempComponents = _.map(allComponents,function(item){
        return $.extend(true,{},item);
    });

    return _.filter(tempComponents,function(item){
        return item.name.toLowerCase().indexOf(filter) >= 0;
    })
}
