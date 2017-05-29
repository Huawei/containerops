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

import { mouseoverRelevantWorkflow, mouseoutRelevantWorkflow } from "../relation/lineHover";
import { clickAction } from "../action/clickAction";
import { workflowData } from "./main";
import { initLine } from "./initLine";
import { initWorkflow } from "./initWorkflow";
import { deleteAction } from "../action/addOrDeleteAction";
import * as initButton from "./initButton";
import * as config from "../common/config";

export function animationForRemoveAction(parentId, itemId, itemIndex) {
    var actionViewId = "action-" + parentId;
    /* make target action and reference items disappear */
    var target = "#" + itemId;
    var inputCircle = "#action-self-line-input-" + itemId;
    var outputCircle = "#action-self-line-output-" + itemId;
    var linkPath = "#action-self-line-path-" + itemId;
    var dispappearArray = [target, inputCircle, outputCircle, linkPath];
    var relatedLines = util.findAllRelatedLines(itemId);
    _.each(relatedLines, function(item) {
        var selector = "#" + item.id;
        dispappearArray.push(selector);
    });
    util.disappearAnimation(dispappearArray);

    /* make sibling actions and reference items transform  */
    var siblings = "#" + actionViewId + ">image";
    var siblingInputCircle = "#" + "action-self-line-" + parentId + " > .action-self-line-input";
    var siblingOutputCircle = "#" + "action-self-line-" + parentId + " > .action-self-line-output";
    var siblingLinkPath = "#" + "action-self-line-" + parentId + "> path";
    var transformArray = [{ "selector": siblings, "type": "siblings", "itemIndex": itemIndex }, { "selector": siblingInputCircle, "type": "others", "itemIndex": itemIndex },
        { "selector": siblingOutputCircle, "type": "others", "itemIndex": itemIndex }, { "selector": siblingLinkPath, "type": "others", "itemIndex": itemIndex }
    ]
    util.transformAnimation(transformArray, "action");
}

export function initAction() {
    constant.actionsView.selectAll("g").remove();

    /* draw actions in actionView , data source is stage.actions */
    constant.workflowView.selectAll("image").each(function(d, i) {
        if (d.type == constant.WORKFLOW_STAGE && d.actions != null && d.actions.length > 0) {
            var actionViewId = "action" + "-" + d.id;
            
            constant.actionView[actionViewId] = constant.actionsView.append("g")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight)
                .attr("id", actionViewId);
            var actionStartX = d.translateX + (constant.svgStageWidth - constant.svgActionWidth) / 2;
            var actionStartY = d.translateY;

            constant.actionView[actionViewId].selectAll("image")
                .data(d.actions).enter()
                .append("image")
                .attr("xlink:href", function(ad, ai) {
                    if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "action" && constant.currentSelectedItem.data.id == ad.id) {
                        return config.getSVG(config.SVG_ACTION_SELECTED);
                    } else {
                        return config.getSVG(config.SVG_ACTION);
                    }

                })
                .attr("id", function(ad, ai) {
                    return ad.id;
                })
                .attr("data-index", function(ad, ai) {
                    return ai;
                })
                .attr("data-parent", i)
                .attr("width", function(ad, ai) {
                    return constant.svgActionWidth;
                })
                .attr("height", function(ad, ai) {
                    return constant.svgActionHeight;
                })
                .style("cursor", "pointer")
                .attr("translateX", actionStartX)
                .attr("translateY", function(ad, ai) {
                    /* draw difference distance between action and stage grouped by stage index */
                    if (i % 2 == 0) {
                        ad.translateY = actionStartY + constant.svgStageHeight - 55 + constant.ActionNodeSpaceSize * (ai + 1);
                    } else {
                        ad.translateY = actionStartY + constant.svgStageHeight - 10 + constant.ActionNodeSpaceSize * (ai + 1);
                    }
                    return ad.translateY;
                })
                .attr("transform", function(ad, ai) {
                    ad.width = constant.svgActionWidth;
                    ad.height = constant.svgActionHeight;
                    if (i % 2 == 0) {
                        ad.translateX = actionStartX;
                        ad.translateY = actionStartY + constant.svgStageHeight - 55 + constant.ActionNodeSpaceSize * (ai + 1);
                    } else {
                        ad.translateX = actionStartX;
                        ad.translateY = actionStartY + constant.svgStageHeight - 10 + constant.ActionNodeSpaceSize * (ai + 1);
                    }

                    return "translate(" + ad.translateX + "," + ad.translateY + ")";
                })
                .on("click", function(ad, ai) {
                    clickAction(ad, ai);
                    util.changeCurrentElement(constant.currentSelectedItem);
                    constant.setCurrentSelectedItem({ "data": ad, "parentData": d, "type": "action" });
                    initButton.updateButtonGroup("action");
                    d3.select("#" + ad.id).attr("href", config.getSVG(config.SVG_ACTION_SELECTED));
                    util.cleanToolTip(constant.workflowView, "#workflow-element-popup");
                })
                .on("mouseout", function(ad, ai) {
                    util.cleanToolTip(constant.workflowView, "#workflow-element-popup");
                })
                .on("mouseover", function(ad, ai) {
                    var x = ad.translateX;
                    var y = ad.translateY + constant.svgActionHeight;
                    let text = "Click to Assign Component for Action";
                    let width = 250;
                    let height = null;
                    // if (ad.component.name && ad.setupData && ad.setupData.action && ad.setupData.action.name && ad.setupData.action.name != "") {
                    if (ad.component && ad.component.name && ad.setupData && ad.setupData.action ) {
                        text = ["Component: " + ad.component.name + "[" + ad.component.versionname+"]",
                               "Name: " +ad.setupData.action.name,"Timeout: " + ad.setupData.action.timeout+"(S)",
                               "Image: " + ad.setupData.action.image.name + ":"+ad.setupData.action.image.tag,
                               "Limits:[cpu:" +ad.setupData.pod.spec.containers[0].resources.limits.cpu+", memory:" + ad.setupData.pod.spec.containers[0].resources.limits.memory+"(Mi)]",
                               "Requests:[cpu:" +ad.setupData.pod.spec.containers[0].resources.requests.cpu+", memory:" + ad.setupData.pod.spec.containers[0].resources.requests.memory+"(Mi)]"];
                        // width = text.length * 7 + 20;
                        width = 300;
                        height = text.length * constant.popupHeight;
                    }
                    let options = {
                        "x": x,
                        "y": y,
                        "text": text,
                        "popupId": "workflow-element-popup",
                        "parentView": constant.workflowView,
                        "width": width,
                        "height": height
                    };
                    util.showToolTip(options);
                })
        }

    });

    initLine();
}
