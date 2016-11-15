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

import { resizeWidget } from "../theme/widget";
import * as historyDataService from "./historyData";
import { notify } from "../common/notify";
import { loading } from "../common/loading";


export function getActionHistory(pipelineName,stageName,actionName,actionLogID) {
    var promise = historyDataService.getActionRunHistory(pipelineName,stageName,actionName,actionLogID);
    promise.done(function(data) {
        loading.hide();
        showActionHistoryView(data.result,actionName);
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

function showActionHistoryView(history,actionname) {
    $.ajax({
        url: "../../templates/history/actionHistory.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#history-pipeline-detail").html($(data));

            $("#actionHistoryTitle").text("Action history -- " + actionname);

            var inputStream = JSON.stringify(history.data.input,undefined,2);
            $("#action-input-stream").val(inputStream);

            var outputStream = JSON.stringify(history.data.output,undefined,2);
            $("#action-output-stream").val(outputStream);

            _.each(history.logList,function(log,index){
                var row = `<p class="history-log" data-index="`+ index +`"></p>`;
                $("#logs").append(row);
                $('.history-log[data-index="'+index+'"]').text(log);
            });

            resizeWidget();
        }
    });
}
