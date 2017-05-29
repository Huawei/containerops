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
import * as config from "../common/config";

import { initAction } from "./initAction";
import { addAction, deleteAction } from "../action/addOrDeleteAction";
import { workflowData,workflowSettingData } from "./main";
import { animationForRemoveStage, initWorkflow } from "./initWorkflow";
import { addStage, deleteStage } from "../stage/addOrDeleteStage";
import { animationForRemoveAction } from "./initAction";
import { setPath } from "../relation/setPath";
import { initLine } from "./initLine";
import { initWorkflowSetting } from "./workflowSetting";

let rectBackgroundY = 15;

export function initButton() {
    let scaleObj = { "zoomScale": 1, "zoomTargetScale": 1 };
    constant.buttonView
        .append("rect")
        .attr("width", constant.svgWidth)
        .attr("height", rectBackgroundY)
        .style({
            "fill": "#ffffff"
        });
    constant.buttonView
        .append("rect")
        .attr("width", constant.svgWidth)
        .attr("height", 2 * constant.buttonVerticalSpace + constant.buttonHeight)
        .attr("y", rectBackgroundY)
        .style({
            "fill": "#f7f7f7"
        });
    util.showZoomBtn(1, "zoomin", constant.buttonView, constant.workflowView, scaleObj);
    util.showZoomBtn(2, "zoomout", constant.buttonView, constant.workflowView, scaleObj);
    showOptBtn(3, "setting", "wfSettingBtn");
}
export function updateButtonGroup(currentItemType) {
    cleanOptBtn();
    if (constant.currentSelectedItem != null) {
        if (currentItemType == "stage") {
            showOptBtn(4, "add");
            showOptBtn(5, "delete");
        } else if (currentItemType == "action") {
            showOptBtn(4, "delete");
        } else if (currentItemType == "line") {
            showOptBtn(4, "removeLink");
        } else {
            cleanOptBtn();
        }
    }
}

function showOptBtn(index, type, css_class) {
    css_class = css_class || "optBtn";
    constant.buttonView
        .append("image")
        .attr("xlink:href", function(ad, ai) {
            if (type == "add") {
                return config.getSVG(config.SVG_ADD_ACTION);
            } else if (type == "delete") {
                return config.getSVG(config.SVG_DELETE);
            } else if (type == "removeLink") {
                return config.getSVG(config.SVG_REMOVE_LINK);
            } else if(type == 'setting'){
                return config.getSVG(config.SVG_WORKFLOW_SET);
            }

        })
        .attr("translateX", function(d, i) {
            return index * constant.buttonHorizonSpace + (index - 1) * constant.buttonWidth;
        })
        .attr("translateY", function(d, i) {
            return constant.buttonVerticalSpace + rectBackgroundY;
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
            } else if(type == "setting"){
                return "wfSettingBtn";
            }

        })
        .attr("class", css_class)
        .attr("width", constant.buttonWidth)
        .attr("height", constant.buttonHeight)
        .on("mouseover", function(d, i) {
            d3.select(this).style("cursor", "pointer");
            let content = "";
            let href = "";
            let width = null;
            if (type == "add") {
                content = "Add Action";
                href = config.getSVG(config.SVG_ADD_ACTION_SELECTED);

            } else if (type == "delete") {
                content = "Delete";
                 href = config.getSVG(config.SVG_DELETE_SELECTED);

            } else if (type == "removeLink") {
                content = "Remove Link";
                 href = config.getSVG(config.SVG_REMOVE_LINK_SELECTED);
            } else if (type == "setting") {
                content = "Workflow Setting";
                 href = config.getSVG(config.SVG_WORKFLOW_SET_SELECTED);
                 width = 130;
            }
            d3.select(this).attr("href", href);
            let options = {
                "x": Number(d3.select(this).attr("translateX")),
                "y": Number(d3.select(this).attr("translateY")) + constant.buttonHeight,
                "text": content,
                "popupId": "button-element-popup",
                "parentView": constant.buttonView,
                'width':width
            };
            util.showToolTip(options);
        })
        .on("mouseout", function(d, i) {
            util.cleanToolTip(constant.buttonView, "#button-element-popup");
            let href = "";
            if (type == "add") {
                href = config.getSVG(config.SVG_ADD_ACTION);

            } else if (type == "delete") {
                href = config.getSVG(config.SVG_DELETE);

            } else if (type == "removeLink") {
                href = config.getSVG(config.SVG_REMOVE_LINK);
            } else if(type == "setting"){
                href = config.getSVG(config.SVG_WORKFLOW_SET);
            }
            d3.select(this).attr("href", href);
        })
        .on("click", function(d, i) {
            util.cleanToolTip(constant.buttonView, "#button-element-popup");
            if (type == "add") {
                addAction(constant.currentSelectedItem.data.actions);
                initAction();
            } else if (type == "delete") {
                $("#workflow-info-edit").html("");
                var timeout = 0;
                var index = d3.select("#" + constant.currentSelectedItem.data.id).attr("data-index");
                /* if remove the node is not the last one, add animation to action */
                if (constant.currentSelectedItem.type == "stage") {
                    if (i < workflowData.length - 1) {
                        timeout = 400;
                        animationForRemoveStage(constant.currentSelectedItem.data.id, index);
                    }
                    setTimeout(function() {
                        deleteStage(constant.currentSelectedItem.data, index);
                        constant.setCurrentSelectedItem(null);
                        initWorkflow();
                    }, timeout)
                } else if (constant.currentSelectedItem.type == "action") {
                    $("#workflow-info-edit").html("");
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

                    /* reload workflow after the animation */
                    setTimeout(function() {
                        deleteAction(actionData, index);
                        constant.setCurrentSelectedItem(null);
                        initWorkflow();
                    }, timeout)
                }
                cleanOptBtn();


            } else if (type == "removeLink") {
                $("#workflow-info-edit").html("");
                var id = constant.currentSelectedItem.data.attr("id");
                constant.currentSelectedItem.data.remove();
                var lineData = _.find(constant.linePathAry, function(item) {
                    return item.id == id;
                })
                var index = _.indexOf(constant.linePathAry, lineData);
                constant.linePathAry.splice(index, 1);
                cleanOptBtn();
                constant.setCurrentSelectedItem(null);
            } else if(type == 'setting'){
                initWorkflowSetting(workflowSettingData);
            }

        })
}

function cleanOptBtn() {
    constant.buttonView.selectAll("image.optBtn").remove();
}
