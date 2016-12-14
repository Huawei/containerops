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

// function(workflowName,versionName,workflowRunSequence,stageName,actionName){
export function getActionHistory(workflowName,versionName,workflowRunSequence,stageName,actionName) {
    // loading.show();
    var promise = historyDataService.getActionRunHistory(workflowName,versionName,workflowRunSequence,stageName,actionName);
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

let eventStatusData;
let sequenceLogDetailData = [];
let sequencResultLogsData = [];
function showActionHistoryView(history,actionname) {
    $.ajax({
        url: "../../templates/history/actionHistory.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#history-workflow-detail").html($(data));
            $("#actionHistoryTitle").text("Action history -- " + actionname);

            var inputStream = JSON.stringify(history.data.input,undefined,2);
            $("#action-input-stream").val(inputStream);

            var outputStream = JSON.stringify(history.data.output,undefined,2);
            $("#action-output-stream").val(outputStream);

            var eventLogs = history.logList;
            _.each(eventLogs,function(log,index){

                let allLogs = log.substr(23);
                allLogs = allLogs.replace(/\\n/g , "\\u003cbr /\\u003e")
                eventStatusData = JSON.parse(allLogs);
                let lineNo = index + 1;
                let eventTime = log.substr(0,19);

                sequenceLogDetailData.push(eventStatusData.INFO);
                getEventStatus(eventStatusData,eventTime,lineNo);

                if ( eventStatusData.EVENT == "CO_TASK_RESULT" && eventStatusData.INFO != null ){
                    let uResultLog = JSON.stringify(eventStatusData.INFO); 
                     sequencResultLogsData.push(uResultLog);
                }
                
            })

            getResultLog(sequencResultLogsData);
            resizeWidget()

            $(".sequencelog-detail").on("click",function(e){
                let target = $(e.target);
                var tempLogIdArray = (target.attr("data-logid")).split("_");
                if(null != tempLogIdArray && tempLogIdArray.length > 1){
                    var logjsoin =  sequenceLogDetailData[tempLogIdArray[1]];
                    let detailData = "";
                    for( let prop in logjsoin){
                        var showLogJson = logjsoin[prop];
                        if( typeof(showLogJson) == "object"){
                            for( let subProp in showLogJson){
                                let showStr =showLogJson[subProp].replace(/\\n/g,"<br/>");
                                detailData += subProp + ":" + showStr;
                                detailData += "<br />"
                            }
                        }else {
                            detailData += prop + ":" + showLogJson;
                            detailData += "<br />";
                        }
                    }
                    $(".dialogContant").html(detailData);   
                }

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


function getEventStatus(eventData,eventTime,lineNo){
    var row = `<tr class="log-item"><td>`
            + lineNo +`</td><td>`
            + eventTime +`</td><td>`
            + eventData.EVENT +`</td><td>`
            + eventData.EVENT_ID +`</td><td>`
            + eventData.RUN_ID +`</td><td>`
            + eventData.INFO.status +`</td><td>`
            + eventData.INFO.result +`</td>`
            + `<td><button data-logid="`
            + "info_" + (lineNo - 1) + `" type="button" class="btn btn-success sequencelog-detail"><i class="glyphicon glyphicon-list-alt" style="font-size:14px"></i>&nbsp;&nbsp;Detail</button></td></tr>`;
    
    $("#logs-tr").append(row);
}

function getResultLog(resultData){
    $("#resultLog_list").html();
    if( resultData != null && resultData.length>0) {
        _.each(resultData,function(rd,i){
            var row = `<tr><td class="loglist-td">`+i+`</td><td>`
                + rd +`</td></tr>`; 
            $("#resultLog_list").append(row);
        })
        sequencResultLogsData = [];
    }
}

