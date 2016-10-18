import * as constant from "../common/constant";
import * as util from "../common/util";

import { drag } from "../common/drag";
import { mouseoverRelevantPipeline, mouseoutRelevantPipeline } from "../relation/lineHover";
import { clickAction } from "../action/clickAction";
import { dragDropSetPath } from "../relation/dragDropSetPath";
import { removeLinkArray } from "../relation/removeLinkArray";
import { pipelineData } from "./main";
import { initLine } from "./initLine";
import { initPipeline } from "./initPipeline";
import { deleteAction} from "../action/addOrDeleteAction";

var animationForRemove = function(parentId, parentIndex, itemId, itemIndex) {
    var actionViewId = "action-" + parentId;
    /* make target action and reference items disappear */
    var target = "#" + itemId;
    var inputCircle = "#action-self-line-input-" + itemId;
    var outputCircle = "#action-self-line-output-" + itemId;
    var linkPath = "#action-self-line-path-" + itemId;
    var dispappearArray = [target, inputCircle, outputCircle, linkPath];
    var relatedLines = util.findAllRelatedLines(itemId);
    _.each(relatedLines, function(item) {
        var selector = "." + item.defaultClass;
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
    constant.pipelineView.selectAll("image").each(function(d, i) {
        if (d.type == constant.PIPELINE_STAGE && d.actions != null && d.actions.length > 0) {
            var actionViewId = "action" + "-" + d.id;

            constant.actionView[actionViewId] = constant.actionsView.append("g")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight)
                .attr("id", actionViewId);

            var actionStartX = d.translateX + 7.5;
            var actionStartY = d.translateY;

            constant.actionView[actionViewId].selectAll("image")
                .data(d.actions).enter()
                .append("image")
                .attr("xlink:href", function(ad, ai) {
                    return "../../assets/svg/action-bottom.svg";
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
                .on("mouseover", function(ad, ai) {
                    // mouseoverRelevantPipeline(ad);
                    if (d3.event.movementX != 0 || d3.event.movementY != 0) {
                        constant.buttonView.selectAll("#button-group").remove();
                        var editBtnShow = false,
                            deleteBtnShow = false;
                        /* add g as button group to buttonView when mouse over nodes of pipeline */
                        constant.buttonView
                            .append("g")
                            .attr("id", "button-group")
                            .attr("width", constant.svgActionWidth * 2)
                            .attr("height", constant.svgActionHeight * 2)
                            .on("mouseout", function(rd, ri) {
                                setTimeout(function() {
                                    if (!(editBtnShow || deleteBtnShow)) {
                                        constant.buttonView.selectAll("#button-group").remove();
                                    }

                                }, 100)
                            })
                            /* add edit image to edit setup data and input output data*/
                        constant.buttonView.select("#button-group")
                            .append("image")
                            .attr("id", "edit-image-" + ad.id)
                            .attr("xlink:href", function(ed, ei) {
                                return "../../assets/svg/edit-mouseout-2.svg";
                            })
                            .attr("width", function(ed, ei) {
                                return constant.svgActionWidth * 2;
                            })
                            .attr("height", function(ed, ei) {
                                return constant.svgActionHeight * 2 * 0.5;
                            })
                            .attr("transform", function(ed, ei) {
                                let translateX = ad.translateX - constant.svgActionWidth * 0.5;
                                let translateY = ad.translateY - constant.svgActionHeight * 0.5 + 1;
                                return "translate(" + translateX + "," + translateY + ")";
                            })
                            .on("mouseover", function(ed, ei) {
                                d3.event.stopPropagation();
                                editBtnShow = true;
                                constant.buttonView.select("#edit-image-" + ad.id)
                                    .attr("xlink:href", function(ed, ei) {
                                        return "../../assets/svg/edit-mouseover-2.svg";
                                    })
                            })
                            .on("mouseout", function(ed, ei) {
                                editBtnShow = false;
                                constant.buttonView.select("#edit-image-" + ad.id)
                                    .attr("xlink:href", function(ed, ei) {
                                        return "../../assets/svg/edit-mouseout-2.svg";
                                    })
                            })
                            .on("click", function(ed, ei) {
                                constant.buttonView.selectAll("#button-group").remove();
                                clickAction(ad, ai);

                            })
                        constant.buttonView.select("#button-group")
                            .append("image")
                            .attr("id", "delete-image-" + ad.id)
                            .attr("xlink:href", function(dd, di) {
                                return "../../assets/svg/delete-mouseout.svg";
                            })
                            .attr("width", function(dd, di) {
                                return constant.svgActionWidth * 2;
                            })
                            .attr("height", function(dd, di) {
                                return constant.svgActionHeight * 2 * 0.5;
                            })
                            .attr("transform", function(dd, di) {
                                let translateX = ad.translateX - constant.svgActionWidth * 0.5;
                                let translateY = ad.translateY + constant.svgActionHeight * 0.5 - 1;
                                return "translate(" + translateX + "," + translateY + ")";
                            })
                            .on("mouseover", function(dd, di) {
                                d3.event.stopPropagation();
                                deleteBtnShow = true;
                                constant.buttonView.select("#delete-image-" + ad.id)
                                    .attr("xlink:href", function(d, i) {
                                        return "../../assets/svg/delete-mouseover.svg";
                                    })
                            })
                            .on("mouseout", function(dd, di) {
                                deleteBtnShow = false;
                                constant.buttonView.select("#delete-image-" + ad.id)
                                    .attr("xlink:href", function(dd, di) {
                                        return "../../assets/svg/delete-mouseout.svg";
                                    })
                            })
                            .on("click", function(dd, di) {
                                constant.buttonView.selectAll("#button-group").remove();
                                $("#pipeline-info-edit").html("");

                                var timeout = 0;
                                /* if remove the node is not the last one, add animation to action */
                                if (ai < d.actions.length - 1) {
                                    timeout = 400;
                                    animationForRemove(d.id, i, ad.id, ai);
                                }

                                /* reload pipeline after the animation */
                                setTimeout(function() {
                                    deleteAction(ad, ai);
                                    initPipeline();
                                    initAction();
                                }, timeout)
                                constant.buttonView.selectAll("#button-group").remove();

                            })
                    }

                })

            // .call(drag);


        }

    });

    initLine();
}
