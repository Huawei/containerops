

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
    var finalRelation = [];
    if (needDel) {
        relation = delRelation(relation,fromPath,toPath);
    }

    if (fromPath.split(".").length == 2) {
        // 如果是根节点,则直接添加一个即可,

        finalRelation = relation.concat(calcPipelineInfo(fromPath,toPath,visibleFromNode,visibleToNode));
        
        // getVisibleNode(finalRelation,fromPath,toPath,visibleFromNode,visibleToNode)

        
    }else{
        for (var i = 0; i < relation.length; i ++ ) {
            var tempRelation = relation[i];
            if (tempRelation.from + "." + fromPath.split(".")[fromPath.split(".").length - 1] == fromPath) {
                // 如果当前点 + 目标点的最后一段 == 目标点 则代表找到了目标点的直接父节点,直接在当前点的child中加入指定的关系即可
                if (!tempRelation.child) {
                    tempRelation.child = [];
                }

                tempRelation.child = tempRelation.child.concat(calcPipelineInfo(fromPath,toPath,visibleFromNode,visibleToNode));
            } else if (tempRelation.child && fromPath.indexOf(tempRelation.from + ".") == 0) {
                // 如果当前点存在子节点,并且是目标点的父节点,则进入寻找
                tempRelation.child = addRelation(tempRelation.child,false,fromPath,toPath,visibleFromNode,visibleToNode)
            }

            finalRelation = finalRelation.concat(tempRelation)
        }
    }
    return finalRelation;
}

export function delRelation(relation,fromPath) {
    var finalRelation = [];
    for (var i = 0; i < relation.length; i ++ ) {
        var tempRelation = relation[i];

        if (tempRelation.from == fromPath) {
            // 如果开始路径相同,则直接删除即可
            continue;
        } else if (tempRelation.child && fromPath.indexOf(tempRelation.from+".") == 0) {
            // 如果当前起始路径含有子节点,则在子节点里找
            // 这里添加一个判断,只有当前节点是fromPath的父级时才进入判断,这样可以省略很多无用的子节点判断
            tempRelation.child = delRelation(tempRelation.child,fromPath);
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
                pipelineInfo['child'] = [];
                for (var j = 0; j < fromNode.childNode.length; j ++) {
                    var childResult = getPipelineMap(fromNode.childNode[j],toNodes[i].childNode,visibleFromNode,visibleToNode);
                    pipelineInfo.child = pipelineInfo.child.concat(childResult);
                }

                resultMap = resultMap.concat(pipelineInfo);
            } else {
                var pipelineInfo = calcPipelineInfo(fromNode.path,toNodes[i].path,visibleFromNode,visibleToNode);
                resultMap = resultMap.concat(pipelineInfo);
                break;
            }
        }

        // // 如果toNodes存在子节点,则寻找子节点有没有可以匹配当前节点的
        // if (toNodes[i].childNode) {
        //     var tempResult = getPipelineMap(fromNode,toNodes[i].childNode,visibleFromNode,visibleToNode)
        //     if (tempResult) {
        //         resultMap = tempResult;
        //         break;
        //     }
        // }
    }


    if (resultMap.length > 0) {
        // console.log(resultMap);
        return resultMap;
    }else {
        return null;
    }
}

function getVisibleNode(relation,fromPath,toPath,visibleFromNode,visibleToNode){
    var _fromChildNode = [],
        _toChildNode = [],
        _fromNodeArray = visibleFromNode.split(";"),
        _toNodeArray = visibleToNode.split(";");
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
                
            }
        }

    }



}




