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


import { isObject, isArray, isBoolean, isNumber, isString, judgeType } from '../common/util';
import * as conflictUtil from "../relation/conflict";

import { getPathData } from "../relation/setPath";

import { notify } from "../common/notify";

export function getConflict(actionId) {
    $("#conflictTreeView").empty();
    let conflict = conflictUtil.getActionConflict(actionId);
    console.log(conflict);
    if (_.isEmpty(conflict)) {
        let noconflict = "<h4 class='pr'>" +
            "<em>No Conflict</em>" +
            "</h4>";
        $("#conflictTreeView").append(noconflict);
    } else {
        svgTree(d3.select("#conflictTreeView"), conflict, actionId);
    }
}

export function svgTree(container, data, actionId) {

    let svgWidth = ($("#actionTabsContent").width() - 110) / 2;
    let conflictActions = _.filter(data.node, function(node) {
        return node.id != actionId;
    });
    let curAction = _.filter(data.node, function(node) {
        return node.id == actionId;
    });
    let conflictArray = transformJson(conflictActions, 60);
    let curActionArray = transformJson(curAction, svgWidth + 60);
    let svg = container.append("svg")
        .attr("width", "100%")
        .attr("height", 600)
        .style("fill", "white");
    for (let i = 0; i < conflictArray.length; i++) {
        construct(svg, conflictArray[i], "conflict-source");
    }

    for (let i = 0; i < curActionArray.length; i++) {
        construct(svg, curActionArray[i], "conflict-base");
    }

    drawLine(data.line);

    function transformJson(data, initX) {

        let jsonArray = [];
        let depthY = 0;
        let depthX = 1;

        for (let i = 0; i < data.length; i++) {
            depthY++;
            jsonArray.push({
                depthX: depthX,
                depthY: depthY,
                type: "object",
                initX: initX,
                name: data[i].name,
                class: "",
                category: "action",
                parentActionId:data[i].id
            });

            for (let j = 0; j < data[i].conflicts.length; j++) {

                let conflicts = data[i].conflicts[j];

                for (let key in conflicts) {
                    depthY++;

                    jsonArray.push({
                        depthX: 2,
                        depthY: depthY,
                        type: judgeType(conflicts[key]),
                        initX: initX,
                        parentActionId:data[i].id,
                        name: key,
                        class: key,
                        category: "property"
                    });

                    getChildJson(conflicts[key], 3, data[i].id, key);

                }
            }
        }

        function getChildJson(data, depthX, parentActionId, parentPath) {
            if (isObject(data)) {
                for (let key in data) {
                    depthY++;
                    jsonArray.push({
                        depthX: depthX,
                        depthY: depthY,
                        type: judgeType(data[key]),
                        initX: initX,
                        parentActionId:parentActionId,
                        name: key,
                        class: parentPath + "." + key,
                        category: "property",

                    });
                    getChildJson(data[key], depthX + 1, parentActionId, parentPath + "." + key);
                }
            }

            if (isArray(data) && data.length > 0) {

            }
        }

        return jsonArray;

    }


    function construct(svg, options, type) {
        let gLine = svg.append("g")
            .attr("id", "conflictLine");

        let g = svg.append("g")
            .attr("transform", "translate(" + (options.depthX * 20 + options.initX) + "," + (options.depthY * 28) + ")")
            .attr("id",options.parentActionId+"_"+options.class.replace(/\./g, "_"))
            .attr("tx", (options.depthX * 20 + options.initX))
            .attr("ty", (options.depthY * 28))
            .attr("data-type", type)
            .attr("data-class", options.class.replace(/\./g, "_"))
            .attr("data-clean",options.class)
            .style("cursor", "pointer")

        .on("mouseover", function() {
                // var selectedConflict = callFunction(options.path);
                if (options.category == "property") {
                    // d3.selectAll("[data-class=" + options.class + "]").select("rect")
                    //     .attr("fill", "#333");
                    d3.selectAll("[data-class=" + options.class.replace(/\./g, "_") + "]").each(function(d, i) {
                        var d3DOM = d3.select(this);
                        d3DOM.select("rect").attr("fill", "#333");
                        var elementType = d3DOM.attr("data-type");
                        if (elementType == "conflict-line") {
                            d3DOM.attr("stroke", "#333");
                        } else {
                            d3DOM.select(".conflict-image").attr("xlink:href", function() {
                                if (elementType == "conflict-source") {
                                    return "../../assets/svg/remove-conflict.svg";
                                } else if (elementType == "conflict-base") {
                                    return "../../assets/svg/highlight-conflict.svg";
                                }
                            });
                        }

                    })

                }

            })
            .on("mouseout", function() {
                if (options.category == "property") {
                    d3.selectAll("[data-class=" + options.class.replace(/\./g, "_") + "]").select("rect")
                        .attr("fill", function() {
                            switch (options.type) {
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
                        });

                    // d3.selectAll("[data-class=" + options.class + "]").select(".conflict-image")
                    //     .attr("xlink:href", "../../assets/svg/conflict.svg");
                    d3.selectAll("[data-class=" + options.class.replace(/\./g, "_") + "]").each(function() {
                        var d3DOM = d3.select(this);
                        var type = d3DOM.attr("data-type");
                        if (type == "conflict-line") {
                            d3DOM.attr("stroke", "#f9f065")
                        } else {
                            d3DOM.select(".conflict-image")
                                .attr("xlink:href", "../../assets/svg/conflict.svg");
                        }
                    })
                }

            })
        let rect = g.append('rect')
            .attr("ry", 4)
            .attr("rx", 4)
            .attr("y", 0)
            .attr("width", 135)
            .attr("height", 24)
            .attr("fill", function() {
                switch (options.type) {
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
                }
            });

        let clashImage = g.append('image')
            .attr("transform", "translate(0,0)")
            .attr("xlink:href", "../../assets/svg/conflict.svg")
            .attr("x", 2)
            .attr("y", 2)
            .attr("width", 20)
            .attr("height", 20)
            .attr("class", "conflict-image")
            .on("click", function() {
                if (type == "conflict-source") {
                    
                    if(options.parentActionId){
                        var conflictSourceId = options.parentActionId;
                        conflictUtil.cleanConflict(conflictSourceId,actionId,options.class);
                        getConflict(actionId);
                        notify("Remove conflict", "success");
                        return false;
                    }
                    
                }

            });
        let typeImage = g.append('image')
            .attr("transform", "translate(115,0)")
            .attr("xlink:href", function() {
                switch (options.type) {
                    case "string":
                        return "../../assets/images/string.png";
                        break;
                    case "object":
                        return "../../assets/images/object.png";
                        break;
                    case "number":
                        return "../../assets/images/number.png";
                        break;
                    case "array":
                        return "../../assets/images/array.png";
                        break;
                    case "boolean":
                        return "../../assets/images/boolean.png";
                        break;
                    default:
                        return "";
                }
            })
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

function drawLine(lineArray) {
    for (let i = 0; i < lineArray.length; i++) {
        let start = $("#" + (lineArray[i].fromData).replace(/\./g, "_"));
        let end = $("#" + (lineArray[i].toData).replace(/\./g, "_"));
        let point = [start.attr("tx"), start.attr("ty"), end.attr("tx"), end.attr("ty")];
        var dataClass = _.rest(lineArray[i].fromData.split(".")).join("_");
        // drawLinePath(point, "", "");
        drawLinePath(point, dataClass, "conflict-line");
    }
}


function drawLinePath(point, dataClass, type, fromPath, toPath) {
    // var offsetTop = $("#bipatiteLineSvg").offset().top;
    // var offsetLeft = $("#bipatiteLineSvg").offset().left;
    var x1 = parseInt(point[0]) + 70;
    var y1 = parseInt(point[1]);
    var x2 = parseInt(point[2]);
    var y2 = parseInt(point[3]);
    var d = getPathData({ x: x1, y: y1 }, { x: x2, y: y2 });

    d3.select("#conflictLine")
        .append("path")
        .attr("d", d)
        .attr("stroke", "#f9f065")
        .attr("stroke-width", 6)
        .attr("fill", "none")
        .attr("stroke-opacity", "0.8")
        .attr("class", "cursor")
        .attr("data-class", dataClass)
        .attr("data-type", type);

}
