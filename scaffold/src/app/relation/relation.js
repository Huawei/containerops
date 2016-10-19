

export function initPipeline(fromNodes,toNodes,visibleFromNode,visibleToNode) {
        var result = [];
        
        visibleFromNode.sort().reverse();
        visibleToNode.sort().reverse();

        for (var i = 0; i < fromNodes.length; i ++){
            var tempFromNode = [];

            var relation = getPipelineMap(fromNodes[i],toNodes,visibleFromNode,visibleToNode);
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


    relation = relation.concat(calcPipelineInfo(fromPath,toPath,visibleFromNode,visibleToNode));


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

               relation = relation.concat(calcPipelineInfo(_fromChildNode[i],_toChildNode[j],visibleFromNode,visibleToNode));
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

        // 如果当前路径不是要删除的,则保留下来当前路径
        finalRelation = finalRelation.concat(tempRelation);
    }

    return finalRelation;
}



function calcPipelineInfo(fromPath,toPath,visibleFromNode,visibleToNode) {
    var pipelineInfo = {};
    for (var i = 0; i < visibleToNode.length; i ++ ) {
        // 通过正则匹配当前路径可以匹配上的,
        // 因为已经进行过排序+倒序了,所以第一个匹配上的肯定是距离当前节点最近的可见节点
        var regx = new RegExp('^' + visibleToNode[i]);
        var rs = regx.exec(toPath);
        if (rs) {
            pipelineInfo['to'] = toPath;
            pipelineInfo['toShow'] = visibleToNode[i];

            if (toPath== visibleToNode[i]) {
                pipelineInfo["isToEqual"] = true;
            }else {
                pipelineInfo["isToEqual"] = false;
            }
            break;
        }
    }

    for (var i = 0; i < visibleFromNode.length; i ++ ) {
        var regx = new RegExp('^' + visibleFromNode[i]);
        var rs = regx.exec(fromPath);
        if (rs) {
            pipelineInfo['from'] = fromPath;
            pipelineInfo['fromShow'] = visibleFromNode[i];
            if (fromPath == visibleFromNode[i]) {
                pipelineInfo["isFromEqual"] = true;
            }else {
                pipelineInfo["isFromEqual"] = false;
            }
            break;
        }
    }

    return pipelineInfo;
}


function getPipelineMap(fromNode,toNodes,visibleFromNode,visibleToNode) {
    var resultMap = [];

    for (var i = 0; i < toNodes.length; i ++) {
        // 只有类型和名字相等才可以自动匹配上
        if (fromNode.key == toNodes[i].key && fromNode.type == toNodes[i].type){
            // 如果是对象,则匹配其所有子子节点
            if (fromNode.type == "object" && fromNode.childNode && toNodes[i].childNode) {
                var pipelineInfo = calcPipelineInfo(fromNode.path,toNodes[i].path,visibleFromNode,visibleToNode);
                resultMap = resultMap.concat(pipelineInfo);

                for (var j = 0; j < fromNode.childNode.length; j ++) {
                    var childResult = getPipelineMap(fromNode.childNode[j],toNodes[i].childNode,visibleFromNode,visibleToNode);
                    resultMap = resultMap.concat(childResult);
                }

            } else {
                var pipelineInfo = calcPipelineInfo(fromNode.path,toNodes[i].path,visibleFromNode,visibleToNode);
                resultMap = resultMap.concat(pipelineInfo);
                break;
            }
        }
    }

    if (resultMap.length > 0) {
        return resultMap;
    }else {
        return null;
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
               _relation = _relation.concat(calcPipelineInfo(_fromChildNode[i],_toChildNode[j],visibleFromNode,visibleToNode));
            }
        }

    }

    return _relation;

}




