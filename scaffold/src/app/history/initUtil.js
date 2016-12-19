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

import * as constant from "../common/constant";
import * as sequenceUtil from "./initUtil";

export function draged(d) {
    if (d && d.name && d.name == "conflictTree") {
        var scale = Number(d3.select(this).attr("scale"));
        d3.select(this).attr("translateX", d3.event.x)
            .attr("translateY", d3.event.y)
            .attr("scale", scale)
            .attr("transform", "translate(" + d3.event.x + "," + d3.event.y + ") scale(" + scale + ")");

    } else {
        var scale = Number(constant.sequenceWorkflowView.attr("scale"));
        var translate = "translate(" + (d3.event.x) + "," + (d3.event.y) + ") scale(" + scale + ")";
        var targetCollection = [constant.sequenceWorkflowView, constant.sequenceActionsView, constant.sequenceLinesView];
        _.each(targetCollection, function(target) {
            target
                .attr("translateX", d3.event.x)
                .attr("translateY", d3.event.y)
                .attr("transform", translate)
                .attr("scale", scale)
        })
    }
}

let zoom = d3.behavior.zoom()
    .on("zoom", redraw)
    .scaleExtent([constant.zoomMinimum, constant.zoomMaximum]);

function redraw(d) {
    if (d && d.name && d.name == "conflictTree") {
        d3.select(this)
            .attr("translateX", d3.event.translate[0])
            .attr("translateY", d3.event.translate[1])
            .attr("transform", "translate(" + d3.event.translate + ")" + " scale(" + d3.event.scale + ")")
            .attr("scale", d3.event.scale)
    } else {
        var targetCollection = [constant.sequenceWorkflowView, constant.sequenceActionsView, constant.sequenceLinesView];
        _.each(targetCollection, function(target) {
            target
                .attr("translateX", d3.event.translate[0])
                .attr("translateY", d3.event.translate[1])
                .attr("transform", "translate(" + d3.event.translate + ")" + " scale(" + d3.event.scale + ")")
                .attr("scale", d3.event.scale)
        })
    }

}
export function zoomed(type, target, scaleObj) {
    var currentTranslateX = Number(target.attr("translateX"));
    var currentTranslateY = Number(target.attr("translateY"));
    var currentTranslate = [currentTranslateX, currentTranslateY];
    // zoom.scale(scale).translate(currentTranslate).event(constant.workflowView);
    d3.transition().duration(constant.zoomDuration).tween("zoom", function() {
        if (type == "zoomin") {
            scaleObj.zoomTargetScale = (scaleObj.zoomScale + constant.zoomFactor) <= constant.zoomMaximum ? (scaleObj.zoomScale + constant.zoomFactor) : constant.zoomMaximum;
        } else if (type == "zoomout") {
            scaleObj.zoomTargetScale = (scaleObj.zoomScale - constant.zoomFactor) >= constant.zoomMinimum ? (scaleObj.zoomScale - constant.zoomFactor) : constant.zoomMinimum;
        }
        var interpolate_scale = d3.interpolate(scaleObj.zoomScale, scaleObj.zoomTargetScale),
            interpolate_trans = d3.interpolate(currentTranslate, currentTranslate);
        return function(t) {
            zoom.scale(interpolate_scale(t))
                .translate(interpolate_trans(t))
                .event(target);
            scaleObj.zoomScale = scaleObj.zoomTargetScale;
        };
    });
}
export function showZoomBtn(index, type, containerView, target, scaleObj, options) {
    options = options || {};
    var horizonSpace = options.horizonSpace || constant.buttonHorizonSpace;
    var verticalSpace = options.verticalSpace || constant.buttonVerticalSpace;
    var backgroundY = options.backgroundY || constant.rectBackgroundY;

    containerView
        .append("image")
        .attr("xlink:href", function(ad, ai) {
            if (type == "zoomin") {
                return "../../assets/svg/zoomin.svg";
            } else if (type == "zoomout") {
                return "../../assets/svg/zoomout.svg";
            }

        })
        .attr("translateX", function(d, i) {
            return index * horizonSpace + (index - 1) * constant.buttonWidth;
        })
        .attr("translateY", function(d, i) {

            return verticalSpace + backgroundY;
        })
        .attr("transform", function(d, i) {
            let translateX = d3.select(this).attr("translateX");
            let translateY = d3.select(this).attr("translateY");
            return "translate(" + translateX + "," + translateY + ")";
        })
        .attr("width", constant.buttonWidth)
        .attr("height", constant.buttonHeight)
        .style("cursor", "pointer")
        .on("mouseover", function(d, i) {
            let content = "";
            let href = "";
            if (type == "zoomin") {
                content = "Zoomin";
                href = "../../assets/svg/zoomin.svg";
            } else if (type == "zoomout") {
                content = "Zoomout";
                href = "../../assets/svg/zoomout.svg";
            }
            d3.select(this).attr("href", href);
            let options = {
                "x": Number(d3.select(this).attr("translateX")),
                "y": Number(d3.select(this).attr("translateY")) + constant.buttonHeight,
                "text": content,
                "popupId": "button-element-popup",
                "parentView": containerView
            };
            showToolTip(options);
        })
        .on("mouseout", function(d, i) {
            cleanToolTip(containerView, "#button-element-popup");
        })
        .on("click", function(d, i) {
            zoomed(type, target, scaleObj);
        })
}
export function showToolTip(options) {
    var x = options.x,
        y = options.y,
        text = options.text,
        popupId = options.popupId,
        parentView = options.parentView,
        width = options.width || constant.popupWidth,
        height = options.height || constant.popupHeight;

    parentView
        .append("g")
        .attr("id", popupId);
    parentView.selectAll("#" + popupId)
        .append("rect")
        .attr("width", width)
        .attr("height", height)
        .attr("x", function(pd, pi) {
            return x;
        })
        .attr("y", function(pd, pi) {
            return y;
        })
        .attr("rx", 3)
        .attr("ry", 3)
        .style("fill", constant.toolTipBackground)
        .style("opacity", 0.9)
    parentView.selectAll("#" + popupId)
        .append("text")
        .attr("x", x + 10)
        .attr("y", y + height / 2 + 4)
        .style("fill", "white")
        .style("opacity", 0.9)
        .text(text)
}

export function cleanToolTip(containerView, id) {
    containerView.selectAll(id).remove();
}

 let rectBackgroundY = 15;
export function initButton() {
    let scaleObj = { "zoomScale": 1, "zoomTargetScale": 1 };
    constant.sequenceButtonView
        .append("rect")
        .attr("width", constant.svgWidth)
        .attr("height", rectBackgroundY)
        .style({
            "fill": "#ffffff"
        });
    constant.sequenceButtonView
        .append("rect")
        .attr("width", constant.svgWidth)
        .attr("height", 2 * constant.buttonVerticalSpace + constant.buttonHeight)
        .attr("y", rectBackgroundY)
        .style({
            "fill": "#f7f7f7"
        });
    sequenceUtil.showZoomBtn(1, "zoomin", constant.sequenceButtonView, constant.sequenceWorkflowView, scaleObj);
    sequenceUtil.showZoomBtn(2, "zoomout", constant.sequenceButtonView, constant.sequenceWorkflowView, scaleObj);
}