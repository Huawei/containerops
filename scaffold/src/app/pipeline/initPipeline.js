import * as constant from "../common/constant";
import * as util from "../common/util";

import { pipelineData } from "./main";
import { drag } from "../common/drag";
import { clickStart } from "../stage/clickStart";
import { addStage, deleteStage } from "../stage/addOrDeleteStage";
import { clickStage } from "../stage/clickStage";
import { initAction } from "../pipeline/initAction";
import { mouseoverRelevantPipeline, mouseoutRelevantPipeline } from "../relation/lineHover";
import { dragDropSetPath } from "../relation/dragDropSetPath";
import { removeLinkArray } from "../relation/removeLinkArray";
import { addAction } from "../action/addOrDeleteAction";

var animationForRemove = function(itemId, itemIndex) {
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
    // console.log("pipelineData");
    // console.log(pipelineData);
    constant.pipelineView.selectAll("image").remove();
    constant.pipelineView.selectAll("image")
        .data(pipelineData)
        .enter()
        .append("image")
        .attr("xlink:href", function(d, i) {
            // console.log(d.type);
            if (d.type == constant.PIPELINE_START) {
                return "../../assets/svg/start.svg";
            } else if (d.type == constant.PIPELINE_ADD_STAGE) {
                return "../../assets/svg/addStage.svg";
            } else if (d.type == constant.PIPELINE_END) {
                return "../../assets/svg/end.svg";
            } else if (d.type == constant.PIPELINE_STAGE) {
                return "../../assets/svg/stage.svg";
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
        .on("click", function(d, i) {
            if (d.type == constant.PIPELINE_ADD_STAGE) {
                addStage(d, i);
                initPipeline();
                initAction();
            }
        })
        .on("mouseout", function(d, i) {
            d3.event.stopPropagation();
            if (d.type == constant.PIPELINE_ADD_STAGE) {
                d3.select(this)
                    .attr("xlink:href", function(d, i) {
                        return "../../assets/svg/addStage.svg";
                    })
            }
        })
        .on("mouseover", function(d, i) {
            // console.log(d3.event.movementX);
            // console.log(d3.event.movementY);

            if (d.type == constant.PIPELINE_ADD_STAGE) {
                d3.select(this)
                    .attr("xlink:href", function(d, i) {
                        return "../../assets/svg/addStage-mouseover.svg";
                    })

            } else if (d.type == constant.PIPELINE_STAGE || d.type == constant.PIPELINE_START) {
                if (d3.event.movementX != 0 || d3.event.movementY != 0) {
                    constant.buttonView.selectAll("#button-group").remove();
                    var editBtnShow = false,
                        addActionBtnShow = false,
                        deleteBtnShow = false;
                    /**
                       Add g element to buttonView as the parent group of action buttons
                    */
                    constant.buttonView
                        .append("g")
                        .attr("id", "button-group")
                        .attr("width", constant.svgStageWidth * 2)
                        .attr("height", constant.svgStageHeight * 2)
                        .on("mouseout", function(rd, ri) {
                            setTimeout(function() {
                                if (d.type == constant.PIPELINE_START) {
                                    if (!(editBtnShow)) {
                                        constant.buttonView.selectAll("#button-group").remove();
                                    }
                                } else if (d.type == constant.PIPELINE_STAGE) {
                                    if (!(editBtnShow || addActionBtnShow || deleteBtnShow)) {
                                        constant.buttonView.selectAll("#button-group").remove();
                                    }
                                }


                            }, 100)
                        })

                    /**
                       Add edit image to button-group g as the first child
                    */
                    constant.buttonView.select("#button-group")
                        .append("image")
                        .attr("id", "edit-image-" + d.id)
                        .attr("xlink:href", function(ed, ei) {
                            if (d.type == constant.PIPELINE_START) {
                                return "../../assets/svg/edit-mouseover-1.svg";
                            } else if (d.type == constant.PIPELINE_STAGE) {
                                return "../../assets/svg/edit-mouseout-3.svg";
                            }

                        })
                        .attr("width", function(ed, ei) {
                            if (d.type == constant.PIPELINE_START) {
                                return constant.svgStageWidth * 2;
                            } else if (d.type == constant.PIPELINE_STAGE) {
                                return constant.svgStageWidth * 2 * 0.5;
                            }
                        })
                        .attr("height", function(ed, ei) {
                            if (d.type == constant.PIPELINE_START) {
                                return constant.svgStageHeight * 2;
                            } else if (d.type == constant.PIPELINE_STAGE) {
                                return constant.svgStageHeight * 2 * 0.5;;
                            }
                        })
                        .attr("transform", function(ed, ei) {
                            let translateX = 0;
                            let translateY = 0;
                            if (d.type == constant.PIPELINE_START) {
                                translateX = d.translateX - constant.svgStageWidth * 0.5;
                                translateY = d.translateY - constant.svgStageHeight * 0.5;
                            } else if (d.type == constant.PIPELINE_STAGE) {
                                translateX = d.translateX - constant.svgStageWidth * 0.5 + 1.9;
                                translateY = d.translateY - constant.svgStageHeight * 0.5 + 1.9;
                            }

                            return "translate(" + translateX + "," + translateY + ")";
                        })
                        .on("mouseover", function(ed, ei) {
                            d3.event.stopPropagation();
                            editBtnShow = true;
                            constant.buttonView.select("#edit-image-" + d.id)
                                .attr("xlink:href", function(ed, ei) {
                                    if (d.type == constant.PIPELINE_START) {
                                        return "../../assets/svg/edit-mouseover-1.svg";
                                    } else if (d.type == constant.PIPELINE_STAGE) {
                                        return "../../assets/svg/edit-mouseover-3.svg";
                                    }
                                })
                        })
                        .on("mouseout", function(ed, ei) {
                            editBtnShow = false;
                            constant.buttonView.select("#edit-image-" + d.id)
                                .attr("xlink:href", function(ed, ei) {
                                    // return "../../assets/svg/edit-grey.svg";
                                    if (d.type == constant.PIPELINE_STAGE) {
                                        return "../../assets/svg/edit-mouseout-3.svg";
                                    }
                                })
                        })
                        .on("click", function(ed, ei) {
                            constant.buttonView.selectAll("#button-group").remove();
                            if (d.type == constant.PIPELINE_START) {
                                clickStart(d, i);
                            } else if (d.type == constant.PIPELINE_STAGE) {
                                clickStage(d, i);
                            }

                        })


                    if (d.type == constant.PIPELINE_STAGE) {

                        /**
                            Add add action image to button-group g as the second child (Only for stage node)
                         */
                        constant.buttonView.select("#button-group")
                            .append("image")
                            .attr("id", "add-image-" + d.id)
                            .attr("xlink:href", function(ad, ai) {
                                return "../../assets/svg/add-mouseout.svg";
                            })
                            .attr("width", function(ad, ai) {
                                return constant.svgStageWidth * 2 * 0.5;
                            })
                            .attr("height", function(ad, ai) {
                                return constant.svgStageHeight * 2 * 0.5;
                            })
                            .attr("transform", function(ad, ai) {
                                let translateX = d.translateX + constant.svgStageWidth * 0.5 - 1.9;
                                let translateY = d.translateY - constant.svgStageHeight * 0.5 + 1.9;
                                return "translate(" + translateX + "," + translateY + ")";
                            })
                            .on("mouseover", function(ad, ai) {
                                d3.event.stopPropagation();
                                addActionBtnShow = true;
                                constant.buttonView.select("#add-image-" + d.id)
                                    .attr("xlink:href", function(ad, ai) {
                                        return "../../assets/svg/add-mouseover.svg";
                                    })
                            })
                            .on("mouseout", function(ad, ai) {
                                addActionBtnShow = false;
                                constant.buttonView.select("#add-image-" + d.id)
                                    .attr("xlink:href", function(d, ai) {
                                        return "../../assets/svg/add-mouseout.svg";
                                    })
                            })
                            .on("click", function(ad, ai) {
                                constant.buttonView.selectAll("#button-group").remove();
                                addAction(d.actions);
                                initAction();
                            })

                        /**
                           Add delete image to button-group g as the third child (Only for stage node)
                        */
                        constant.buttonView.select("#button-group")
                            .append("image")
                            .attr("id", "delete-image-" + d.id)
                            .attr("xlink:href", function(dd, di) {
                                return "../../assets/svg/delete-mouseout.svg";
                            })
                            .attr("width", function(dd, di) {
                                return constant.svgStageWidth * 2 - 0.6;
                            })
                            .attr("height", function(dd, di) {
                                return constant.svgStageHeight * 2 * 0.5;
                            })
                            .attr("transform", function(dd, di) {
                                let translateX = d.translateX - constant.svgStageWidth * 0.5 + 0.2;
                                let translateY = d.translateY + constant.svgStageHeight * 0.5 + 0.2;
                                return "translate(" + translateX + "," + translateY + ")";
                            })
                            .on("mouseover", function(dd, di) {
                                d3.event.stopPropagation();
                                deleteBtnShow = true;
                                constant.buttonView.select("#delete-image-" + d.id)
                                    .attr("xlink:href", function(d, i) {
                                        return "../../assets/svg/delete-mouseover.svg";
                                    })
                            })
                            .on("mouseout", function(dd, di) {
                                deleteBtnShow = false;
                                constant.buttonView.select("#delete-image-" + d.id)
                                    .attr("xlink:href", function(dd, di) {
                                        return "../../assets/svg/delete-mouseout.svg";
                                    })
                            })
                            .on("click", function(dd, di) {
                                constant.buttonView.selectAll("#button-group").remove();
                                $("#pipeline-info-edit").html("");
                                var timeout = 0;
                                /* if remove the node is not the last one, add animation to action */
                                if (i < pipelineData.length - 1) {
                                    timeout = 400;
                                    animationForRemove(d.id, i);
                                }
                                setTimeout(function() {
                                    deleteStage(d, i);
                                    initPipeline();
                                    initAction();
                                }, timeout)

                            })
                    }
                }

            }


        })

    .call(drag);



}
