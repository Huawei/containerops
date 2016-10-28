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

import { linePathAry } from "../common/constant";

export  function hasConflict(startActionID,endActionID) {
    let result = false;
    let receiveData = {};

    for(let i = 0; i < linePathAry.length; i ++) {
        let lineInfo = linePathAry[i];
        if (lineInfo.endData.id == endActionID && lineInfo.startData.id == startActionID) {
            for (let j = 0; j < lineInfo.relation.length; j ++) {
                let currentRelation = lineInfo.relation[j];

                receiveData[currentRelation.to] = true;
            }
            break;
        }
    }

    for (let i =0; i < linePathAry.length; i ++) {
        let lineInfo = linePathAry[i];
        if (lineInfo.endData.id == endActionID && lineInfo.startData.id != startActionID) {
            for (let j = 0; j < lineInfo.relation.length; j ++) {
                let currentRelation = lineInfo.relation[j];

                if (receiveData[currentRelation.to]) {
                    result = true;
                    break;
                }
            }
        }
    }

    return result;
}

export function getActionConflict(actionID) {
    let result = {};
    let conflicts = {};
    let actionReceiveData = {};

    for (let i = 0; i < linePathAry.length; i ++) {
        let lineInfo = linePathAry[i];
        if (lineInfo.endData.id == actionID) {
            for (let j = 0; j < lineInfo.relation.length; j ++) {
                let currentRelation = lineInfo.relation[j];
                let currentPath = lineInfo.startData.id + currentRelation.from

                if (!actionReceiveData[currentRelation.to]) { actionReceiveData[currentRelation.to] = [];}

                actionReceiveData[currentRelation.to].push(currentPath);
            }
        }
    }

    for (let p in actionReceiveData) {
        if (actionReceiveData[p].length > 1) {
            for (let i = 0; i < actionReceiveData[p].length; i ++) {
                let fromPath = actionReceiveData[p][i];
                let fromActionId = fromPath.split(".")[0]
                let fromNodePath = fromPath.substring(fromPath.indexOf("."))
                let line = {};

                if (!result.line) {result.line = [];}
                if (!conflicts[actionID]) {conflicts[actionID] = {};}
                if (!conflicts[fromActionId]) {conflicts[fromActionId] = {};}

                conflicts[fromActionId] = setConflictPath(conflicts[fromActionId], fromNodePath);
                conflicts[actionID] = setConflictPath(conflicts[actionID], p);

                line.fromData = fromPath;
                line.toData = actionID + p;
                result.line.push(line);
            }
        }
    }

    for (let p in conflicts) {
        let node = {};
        let nodeConflicts = [];
        if (!result.node) {result.node = [];}

        for (let prop in conflicts[p]) {
            let nodeConflict = {};
            nodeConflict[prop] = conflicts[p][prop];

            nodeConflicts.push(nodeConflict);
        }

        node.id = p;
        node.conflicts = nodeConflicts;

        result.node.push(node);
    }

    return result;
}

function setConflictPath(obj,path) {
    path = path.substring(1);
    currentProp = path.split(".")[0];

    if (path.split(".").length > 1) {
        if (!obj[currentProp]) {obj[currentProp] = {};}
        obj[currentProp] = setConflictPath(obj[currentProp], path.substring(path.indexOf(".")));            
    } else {
        if (!obj[currentProp]) {obj[currentProp] = "";}
    }

    return obj;
}
