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

let linesView, actionsView, workflowView, buttonView;

export function initDesigner() {
    let $div = $("#div-d3-main-svg").height($("main").height() * 2 / 3);
    // let zoom = d3.behavior.zoom().on("zoom", zoomed);
    let drag = d3.behavior.drag()
        .origin(function() {
            return { "x": 0, "y": 0 };
        })
        .on("dragstart", dragStart)
        .on("drag", util.draged);

    function dragStart() {
        d3.event.sourceEvent.stopPropagation();
        drag.origin(function() {
            return { "x": constant.workflowView.attr("translateX"), "y": constant.workflowView.attr("translateY") }
        });
    }
    constant.setSvgWidth("100%");
    constant.setSvgHeight("100%");
    constant.setWorkflowNodeStartX(50);
    constant.setWorkflowNodeStartY(($div.height() + 2 * constant.buttonVerticalSpace + constant.buttonHeight) * 0.2);

    $div.empty();

    let svg = d3.select("#div-d3-main-svg")
        .on("touchstart", nozoom)
        .on("touchmove", nozoom)
        .append("svg")
        .attr("width", constant.svgWidth)
        .attr("height", constant.svgHeight)
        .style("fill", "white");

    let g = svg.append("g")
        // .call(zoom)
        .call(drag)
        // .on("dblclick.zoom", null)
        // .on("wheel.zoom", null);

    let svgMainRect = g.append("rect")
        .attr("width", constant.svgWidth)
        .attr("height", constant.svgHeight)
        .on("click", clicked);

    linesView = g.append("g")
        .attr("width", constant.svgWidth)
        .attr("height", constant.svgHeight)
        .attr("id", "linesView")
        .attr("translateX", 0)
        .attr("translateY", 0)
        .attr("transform", "translate(0,0) scale(1)")
        .attr("scale", 1);

    actionsView = g.append("g")
        .attr("width", constant.svgWidth)
        .attr("height", constant.svgHeight)
        .attr("id", "actionsView")
        .attr("translateX", 0)
        .attr("translateY", 0)
        .attr("transform", "translate(0,0) scale(1)")
        .attr("scale", 1);

    workflowView = g.append("g")
        .attr("width", constant.svgWidth)
        .attr("height", constant.svgHeight)
        .attr("id", "workflowView")
        .attr("translateX", 0)
        .attr("translateY", 0)
        .attr("transform", "translate(0,0) scale(1)")
        .attr("scale", 1);

    buttonView = g.append("g")
        .attr("width", constant.svgWidth)
        .attr("height", constant.svgHeight)
        .attr("id", "buttonView");


    // let actionLinkView = g.append("g")
    //     .attr("width", constant.svgWidth)
    //     .attr("height", constant.svgHeight)
    //     .attr("id", "actionLinkView");

    constant.setSvg(svg);
    constant.setG(g);
    constant.setSvgMainRect(svgMainRect);
    constant.setLinesView(linesView);
    constant.setActionsView(actionsView);
    constant.setWorkflowView(workflowView);
    constant.setButtonView(buttonView);
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
