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

import * as constant from "../common/constant";
import * as util from "../common/util";
import { makeFrontLayer, makeBackLayer, clickWorkflowLine } from "./lineHover";

/* draw real line between from node and to node */
export function setPath(options) {
    var fromDom = $("#" + options.startData.id)[0].__data__;
    var toDom = $("#" + options.endData.id)[0].__data__;

    /* line start point(x,y) is the circle(x,y) */
    var startPoint = {},
        endPoint = {};
    if (fromDom.type == constant.WORKFLOW_START) {
        startPoint = { x: fromDom.translateX + 1, y: fromDom.translateY + 57 };
    } else if (fromDom.type == constant.WORKFLOW_ACTION) {
        startPoint = { x: fromDom.translateX + 19, y: fromDom.translateY + 4 };
    }
    endPoint = { x: toDom.translateX - 12, y: toDom.translateY + 4 };
    
    constant.lineView[options.workflowLineViewId]
        .append("path")
        .attr("d", getPathData(startPoint, endPoint))
        .attr("fill", "none")
        .attr("stroke-opacity", "1")
        .attr("stroke", function(d, i) {
            if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "line" && constant.currentSelectedItem.data.attr("id") == options.id) {
                makeFrontLayer(this);
                return "#81D9EC";
            } else {
                makeBackLayer(this);
                return "#E6F3E9";
            }
        })
        .attr("stroke-width", 10)
        .attr("data-index", options.index)
        .attr("id", options.id)
        .style("cursor", "pointer")
        .on("click", function(d) {
            clickWorkflowLine(this);
        })
}

export function getPathData(startPoint, endPoint) {
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
