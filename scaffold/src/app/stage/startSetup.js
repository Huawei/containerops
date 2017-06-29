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

import {initStartIO,initTreeEdit,initFromEdit,initFromView,getOutputForEvent} from "./startIO";
import {getWorkflowToken} from "../workflow/main";
import { notify } from "../common/notify";
import { loading } from "../common/loading";

let startData;
let urlcopy,tokencopy;

export function initStartSetup(start){
    startData = start;

    // url and token
    showWorkflow_URL_Token();

    if(_.isUndefined(urlcopy)){
        urlcopy = new Clipboard('#copyUrl');
        urlcopy.on('success', function(e) {
            notify("Url copied.","info");
            e.clearSelection();
        });
        urlcopy.on('error', function(e) {
            notify("Copy url failed.","info");
        });
    }
    
    if(_.isUndefined(tokencopy)){
        tokencopy = new Clipboard('#copyToken');
        tokencopy.on('success', function(e) {
            notify("Token copied.","info");
            e.clearSelection();
        });
        tokencopy.on('error', function(e) {
            notify("Copy token failed.","info");
        });
    } 

    initStartIO(start);
}

function showWorkflow_URL_Token(){
    var promise = getWorkflowToken();
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
