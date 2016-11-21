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

import * as startSetupData from "./startSetupData";
import {initStartIO,initTreeEdit,initFromEdit,initFromView,getOutputForEvent} from "./startIO";
import {getPipelineToken} from "../pipeline/main";
import { notify } from "../common/notify";
import { loading } from "../common/loading";

let startData;
export function initStartSetup(start){
    startData = start;

    // url and token
    showPipeline_URL_Token();
    var urlcopy = new Clipboard('#copyUrl');
    var tokencopy = new Clipboard('#copyToken');
    urlcopy.on('success', function(e) {
        notify("Url copied.","info");
        e.clearSelection();
    });
    urlcopy.on('error', function(e) {
        notify("Copy url failed.","info");
    });
    tokencopy.on('success', function(e) {
        notify("Token copied.","info");
        e.clearSelection();
    });
    tokencopy.on('error', function(e) {
        notify("Copy token failed.","info");
    });

    startSetupData.getStartSetupData(start);
    initStartIO(start);

    // type select
    $("#type-select").val(startSetupData.getTypeSelect());
    selectType(startSetupData.getTypeSelect());

    $("#type-select").on("change",function(){
        startSetupData.setTypeSelect();
        selectType(startSetupData.getTypeSelect(),true);
    });

    $("#type-select").select2({
        minimumResultsForSearch: Infinity
    });

    // event select
    $("#event-select").on("change",function(){
        startSetupData.setEventSelect();
        getOutputForEvent(startSetupData.getEventSelect());
    });
}

function selectType(pipelineType,isTypeChange){
    if(pipelineType == "github" || pipelineType == "gitlab"){
        $("#event_select").show();
        $("#outputTreeViewer").show();
        $("#outputTreeDesigner").hide();
        
        $("#event-select").val(startSetupData.getEventSelect());
        $("#event-select").select2({
            minimumResultsForSearch: Infinity
        });
        getOutputForEvent(startSetupData.getEventSelect()); 
    }else{
        $("#event_select").hide();
        $("#outputTreeViewer").hide();
        $("#outputTreeDesigner").show();

        if(isTypeChange){
            startData.outputJson = {};
        } 
        initTreeEdit();
        initFromEdit("output");
    }
}

function showPipeline_URL_Token(){
    var promise = getPipelineToken();
    promise.done(function(data) {
        loading.hide();
        $("#pp-url").val(data.url);
        $("#pp-token").val(data.token);
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
