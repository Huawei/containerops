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

import { workflowData } from "./main";
import { clickStart } from "../stage/clickStart";
import { addStage, deleteStage } from "../stage/addOrDeleteStage";
import { clickStage } from "../stage/clickStage";
import { initAction } from "./initAction";
import { mouseoverRelevantWorkflow, mouseoutRelevantWorkflow } from "../relation/lineHover";
import { addAction } from "../action/addOrDeleteAction";
import * as initButton from "./initButton";
import * as config from "../common/config";

export function animationForRemoveStage(itemId, itemIndex) {
    var target = "#" + itemId;
    var actions = "#action" + "-" + itemId + "> image";
    var actionReference = "#action-self-line-" + itemId;
    var dispappearArray = [target, actions, actionReference];
    util.disappearAnimation(dispappearArray);
    var siblings = "#workflowView" + ">image";
    var transformArray = [{ "selector": siblings, "type": "siblings", "itemIndex": itemIndex }]
    util.transformAnimation(transformArray, "stage");
}

export function initWorkflow() {
    constant.workflowView.selectAll("image").remove();
    constant.workflowView.selectAll("image")
        .data(workflowData)
        .enter()
        .append("image")
        .attr("xlink:href", function(d, i) {
            if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "stage" && constant.currentSelectedItem.data.id == d.id) {
                if (d.type == constant.WORKFLOW_START) {
                    return config.getSVG(config.SVG_START_SELECTED);
                } else if (d.type == constant.WORKFLOW_STAGE) {
                   return config.getSVG(config.SVG_STAGE_SELECTED);
                }

            } else {
                if (d.type == constant.WORKFLOW_START) {

                    return config.getSVG(config.SVG_START);
                } else if (d.type == constant.WORKFLOW_ADD_STAGE) {
                    return config.getSVG(config.SVG_ADD_STAGE);
                } else if (d.type == constant.WORKFLOW_END) {
                    return config.getSVG(config.SVG_END);
                } else if (d.type == constant.WORKFLOW_STAGE) {
                    return config.getSVG(config.SVG_STAGE);
                }
            }
        })
        .attr("id", function(d, i) {
            return d.id;
        })
        .attr("data-index", function(d, i) {
            return i;
        })
        .attr("width", function(d, i) {
            return constant.svgStageWidth;
        })
        .attr("height", function(d, i) {
            return constant.svgStageHeight;
        })
        .attr("transform", function(d, i) {
            d.width = constant.svgStageWidth;
            d.height = constant.svgStageHeight;
            d.translateX = i * constant.WorkflowNodeSpaceSize + constant.workflowNodeStartX;
            d.translateY = constant.workflowNodeStartY;
            return "translate(" + d.translateX + "," + d.translateY + ")";
        })
        .attr("translateX", function(d, i) {
            return i * constant.WorkflowNodeSpaceSize + constant.workflowNodeStartX;
        })
        .attr("translateY", constant.workflowNodeStartY)
        .attr("class", function(d, i) {
            if (d.type == constant.WORKFLOW_START) {
                return constant.WORKFLOW_START;
            } else if (d.type == constant.WORKFLOW_ADD_STAGE) {
                return constant.WORKFLOW_ADD_STAGE;
            } else if (d.type == constant.WORKFLOW_END) {
                return constant.WORKFLOW_END;
            } else if (d.type == constant.WORKFLOW_STAGE) {
                return constant.WORKFLOW_STAGE;
            }
        })
        .style("cursor", "pointer")
        .on("click", function(d, i) {
            util.cleanToolTip(constant.workflowView, "#workflow-element-popup");
            if (d.type == constant.WORKFLOW_ADD_STAGE) {
                addStage(d, i);
                initWorkflow();
            } else if (d.type == constant.WORKFLOW_STAGE) {
                clickStage(d, i);
                util.changeCurrentElement(constant.currentSelectedItem); /* remove previous selected item style before set current item */
                constant.setCurrentSelectedItem({ "data": d, "type": "stage" }); /* save current item to constant.currentSelectedItem */
                initButton.updateButtonGroup("stage"); /* update the buttons on left top according to current item */
                d3.select("#" + d.id).attr("href", config.getSVG(config.SVG_STAGE_SELECTED)); /* set current item to selected style */
            } else if (d.type == constant.WORKFLOW_START) {
                clickStart(d, i);
                util.changeCurrentElement(constant.currentSelectedItem);
                constant.setCurrentSelectedItem({ "data": d, "type": "start" });
                initButton.updateButtonGroup("start");
                d3.select("#" + d.id).attr("href", config.getSVG(config.SVG_START_SELECTED));
            }
        })

    .on("mouseover", function(d, i) {
            // console.log(d3.event.movementX);
            // console.log(d3.event.movementY);
            var options = {};
            if (d.type == constant.WORKFLOW_ADD_STAGE) {
                d3.select(this)
                    .attr("xlink:href", function(d, i) {
                        return config.getSVG(config.SVG_ADD_STAGE_SELECTED);
                    })
                    .style({
                        "cursor": "pointer"
                    })
                options = {
                    "x": i * constant.WorkflowNodeSpaceSize + constant.workflowNodeStartX,
                    "y": constant.workflowNodeStartY + constant.svgStageHeight,
                    "text": "Add Stage",
                    "popupId": "workflow-element-popup",
                    "parentView": constant.workflowView
                };
                util.showToolTip(options);

            } else if (d.type == constant.WORKFLOW_STAGE || d.type == constant.WORKFLOW_START) {
                let text = d.type==constant.WORKFLOW_START?"Click to Edit Output":"Click to Edit Name and Timeout";
                let width = 220;
                let height = null;
                if (d.setupData && ((d.setupData.name && d.setupData.name != "") || (d.setupData.timeout && d.setupData.timeout != ""))) {
                    text = ["Name: " + d.setupData.name, "Timeout: "+d.setupData.timeout+"(S)"];
                    // width = text.length * 8 + 20;
                    width = 300;
                    height = text.length * constant.popupHeight;
                }
                options = {
                    "x": i * constant.WorkflowNodeSpaceSize + constant.workflowNodeStartX,
                    "y": constant.workflowNodeStartY + constant.svgStageHeight,
                    "text": text,
                    "popupId": "workflow-element-popup",
                    "parentView": constant.workflowView,
                    "width": width,
                    "height":height
                };
                util.showToolTip(options);
            }


        })
        .on("mouseout", function(d, i) {
            d3.event.stopPropagation();
            if (d.type == constant.WORKFLOW_ADD_STAGE) {
                d3.select(this)
                    .attr("xlink:href", function(d, i) {
                        return config.getSVG(config.SVG_ADD_STAGE);
                    })
            }
            util.cleanToolTip(constant.workflowView, "#workflow-element-popup");

        })

    initAction();

}
