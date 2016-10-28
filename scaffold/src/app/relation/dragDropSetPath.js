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

import { setPath, getPathData } from "./setPath";
import { linePathAry } from "../common/constant";
import { bipatiteView } from "./bipatiteView";
import { notify } from "../common/notify";

export function dragDropSetPath(options) {

    var fromNodeData = options.data; /* from node data */
    // var _path = d3.select("svg>g").insert("path", ":nth-child(2)").attr("class", "drag-drop-line"),
    // var _path = d3.select("svg>g").append("path").attr("class", "drag-drop-line"),
    var _path = d3.select("#pipeline-line-view").append("path").attr("class", "drag-drop-line"),
        _offsetX = $("main").offset().left,
        _offsetY = $("#designerMenubar").height(),
        _startX = $(window.event.target).offset().left - _offsetX,
        _startY = $(window.event.target).offset().top - _offsetY - 12,
        _pageTitleHeight = $(".page-title").height();

    /* draw temporary line by mouse move*/
    document.onmousemove = function(e) {
            var diffX = e.pageX - _startX - _offsetX,
                diffY = e.pageY - _startY - _offsetY;
            _path.attr("d", getPathData({ x: _startX - 60, y: _startY - (105 + _pageTitleHeight) }, { x: _startX + diffX - 40, y: _startY + diffY - (130 + _pageTitleHeight) }))
                .attr("fill", "none")
                .attr("stroke-opacity", "1")
                .attr("stroke", "#81D9EC")
                .attr("stroke-width", 10);
        }
        /* remove temporary line and draw the real line between nodes with data */
    document.onmouseup = function(e) {

        document.onmousemove = null;
        document.onmouseup = null;
        d3.select(".drag-drop-line").remove();

        try {
            var toNodeData = d3.select(e.target)[0][0].__data__; /* target node(action) data */
            var _id = fromNodeData.id + "-" + toNodeData.id; /* id is set to from data id add target id */
            if (d3.selectAll("#" + _id)[0].length > 0) {
                notify("Duplicate addition is prohibited", "error");
                return false;
            }
        } catch (e) {

        }

        if (toNodeData != undefined && toNodeData.translateX > fromNodeData.translateX && toNodeData.type === "pipeline-action") {
            
            let dataJson = {
                pipelineLineViewId: "pipeline-line-view",
                startData: fromNodeData,
                endData: toNodeData,
                startPoint: { x: fromNodeData.translateX, y: fromNodeData.translateY },
                endPoint: { x: toNodeData.translateX, y: toNodeData.translateY },
                id: _id
            };

            setPath(dataJson);
            linePathAry.push(dataJson);

            bipatiteView(fromNodeData.outputJson,toNodeData.inputJson,dataJson);

            // if(checkConflict(fromNodeData.id, toNodeData.id)){
            //     notify("Conflict with other inputs, please click target action to resolve conflict first", "error");
            //     return false;
            // }
        }
    }
}
