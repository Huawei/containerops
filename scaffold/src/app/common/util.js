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

import * as constant from "./constant";

export function isObject(o) {
    return Object.prototype.toString.call(o) == '[object Object]';
}
export function isArray(o) {
    return Object.prototype.toString.call(o) == '[object Array]';
}
export function isBoolean(o) {
    return Object.prototype.toString.call(o) == '[object Boolean]';
}
export function isNumber(o) {
    return Object.prototype.toString.call(o) == '[object Number]';
}
export function isString(o) {
    return Object.prototype.toString.call(o) == '[object String]';
}

export function findAllRelatedLines(itemId) {
    var relatedLines = _.filter(constant.linePathAry, function(item) {
        return (item.startData != undefined && item.endData != undefined) && (item.startData.id == itemId || item.endData.id == itemId)
    });
    return relatedLines;
}
export function findInputLines(itemId) {
    var relatedLines = _.filter(constant.linePathAry, function(item) {
        return (item.endData != undefined) && (item.endData.id == itemId)
    });
    return relatedLines;
}
export function findOutputLines(itemId) {
    var relatedLines = _.filter(constant.linePathAry, function(item) {
        return (item.startData != undefined) && (item.startData.id == itemId)
    });
    return relatedLines;
}
export function removeRelatedLines(args) {
    if (isString(args)) {
        var relatedLines = findAllRelatedLines(args);
        constant.setLinePathAry(_.difference(constant.linePathAry, relatedLines));
    } else {
        _.each(args, function(item) {
            removeRelatedLines(item.id);
        })
    }

}
export function findAllActionsOfStage(stageId) {
    var groupId = "#action" + "-" + stageId;
    var selector = groupId + "> image";
    return $(selector);
}
export function disappearAnimation(args) {
    if (isString(args)) {
        d3.selectAll(args)
            .transition()
            .duration(200)
            .style("opacity", 0);
    } else {
        _.each(args, function(selector) {
            disappearAnimation(selector);
        })
    }

}
export function transformAnimation(args, type) {
    _.each(args, function(item) {
        d3.selectAll(item.selector)
            .filter(function(d, i) {
                return i > item.itemIndex
            })
            .transition()
            .delay(200)
            .duration(200)
            .attr("transform", function(d, i) {
                var translateX = 0,
                    translateY = 0;
                if (type == "action") {
                    translateX = item.type == "siblings" ? d.translateX : 0;
                    translateY = item.type == "siblings" ? (d.translateY - constant.ActionNodeSpaceSize) : (0 - constant.ActionNodeSpaceSize);

                } else if (type == "stage") {
                    translateX = item.type == "siblings" ? (d.translateX - constant.PipelineNodeSpaceSize) : (0 - constant.PipelineNodeSpaceSize);
                    translateY = item.type == "siblings" ? d.translateY : 0;

                }
                return "translate(" + translateX + "," + translateY + ")";

            });
    })

}

export function judgeType(target) {
    if (isObject(target)) {
        return "object";
    } else if (isArray(target)) {
        return "array";
    } else if (isBoolean(target)) {
        return "boolean";
    } else if (isString(target)) {
        return "string";
    } else if (isNumber(target)) {
        return "number";
    } else {
        return "null";
    }
}

export function changeCurrentElement(previousData) {
    if (previousData != null) {
        switch (previousData.type) {
            case "stage":
                d3.select("#" + previousData.data.id).attr("href", "../../assets/svg/stage-latest.svg");
                break;
            case "start":
                d3.select("#" + previousData.data.id).attr("href", "../../assets/svg/start-latest.svg");
                break;
            case "action":
                d3.select("#" + previousData.data.id).attr("href", "../../assets/svg/action-bottom.svg");
                break;
            case "line":
                d3.select("#" + previousData.data.attr("id")).attr("stroke", "#E6F3E9");
                break;

        }

    }

}
