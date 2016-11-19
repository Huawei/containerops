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

// function(pipelineName,versionName,pipelineRunSequence,stageName,actionName){
export function getActionHistory(pipelineName,versionName,pipelineRunSequence,stageName,actionName) {
    loading.show();
    var promise = historyDataService.getActionRunHistory(pipelineName,versionName,pipelineRunSequence,stageName,actionName);
    promise.done(function(data) {
        loading.hide();
        showActionHistoryView(data.result,actionName);
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else {
            notify("Server is unreachable", "error");
        }
    });
}

let sequenceLogDetail = [];
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
                let allLogs = log.substr(23);
                let logJson = JSON.parse(allLogs);
                let num = index + 1;

                if(!logJson.data && !logJson.resp){
                    sequenceLogDetail[index] = logJson.INFO.output;
                    let logTime = log.substr(0,19);

                    var row = `<tr class="log-item"><td>`
                            + num +`</td><td>`
                            + logTime +`</td><td>`
                            + logJson.EVENT +`</td><td>`
                            + logJson.EVENTID +`</td><td>`
                            + logJson.RUN_ID +`</td><td>`
                            + logJson.INFO.status +`</td><td>`
                            + logJson.INFO.result +`</td><td></td><td></td><td><button data-logid="`
                            + index + `" type="button" class="btn btn-success sequencelog-detail"><i class="glyphicon glyphicon-list-alt" style="font-size:14px"></i>&nbsp;&nbsp;Detail</button></td></tr>`;
                    $("#logs-tr").append(row);

                } else {
                    var row = `<tr class="log-item"><td>`
                                    + num + `</td><td></td><td></td><td></td><td></td><td></td><td></td><td>`
                                    + logJson.data +`</td><td>`
                                    + logJson.resp +`</td><td></td></tr>`;
                    $("#logs-tr").append(row);    
                }
            })

             resizeWidget()

            $(".sequencelog-detail").on("click",function(e){
                let target = $(e.target);

                _.each(sequenceLogDetail,function(d,i){
                    if(target.data("logid") == i){
                        let detailData = "";
                        for( let prop in d){
                            detailData += prop + ":" + d[prop].replace(/\n/g, "<br />");
                            detailData += "<br />";
                        }

                        $(".dialogContant").html(detailData);
                    }
                })
  
                $(".dialogWindow").css("height","auto");
                $("#dialog").show();
                if( $(".dialogWindow").height() < $("#dialog").height() * 0.75 ){
                    
                    $(".dialogWindow").css("height","auto");

                } else {
                    
                    $(".dialogWindow").css("height","80%");
                    $(".dialogContant").css("height","100%");
                }

                $("#detailClose").on("click", function(){
                   $("#dialog").hide(); 
                })
            })
        }
    });
}
