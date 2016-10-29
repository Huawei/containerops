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

export function initPipeline(fromNodes,toNodes) {
        var result = [];
        
        // visibleFromNode.sort().reverse();
        // visibleToNode.sort().reverse();

        for (var i = 0; i < fromNodes.length; i ++){
            var tempFromNode = [];

            var tempResult = getPipelineMap(fromNodes[i],toNodes);
            var relation = tempResult.resultMap
            var isAllEqual = tempResult. isAllEqual
            if (relation) {
                tempFromNode = tempFromNode.concat(relation);
            }

            result = result.concat(tempFromNode);
        }

        return result;
}


export function addRelation(relation,needDel,fromPath,toPath,visibleFromNode,visibleToNode) {
   
    if (needDel) {
        relation = delRelation(relation,fromPath);
    }


    relation = relation.concat(calcPipelineInfo(fromPath,toPath));


    var _fromChildNode = [],
        _toChildNode = [],
        _fromNodeArray = visibleFromNode.split(";"),
        _toNodeArray = visibleToNode.split(";"),
        _relation = [];

    for(var i=0;i<_fromNodeArray.length;i++){
        if(_fromNodeArray[i].indexOf(fromPath) == 0){
           _fromChildNode.push(_fromNodeArray[i]);
        }
    }

    for(var i=0;i<_toNodeArray.length;i++){
        if(_toNodeArray[i].indexOf(toPath) == 0){
           _toChildNode.push(_toNodeArray[i]);
        }
    }
    
    
    for(var i =1;i<_fromChildNode.length;i++){
        for(var j =1;j<_toChildNode.length;j++){
            
            if(_fromChildNode[i].replace(fromPath,"") == _toChildNode[j].replace(toPath,"")){

               relation = relation.concat(calcPipelineInfo(_fromChildNode[i],_toChildNode[j]));
            }
        }

    }

    return relation;
}

export function delRelation(relation,fromPath) {
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



function calcPipelineInfo(fromPath,toPath) {
    var pipelineInfo = {};
    // for (var i = 0; i < visibleToNode.length; i ++ ) {
        
    //     var regx = new RegExp('^' + visibleToNode[i]);
    //     var rs = regx.exec(toPath);
    //     if (rs) {
            pipelineInfo['to'] = toPath;
    //         pipelineInfo['toShow'] = visibleToNode[i];

    //         if (toPath== visibleToNode[i]) {
    //             pipelineInfo["isToEqual"] = true;
    //         }else {
    //             pipelineInfo["isToEqual"] = false;
    //         }
    //         break;
    //     }
    // }

    // for (var i = 0; i < visibleFromNode.length; i ++ ) {
    //     var regx = new RegExp('^' + visibleFromNode[i]);
    //     var rs = regx.exec(fromPath);
    //     if (rs) {
            pipelineInfo['from'] = fromPath;
    //         pipelineInfo['fromShow'] = visibleFromNode[i];
    //         if (fromPath == visibleFromNode[i]) {
    //             pipelineInfo["isFromEqual"] = true;
    //         }else {
    //             pipelineInfo["isFromEqual"] = false;
    //         }
    //         break;
    //     }
    // }

    return pipelineInfo;
}


function getPipelineMap(fromNode,toNodes) {
    var resultMap = [];
    var isAllEqual = true;

    for (var i = 0; i < toNodes.length; i ++) {
        // 只有类型和名字相等才可以自动匹配上
        if (fromNode.key == toNodes[i].key && fromNode.type == toNodes[i].type){
            // 如果是对象,则匹配其所有子子节点
            if (fromNode.type == "object" && fromNode.childNode && toNodes[i].childNode) {
                var isChildAllEqual = true;

                for (var j = 0; j < fromNode.childNode.length; j ++) {
                    var result = getPipelineMap(fromNode.childNode[j],toNodes[i].childNode);
                    var childResult=result.resultMap;
                    var isChildEqual = result.isAllEqual;
                    if (childResult) {
                        resultMap = resultMap.concat(childResult);
                    }

                    if (!isChildEqual) {
                        isAllEqual = false;
                        isChildAllEqual = false;
                    }
                }

                if (isChildAllEqual) {
                    var pipelineInfo = calcPipelineInfo(fromNode.path,toNodes[i].path);
                    resultMap = resultMap.concat(pipelineInfo);
                }

            } else {
                var pipelineInfo = calcPipelineInfo(fromNode.path,toNodes[i].path);
                resultMap = resultMap.concat(pipelineInfo);
                break;
            }
        } else {
            isAllEqual = false;
        }
    }

    if (resultMap.length > 0) {
        var result = {"resultMap":resultMap,"isAllEqual":isAllEqual};
        return result;
    }else {
        var result = {"resultMap":resultMap,"isAllEqual":isAllEqual};
        return null, isAllEqual;
    }
}

function childNodeRelation(fromPath,toPath,visibleFromNode,visibleToNode){
    var _fromChildNode = [],
        _toChildNode = [],
        _fromNodeArray = visibleFromNode.split(";"),
        _toNodeArray = visibleToNode.split(";"),
        _relation = [];

    for(var i=0;i<_fromNodeArray.length;i++){
        if(_fromNodeArray[i].indexOf(fromPath) == 0){
           _fromChildNode.push(_fromNodeArray[i]);
        }
    }

    for(var i=0;i<_toNodeArray.length;i++){
        if(_toNodeArray[i].indexOf(toPath) == 0){
           _toChildNode.push(_toNodeArray[i]);
        }
    }
    
    
    for(var i =0;i<_fromChildNode.length;i++){
        for(var j =0;j<_toChildNode.length;j++){
            if(_fromChildNode[i] == _toChildNode[j]){
               _relation = _relation.concat(calcPipelineInfo(_fromChildNode[i],_toChildNode[j]));
            }
        }

    }

    return _relation;

}




