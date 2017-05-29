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

import * as util from "../common/util";
import * as constant from "../common/constant";
import * as initButton from "../workflow/initButton";
import { editLine } from "./editLine";

export function mouseoverRelevantWorkflow(param) {
    var outputLines = util.findOutputLines(param.id);
    _.each(outputLines, function(line) {
        d3.select("#" + line.id).attr("stroke", function() {
            makeFrontLayer(this);
            return "#81D9EC";
        });
    });
}


export function mouseoutRelevantWorkflow(param) {
    var outputLines = util.findOutputLines(param.id);
    var tempLines = constant.linePathAry;
    if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "line") {
        var id = constant.currentSelectedItem.data.attr("id");
        var currentLineData = _.find(constant.linePathAry, function(line) {
            return id == line.id;
        })
        tempLines = _.without(outputLines, currentLineData);
    }

    _.each(tempLines, function(line) {
        d3.select("#" + line.id).attr("stroke", function() {
            makeBackLayer(this);
            return "#E6F3E9";
        });
    });
}

export function makeFrontLayer(element) {
    element.parentNode.appendChild(element);
}

export function makeBackLayer(element) {
    var firstChild = element.parentNode.firstChild;
    if (firstChild) {
        element.parentNode.insertBefore(element, firstChild);
    }
}

export function showOutputLines(param, i) {
    $.ajax({
        url: "../../templates/relation/showOutputLine.html",
        type: "GET",
        cache: false,
        success: function(data) {
            var outputLines = util.findOutputLines(param.id);
            var template = _.template(data)({ "lines": outputLines });
            $("#div-d3-lines-table").html($(template));
            $("#div-d3-lines-table").css("top",$("#div-d3-main-svg > svg").offset().top+ 2 * constant.buttonVerticalSpace + constant.buttonHeight + 15);
            $(".output-line-tr").css('cursor', 'pointer');
            $(".output-line-tr").on("mouseover", function() {
                    $(this).css("background-color", '#48a746');
                    $(this).css("color", 'white');
                    highlightSelectedLine($(this).data("lineid"), true);
                })
                .on("mouseout", function() {
                    $(this).css("background-color", 'white');
                    $(this).css("color", '#555');
                    highlightSelectedLine($(this).data("lineid"), false);
                })
                .on("click", function() {
                    // $(this).css("color", '#555');
                    $(this).css("background-color", '#81D9EC');
                    clickWorkflowLine(d3.select("#" + $(this).data("lineid"))[0][0])
                })
        }
    });
}
export function clickWorkflowLine(element) {
    element.parentNode.appendChild(element); // make this line to front layer
    var self = $(element);
    util.changeCurrentElement(constant.currentSelectedItem);
    constant.setCurrentSelectedItem({ "data": self, "type": "line" });
    initButton.updateButtonGroup("line");
    d3.select(element).attr("stroke", "#81D9EC");
    $.ajax({
        url: "../../templates/relation/editLine.html",
        type: "GET",
        cache: false,
        success: function(data) {
            editLine(data, self);
        }
    });
}

function highlightSelectedLine(lineId, highlight) {
    if (highlight) {
        makeFrontLayer(d3.select("#" + lineId)[0][0]);
        d3.select("#" + lineId).attr("stroke", "#48a746");
    } else {
        
        if (constant.currentSelectedItem != null && constant.currentSelectedItem.type == "line" && constant.currentSelectedItem.data.attr("id") == lineId) {
            makeFrontLayer(d3.select("#" + lineId)[0][0]);
            d3.select("#" + lineId).attr("stroke", "#81D9EC");
        } else {
            makeBackLayer(d3.select("#" + lineId)[0][0]);
            d3.select("#" + lineId).attr("stroke", "#E6F3E9");
        }
    }
}
