import {jsonEditor} from "../../vendor/jquery.jsoneditor";
import {bipatiteJson} from "./bipatiteJson";
import {bipatiteLine} from "./bipatiteLine";
import {bipatiteView} from "./bipatiteView";
import {resizeWidget} from "../theme/widget";
import {pipelineData} from "../pipeline/main";
import * as constant from "../common/constant";

export var importJson = {};
   
export var outputJson = {};

// export var bipatiteViewOn = false;

export function pipelineEdit(data,linkDom){
    
    var fromParent = linkDom.attr("from-parent"),
        fromIndex = linkDom.attr("from-index"),
        toParent = linkDom.attr("to-parent"),
        toIndex = linkDom.attr("to-index"),
        index = linkDom.attr("data-index");

    $("#pipeline-info-edit").html($(data));


    $("#removeLink").click(function(){
        linkDom.remove();
        constant.linePathAry.splice(index, 1);
        $("#pipeline-info-edit").html("");
    })

    if(fromParent != -1){
        if(pipelineData[fromParent].actions[fromIndex].inputJson!= undefined){
            $("#importDiv").html("");
            importJson = pipelineData[fromParent].actions[fromIndex].outputJson;
        }else{
            $("#importDiv").html("no data");
            importJson = {};
        }
    }else{
        $("#importDiv").html("no data");
        importJson = {};
    }

    if(pipelineData[toParent].actions[toIndex].outputJson != undefined){
        $("#outputDiv").html("");
        outputJson = pipelineData[toParent].actions[toIndex].inputJson;
    }else{
        $("#outputDiv").html("no data");
        outputJson = {};
    }

    bipatiteView(importJson,outputJson,constant.linePathAry[index]);
    // console.log("line")
    // console.log(constant.linePathAry[index]);
    resizeWidget();
}