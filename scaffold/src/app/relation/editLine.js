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

import { jsonEditor } from "../../vendor/jquery.jsoneditor";
import { bipatiteView } from "./bipatiteView";
import { resizeWidget } from "../theme/widget";
import { pipelineData } from "../pipeline/main";
import * as constant from "../common/constant";


export var lineInputJSON = {};

export var lineOutputJSON = {};

export function editLine(editPage, currentLine) {

    var index = currentLine.attr("data-index");
    var id = currentLine.attr("id");
    $("#pipeline-info-edit").html($(editPage));

    $("#removeLink").click(function() {
        currentLine.remove();
        constant.linePathAry.splice(index, 1);
        $("#pipeline-info-edit").html("");
    })

    $("#importDiv").html("");
    $("#outputDiv").html("");
    var currentLineData = _.find(constant.linePathAry, function(line) {
        return id == line.id;
    })
    lineInputJSON = currentLineData.startData.outputJson;
    lineOutputJSON = currentLineData.endData.inputJson;
    if (_.isEmpty(lineInputJSON)) {
        $("#importDiv").html("no data");
    }
    if (_.isEmpty(lineOutputJSON)) {
        $("#outputDiv").html("no data");
    }
    bipatiteView(lineInputJSON, lineOutputJSON, constant.linePathAry[index]);


    $("#afreshRelation").click(function() {
        constant.linePathAry[index].relation = undefined;
        bipatiteView(lineInputJSON, lineOutputJSON, constant.linePathAry[index]);
        
    })


    resizeWidget();
}
