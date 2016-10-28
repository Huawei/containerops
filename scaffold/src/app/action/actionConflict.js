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
import {getPathData} from "../relation/setPath";

export function checkConflict(fromNodeId, toNodeId) {

}
// var conflict = {
//     "node": [{
//         "name": "action1",
//         "id": "action1ID",
//         "conflicts": [
//             { "aaa": "" },
//             { "bbb": { "name": { "ddd": "" } } },
//             { "ccc": [] }
//         ]
//     }, {
//         "name": "action2",
//         "id": "action2ID",
//         "conflicts": [
//             { "ccc": "" },
//             { "eee": [] }
//         ]
//     }, {
//         "name": "action3",
//         "id": "action3ID",
//         "conflicts": [
//             { "aaa": "" },
//             { "bbb": {} },
//             { "eee": [] }
//         ]
//     }, {
//         "name": "currentAction",
//         "id": "currentActionID",
//         "conflicts": [
//             { "aaa": "" },
//             { "bbb": {} },
//             { "ccc": "" },
//             { "eee": [] }
//         ]
//     }, {
//         "name": "currentAction",
//         "id": "currentActionID",
//         "conflicts": [
//             { "aaa": "" },
//             { "bbb": {} },
//             { "ccc": "" },
//             { "eee": [] }
//         ]
//     }, {
//         "name": "currentAction",
//         "id": "currentActionID",
//         "conflicts": [
//             { "aaa": "" },
//             { "bbb": {} },
//             { "ccc": "" },
//             { "eee": [] }
//         ]
//     }],
//     "line": [{
//         "fromData": "action1ID.bbb.name.ddd",
//         "toData": "action1ID.aa.bb"
//     }]
// };

export function getConflict(targetAction) {
    let actionId = targetAction.id;
    $("#conflictTreeView").empty();
    let conflict = conflictUtil.getActionConflict(actionId);
    if (_.isEmpty(conflict)) {
        let noconflict = "<h4 class='pr'>" +
            "<em>No Conflict</em>" +
            "</h4>";
        $("#conflictTreeView").append(noconflict);
    } else {
        svgTree(d3.select("#conflictTreeView"), conflict,actionId);
    }
}


export function svgTree(container, data,actionId) {
    
    let svgWidth = ($("#actionTabsContent").width()-110)/2;
    

    let conflictActions = _.filter(data.node,function(node){
        return node.id != actionId;
    });

    let curAction = _.filter(data.node,function(node){
        return node.id == actionId;
    });


    let conflictArray = transformJson(conflictActions,60);
    let curActionArray = transformJson(curAction,svgWidth+60);

    let svg = container.append("svg")
        .attr("width", "100%")
        .attr("height", 600)
        .style("fill", "white");


    for (let i = 0; i < conflictArray.length; i++) {
        construct(svg, conflictArray[i]);
    }

    for (let i = 0; i < curActionArray.length; i++) {
        construct(svg, curActionArray[i]);
    }

    drawLine(data.line);

    function transformJson(data,initX) {

        let jsonArray = [];
        let depthY = 0;
        let depthX = 1;

        for(let i =0;i<data.length;i++){
            depthY++;
            jsonArray.push({
                depthX: depthX,
                depthY: depthY,
                type: "object",
                initX : initX,
                name: data[i].name
            });

            for (let j = 0; j < data[i].conflicts.length; j++) {

                let conflicts = data[i].conflicts[j];

                for (let key in conflicts) {
                    depthY++;

                    jsonArray.push({
                        depthX: 2,
                        depthY: depthY,
                        type: judgeType(conflicts[key]),
                        initX:initX,
                        path:data[i].name+"_"+key,
                        name: key
                    });

                    getChildJson(conflicts[key], 3,data[i].name+"_"+key);
                }
            }
        }

        function getChildJson(data, depthX,path){
             if (isObject(data)) {
                for (let key in data) {
                    depthY++;
                    jsonArray.push({
                        depthX: depthX,
                        depthY: depthY,
                        type: judgeType(data[key]),
                        initX:initX,
                        path:path+"_"+key,
                        name: key
                    });
                    getChildJson(data[key], depthX + 1);
                }
            }

            if (isArray(data) && data.length > 0) {

            }
        }

        return jsonArray;

    }



    function construct(svg, options) {
        let gLine = svg.append("g")
            .attr("id","conflictLine");


        let g = svg.append("g")
            .attr("transform", "translate(" + (options.depthX * 20 + options.initX) + "," + (options.depthY * 28) + ")")
            .attr("class",options.path)
            .attr("tx",(options.depthX * 20 + options.initX))
            .attr("ty",(options.depthY * 28));

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
            .attr("height", 20);



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

function drawLine(lineArray){
    for(let i=0;i<lineArray.length;i++){
        let start = $("."+(lineArray[i].fromData).replace(/\./g,"_") );
        let end = $("."+(lineArray[i].toData).replace(/\./g,"_"));      
        let point = [start.attr("tx"),start.attr("ty"),end.attr("tx"),end.attr("ty")];
        
        drawLinePath(point,"","");           
    }
}


function drawLinePath(point,fromPath,toPath){
    // var offsetTop = $("#bipatiteLineSvg").offset().top;
    // var offsetLeft = $("#bipatiteLineSvg").offset().left;
    var x1 = parseInt(point[0])+70;
    var y1 = parseInt(point[1]);
    var x2 = parseInt(point[2]);
    var y2 = parseInt(point[3]);
    var d = getPathData({x:x1,y:y1},{x:x2,y:y2});

    d3.select("#conflictLine")
    .append("path")
    .attr("d",d)
    .attr("stroke", "#f9f065")
    .attr("stroke-width", 6)
    .attr("fill","none")
    .attr("stroke-opacity", "0.8")
    .attr("class","cursor");

}