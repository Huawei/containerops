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

export function drawTreeNode(x, y, text, type, opacity, conflictNodes) {
    var p = d3.scale.category10();
    var r = p.range();
    var color = d3.scale.ordinal().domain(["string", "object", "boolean", "number"]).range(r)(type);
    var strokeColor = d3.scale.ordinal().domain(["string", "object", "boolean", "number"]).range(r)(type);
    var rectWidth = 180,
        rectHeight = 30,
        imageWidth = 26,
        imageHeight = 26,
        buttonWidth = 15,
        buttonHeight = 15,
        radius = 8,
        space = 5,
        textSpace2Top = 4,
        popupWidth = 180,
        popupHeight = 30,
        popupFill = "#333",
        textColor = "white";
    var randomID = uuid.v1();
    constant.treeView.append("g")
        .attr("id", "tree-node-group-" + randomID)
        .on("mouseover", function(d, i) {
            highlightConflict(conflictNodes);
        });
    constant.treeView.selectAll("#tree-node-group-" + randomID)
        .append("rect")
        .attr({
            "width": rectWidth,
            "height": rectHeight,
            "x": x,
            "y": y,
            "rx": radius,
            "ry": radius
        })
        .style({
            "fill": color,
            "fill-opacity": opacity,
            "stroke": strokeColor,
            "stroke-width": 2
        })

    constant.treeView.selectAll("#tree-node-group-" + randomID)
        .append("image")
        .attr({
            "xlink:href": "../../assets/svg/add.svg",
            "width": imageWidth,
            "height": imageHeight,
            "x": x + space,
            "y": y + (rectHeight - imageHeight) / 2
        })
        .style({
            "cursor": "pointer"
        })
        .on("mouseover", function(d, i) {
            constant.treeView.selectAll("#tree-element-popup").remove();
            constant.treeView
                .append("g")
                .attr("id", "tree-element-popup");
            constant.treeView.selectAll("#tree-element-popup")
                .append("rect")
                .attr("width", popupWidth)
                .attr("height", popupHeight)
                .attr("x", function(d, i) {
                    return x + space;
                })
                .attr("y", y + rectHeight)
                .attr("id", "tree-element-popup")
                .attr("rx", radius)
                .attr("ry", radius)
                .style("fill", popupFill)
                .style("opacity", 1)
            constant.treeView.selectAll("#tree-element-popup")
                .append("text")
                .attr("x", x + 3 * space)
                .attr("y", y + rectHeight + popupHeight / 2 + textSpace2Top)
                .style("fill", textColor)
                .style("opacity", 1)
                .text("popup for" + text)

        })
        .on("mouseout", function(d, i) {
            d3.select("#tree-element-popup").remove();
        })
        .on("click", function(d, i) {
            resolveConflict();
        })
    constant.treeView.selectAll("#tree-node-group-" + randomID)
        .append("text")
        .style({
            "fill": textColor,
            "font-weight": "bold"
        })
        .attr("x", function() {
            return x + space * 2 + imageWidth;
        })
        .attr("y", function() {
            return y + rectHeight / 2 + textSpace2Top;
        })
        .text(text + "-" + type)
        .on("mouseover", function(d, i) {

        })
        .on("mouseout", function(d, i) {

        })
        .on("click", function(d, i) {

        });
    constant.treeView.selectAll("#tree-node-group-" + randomID)
        .append("image")
        .attr({
            "xlink:href": "../../assets/svg/end.svg",
            "width": imageWidth,
            "height": imageHeight,
            "x": x + rectWidth - space - imageWidth,
            "y": y + (rectHeight - imageHeight) / 2
        })
        .style({
            "cursor": "pointer"
        })
        .on("mouseover", function(d, i) {
            constant.treeView.selectAll("#tree-element-popup").remove();
            constant.treeView
                .append("g")
                .attr("id", "tree-element-popup");
            constant.treeView.selectAll("#tree-element-popup")
                .append("rect")
                .attr("width", popupWidth)
                .attr("height", popupHeight)
                .attr("x", function(d, i) {
                    return x + rectWidth - space - imageWidth;
                })
                .attr("y", y + rectHeight)
                .attr("id", "tree-element-popup")
                .attr("rx", radius)
                .attr("ry", radius)
                .style("fill", popupFill)
                .style("opacity", 1)
            constant.treeView.selectAll("#tree-element-popup")
                .append("text")
                .attr("x", x + rectWidth - imageWidth + space)
                .attr("y", y + rectHeight + popupHeight / 2 + textSpace2Top)
                .style("fill", textColor)
                .style("opacity", 1)
                .text("popup for" + text)

        })
        .on("mouseout", function(d, i) {
            d3.select("#tree-element-popup").remove();

        })
        .on("click", function(d, i) {

        })

}

function showTooltip() {

}

function highlightConflict(conflictNodes) {

}

function resolveConflict() {

}
