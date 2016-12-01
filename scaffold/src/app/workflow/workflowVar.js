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

import * as workflowDataService from "./workflowData";
import { notify, confirm } from "../common/notify";
import { loading } from "../common/loading";
import {hideWorkflowEnv} from "./workflowEnv";

let workflowName, workflowVersionID;
export let workflowVars;

export function initWorkflowVar(name,versionid){
    workflowName = name;
    workflowVersionID = versionid;

    getVarList();
}

export function showWorkflowVar() {
    if ($("#env-setting").hasClass("env-setting-closed")) {
        $("#env-setting").removeClass("env-setting-closed");
        $("#env-setting").addClass("env-setting-opened");
        $("#close_pp_env").removeClass("workflow-open-env");
        $("#close_pp_env").addClass("workflow-close-env");

        $.ajax({
            url: "../../templates/workflow/varSetting.html",
            type: "GET",
            cache: false,
            success: function(data) {
                $("#env-setting").html($(data));

                $(".add-var").on('click', function() {
                    workflowVars.push(["", ""]);
                    showVarKVs();
                });

                $(".workflow-close-env").on('click', function() {
                    hideWorkflowEnv();
                });

                $(".save-var").on('click', function() {
                    saveWorkflowVars();
                });

                showVarKVs();
            }
        });

    } else {
        hideWorkflowEnv();
    }
}

function getVarList() {
    var promise = workflowDataService.getVars(workflowName, workflowVersionID);
    promise.done(function(data) {
        loading.hide();
        workflowVars = _.pairs(data.var);
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

function showVarKVs() {
    $("#vars").empty();
    _.each(workflowVars,function(item,index){
         var row = '<div class="env-row"><div class="env-key-div">'
                        +'<div>'
                            +'<label for="normal-field" class="col-sm-3 control-label" style="margin-top:5px">'
                                +'KEY'
                            +'</label>'
                            +'<div class="col-sm-9" data-index="' + index + '">'
                                +'<input type="text" value="' + item[0] + '" class="form-control pp-env-input pp-var-key" required>'
                            +'</div>'
                        +'</div>'
                    +'</div>'
                    +'<div class="env-value-div" style="margin-left:0">'
                        +'<div>'
                            +'<label for="normal-field" class="col-sm-3 control-label" style="margin-top:5px">'
                                +'VALUE'
                            +'</label>'
                            +'<div class="col-sm-9" data-index="' + index + '">' 
                                +'<input type="text" class="form-control pp-env-input pp-var-value" required>'
                            +'</div>'
                        +'</div>'
                    +'</div>'
                    +'<div class="env-remove-div pp-rm-vkv" data-index="' + index + '">'
                        +'<span class="glyphicon glyphicon-remove"></span>'
                    +'</div></div>';
        $("#envs").append(row);
        $("#envs").find("div[data-index="+index+"]").find(".pp-var-value").val(item[1]);
    });

    $(".pp-var-key").on('input',function(event){
        var key = $(event.currentTarget).val();
        $(event.currentTarget).val(key.toUpperCase());
    });

    $(".pp-var-key").on('blur',function(event){
        var index = $(event.currentTarget).parent().data("index");
        workflowVars[index][0] = $(event.currentTarget).val();
    });

    $(".pp-var-value").on('blur',function(event){
        var index = $(event.currentTarget).parent().data("index");
        workflowVars[index][1] = $(event.currentTarget).val();
    });

    $(".pp-rm-vkv").on('click',function(event){
        var index = $(event.currentTarget).data("index");
        workflowVars.splice(index, 1);
        showVarKVs();
    }); 
}

function saveWorkflowVars() {
    var promise = workflowDataService.setVars(workflowName, workflowVersionID, workflowVars);
    if (promise) {
        promise.done(function(data) {
            loading.hide();
            notify(data.message, "success");
            hideWorkflowEnv();
        });
        promise.fail(function(xhr, status, error) {
            loading.hide();
            if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
                notify(xhr.responseJSON.errMsg, "error");
            } else if(xhr.statusText != "abort") {
                notify("Server is unreachable", "error");
            }
            hideWorkflowEnv();
        });
    }
}