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

import { setPath, getPathData } from "./setPath";
import { linePathAry } from "../common/constant";
import { bipatiteView } from "./bipatiteView";
import { notify } from "../common/notify";
import * as conflictUtil from "../relation/conflict";

export function dragDropSetPath(options) {
   
    var fromNodeData = options.data; /* from node data */
    // var _path = d3.select("svg>g").insert("path", ":nth-child(2)").attr("class", "drag-drop-line"),
    // var _path = d3.select("svg>g").append("path").attr("class", "drag-drop-line"),
    var _path = d3.select("#workflow-line-view").append("path").attr("class", "drag-drop-line");
    var svgDOM = $("#div-d3-main-svg > svg"),
        svgOffsetX = svgDOM.offset().left,
        svgOffsetY = svgDOM.offset().top,
        _startX = $(window.event.target).offset().left - svgOffsetX,
        _startY = $(window.event.target).offset().top - svgOffsetY;

    /* draw temporary line by mouse move*/
    document.onmousemove = function(e) {
            var diffX = e.pageX - svgOffsetX,
                diffY = e.pageY - svgOffsetY;
            var translateX = parseInt(d3.select("#linesView")[0][0].attributes["translateX"] ? d3.select("#linesView")[0][0].attributes["translateX"].value : 0);
            var translateY = parseInt(d3.select("#linesView")[0][0].attributes["translateY"] ? d3.select("#linesView")[0][0].attributes["translateY"].value : 0);
            var scale = parseFloat(d3.select("#linesView")[0][0].attributes["scale"] ? d3.select("#linesView")[0][0].attributes["scale"].value : 1);
            _path.attr("d", getPathData({ x: (_startX - 15 * scale - translateX) / scale, y: (_startY - 8 * scale - translateY) / scale}, { x: (diffX - translateX) / scale, y: (diffY - 10 * scale - translateY)/scale}))
                .attr("fill", "none")
                .attr("stroke-opacity", "1")
                .attr("stroke", "#81D9EC")
                .attr("stroke-width", 10)
                .style("cursor","pointer");
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

        if (toNodeData != undefined && toNodeData.translateX > fromNodeData.translateX && toNodeData.type === "workflow-action") {
            
            let dataJson = {
                workflowLineViewId: "workflow-line-view",
                startData: fromNodeData,
                endData: toNodeData,
                startPoint: { x: fromNodeData.translateX, y: fromNodeData.translateY },
                endPoint: { x: toNodeData.translateX, y: toNodeData.translateY },
                id: _id
            };

            setPath(dataJson);
            linePathAry.push(dataJson);

            bipatiteView(fromNodeData.outputJson,toNodeData.inputJson,dataJson);

            if(conflictUtil.hasConflict(fromNodeData.id, toNodeData.id)){
                notify("Conflict with other inputs, please click target action to resolve conflict first", "error",10);
                return false;
            }
        }
    }
}
