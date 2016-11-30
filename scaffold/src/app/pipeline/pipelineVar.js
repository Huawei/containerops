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

import * as pipelineDataService from "./pipelineData";
import { notify, confirm } from "../common/notify";
import { loading } from "../common/loading";
import {hidePipelineEnv} from "./pipelineEnv";

let pipelineName, pipelineVersionID;
export let pipelineVars;

export function initPipelineVar(name,versionid){
    pipelineName = name;
    pipelineVersionID = versionid;

    getVarList();
}

export function showPipelineVar() {
    if ($("#env-setting").hasClass("env-setting-closed")) {
        $("#env-setting").removeClass("env-setting-closed");
        $("#env-setting").addClass("env-setting-opened");
        $("#close_pp_env").removeClass("pipeline-open-env");
        $("#close_pp_env").addClass("pipeline-close-env");

        $.ajax({
            url: "../../templates/pipeline/varSetting.html",
            type: "GET",
            cache: false,
            success: function(data) {
                $("#env-setting").html($(data));

                $(".add-var").on('click', function() {
                    pipelineVars.push(["", ""]);
                    showVarKVs();
                });

                $(".pipeline-close-env").on('click', function() {
                    hidePipelineEnv();
                });

                $(".save-var").on('click', function() {
                    savePipelineVars();
                });

                showVarKVs();
            }
        });

    } else {
        hidePipelineEnv();
    }
}

function getVarList() {
    var promise = pipelineDataService.getVars(pipelineName, pipelineVersionID);
    promise.done(function(data) {
        loading.hide();
        pipelineVars = _.pairs(data.var);
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
    _.each(pipelineVars,function(item,index){
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
        pipelineVars[index][0] = $(event.currentTarget).val();
    });

    $(".pp-var-value").on('blur',function(event){
        var index = $(event.currentTarget).parent().data("index");
        pipelineVars[index][1] = $(event.currentTarget).val();
    });

    $(".pp-rm-vkv").on('click',function(event){
        var index = $(event.currentTarget).data("index");
        pipelineVars.splice(index, 1);
        showVarKVs();
    }); 
}

function savePipelineVars() {
    var promise = pipelineDataService.setVars(pipelineName, pipelineVersionID, pipelineVars);
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