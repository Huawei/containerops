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
import * as util from "../common/util";

import { initAction } from "./initAction";
import { addAction, deleteAction } from "../action/addOrDeleteAction";
import { pipelineData } from "./main";
import { animationForRemoveStage, initPipeline } from "./initPipeline";
import { addStage, deleteStage } from "../stage/addOrDeleteStage";
import { animationForRemoveAction } from "./initAction";
import { setPath } from "../relation/setPath";
import { initLine } from "./initLine";

export var buttonWidth = 23,
    buttonHeight = 23,
    buttonVerticalSpace = 6,
    background = "#555",
    buttonHorizonSpace = 20,
    rectBackgroundY = 15;

export function initButton() {
    constant.zoomScale = 1;
    constant.zoomTargetScale = 1;
    constant.buttonView
        .append("rect")
        .attr("width", constant.svgWidth)
        .attr("height", 2 * buttonVerticalSpace + buttonHeight)
        .attr("y", rectBackgroundY)
        .style({
            "fill": "#f7f7f7"
        });
    showZoomBtn(1, "zoomin");
    showZoomBtn(2, "zoomout");
    // showSeperateLine(3);
}
export function updateButtonGroup(currentItemType) {
    cleanOptBtn();
    if (constant.currentSelectedItem != null) {
        if (currentItemType == "stage") {
            showOptBtn(3, "add");
            showOptBtn(4, "delete", currentItemType);
        } else if (currentItemType == "action") {
            showOptBtn(3, "delete", currentItemType);
        } else if (currentItemType == "line") {
            showOptBtn(3, "removeLink");
        } else {
            cleanOptBtn();
        }
    }
}

export function showToolTip(x, y, text, popupId, parentView, width, height) {
    parentView
        .append("g")
        .attr("id", popupId);
    parentView.selectAll("#" + popupId)
        .append("rect")
        .attr("width", width || constant.popupWidth)
        .attr("height", height || constant.popupHeight)
        .attr("x", function(pd, pi) {
            return x;
        })
        .attr("y", function(pd, pi) {
            return y;
        })
        .attr("rx", 3)
        .attr("ry", 3)
        .style("fill", background)
        .style("opacity", 0.9)
    parentView.selectAll("#" + popupId)
        .append("text")
        .attr("x", x + 10)
        .attr("y", y + constant.popupHeight / 2 + 4)
        .style("fill", "white")
        .style("opacity", 0.9)
        .text(text)
}

function showOptBtn(index, type) {
    constant.buttonView
        .append("image")
        .attr("xlink:href", function(ad, ai) {
            if (type == "add") {
                return "../../assets/svg/add-action-latest.svg";
            } else if (type == "delete") {
                return "../../assets/svg/delete-latest.svg";
            } else if (type == "removeLink") {
                return "../../assets/svg/remove-link-latest.svg";
            }

        })
        .attr("translateX", function(d, i) {
            return index * buttonHorizonSpace + (index - 1) * buttonWidth;
        })
        .attr("translateY", function(d, i) {
            return buttonVerticalSpace + rectBackgroundY;
        })
        .attr("transform", function(d, i) {
            let translateX = d3.select(this).attr("translateX");
            let translateY = d3.select(this).attr("translateY");
            return "translate(" + translateX + "," + translateY + ")";
        })
        .attr("id", function(d, i) {
            if (type == "add") {
                return "addActionBtn";
            } else if (type == "delete") {
                return "deleteBtn";
            } else if (type == "removeLink") {
                return "removeLinkBtn";
            }

        })
        .attr("class", "optBtn")
        .attr("width", buttonWidth)
        .attr("height", buttonHeight)
        .on("mouseover", function(d, i) {
            d3.select(this).style("cursor", "pointer");
            let content = "";
            let href = "";
            if (type == "add") {
                content = "Add Action";
                href = "../../assets/svg/add-action-selected-latest.svg";

            } else if (type == "delete") {
                content = "Delete";
                href = "../../assets/svg/delete-selected-latest.svg";

            } else if (type == "removeLink") {
                content = "Remove Link";
                href = "../../assets/svg/remove-link-selected-latest.svg";
            }
            d3.select(this).attr("href", href);
            showToolTip(Number(d3.select(this).attr("translateX")), Number(d3.select(this).attr("translateY")) + buttonHeight, content, "button-element-popup", constant.buttonView);
        })
        .on("mouseout", function(d, i) {
            cleanToolTip();
            let href = "";
            if (type == "add") {
                href = "../../assets/svg/add-action-latest.svg";

            } else if (type == "delete") {
                href = "../../assets/svg/delete-latest.svg";

            } else if (type == "removeLink") {
                href = "../../assets/svg/remove-link-latest.svg";
            }
            d3.select(this).attr("href", href);
        })
        .on("click", function(d, i) {
            cleanToolTip();
            if (type == "add") {
                addAction(constant.currentSelectedItem.data.actions);
                initAction();
            } else if (type == "delete") {
                $("#pipeline-info-edit").html("");
                var timeout = 0;
                var index = d3.select("#" + constant.currentSelectedItem.data.id).attr("data-index");
                /* if remove the node is not the last one, add animation to action */
                if (constant.currentSelectedItem.type == "stage") {
                    if (i < pipelineData.length - 1) {
                        timeout = 400;
                        animationForRemoveStage(constant.currentSelectedItem.data.id, index);
                    }
                    setTimeout(function() {
                        deleteStage(constant.currentSelectedItem.data, index);
                        constant.setCurrentSelectedItem(null);
                        initPipeline();
                    }, timeout)
                } else if (constant.currentSelectedItem.type == "action") {
                    $("#pipeline-info-edit").html("");
                    var timeout = 0;
                    // TODO
                    var index = d3.select("#" + constant.currentSelectedItem.data.id).attr("data-index");
                    var stageData = constant.currentSelectedItem.parentData;
                    var actionData = constant.currentSelectedItem.data;
                    /* if remove the node is not the last one, add animation to action */
                    if (index < stageData.actions.length - 1) {
                        timeout = 400;
                        animationForRemoveAction(stageData.id, actionData.id, index);
                    }

                    /* reload pipeline after the animation */
                    setTimeout(function() {
                        deleteAction(actionData, index);
                        constant.setCurrentSelectedItem(null);
                        initPipeline();
                    }, timeout)
                }
                cleanOptBtn();


            } else if (type == "removeLink") {
                $("#pipeline-info-edit").html("");
                var id = constant.currentSelectedItem.data.attr("id");
                constant.currentSelectedItem.data.remove();
                var lineData = _.find(constant.linePathAry, function(item) {
                    return item.id == id;
                })
                var index = _.indexOf(constant.linePathAry, lineData);
                constant.linePathAry.splice(index, 1);
                cleanOptBtn();
                constant.setCurrentSelectedItem(null);
            }

        })
}

function showZoomBtn(index, type) {
    constant.buttonView
        .append("image")
        .attr("xlink:href", function(ad, ai) {
            if (type == "zoomin") {
                return "../../assets/svg/zoomin.svg";
            } else if (type == "zoomout") {
                return "../../assets/svg/zoomout.svg";
            }

        })
        .attr("translateX", function(d, i) {
            return index * buttonHorizonSpace + (index - 1) * buttonWidth;
        })
        .attr("translateY", function(d, i) {
            return buttonVerticalSpace + rectBackgroundY;
        })
        .attr("transform", function(d, i) {
            let translateX = d3.select(this).attr("translateX");
            let translateY = d3.select(this).attr("translateY");
            return "translate(" + translateX + "," + translateY + ")";
        })
        .attr("id", function(d, i) {
            if (type == "zoomin") {
                return "pipeline-zoomin";
            } else if (type == "zoomout") {
                return "pipeline-zoomout";
            }
        })
        .attr("width", buttonWidth)
        .attr("height", buttonHeight)
        .style("cursor", "pointer")
        .on("mouseover", function(d, i) {
            let content = "";
            let href = "";
            // let tooltipPosX = index * buttonHorizonSpace + (index - 1) * buttonWidth;
            if (type == "zoomin") {
                content = "Zoomin";
                href = "../../assets/svg/zoomin.svg";
            } else if (type == "zoomout") {
                content = "Zoomout";
                href = "../../assets/svg/zoomout.svg";
            }
            d3.select(this).attr("href", href);
            showToolTip(Number(d3.select(this).attr("translateX")), Number(d3.select(this).attr("translateY")) + buttonHeight, content, "button-element-popup", constant.buttonView);
        })
        .on("mouseout", function(d, i) {
            cleanToolTip();
        })
        .on("click", function(d, i) {
            util.zoomed(type);
        })
}
function showSeperateLine(index){
    constant.buttonView
        .append("image")
        .attr("xlink:href", function(ad, ai) {
                return "../../assets/svg/seperate-line.svg";
        })
        .attr("translateX", function(d, i) {
            return index * buttonHorizonSpace + (index - 1) * buttonWidth;
        })
        .attr("translateY", function(d, i) {
            return rectBackgroundY;
        })
        .attr("transform", function(d, i) {
            let translateX = d3.select(this).attr("translateX");
            let translateY = d3.select(this).attr("translateY");
            return "translate(" + translateX + "," + translateY + ")";
        })
        .attr("width", 2 * buttonVerticalSpace + buttonHeight)
        .attr("height", 2 * buttonVerticalSpace + buttonHeight)
}
function cleanOptBtn() {
    constant.buttonView.selectAll("image.optBtn").remove();
}

function cleanToolTip() {
    constant.buttonView.selectAll("#button-element-popup").remove();
}
