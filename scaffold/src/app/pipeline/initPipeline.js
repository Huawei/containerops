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

import { pipelineData } from "./main";
import { clickStart } from "../stage/clickStart";
import { addStage, deleteStage } from "../stage/addOrDeleteStage";
import { clickStage } from "../stage/clickStage";
import { initAction } from "./initAction";
import { mouseoverRelevantPipeline, mouseoutRelevantPipeline } from "../relation/lineHover";
import { addAction } from "../action/addOrDeleteAction";
import * as initButton from "./initButton";

export function animationForRemoveStage(itemId, itemIndex) {
    var target = "#" + itemId;
    var actions = "#action" + "-" + itemId + "> image";
    var actionReference = "#action-self-line-" + itemId;
    var dispappearArray = [target, actions, actionReference];
    util.disappearAnimation(dispappearArray);
    var siblings = "#pipelineView" + ">image";
    var transformArray = [{ "selector": siblings, "type": "siblings", "itemIndex": itemIndex }]
    util.transformAnimation(transformArray, "stage");
}


export function initPipeline() {
    constant.pipelineView.selectAll("image").remove();
    constant.pipelineView.selectAll("image")
        .data(pipelineData)
        .enter()
        .append("image")
        .attr("xlink:href", function(d, i) {
            if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "stage" && constant.currentSelectedItem.data.id == d.id) {
                if (d.type == constant.PIPELINE_START) {
                    return "../../assets/svg/start-selected-latest.svg";
                } else if (d.type == constant.PIPELINE_STAGE) {
                    return "../../assets/svg/stage-selected-latest.svg";
                }

            } else {
                if (d.type == constant.PIPELINE_START) {
                    return "../../assets/svg/start-latest.svg";
                } else if (d.type == constant.PIPELINE_ADD_STAGE) {
                    return "../../assets/svg/add-stage-latest.svg";
                } else if (d.type == constant.PIPELINE_END) {
                    return "../../assets/svg/end-latest.svg";
                } else if (d.type == constant.PIPELINE_STAGE) {
                    return "../../assets/svg/stage-latest.svg";
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
            d.translateX = i * constant.PipelineNodeSpaceSize + constant.pipelineNodeStartX;
            d.translateY = constant.pipelineNodeStartY;
            return "translate(" + d.translateX + "," + d.translateY + ")";
        })
        .attr("translateX", function(d, i) {
            return i * constant.PipelineNodeSpaceSize + constant.pipelineNodeStartX;
        })
        .attr("translateY", constant.pipelineNodeStartY)
        .attr("class", function(d, i) {
            if (d.type == constant.PIPELINE_START) {
                return constant.PIPELINE_START;
            } else if (d.type == constant.PIPELINE_ADD_STAGE) {
                return constant.PIPELINE_ADD_STAGE;
            } else if (d.type == constant.PIPELINE_END) {
                return constant.PIPELINE_END;
            } else if (d.type == constant.PIPELINE_STAGE) {
                return constant.PIPELINE_STAGE;
            }
        })
        .style("cursor", "pointer")
        .on("click", function(d, i) {
            util.cleanToolTip(constant.pipelineView, "#pipeline-element-popup");
            if (d.type == constant.PIPELINE_ADD_STAGE) {
                addStage(d, i);
                initPipeline();
            } else if (d.type == constant.PIPELINE_STAGE) {
                clickStage(d, i);
                util.changeCurrentElement(constant.currentSelectedItem); /* remove previous selected item style before set current item */
                constant.setCurrentSelectedItem({ "data": d, "type": "stage" }); /* save current item to constant.currentSelectedItem */
                initButton.updateButtonGroup("stage"); /* update the buttons on left top according to current item */
                d3.select("#" + d.id).attr("href", "../../assets/svg/stage-selected-latest.svg"); /* set current item to selected style */
            } else if (d.type == constant.PIPELINE_START) {
                clickStart(d, i);
                util.changeCurrentElement(constant.currentSelectedItem);
                constant.setCurrentSelectedItem({ "data": d, "type": "start" });
                initButton.updateButtonGroup("start");
                d3.select("#" + d.id).attr("href", "../../assets/svg/start-selected-latest.svg");
            }
        })

    .on("mouseover", function(d, i) {
            // console.log(d3.event.movementX);
            // console.log(d3.event.movementY);
            var options = {};
            if (d.type == constant.PIPELINE_ADD_STAGE) {
                d3.select(this)
                    .attr("xlink:href", function(d, i) {
                        return "../../assets/svg/add-stage-selected-latest.svg";
                    })
                    .style({
                        "cursor": "pointer"
                    })
                options = {
                    "x": i * constant.PipelineNodeSpaceSize + constant.pipelineNodeStartX,
                    "y": constant.pipelineNodeStartY + constant.svgStageHeight,
                    "text": "Add Stage",
                    "popupId": "pipeline-element-popup",
                    "parentView": constant.pipelineView
                };
                 util.showToolTip(options);

            } else if (d.type == constant.PIPELINE_STAGE || d.type == constant.PIPELINE_START) {
                let text = "Click to Edit";
                let width = null;
                if (d.setupData && d.setupData.name && d.setupData.name != "") {
                    text = d.setupData.name;
                    width = text.length * 8 + 20;
                }
                options = {
                    "x": i * constant.PipelineNodeSpaceSize + constant.pipelineNodeStartX,
                    "y": constant.pipelineNodeStartY + constant.svgStageHeight,
                    "text": text,
                    "popupId": "pipeline-element-popup",
                    "parentView": constant.pipelineView,
                    "width": width
                };
                 util.showToolTip(options);
            }
           

        })
        .on("mouseout", function(d, i) {
            d3.event.stopPropagation();
            if (d.type == constant.PIPELINE_ADD_STAGE) {
                d3.select(this)
                    .attr("xlink:href", function(d, i) {
                        return "../../assets/svg/add-stage-latest.svg";
                    })
            }
            util.cleanToolTip(constant.pipelineView, "#pipeline-element-popup");

        })

    initAction();

}
