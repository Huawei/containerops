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
import { pipelineData } from "../pipeline/main";

export  function hasConflict(startActionID,endActionID) {
    let result = false;
    let receiveData = {};

    for(let i = 0; i < linePathAry.length; i ++) {
        let lineInfo = linePathAry[i];
        if (lineInfo.endData.id == endActionID && lineInfo.startData.id == startActionID && lineInfo.relation) {
            for (let j = 0; j < lineInfo.relation.length; j ++) {
                let currentRelation = lineInfo.relation[j];

                receiveData[currentRelation.to] = true;
            }
            break;
        }
    }

    for (let i =0; i < linePathAry.length; i ++) {
        let lineInfo = linePathAry[i];
        if (lineInfo.endData.id == endActionID && lineInfo.startData.id != startActionID && lineInfo.relation) {
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
        if (lineInfo.endData.id == actionID && lineInfo.relation) {
            actionReceiveData = setReceiveData(actionReceiveData, lineInfo.startData.id, lineInfo.relation);
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

                let fromAction = getAction(fromActionId);
                let toAction = getAction(actionID);

                let fromActionValue;
                if (fromAction.outputJson) {
                    getObjValue(fromAction.outputJson,fromNodePath)
                } else {
                    fromActionValue = null;
                }

                let toActionValue;
                if (toAction.inputJson) {
                    getObjValue(toAction.inputJson,p)
                } else {
                    toActionValue = null;
                }

                conflicts[fromActionId] = setConflictPath(conflicts[fromActionId], fromNodePath, fromActionValue);
                conflicts[actionID] = setConflictPath(conflicts[actionID], p, toActionValue);

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

        let action = getAction(p);
        let actionName = "";
        if (action.setupData.action.name) {
            actionName = action.setupData.action.name;
        }

        node.id = p;
        node.name = actionName;
        node.conflicts = nodeConflicts;

        result.node.push(node);
    }

    return result;
}

function setReceiveData(actionReceiveData, actionId, relationList) {
    let allLeafNodes = [];
    for (let i = 0; i < relationList.length; i ++ ) {
        let relation = relationList[i];
        let isLeafNode = true;

        for (let j = 0; j < relationList.length; j ++) {
            if ((relation.from+".").indexOf(relationList[j].from+".") == -1) {
                isLeafNode = false;
                break;
            }
        }

        if (isLeafNode) {
            relation.finalPath = actionId + relation.from;
            allLeafNodes.push(relation);
        }
    }

    for (let i = 0; i < allLeafNodes.length; i ++ ) {
        let currentRelation = allLeafNodes[i];

        if (!actionReceiveData[currentRelation.to]) { actionReceiveData[currentRelation.to] = [];}

        actionReceiveData[currentRelation.to].push(currentRelation.finalPath);
    }

    return actionReceiveData;
}

export function cleanConflict(fromActionId,toActionId,path){
    let ary = path.split("_");
    let formPath = "";
    for(var i=1;i<ary.length;i++){
      formPath = formPath+"."+ary[i];
    }

    let line = _.find(linePathAry, function(line) {
        return (fromActionId+""+toActionId) == line.id;
    })

    line.relation = delRelation(line.relation,formPath);
}

function setConflictPath(obj,path,info) {
    path = path.substring(1);
    let currentProp = path.split(".")[0];

    if (path.split(".").length > 1) {
        if (!obj[currentProp]) {obj[currentProp] = {};}
        obj[currentProp] = setConflictPath(obj[currentProp], path.substring(path.indexOf(".")), info);            
    } else {
        if (!obj[currentProp]) {obj[currentProp] = info;}
    }

    return obj;
}

function delRelation(relation,fromPath) {
    var finalRelation = [];
    for (var i = 0; i < relation.length; i ++ ) {
        var tempRelation = relation[i];
        
        if(tempRelation.from.indexOf(fromPath) == 0){
            continue
        }

        finalRelation = finalRelation.concat(tempRelation);
    }

    return finalRelation;
}

function getAction(actionId) {
    for (let i = 0; i < pipelineData.length; i ++) {
        stage = pipelineData[i];
        if (stage.actions) {
            for (let j = 0; j < stage.actions.length; j ++ ) {
                let action = stage.actions[j];
                if (action.id == actionId) {
                    return action;
                }
            }
        }
    }

    return "";
}

function getObjValue(obj, path) {
    path = path.substring(1);
    let value;
    let currentProp = path.split(".")[0];

    if (path.split(".").length > 1) {
        if (typeof(obj[currentProp]) == "undefined") {return null;}
        value = getObjValue(obj[currentProp], path.substring(path.indexOf(".")));
    } else {
        if (typeof(obj[currentProp]) == "undefined") {return null;}
        value = obj[currentProp];
    }

    return value;
}
