import * as constant from "../common/constant";
import { setPath } from "../relation/setPath";
import { drag } from "../common/drag";
import { mouseoverRelevantPipeline, mouseoutRelevantPipeline } from "../relation/lineHover";
import { dragDropSetPath } from "../relation/dragDropSetPath";

export function initLine() {

    constant.linesView.selectAll("g").remove();

    var diagonal = d3.svg.diagonal();

    var pipelineLineViewId = "pipeline-line-view";

    constant.lineView[pipelineLineViewId] = constant.linesView.append("g")
        .attr("width", constant.svgWidth)
        .attr("height", constant.svgHeight)
        .attr("id", pipelineLineViewId);

    constant.pipelineView.selectAll("image").each(function(d, i) {

        /* draw the main line of pipeline */
        if (i != 0) {
            constant.lineView[pipelineLineViewId]
                .append("path")
                .attr("d", function() {
                    return diagonal({
                        source: { x: d.translateX - constant.PipelineNodeSpaceSize, y: constant.pipelineNodeStartY + 22.5 },
                        target: { x: d.translateX + 2, y: constant.pipelineNodeStartY + 22.5 }
                    });
                })
                .attr("fill", "none")
                .attr("stroke", "#333")
                .attr("stroke-width", 2);
        }
        if (d.type == constant.PIPELINE_START) {
            /* draw the vertical line and circle for start node  in lineView -> pipeline-line-view */
            constant.lineView[pipelineLineViewId]
                .append("path")
                .attr("d", function() {
                    return diagonal({
                        source: { x: d.translateX + constant.svgStageWidth / 2, y: constant.pipelineNodeStartY + constant.svgStageWidth / 2 },
                        target: { x: d.translateX + constant.svgStageWidth / 2, y: constant.pipelineNodeStartY + constant.svgStageWidth + 6 }
                    })
                })
                .attr("fill", "none")
                .attr("stroke", "#aaa")
                .attr("stroke-width", 1);
            constant.lineView[pipelineLineViewId]
                .append("circle")
                .attr("cx", function(cd, ci) {
                    return d.translateX + constant.svgStageWidth / 2;
                })
                .attr("cy", function(cd, ci) {
                    return constant.pipelineNodeStartY + constant.svgStageWidth + 12;
                })
                .attr("r", function(cd, ci) {
                    return 8;
                })
                .attr("fill", "#fff")
                .attr("stroke", "#aaa")
                .attr("stroke-width", 2)
                /* mouse over the circle show relevant lines of start stage */
                .on("mouseover", function(cd, ci) {
                    mouseoverRelevantPipeline(d);
                })
                /* mouse over the circle to draw line from start stage */
                .on("mousedown", function(cd, ci) {
                    d3.event.stopPropagation();
                    dragDropSetPath({
                        "data": d,
                        "node": i
                    });
                })
                .on("mouseout", function(cd, ci) {
                    mouseoutRelevantPipeline(d);
                })

        }
        /* draw line from action 2 stage and circle of action self to accept and emit lines  */
        if (d.type == constant.PIPELINE_STAGE && d.actions != null && d.actions.length > 0) {

            var actionLineViewId = "action-line" + "-" + d.id;
            var action2StageLineViewId = "action-2-stage-line" + "-" + d.id;
            var actionSelfLine = "action-self-line" + "-" + d.id
                /* Action 2 Stage */
            constant.lineView[action2StageLineViewId] = constant.linesView.append("g")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight)
                .attr("id", action2StageLineViewId);

            constant.lineView[action2StageLineViewId].selectAll("path")
                .data(d.actions).enter()
                .append("path")
                .attr("d", function(ad, ai) {
                    /* draw the tail line of action */
                    constant.lineView[action2StageLineViewId]
                        .append("path")
                        .attr("d", function(fd, fi) {
                            return diagonal({
                                source: { x: ad.translateX + 15, y: ad.translateY + 28 },
                                target: { x: ad.translateX + 15, y: ad.translateY + 40 }
                            });
                        })
                        .attr("fill", "none")
                        .attr("stroke", "black")
                        .attr("stroke-width", 1)
                        .attr("stroke-dasharray", "2,2");
                    /* draw different length line group by stage index */
                    if (i % 2 == 0) {
                        return diagonal({
                            source: { x: ad.translateX + 15, y: ad.translateY },
                            target: { x: ad.translateX + 15, y: ad.translateY - 40 }
                        });
                    } else {
                        return diagonal({
                            source: { x: ad.translateX + 15, y: ad.translateY },
                            target: { x: ad.translateX + 15, y: ad.translateY - 70 }
                        });
                    }
                })
                .attr("fill", "none")
                .attr("stroke", "black")
                .attr("stroke-width", 1)
                .attr("stroke-dasharray", "2,2");

            /* line across action to connect two circles */
            constant.lineView[actionSelfLine] = constant.linesView.append("g")
                .attr("width", constant.svgWidth)
                .attr("height", constant.svgHeight)
                .attr("id", actionSelfLine);

            constant.lineView[actionSelfLine].selectAll("path")
                .data(d.actions).enter()
                .append("path")
                .attr("d", function(ad, ai) {
                    return diagonal({
                        source: { x: ad.translateX - 8, y: ad.translateY + constant.svgActionHeight / 2 },
                        target: { x: ad.translateX + constant.svgActionWidth + 8, y: ad.translateY + constant.svgActionHeight / 2 }
                    })
                })
                .attr("id", function(ad, ai) {
                    return "action-self-line-path-" + ad.id;
                })
                // .attr("transform", function(ad, ai){
                //     return "translate(0,0)"
                // })
                .attr("fill", "none")
                .attr("stroke", "#aaa")
                .attr("stroke-width", 1);

            /* circle on the left */
            constant.lineView[actionSelfLine].selectAll(".action-self-line-input")
                .data(d.actions).enter()
                .append("circle")
                .attr("class", "action-self-line-input")
                .attr("cx", function(ad, ai) {
                    return ad.translateX - 16;
                })
                .attr("cy", function(ad, ai) {
                    return ad.translateY + constant.svgActionHeight / 2;
                })
                .attr("r", function(ad, ai) {
                    return 8;
                })
                .attr("id", function(ad, ai) {
                    return "action-self-line-input-" + ad.id;
                })
                // .attr("transform", function(ad, ai){
                //      return "translate(0,0)"
                // })
                .attr("fill", "#fff")
                .attr("stroke", "#aaa")
                .attr("stroke-width", 2)
                /* mouse down to drag lines */
                .on("mousedown", function(ad, ai) {
                    d3.event.stopPropagation();
                    dragDropSetPath({
                        "data": ad,
                        "node": ai
                    });
                })
                // .on("mouseover", function(ad, ai){
                //     var translateX = 0 ;
                //     var translateY = 0;
                //     d3.select(this).attr("transform","translate(-250,-190) scale(2)");
                // })
                // .on("mouseup", function(ad, ai){
                //     var translateX = ad.translateX - 16;
                //     var translateY =  ad.translateY + constant.svgActionHeight/2;
                //     d3.select(this).attr("transform","translate(0,0) scale(1)");
                // })
                // .on("mouseout", function(ad, ai){
                //     var translateX = ad.translateX - 16;
                //     var translateY =  ad.translateY + constant.svgActionHeight/2;
                //     d3.select(this).attr("transform","translate(0,0) scale(1)");
                // })
                /* circle on the right */
            constant.lineView[actionSelfLine].selectAll(".action-self-line-output")
                .data(d.actions).enter()
                .append("circle")
                .attr("class", "action-self-line-output")
                .attr("cx", function(ad, ai) {
                    return ad.translateX + constant.svgActionWidth + 16;
                })
                .attr("cy", function(ad, ai) {
                    return ad.translateY + constant.svgActionHeight / 2;
                })
                .attr("r", function(ad, ai) {
                    return 8;
                })
                .attr("id", function(ad, ai) {
                    return "action-self-line-output-" + ad.id
                })
                // .attr("transform", function(ad, ai){                     
                //      return "translate(0,0)"
                // })
                .attr("fill", "#fff")
                .attr("stroke", "#aaa")
                .attr("stroke-width", 2)
                .on("mouseover", function(ad, ai) {
                    mouseoverRelevantPipeline(ad);
                })
                .on("mousedown", function(ad, ai) {
                    d3.event.stopPropagation();
                    dragDropSetPath({
                        "data": ad,
                        "node": ai
                    });
                })
                .on("mouseout", function(ad, ai) {
                    mouseoutRelevantPipeline(ad);
                })
        }

    });

    constant.linePathAry.forEach(function(i) {
        setPath(i);
    })
}
