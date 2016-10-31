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

import * as historyDataService from "./historyData";
import * as constant  from "../common/constant";
import { loading } from "../common/loading";
import * as historyMain from "./main";

export function initHistoryList() {
    loading.show();
    var promise = historyDataService.sequenceList();
    promise.done(function(data) {
        loading.hide();
        constant.sequenceAllList = data.pipelineList;
        
        if (constant.sequenceAllList.length > 0) {
            getHistoryList(constant.sequenceAllList);
        } else {
            notify("Server is unreachable", "error");
        }
    });
    promise.fail(function(xhr, status, error) {
        loading.hide();
        if (xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg, "error");
        } else {
            notify("Server is unreachable", "error");
        }
    });
}



function getHistoryList (pipelineData) {
    $.ajax({
        url: "../../templates/history/historyList.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#history-pipeline-list").html($(data));
            $("#historyPipelinelist").show();
            $(".pipelinelist_body").empty();

            var hppItem = $(".pipelinelist_body");

            _.each(pipelineData, function (pd){
                var hpRow = `<tr data-id="pd`+ pd.id +`" class="pp-row">
                                <td class="pptd">
                                    <span class="glyphicon glyphicon-menu-down treeclose treecontroller" data-name=` 
                                    + pd.name + `></span>&nbsp;&nbsp;&nbsp;&nbsp;` 
                                    + pd.name + `</td><td></td><td></td>
                                    <td data-btnId="pd`+ pd.id+`"></td></tr>`;
                hppItem.append(hpRow);

                _.each(pd.versionList, function(vd){
                    var hvRow =`<tr data-id="vd`+ vd.id +` data-pname=` + vd.name + ` data-version=` + vd.id + ` data-versionid=`+ vd.id + ` class="ppversion-row">
                                    <td></td>`;

                    if(_.isUndefined(vd.status) && vd.sequenceList.length <= 0){

                        hvRow += `<td class="pptd">`+ vd.name + `</td>
                                    <td><div class="state-list"><div class="state-icon-list state-norun"></div></div></td><td></td>`;

                        hppItem.append(hvRow);

                    } else {

                        hvRow += `<td class="pptd"><span class="glyphicon glyphicon-menu-down treeclose treecontroller" data-name=` 
                        + vd.name + `></span>&nbsp;&nbsp;&nbsp;&nbsp;` + vd.name + `</td>`;

                        // version is close  
                        // if( isTreeOpen == false ){

                        //     if(vd.status == true){
                        //         hvRow += `<td><div class="state-list">
                        //                 <div class="state-icon-list state-success"></div>
                        //                 <span class="state-label-list">` + vd.time + `</span>
                        //             </div></td>`;
                        //     } else {
                        //         hvRow += `<td><div class="state-list">
                        //                 <div class="state-icon-list state-fail"></div>
                        //                 <span class="state-label-list">` + vd.time + `</span>
                        //             </div></td>`
                        //     }

                        // } else {

                            hvRow += `<td class="pptd">`+ vd.info +`</td>`;
                           
                        // }


                        hvRow += `<td data-btnId="vd`+ vd.id +`"></td></tr> `

                        hppItem.append(hvRow);

                        if( vd.sequenceList.length > 0){
                            _.each(vd.sequenceList, function(sd){
                                var hsRow =`<tr data-id="sd`+ sd.pipelineSequenceID +` class="ppversion-row"><td></td><td></td>`;

                                if(sd.status == true){
                                    hsRow += `<td><div class="state-list">
                                            <div class="state-icon-list state-success"></div>
                                            <span class="state-label-list">` + sd.time + `</span>
                                        </div></td>`;
                                } else {
                                    hsRow += `<td><div class="state-list">
                                            <div class="state-icon-list state-fail"></div>
                                            <span class="state-label-list">` + sd.time + `</span>
                                        </div></td>`
                                }
                                
                                hsRow += `<td  data-name=`  + pd.name + ` data-btnId=`+ sd.pipelineSequenceID +`><button type="button" class="btn btn-success ppview"><i class="glyphicon glyphicon-list-alt" style="font-size:16px"></i>&nbsp;&nbsp;detail</button></td></tr> `

                                hppItem.append(hsRow)
                            });
                        }
                    }
                });
            });

                $(".ppview").on("click",function(event){
                    var pname = $(event.currentTarget).parent().data("name");
                    var sid = $(event.currentTarget).parent().data("btnid");
                    historyMain.initHistoryPage(pname,sid)
                })
        }
    });
}



