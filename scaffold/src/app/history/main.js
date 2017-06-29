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

import { loading } from "../common/loading";
import * as constant from "../common/constant";
import * as historyDataService from "./historyData";
import { setPath } from "../relation/setPath";
import { notify } from "../common/notify";
import { getActionHistory } from "./actionHistory";
import { getContainerLogs } from "./actionHistory";
import { getLineHistory } from "./lineHistory";
import { getHistoryList,addFilterWorklowEvent,getSequenceNum } from "./historyList";
import {changeCurrentElement} from "../common/util";
import * as sequenceUtil from "./initUtil";

var filterType = 'fuzzy';
export function initHistoryPage() {
    var type = arguments.length===0 ? 'fuzzy': 'exact';
    var keywords = arguments.length===0 ? '': arguments[0];
    clearTimeout(timer);
    setFilterType(type);
    getHistory(keywords,filterType);  
}

function getHistory(keywords,filterType) {
    $.ajax({
        url: "../../templates/history/historyList.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#main").html($(data));
            addFilterWorklowEvent();
            getSequenceNum();
            getHistoryList(keywords,filterType);
        }
    });
}

function setFilterType(type){
    filterType = type;
}

function forVdSequenceList(vd,index,length,pdId,pdName,vdId,vdName){

    var hsRowArray = [];
    var tempLength = (length > 0) ? length : vd.length;


    for(var i = index ; i < tempLength ; i++){
        var sd = vd[i] ;
        var hsRow = `<tr data-id=` + sd.sequence + ` data-status=` + sd.status + ` data-pname=` + pdName + ` data-version=` + vdName + ` data-versionid=` + vdId + ` class="sequence-row"><td></td><td></td>`;
        // var hsRow = `<tr data-id=` + sd.sequence + ` data-status=` + sd.status + ` data-pname=` + pd.name + ` data-version=` + vd.name + ` data-versionid=` + vd.id + ` class="sequence-row"><td></td><td></td>`;
        let sdTime = sd.time;

        if (sd.status == 1 || sd.status == 0) {

            hsRow += `<td><div class="state-list"><div class="state-icon-list state-waitStart"></div><span class="state-label-list">` + sd.time + `</span></div></td>`;

        } else if (sd.status == 2) {

            hsRow += `<td><div class="state-list"><div class="state-icon-list state-running"></div><span class="state-label-list">` + sd.time + `</span></div></td>`;

        } else if (sd.status == 3) {

            hsRow += `<td><div class="state-list"><div class="state-icon-list state-success"></div><span class="state-label-list">` + sd.time + `</span></div></td>`;

        } else if(sd.status == 4) {

            hsRow += `<td><div class="state-list"><div class="state-icon-list state-fail"></div><span class="state-label-list">` + sd.time + `</span></div></td>`
        }

        hsRow += `<td><button type="button" class="btn btn-success sequence-detail"><i class="glyphicon glyphicon-list-alt" style="font-size:16px"></i><span style="margin-left:5px">Detail</span></button></td></tr> `

         if(i >= index){
            hsRowArray[i]=hsRow ;
        }

        if(length > 0 ){
            if(hsRowArray.length == tempLength){
                var btnMore = `<tr data-insertid=` + sd.workflowSequenceID + ` class="btn-more"><td colspan="4" id="btn_`+pdName+`_`+pdId+`" class="pptd btn-showMorm"  >点击查看更多 \<\< </td></tr>`;
                hsRowArray[i+1] = btnMore ;
                // $(".btn-more").css({"font-size":"12px","color":"#7C7C7C","text-align":"center"});
                break ;
            }
        }
    }
    return hsRowArray ;
}

function addMore(vd,index,pdId,pdName,vdId,vdName){

    var tempLength = vd.length;


    for(var i = index ; i < tempLength ; i++){
        var sd = vd[i] ;
        var hsRow = `<tr data-id=` + sd.sequence + ` data-status=` + sd.status + ` data-pname=` + pdName + ` data-version=` + vdName + ` data-versionid=` + vdId + ` class="sequence-row"><td></td><td></td>`;
        // var hsRow = `<tr data-id=` + sd.sequence + ` data-status=` + sd.status + ` data-pname=` + pd.name + ` data-version=` + vd.name + ` data-versionid=` + vd.id + ` class="sequence-row"><td></td><td></td>`;
        let sdTime = sd.time;

        if (sd.status == 1 || sd.status == 0) {

            hsRow += `<td><div class="state-list"><div class="state-icon-list state-waitStart"></div><span class="state-label-list">` + sd.time + `</span></div></td>`;

        } else if (sd.status == 2) {

            hsRow += `<td><div class="state-list"><div class="state-icon-list state-running"></div><span class="state-label-list">` + sd.time + `</span></div></td>`;

        } else if (sd.status == 3) {

            hsRow += `<td><div class="state-list"><div class="state-icon-list state-success"></div><span class="state-label-list">` + sd.time + `</span></div></td>`;

        } else if(sd.status == 4) {

            hsRow += `<td><div class="state-list"><div class="state-icon-list state-fail"></div><span class="state-label-list">` + sd.time + `</span></div></td>`
        }

        hsRow += `<td><button type="button" class="btn btn-success sequence-detail"><i class="glyphicon glyphicon-list-alt" style="font-size:16px"></i><span style="margin-left:5px">Detail</span></button></td></tr> `

         if(i >= index){
            $("#btn_"+pdName+"_"+pdId).parent().before(hsRow);
        }
        $("#btn_"+pdName+"_"+pdId).parent().hide();
    }

}

let historyAbout;
var isAction, actionId;
export function getSequenceDetail(selected_history) {
    // todo: 判断 action 来源
    // 在  initSequenceActionByStage   判断属性
    actionId = selected_history.actionId;
    var actionName = selected_history.actionName;
    var stageId = selected_history.stageId;
    var stageName = selected_history.stageName;
  
    if ( actionId != null || actionId > 0){
        isAction = true;
        actionId = "a-"+ actionId;
        
    }else {
        isAction = false;
    }

    historyAbout = selected_history;
    loading.show();
    constant.sequenceRunStatus = selected_history.sequenceStatus;
    var promise = historyDataService.getWorkflowHistory(selected_history.workflowName, selected_history.versionName, selected_history.sequence);
    promise.done(function(data) {
        loading.hide();
        constant.sequenceRunData = data.define.stageList;
        constant.refreshSequenceRunData = data.define.stageList;
        constant.sequenceLinePathArray = data.define.lineList;
        if (data.define.stageList.length > 0 && !isAction ) {
            initSequenceView(selected_history);
        } else if( data.define.stageList.length > 0 && isAction){
            initSequenceView(selected_history);
            getActionHistory(historyAbout.workflowName,historyAbout.versionName,historyAbout.sequence,historyAbout.stageName,historyAbout.actionName);
            getContainerLogs(historyAbout.workflowName,historyAbout.versionName,historyAbout.sequence,historyAbout.stageName,historyAbout.actionName,"");
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

function initSequenceView(selected_history) {
    $.ajax({
        url: "../../templates/history/historyView.html",
        type: "GET",
        cache: false,
        success: function(data) {
            let zoom = d3.behavior.zoom().on("zoom", zoomed);
            $("#main").html($(data));
            $("#historyView").show("slow");

            $("#selected_workflow").text(selected_history.workflowName + " / " + selected_history.versionName);

            $(".backtolist").on('click', function() {
                initHistoryPage();
                clearTimeout(timer);
            });

            let $div = $("#div-d3-main-svg").height($("main").height() * 3 / 7);
            // let zoom = d3.behavior.zoom().on("zoom", zoomed);
            let drag = d3.behavior.drag()
                .origin(function() {
                    return { "x": 0, "y": 0 };
                })
                .on("dragstart", dragStart)
                .on("drag", sequenceUtil.draged);

            function dragStart() {
                d3.event.sourceEvent.stopPropagation();
                drag.origin(function() {
                    return { "x": constant.sequenceWorkflowView.attr("translateX"), "y": constant.sequenceWorkflowView.attr("translateY") }
                });
            }    

            constant.setSvgWidth("100%");
            constant.setSvgHeight("100%");
            constant.setWorkflowNodeStartX(50);
            constant.setWorkflowNodeStartY($div.height() * 0.2);

            $div.empty();

            let svg = d3.select("#div-d3-main-svg")
                .on("touchstart", nozoom)
                .on("touchmove", nozoom)
                .append("svg")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight)
                .style("fill", "white");

            let g = svg.append("g")
                .call(drag);
            // .on("dblclick.zoom", null);

            let svgMainRect = g.append("rect")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight);

            constant.sequenceLinesView = g.append("g")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight)
                .attr("id", "sequenceLinesView")
                .attr("translateX", 0)
                .attr("translateY", 0)
                .attr("transform", "translate(0,0) scale(1)")
                .attr("scale", 1);

            constant.sequenceActionsView = g.append("g")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight)
                .attr("id", "sequenceActionsView")
                .attr("translateX", 0)
                .attr("translateY", 0)
                .attr("transform", "translate(0,0) scale(1)")
                .attr("scale", 1);

            constant.sequenceWorkflowView = g.append("g")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight)
                .attr("id", "sequenceWorkflowView")
                .attr("translateX", 0)
                .attr("translateY", 0)
                .attr("transform", "translate(0,0) scale(1)")
                .attr("scale", 1);

            constant.sequenceActionLinkView = g.append("g")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight)
                .attr("id", "sequenceActionLinkView");

            constant.sequenceButtonView = g.append("g")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight)
                .attr("id", "buttonView");

            showSequenceView(constant.sequenceRunData);
            sequenceUtil.initButton();
        }
    });
}

function showSequenceView(workflowSequenceData) {
    constant.sequenceWorkflowView.selectAll("image").remove();
    constant.sequenceWorkflowView.selectAll("image")
        .data(workflowSequenceData)
        .enter()
        .append("image")
        .attr("xlink:href", function(d, i) {

            if (d.status == 1 || d.status == 0 ) {

                if (d.type == constant.WORKFLOW_END) {
                    return "../../assets/svg/history-end-waitStart.svg";
                }
                
                if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "stage" && constant.currentSelectedItem.data.id == d.id) {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-selected-waitStart.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-selected-waitStart.svg";
                    }
                } else {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-waitStart.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-waitStart.svg";
                    }
                } 
                   
            } else if (d.status == 2 ) {
                if (d.type == constant.WORKFLOW_END) {
                    return "../../assets/svg/history-end-waitStart.svg";
                }

                if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "stage" && constant.currentSelectedItem.data.id == d.id) {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-selected-running.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-selected-running.svg";
                    }
                } else {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-running.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-running.svg";
                    }
                } 

            } else if (d.status == 3) {

                if (d.type == constant.WORKFLOW_END) {
                    return "../../assets/svg/history-end-success.svg";
                }

                if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "stage" && constant.currentSelectedItem.data.id == d.id) {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-selected-success.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-selected-success.svg";
                    }
                } else {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-success.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-success.svg";
                    }
                }

            } else if(d.status == 4) {

                if (d.type == constant.WORKFLOW_END) {
                    return "../../assets/svg/history-end-fail.svg";
                }

                if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "stage" && constant.currentSelectedItem.data.id == d.id) {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-selected-fail.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-selected-fail.svg";
                    }
                } else {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-fail.svg";
                    } else if (d.type ==  constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-fail.svg";
                    }
                }
            } 
        })
        .attr("id", function(d, i) {
            return d.id;
        })
        .attr("data-index", function(d, i) {
            return i;
        })
        .attr("width", function(d, i) {
            return constant.svgStageWidth;
        })
        .attr("height", function(d, i) {
            return constant.svgStageHeight;
        })
        .attr("transform", function(d, i) {
            d.width = constant.svgStageWidth;
            d.height = constant.svgStageHeight;
            d.translateX = i * constant.WorkflowNodeSpaceSize + constant.workflowNodeStartX;
            d.translateY = constant.workflowNodeStartY;
            return "translate(" + d.translateX + "," + d.translateY + ")";
        })
        .attr("translateX", function(d, i) {
            return i * constant.WorkflowNodeSpaceSize + constant.workflowNodeStartX;
        })
        .attr("translateY", constant.workflowNodeStartY)
        .on("click", function(d, i) {
            if (d.status == 1 || d.status == 0) {

                if (d.type == constant.WORKFLOW_STAGE) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "stage", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-stage-selected-waitStart.svg");
                } else if (d.type == constant.WORKFLOW_START) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "start", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-start-selected-waitStart.svg");
                }

            } else if (d.status == 2) {

                if (d.type == constant.WORKFLOW_STAGE) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "stage", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-stage-selected-running.svg");
                } else if (d.type == constant.WORKFLOW_START) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "start", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-start-selected-running.svg");
                }

            } else if (d.status == 3) {

                if (d.type == constant.WORKFLOW_STAGE) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "stage", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-stage-selected-success.svg");
                } else if (d.type == constant.WORKFLOW_START) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "start", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-start-selected-success.svg");
                }
            }else if (d.status == 4) {

                if (d.type == constant.WORKFLOW_STAGE) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "stage", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-stage-selected-fail.svg");
                } else if (d.type == constant.WORKFLOW_START) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "start", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-start-selected-fail.svg");
                }
            }
        })

    initSequenceStageLine();
    if(constant.sequenceRunStatus == 1 || constant.sequenceRunStatus == 2){
        timerSequenceWorkflowData(historyAbout)
    }
    // initAction();
}

var timer ;
function timerSequenceWorkflowData(refreshSelect_hisotry){
    var promise = historyDataService.getWorkflowHistory(refreshSelect_hisotry.workflowName, refreshSelect_hisotry.versionName, refreshSelect_hisotry.sequence);
    promise.done(function(data) {
        loading.hide();
        constant.refreshSequenceRunData = data.define.stageList;
        constant.sequenceRunStatus = data.define.status; 

        if(constant.refreshSequenceRunData.length > 0){
            showRefreshSequenceView(constant.refreshSequenceRunData);
        }

        if(constant.sequenceRunStatus == 1 || constant.sequenceRunStatus == 2){
            timer = setTimeout(function(){timerSequenceWorkflowData(refreshSelect_hisotry);},5000);
        }else{
            clearTimeout(timer); 
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

function showRefreshSequenceView(refreshWorkflowSequenceData) {
    constant.sequenceWorkflowView.selectAll("image")
        .data(refreshWorkflowSequenceData)
        .attr("xlink:href", function(d, i) {

            if (d.status == 1 || d.status == 0 ) {

                if (d.type == constant.WORKFLOW_END) {
                    return "../../assets/svg/history-end-waitStart.svg";
                }
                
                if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "stage" && constant.currentSelectedItem.data.id == d.id) {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-selected-waitStart.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-selected-waitStart.svg";
                    }
                } else {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-waitStart.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-waitStart.svg";
                    }
                } 
                   
            } else if (d.status == 2 ) {

                if (d.type == constant.WORKFLOW_END) {
                    return "../../assets/svg/history-end-waitStart.svg";
                }

                if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "stage" && constant.currentSelectedItem.data.id == d.id) {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-selected-running.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-selected-running.svg";
                    }
                } else {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-running.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-running.svg";
                    }
                } 

            } else if (d.status == 3) {

                if (d.type == constant.WORKFLOW_END) {
                    return "../../assets/svg/history-end-success.svg";
                }

                if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "stage" && constant.currentSelectedItem.data.id == d.id) {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-selected-success.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-selected-success.svg";
                    }
                } else {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-success.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-success.svg";
                    }
                }

            } else if(d.status == 4) {

                if (d.type == constant.WORKFLOW_END) {
                    return "../../assets/svg/history-end-fail.svg";
                }

                if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "stage" && constant.currentSelectedItem.data.id == d.id) {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-selected-fail.svg";
                    } else if (d.type == constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-selected-fail.svg";
                    }
                } else {
                    if (d.type == constant.WORKFLOW_START) {
                        return "../../assets/svg/history-start-fail.svg";
                    } else if (d.type ==  constant.WORKFLOW_STAGE) {
                        return "../../assets/svg/history-stage-fail.svg";
                    }
                }
            } 
        })
        .attr("id", function(d, i) {
            return d.id;
        })
        .attr("data-index", function(d, i) {
            return i;
        })
        .attr("width", function(d, i) {
            return constant.svgStageWidth;
        })
        .attr("height", function(d, i) {
            return constant.svgStageHeight;
        })
        .attr("transform", function(d, i) {
            d.width = constant.svgStageWidth;
            d.height = constant.svgStageHeight;
            d.translateX = i * constant.WorkflowNodeSpaceSize + constant.workflowNodeStartX;
            d.translateY = constant.workflowNodeStartY;
            return "translate(" + d.translateX + "," + d.translateY + ")";
        })
        .attr("translateX", function(d, i) {
            return i * constant.WorkflowNodeSpaceSize + constant.workflowNodeStartX;
        })
        .attr("translateY", constant.workflowNodeStartY)
        .on("click", function(d, i) {
            if (d.status == 1 || d.status == 0) {

                if (d.type == constant.WORKFLOW_STAGE) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "stage", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-stage-selected-waitStart.svg");
                } else if (d.type == constant.WORKFLOW_START) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "start", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-start-selected-waitStart.svg");
                }

            } else if (d.status == 2) {

                if (d.type == constant.WORKFLOW_STAGE) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "stage", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-stage-selected-running.svg");
                } else if (d.type == constant.WORKFLOW_START) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "start", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-start-selected-running.svg");
                }

            } else if (d.status == 3) {

                if (d.type == constant.WORKFLOW_STAGE) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "stage", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-stage-selected-success.svg");
                } else if (d.type == constant.WORKFLOW_START) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "start", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-start-selected-success.svg");
                }
            }else if (d.status == 4) {

                if (d.type == constant.WORKFLOW_STAGE) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "stage", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-stage-selected-fail.svg");
                } else if (d.type == constant.WORKFLOW_START) {
                    historyChangeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": d, "type": "start", "status": d.status });
                    d3.select("#" + d.id).attr("href", "../../assets/svg/history-start-selected-fail.svg");
                }
            }
        })

    initSequenceStageLine();
    // initAction();
}

function initSequenceStageLine() {

    constant.sequenceLinesView.selectAll("g").remove();

    var diagonal = d3.svg.diagonal();

    var sequenceWorkflowLineViewId = "workflow-line-view";

    constant.sequenceLineView[sequenceWorkflowLineViewId] = constant.sequenceLinesView.append("g")
        .attr("width", constant.svgWidth)
        .attr("height", constant.svgHeight)
        .attr("id", sequenceWorkflowLineViewId);

    constant.sequenceWorkflowView.selectAll("image").each(function(d, i) {

        /* draw the main line of workflow */
        if (i != 0) {
            if (d.status == 0 || d.status == 1 || d.status ==2) {
                constant.sequenceLineView[sequenceWorkflowLineViewId]
                    .append("path")
                    .attr("d", function() {
                        return diagonal({
                            source: { x: d.translateX - constant.WorkflowNodeSpaceSize, y: constant.workflowNodeStartY + constant.svgStageHeight / 2 },
                            target: { x: d.translateX + 2, y: constant.workflowNodeStartY + constant.svgStageHeight / 2 }
                        });
                    })
                    .attr("fill", "none")
                    .attr("stroke", "#54711e")
                    .attr("stroke-width", 2);
            } else if (d.status == 3) {
                constant.sequenceLineView[sequenceWorkflowLineViewId]
                    .append("path")
                    .attr("d", function() {
                        return diagonal({
                            source: { x: d.translateX - constant.WorkflowNodeSpaceSize, y: constant.workflowNodeStartY + constant.svgStageHeight / 2 },
                            target: { x: d.translateX + 2, y: constant.workflowNodeStartY + constant.svgStageHeight / 2 }
                        });
                    })
                    .attr("fill", "none")
                    .attr("stroke", "#00733B")
                    .attr("stroke-width", 2);
            } else if(d.status == 4) {
                constant.sequenceLineView[sequenceWorkflowLineViewId]
                    .append("path")
                    .attr("d", function() {
                        return diagonal({
                            source: { x: d.translateX - constant.WorkflowNodeSpaceSize, y: constant.workflowNodeStartY + constant.svgStageHeight / 2 },
                            target: { x: d.translateX + 2, y: constant.workflowNodeStartY + constant.svgStageHeight / 2 }
                        });
                    })
                    .attr("fill", "none")
                    .attr("stroke", "#7E1101")
                    .attr("stroke-width", 2);
            }
        }

        if (d.type == constant.WORKFLOW_START) {
            /* draw the vertical line and circle for start node  in lineView -> workflow-line-view */
            constant.sequenceLineView[sequenceWorkflowLineViewId]
                .append("path")
                .attr("d", function() {
                    return diagonal({
                        source: { x: d.translateX + constant.svgStageWidth / 2, y: constant.workflowNodeStartY + constant.svgStageHeight / 2 },
                        target: { x: d.translateX + constant.svgStageWidth / 2, y: constant.workflowNodeStartY + constant.svgStageHeight + 10 }
                    })
                })
                .attr("fill", "none")
                .attr("stroke", "#1F6D84")
                .attr("stroke-width", 1);

            constant.sequenceLineView[sequenceWorkflowLineViewId]
                .append("circle")
                .attr("cx", function(cd, ci) {
                    return d.translateX + constant.svgStageWidth / 2;
                })
                .attr("cy", function(cd, ci) {
                    return constant.workflowNodeStartY + constant.svgStageHeight + 19;
                })
                .attr("r", function(cd, ci) {
                    return 8;
                })
                .attr("fill", "#fff")
                .attr("stroke", "#1F6D84")
                .attr("stroke-width", 2)
        }

    });

    initSequenceActionByStage();
    initSequenceAction2StageLine();
    initSequenceActionLinkBase();
    initSequenceActionLinkBasePoint();
    initSequencePath();
}

function initSequenceActionByStage() {
    constant.sequenceActionsView.selectAll("g").remove();
    /* draw actions in actionView , data source is stage.actions */
    constant.sequenceWorkflowView.selectAll("image").each(function(d, i) {
        if (d.type == constant.WORKFLOW_STAGE && d.actions != null && d.actions.length > 0) {

            var actionViewId = "action" + "-" + d.id;
            constant.sequenceActionView[actionViewId] = constant.sequenceActionsView.append("g")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight)
                .attr("id", actionViewId);

            var actionStartX = d.translateX + (constant.svgStageWidth - constant.svgActionWidth) / 2;
            var actionStartY = d.translateY;

            constant.sequenceActionView[actionViewId].selectAll("image")
                .data(d.actions).enter()
                .append("image")
                .attr("xlink:href", function(ad, ai) {
                    if (ad.status == 1 || ad.status == 0 ) {

                        if ( constant.currentSelectedItem != null && constant.currentSelectedItem.type == "action" && constant.currentSelectedItem.data.id == ad.id) {
                            return "../../assets/svg/history-action-selected-waitStart.svg";
                        }else if( isAction && actionId == ad.id){
                            return "../../assets/svg/history-action-selected-waitStart.svg";
                        }else {
                            return "../../assets/svg/history-action-waitStart.svg";
                        }

                    } else if (ad.status == 2) {

                        if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "action" && constant.currentSelectedItem.data.id == ad.id) {
                            return "../../assets/svg/history-action-selected-running.svg";
                        }else if( isAction && actionId == ad.id){
                            return "../../assets/svg/history-action-selected-running.svg";
                        }else {
                            return "../../assets/svg/history-action-running.svg";
                        }

                    } else if (ad.status == 3) {

                        if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "action" && constant.currentSelectedItem.data.id == ad.id) {
                            return "../../assets/svg/history-action-selected-success.svg";
                        }else if( isAction && actionId == ad.id){
                            return "../../assets/svg/history-action-selected-success.svg";
                        }else {
                            return "../../assets/svg/history-action-success.svg";
                        }

                    }else if (ad.status == 4) {

                        if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "action" && constant.currentSelectedItem.data.id == ad.id) {
                            return "../../assets/svg/history-action-selected-fail.svg";
                        }else if( isAction && actionId == ad.id){
                            return "../../assets/svg/history-action-selected-fail.svg";
                        }else {
                            return "../../assets/svg/history-action-fail.svg";
                        }

                    }
                })
                .attr("id", function(ad, ai) {
                    return ad.id;
                })
                .attr("data-index", function(ad, ai) {
                    return ai;
                })
                .attr("data-parent", i)
                .attr("width", function(ad, ai) {
                    return constant.svgActionWidth;
                })
                .attr("height", function(ad, ai) {
                    return constant.svgActionHeight;
                })
                .attr("translateX", actionStartX)
                .attr("translateY", function(ad, ai) {
                    /* draw difference distance between action and stage grouped by stage index */
                    if (i % 2 == 0) {
                        ad.translateY = actionStartY + constant.svgStageHeight - 55 + constant.ActionNodeSpaceSize * (ai + 1);
                    } else {
                        ad.translateY = actionStartY + constant.svgStageHeight - 10 + constant.ActionNodeSpaceSize * (ai + 1);
                    }
                    return ad.translateY;
                })
                .attr("transform", function(ad, ai) {
                    ad.width = constant.svgActionWidth;
                    ad.height = constant.svgActionHeight;
                    if (i % 2 == 0) {
                        ad.translateX = actionStartX;
                        ad.translateY = actionStartY + constant.svgStageHeight - 55 + constant.ActionNodeSpaceSize * (ai + 1);
                    } else {
                        ad.translateX = actionStartX;
                        ad.translateY = actionStartY + constant.svgStageHeight - 10 + constant.ActionNodeSpaceSize * (ai + 1);
                    }

                    return "translate(" + ad.translateX + "," + ad.translateY + ")";
                })
                .style("cursor", "pointer")
                .on("click", function(ad, ai) {
                    // workflowName,versionName,workflowRunSequence,stageName,actionName
                    getActionHistory(historyAbout.workflowName,historyAbout.versionName,historyAbout.sequence, d.setupData.name, ad.setupData.name);
                    getContainerLogs(historyAbout.workflowName,historyAbout.versionName,historyAbout.sequence, d.setupData.name, ad.setupData.name,"");
                    if (ad.status == 1 || ad.status == 0) {
                        historyChangeCurrentElement(constant.currentSelectedItem);
                        constant.setCurrentSelectedItem({ "data": ad, "parentData": d, "type": "action", "status": ad.status });
                        d3.select("#" + ad.id).attr("href", "../../assets/svg/history-action-selected-waitStart.svg");
                    } else if (ad.status == 2) {
                        historyChangeCurrentElement(constant.currentSelectedItem);
                        constant.setCurrentSelectedItem({ "data": ad, "parentData": d, "type": "action", "status": ad.status });
                        d3.select("#" + ad.id).attr("href", "../../assets/svg/history-action-selected-running.svg");
                    } else if (ad.status == 3) {
                        historyChangeCurrentElement(constant.currentSelectedItem);
                        constant.setCurrentSelectedItem({ "data": ad, "parentData": d, "type": "action", "status": ad.status });
                        d3.select("#" + ad.id).attr("href", "../../assets/svg/history-action-selected-success.svg");

                    }else if (ad.status == 4) {
                        historyChangeCurrentElement(constant.currentSelectedItem);
                        constant.setCurrentSelectedItem({ "data": ad, "parentData": d, "type": "action", "status": ad.status });
                        d3.select("#" + ad.id).attr("href", "../../assets/svg/history-action-selected-fail.svg");

                    }
                })
                .on("mouseout", function(ad, ai) {
                    constant.sequenceWorkflowView.selectAll("#workflow-element-popup").remove();
                })
                .on("mouseover", function(ad, ai) {
                    var x = ad.translateX;
                    var y = ad.translateY + constant.svgActionHeight;
                    let text = "";
                    let width = 150;
                    let options = {};
                    if (ad.setupData && ad.setupData.name && ad.setupData.name != "") {
                        text = ad.setupData.name;
                        width = text.length * 12 + 20;
                        options = {
                            "x": x,
                            "y": y,
                            "text": text,
                            "popupId": "workflow-element-popup",
                            "parentView": constant.sequenceWorkflowView,
                            "width": width
                        };
                        sequenceUtil.showToolTip(options);
                    }

                })

        }

    });
}

function initSequenceAction2StageLine() {
    var diagonal = d3.svg.diagonal();

    constant.sequenceWorkflowView.selectAll("image").each(function(d, i) {
        /* draw line from action 2 stage and circle of action self to accept and emit lines  */
        if (d.type == constant.WORKFLOW_STAGE && d.actions != null && d.actions.length > 0) {

            var actionLineViewId = "action-line" + "-" + d.id;
            var action2StageLineViewId = "action-2-stage-line" + "-" + d.id;
            var actionSelfLine = "action-self-line" + "-" + d.id
                /* Action 2 Stage */
            constant.sequenceLineView[action2StageLineViewId] = constant.sequenceLinesView.append("g")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight)
                .attr("id", action2StageLineViewId);

            constant.sequenceLineView[action2StageLineViewId].selectAll("path")
                .data(d.actions).enter()
                .append("path")
                .attr("d", function(ad, ai) {
                    /* draw the tail line of action */
                    constant.sequenceLineView[action2StageLineViewId]
                        .append("path")
                        .attr("d", function(fd, fi) {
                            return diagonal({
                                source: { x: ad.translateX + constant.svgActionWidth / 2, y: ad.translateY + constant.svgActionHeight },
                                target: { x: ad.translateX + constant.svgActionWidth / 2, y: ad.translateY + constant.svgActionHeight + 8 }
                            });
                        })
                        .attr("fill", "none")
                        .attr("stroke", "#1F6D84")
                        .attr("stroke-width", 1)
                        .attr("stroke-dasharray", "2,2");
                    /* draw different length line group by stage index */
                    if (i % 2 == 0) {
                        return diagonal({
                            source: { x: ad.translateX + constant.svgActionWidth / 2, y: ad.translateY },
                            target: { x: ad.translateX + constant.svgActionWidth / 2, y: ad.translateY - 44 }
                        });
                    } else {
                        return diagonal({
                            source: { x: ad.translateX + constant.svgActionWidth / 2, y: ad.translateY },
                            target: { x: ad.translateX + constant.svgActionWidth / 2, y: ad.translateY - 68 }
                        });
                    }
                })
                .attr("fill", "none")
                .attr("stroke", "#1F6D84")
                .attr("stroke-width", 1)
                .attr("stroke-dasharray", "2,2");
        }
    });
}

function initSequenceActionLinkBase() {
    var diagonal = d3.svg.diagonal();

    constant.sequenceWorkflowView.selectAll("image").each(function(d, i) {
        if (d.type == constant.WORKFLOW_STAGE && d.actions != null && d.actions.length > 0) {

            var actionSelfLine = "action-self-line" + "-" + d.id

            /* line across action to connect two circles */
            constant.sequenceLineView[actionSelfLine] = constant.sequenceLinesView.append("g")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight)
                .attr("id", actionSelfLine);

            constant.sequenceLineView[actionSelfLine].selectAll("path")
                .data(d.actions).enter()
                .append("path")
                .attr("d", function(ad, ai) {
                    return diagonal({
                        source: { x: ad.translateX - 8, y: ad.translateY + constant.svgActionHeight / 2 },
                        target: { x: ad.translateX + constant.svgActionWidth + 8, y: ad.translateY + constant.svgActionHeight / 2 }
                    })
                })
                .attr("id", function(ad, ai) {
                    return "action-self-line-path-" + ad.id;
                })
                .attr("fill", "none")
                .attr("stroke", "#1F6D84")
                .attr("stroke-width", 1);
        }
    });
}

function initSequenceActionLinkBasePoint() {
    var diagonal = d3.svg.diagonal();

    constant.sequenceWorkflowView.selectAll("image").each(function(d, i) {
        if (d.type == constant.WORKFLOW_STAGE && d.actions != null && d.actions.length > 0) {

            var actionSelfLine = "action-self-line" + "-" + d.id

            /* circle on the left */
            constant.sequenceLineView[actionSelfLine].selectAll(".action-self-line-input")
                .data(d.actions).enter()
                .append("circle")
                .attr("class", "action-self-line-input")
                .attr("cx", function(ad, ai) {
                    return ad.translateX - 16;
                })
                .attr("cy", function(ad, ai) {
                    return ad.translateY + constant.svgActionHeight / 2;
                })
                .attr("r", function(ad, ai) {
                    return 8;
                })
                .attr("id", function(ad, ai) {
                    return "action-self-line-input-" + ad.id;
                })
                .attr("fill", "#fff")
                .attr("stroke", "#84C1BC")
                .attr("stroke-width", 2)
                .style("cursor", "pointer")

            /* circle on the right */
            constant.sequenceLineView[actionSelfLine].selectAll(".action-self-line-output")
                .data(d.actions).enter()
                .append("circle")
                .attr("class", "action-self-line-output")
                .attr("cx", function(ad, ai) {
                    return ad.translateX + constant.svgActionWidth + 16;
                })
                .attr("cy", function(ad, ai) {
                    return ad.translateY + constant.svgActionHeight / 2;
                })
                .attr("r", function(ad, ai) {
                    return 8;
                })
                .attr("id", function(ad, ai) {
                    return "action-self-line-output-" + ad.id
                })
                .attr("fill", "#fff")
                .attr("stroke", "#84C1BC")
                .attr("stroke-width", 2)
                .style("cursor", "pointer")
        }
    });
}

function initSequencePath() {
    constant.sequenceLinePathArray.forEach(function(i) {
        setSequencePath(i)
    });
}

function setSequencePath(options) {
    var fromDom = $("#" + options.startData.id)[0].__data__;
    var toDom = $("#" + options.endData.id)[0].__data__;
    var lineId = options.id;
    /* line start point(x,y) is the circle(x,y) */

    var startPoint = {},
        endPoint = {};
    if (fromDom.type == constant.WORKFLOW_START) {
        startPoint = { x: fromDom.translateX + 1, y: fromDom.translateY + 57 };
    } else if (fromDom.type == constant.WORKFLOW_ACTION) {
        startPoint = { x: fromDom.translateX + 19, y: fromDom.translateY + 4 };
    }
    endPoint = { x: toDom.translateX - 12, y: toDom.translateY + 4 };
    constant.sequenceLineView[options.workflowLineViewId]
        .append("path")
        .attr("d", getPathData(startPoint, endPoint))
        .attr("fill", "none")
        .attr("stroke-opacity", "1")
        .attr("stroke", function(d, i) {

            if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "line" && constant.currentSelectedItem.data.attr("id") == options.id) {
                return "#81D9EC";
            } else {
                return "#E6F3E9";
            }
        })
        .attr("stroke-width", 10)
        .attr("data-index", options.index)
        .attr("id", options.id)
        .style("cursor", "pointer")
        .on("click", function(d) {
            getLineHistory(historyAbout.workflowName,historyAbout.versionName,historyAbout.sequence,lineId );
            // getLineHistory(historyAbout.workflowName, historyAbout.sequenceID, options.startData.id, options.endData.id);
            var self = $(this);
            historyChangeCurrentElement(constant.currentSelectedItem);
            constant.setCurrentSelectedItem({ "data": self, "type": "line" });
            d3.select(this).attr("stroke", "#81D9EC");
        });
}

function getPathData(startPoint, endPoint) {
    var curvature = .5;
    var x0 = startPoint.x + 30,
        x1 = endPoint.x + 2,
        xi = d3.interpolateNumber(x0, x1),
        x2 = xi(curvature),
        x3 = xi(1 - curvature),
        y0 = startPoint.y + 30 / 2,
        y1 = endPoint.y + 30 / 2;

    return "M" + x0 + "," + y0 + "C" + x2 + "," + y0 + " " + x3 + "," + y1 + " " + x1 + "," + y1;
}

function zoomed() {
    constant.sequenceWorkflowView.attr("transform", "translate(" + d3.event.translate + ") scale(" + d3.event.scale + ")");
    constant.sequenceActionsView.attr("transform", "translate(" + d3.event.translate + ") scale(" + d3.event.scale + ")");
    // buttonView.attr("transform", "translate(" + d3.event.translate + ") scale(" + d3.event.scale + ")");
    constant.sequenceLinesView.attr("transform", "translate(" + d3.event.translate + ") scale(" + d3.event.scale + ")")
        .attr("translateX", d3.event.translate[0])
        .attr("translateY", d3.event.translate[1])
        .attr("scale", d3.event.scale);
}

function clicked(d, i) {
    // constant.buttonView.selectAll("image").remove();
    if (d3.event.defaultPrevented) return; // zoomed
    d3.select(this).transition()
        .transition()
}

function nozoom() {
    d3.event.preventDefault();
}

function historyChangeCurrentElement(previousData) {
    if (previousData != null) {

        if (previousData.status == 3 || previousData.type == "line") {

            switch (previousData.type) {
                case "stage":
                    d3.select("#" + previousData.data.id).attr("href", "../../assets/svg/history-stage-success.svg");
                    break;
                case "start":
                    d3.select("#" + previousData.data.id).attr("href", "../../assets/svg/history-start-success.svg");
                    break;
                case "action":
                    d3.select("#" + previousData.data.id).attr("href", "../../assets/svg/history-action-success.svg");
                    break;
                case "line":
                    d3.select("#" + previousData.data.attr("id")).attr("stroke", "#00733B");
                    break;


            }
        }
    }

    if (previousData != null) {

        if (previousData.status == 4 || previousData.type == "line") {

            switch (previousData.type) {
                case "stage":
                    d3.select("#" + previousData.data.id).attr("href", "../../assets/svg/history-stage-fail.svg");
                    break;
                case "start":
                    d3.select("#" + previousData.data.id).attr("href", "../../assets/svg/history-start-fail.svg");
                    break;
                case "action":
                    d3.select("#" + previousData.data.id).attr("href", "../../assets/svg/history-action-fail.svg");
                    break;
                case "line":
                    d3.select("#" + previousData.data.attr("id")).attr("stroke", "#E6F3E9");
                    break;
            }
        }
    }
}
