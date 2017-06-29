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

import * as util from '../common/util';
import * as conflictUtil from "../relation/conflict";
import * as config from '../common/config';

import { getPathData } from "../relation/setPath";

import { notify } from "../common/notify";

let drag = d3.behavior.drag()
    .origin(function(d) {
        return { "x": 0, "y": 0 };
    })
    .on("dragstart", dragStart)
    .on("drag", util.draged);

function dragStart() {
    d3.event.sourceEvent.stopPropagation();
    drag.origin(function() {
        return { "x": d3.select(this).attr("translateX"), "y": d3.select(this).attr("translateY") }
    });
}

function caclX(id, actionId) {
    let svgWidth = ($("#actionTabsContent").width() - 110) / 2;
    let x = 0;
    if (id != actionId) {
        x = svgWidth / 5;
    } else {
        x = svgWidth + svgWidth / 5;
    }
    return x;
}

function parseNodeType(id, actionId) {
    let type = "";
    if (id == actionId) {
        type = "conflict-base";
    } else {
        type = "conflict-source";
    }
    return type;
}

function generateImage(type) {
    return config.getImage(type);
}

function getDOMData(dom, conflictData) {
    var data = _.find(conflictData, function(item) {
        return dom.attr("id") == item.parentActionId + "_" + item.class.replace(/\./g, "_");
    })
    return data;
}

function drawLine(lineArray) {
    for (let i = 0; i < lineArray.length; i++) {
        let start = $("#" + (lineArray[i].fromData).replace(/\./g, "_"));
        let end = $("#" + (lineArray[i].toData).replace(/\./g, "_"));
        let point = [start.attr("tx"), start.attr("ty"), end.attr("tx"), end.attr("ty")];
        var dataClass = _.rest(lineArray[i].fromData.split(".")).join("_");
        drawLinePath(point, dataClass, "conflict-line");
    }
}

function drawLinePath(point, dataClass, type, fromPath, toPath) {
    var x1 = parseInt(point[0]) + 70;
    var y1 = parseInt(point[1]);
    var x2 = parseInt(point[2]);
    var y2 = parseInt(point[3]);
    var d = getPathData({ x: x1, y: y1 }, { x: x2, y: y2 });

    d3.select("#conflictLine")
        .append("path")
        .attr("d", d)
        .attr("stroke", "#e0e004")
        .attr("stroke-width", 6)
        .attr("fill", "none")
        .attr("stroke-opacity", "0.8")
        .attr("class", "cursor")
        .attr("data-class", dataClass)
        .attr("data-type", type);

}
function drawTree(conflict, actionId){
    let options = { "horizonSpace": 13, "verticalSpace": 10, "backgroundY": 1 };
    let scaleObj = { "zoomScale": 1, "zoomTargetScale": 1 };
    if (_.isEmpty(conflict)) {
        d3.select("#conflictTreeView > svg").append("text")
        .attr("x", 20)
        .attr("y", 50)
        .style("fill", "#555")
        .style("font-size","25px")
        .style("font-weight","700")
        .text("No Conflict");
    } else {
        util.showZoomBtn(1, "zoomin", d3.select("#conflictZoomBtn"), d3.select("#conflictTreeSVGView"), scaleObj, options);
        util.showZoomBtn(2, "zoomout", d3.select("#conflictZoomBtn"), d3.select("#conflictTreeSVGView"), scaleObj, options);
        svgTree(conflict, actionId);
    }
}

export function getConflict(actionId) {
    $("#conflictTreeView").empty();
    // showStartStageCondition(actionId);
    let svg = d3.select("#conflictTreeView").append("svg")
        .attr("width", "100%")
        .attr("height", "100%");
    let conflictTreeMainView = svg.append("g")
        .attr("id", "conflictTreeSVGView")
        .data([{ "name": "conflictTree" }])
        .attr("width", "100%")
        .attr("height", "100%")
        .attr("translateX", 0)
        .attr("translateY", 0)
        .attr("scale", 1)
        .attr("transform", "translate(0,0) scale(1)")
        .call(drag)
    let rect = conflictTreeMainView.append("rect")
        .attr("width", "100%")
        .attr("height", "100%")
        .attr("fill", "white")
        .attr("fill-opacity", 0)
        .style("cursor", "pointer")
        // .data(conflict.node);
    let cLine = conflictTreeMainView.append("g")
        .attr("id", "conflictLine")
        // .data(conflict.line);
    let cNode = conflictTreeMainView.append("g")
        .attr("id", "conflictNode");
    let buttonMainView = svg.append("g")
        .attr({
            "width": "100%",
            "height": "23",
            "fill": "white",
            "id":"conflictZoomBtn"
        })
    let conflict = conflictUtil.getActionConflict(actionId);

    drawTree(conflict, actionId);
}
function showStartStageCondition(actionId){
   var validation = util.hasLinkWithStartStage(actionId);
   if(validation.hasLink){
      var option = _.template('<option value="<%= value %>"><%= value %></option>');
      var options = "";
      _.each(validation.startData.outputJson, function(item){
           options += option({'value':item.event+"_"+item.type});
      })
      // var template = '<select>' + options + '</select>';
      var template = '<div class="row">'+
                        '<label class="control-label">' +
                            'Select Output:' +
                        '</label>' +
                        '<div class="input-group">' +
                            '<select id="start-stage-condition" style="width:100%">' +
                                options +
                            '</select>' +
                        '</div>' +
                     '</div>'
      $("#conflictTreeView").append(template);
      $("#start-stage-condition").select2({
                minimumResultsForSearch: Infinity
       });
   }

}
export function redrawTree(actionId) {
    d3.selectAll("#conflictNode > g").remove();
    d3.selectAll("#conflictLine > path").remove();
    d3.selectAll("#conflictZoomBtn > image").remove();
    let conflict = conflictUtil.getActionConflict(actionId);
    drawTree(conflict, actionId);
}

export function svgTree(data, actionId) {
    let conflictData = transformJson(data.node);
    var container = d3.select("#conflictNode");
    for (let i = 0; i < conflictData.length; i++) {
        var type = conflictData[i].parentActionId == actionId ? "conflict-base" : "conflict-source";
        construct(container, conflictData[i], type);
    }
    drawLine(data.line);

    function transformJson(data) {
        let jsonArray = [];
        let depthYLeft = 0;
        let depthYRight = 0;
        let depthX = 1;
        for (let i = 0; i < data.length; i++) {
            var initX = caclX(data[i].id, actionId);
            var dataType = parseNodeType(data[i].id, actionId);
            dataType == "conflict-source" ? depthYLeft++ : depthYRight++;
            jsonArray.push({
                depthX: depthX,
                depthY: dataType == "conflict-source" ? depthYLeft : depthYRight,
                type: "object",
                initX: initX,
                name: data[i].name,
                class: "",
                category: "action",
                parentActionId: data[i].id
            });

            for (let j = 0; j < data[i].conflicts.length; j++) {

                let conflicts = data[i].conflicts[j];

                for (let key in conflicts) {
                    dataType == "conflict-source" ? depthYLeft++ : depthYRight++;
                    jsonArray.push({
                        depthX: 2,
                        depthY: dataType == "conflict-source" ? depthYLeft : depthYRight,
                        type: util.judgeType(conflicts[key]),
                        initX: initX,
                        parentActionId: data[i].id,
                        name: key,
                        class: key,
                        category: util.judgeType(conflicts[key]) == "object" ? "path" : "property"
                    });

                    getChildJson(conflicts[key], 3, data[i].id, key);

                }
            }
        }

        function getChildJson(data, depthX, parentActionId, parentPath) {
            if (util.isObject(data)) {
                for (let key in data) {
                    // depthY++;
                    dataType == "conflict-source" ? depthYLeft++ : depthYRight++;
                    jsonArray.push({
                        depthX: depthX,
                        depthY: dataType == "conflict-source" ? depthYLeft : depthYRight,
                        type: util.judgeType(data[key]),
                        initX: initX,
                        parentActionId: parentActionId,
                        name: key,
                        class: parentPath + "." + key,
                        category: util.judgeType(data[key]) == "object" ? "path" : "property"

                    });
                    getChildJson(data[key], depthX + 1, parentActionId, parentPath + "." + key);

                }
            }

            if (util.isArray(data) && data.length > 0) {

            }
        }

        return jsonArray;

    }


    function construct(svg, options, type) {

        let g = svg.append("g")
            .data(data.node)
            .attr("transform", "translate(" + (options.depthX * 20 + options.initX) + "," + (options.depthY * 28) + ")")
            .attr("id", options.parentActionId + "_" + options.class.replace(/\./g, "_"))
            .attr("tx", (options.depthX * 20 + options.initX))
            .attr("ty", (options.depthY * 28))
            .attr("data-valuetype", options.type)
            .attr("data-type", type)
            .attr("data-class", options.class.replace(/\./g, "_"))
            .attr("data-clean", options.class)
            .style("cursor", "pointer")
            .on("mouseover", function() {
                mouseoverOrClick(options, "#797979");
            })
            .on("mouseout", function() {
                mouseout(options);
            })
            .on("click", function(d, i) {
                if (options.category == "property") {
                    d3.selectAll("[data-status=after-click").each(function(d, i) {
                        var d3DOM = d3.select(this);
                        var elementType = d3DOM.attr("data-type");
                        if (elementType != "conflict-line") {
                            var valueType = d3DOM.attr("data-valuetype");
                            d3DOM.select("rect").attr("fill", generateColor(valueType));
                            d3DOM.select(".conflict-image").attr("xlink:href", config.getSVG(config.SVG_CONFLICT));
                            var data = getDOMData(d3DOM, conflictData);
                            d3DOM.on("mouseover", function() {
                                    mouseoverOrClick(data, "#797979");
                                })
                                .on("mouseout", function() {
                                    mouseout(data);
                                })
                                .attr("data-status", "");
                        } else {
                            d3DOM.attr("stroke", "#e0e004")
                                .attr("data-status", "");
                        }
                    })

                    d3.selectAll("[data-class=" + options.class.replace(/\./g, "_") + "]").each(function(d, i) {
                        d3.select(this)
                            .on("mouseover", null)
                            .on("mouseout", null)
                            .attr("data-status", "after-click");

                    });

                }
                mouseoverOrClick(options, "#333");


            })
        let rect = g.append('rect')
            .attr("ry", 5)
            .attr("rx", 5)
            .attr("y", 0)
            .attr("width", 135)
            .attr("height", 24)
            .attr("fill", generateColor(options.type));

        let clashImage = g.append('image')
            .attr("transform", "translate(0,0)")
            .attr("xlink:href", config.getSVG(config.SVG_CONFLICT))
            .attr("x", 2)
            .attr("y", 2)
            .attr("width", 20)
            .attr("height", 20)
            .attr("class", "conflict-image")
            .on("click", function() {
                d3.event.stopPropagation();
                if (type == "conflict-source" && options.category == "property") {
                    if (options.parentActionId) {
                        var conflictSourceId = options.parentActionId;
                        conflictUtil.cleanConflict(conflictSourceId, actionId, options.class);
                        redrawTree(actionId);
                        notify("Remove conflict successfully", "success");
                        return false;
                    }
                }

            });
        let typeImage = g.append('image')
            .attr("transform", "translate(115,0)")
            .attr("xlink:href", generateImage(options.type))
            .attr("x", "0")
            .attr("y", "0")
            .attr("width", "20")
            .attr("height", "24")
            .attr("class", "type-image")

        let text = g.append('text')
            .attr("dx", 28)
            .attr("dy", 17)
            .attr("fill", function() {
                if (options.type == "null") {
                    return "#8e8a89";
                } else {
                    return "#fff";
                }
            })
            .text(function() {
                if (options.name.length > 12) {
                    return options.name.substring(0, 10) + "...";
                } else {
                    return options.name;
                }
            });
    }
}

function mouseoverOrClick(options, color) {
    if (options.category == "property") {
        d3.selectAll("[data-class=" + options.class.replace(/\./g, "_") + "]").each(function(d, i) {
            var d3DOM = d3.select(this);
            d3DOM.select("rect").attr("fill", color);
            var elementType = d3DOM.attr("data-type");
            if (elementType == "conflict-line") {
                d3DOM.attr("stroke", color);
            } else {
                d3DOM.select(".conflict-image").attr("xlink:href", function() {
                    if (elementType == "conflict-source") {
                        if (options.type != "object") {
                            return config.getSVG(config.SVG_REMOVE_CONFLICT);
                        } else {
                            return config.getSVG(config.SVG_HIGHLIGHT_CONFLICT);
                        }

                    } else if (elementType == "conflict-base") {
                        return config.getSVG(config.SVG_HIGHLIGHT_CONFLICT);
                    }
                });
            }

        })

    }
}

function mouseout(options) {
    if (options.category == "property") {
        d3.selectAll("[data-class=" + options.class.replace(/\./g, "_") + "]").select("rect")
            .attr("fill", generateColor(options.type));

        d3.selectAll("[data-class=" + options.class.replace(/\./g, "_") + "]").each(function() {

            var d3DOM = d3.select(this);
            var type = d3DOM.attr("data-type");
            if (type == "conflict-line") {
                d3DOM.attr("stroke", "#e0e004")
            } else {
                d3DOM.select(".conflict-image")
                    .attr("xlink:href", config.getSVG(config.SVG_CONFLICT));
            }
        })
    }
}

function generateColor(type) {
    switch (type) {
        case "string":
            return "#13b5b1";
            break;
        case "object":
            return "#eb6876";
            break;
        case "number":
            return "#32b16c";
            break;
        case "array":
            return "#c490c0";
            break;
        case "boolean":
            return "#8fc320";
            break;
        default:
            return "#cfcfcf";
    };
}
