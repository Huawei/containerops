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

import {getPathData} from "./setPath";
import {isObject,isArray,isBoolean,isNumber,isString} from "../common/util";
import {addRelation,delRelation,initWorkflow} from "./relation";
import { notify } from "../common/notify";

var fromParentDom,toParentDom,startKey;

export function bipatiteView(importJson,outputJson,linePathData){
	$("#inputH4").find("select").remove();
	var importTree;
	var outputTree;
	if(isArray(importJson)){
		var select = $('<select></select>');
		startKey = importJson[0].event+'_'+importJson[0].type

		importTree = jsonTransformation(importJson[0].json);
		outputTree = jsonTransformation(outputJson);
		
		if(linePathData.relation == undefined){
			linePathData.relation = {};
			for(var i=0;i<importJson.length;i++){
				var nodeTree = jsonTransformation(importJson[i].json);
				var key = importJson[i].event+'_'+importJson[i].type;
				select.append('<option>'+key +'</option>');
				linePathData.relation[key] = getRelationArray(nodeTree,outputTree);
			}
		}else{
			for(var i=0;i<importJson.length;i++){
				select.append('<option>'+importJson[i].event+'_'+importJson[i].type +'</option>');
			}
		}
		$("#inputH4").append(select);

		
	    

	    initView(importTree,outputTree,linePathData);


	    select.change(function(){
	    	startKey = $(this).val();
	    	var selectIndex = ($(this).get(0).selectedIndex);
	    	importTree = jsonTransformation(importJson[selectIndex].json);
	    	initView(importTree,outputTree,linePathData);
	    })

	}else{
		importTree = jsonTransformation(importJson);
	    outputTree = jsonTransformation(outputJson);

	    if(linePathData.relation == undefined){
			linePathData.relation = getRelationArray(importTree,outputTree);
		}
	}

	initView(importTree,outputTree,linePathData);
}




function getRelationArray(importTree,outputTree){

	return initWorkflow(importTree,outputTree);

}


function initView(importTree,outputTree,linePathData){
	
	

	$("#importDiv").html("");
	$("#outputDiv").html("");
	construct($("#importDiv"),importTree,linePathData.relation);
	construct($("#outputDiv"),outputTree,linePathData.relation);

	if(isArray(linePathData.relation)){
		relationLineInit(linePathData.relation);
	}else{
		relationLineInit(linePathData.relation[startKey]);
	}
	
	var mouseX,mouseY;
	
	$("span.property").mousedown(function(event){

		
		if(event.buttons != 1){
			return false;
		}

		var _startX = $(event.target).offset().left,
	        _startY = $(event.target).offset().top,
	        startClass = $(event.target).parent().attr("class"),
	    	fromPath = $(event.target).parent().attr("data-path").replace(/\-/g,'.'),
	    	fromPath = fromPath.substring(5);
	    	mouseX = event.clientX;
	    	mouseY = event.clientY;
		
	    document.onmousemove = function(event){
	    	dragDropLine([_startX,_startY,event.pageX,event.pageY]);
	    }

	    document.onmouseup = function(event){
	    	if(mouseX == event.clientX && mouseY == event.clientY){
	    		return false;
	    	}
	    	document.onmousemove = null;   
        	document.onmouseup = null; 
        	
	    	var endX = $(event.target).offset().left,
	    		endY = $(event.target).offset().top,
	    		endClass = $(event.target).parent().attr("class"),
	    		toPath = $(event.target).parent().attr("data-path");
	    		
	    	
	    	$("#bipatiteLineSvg .drag-drop-line").remove();
	    	
            
	    	if(toPath != undefined){

	    		if(startClass != endClass){
		    		notify("Different types can't build relationships", "error");
		    		return false;
		    	}

	    		toPath = toPath.replace(/\-/g,'.').substring(5);

	    		if(isArray(linePathData.relation)){
					linePathData.relation = addRelation(linePathData.relation,true,fromPath,toPath,getVisibleInputStr(),getVisibleOutputStr());
	    			relationLineInit(linePathData.relation);
				}else{
					linePathData.relation[startKey] = addRelation(linePathData.relation[startKey],true,fromPath,toPath,getVisibleInputStr(),getVisibleOutputStr());
	    			relationLineInit(linePathData.relation[startKey]);
				}


	    		
	    	}
	    }
	}).mouseup(function(event){
		if(mouseX == event.clientX && mouseY == event.clientY){
    		document.onmousemove = null; 
    		document.onmouseup = null; 
    	}
	});

	
    $('#removeLine').off('click');
	$("#removeLine").on('click', function(){
		var path = $("#bipatiteLineSvg path.active");
		var index = path.attr("data-index"); 
		if(isArray(linePathData.relation)){
			linePathData.relation.splice(index,1);
			relationLineInit(linePathData.relation);
		}else{
			linePathData.relation[startKey].splice(index,1);
			relationLineInit(linePathData.relation[startKey]);
		}
		
		$(this).addClass("hide");
	});




	function construct(root,json){

		for(var i=0;i<json.length;i++){

			var item     = $('<div>',   { 'class': 'item row '+json[i].type, 'data-path': replacePoint(json[i].path) }),
				property =   $('<span>', { 'class': 'property' });

			property.text(json[i].key).attr("title",json[i].key);
			item.append(property);
			root.append(item);

			if(json[i].childNode){
				addExpander(item);
				construct(item,json[i].childNode);
			}
		}	
	}



	function addExpander(item){
		if (item.children('.expander').length == 0) {
	        var expander =   $('<span>',  { 'class': 'expander' });
	        expander.bind('click', function() {
	            var item = $(this).parent();
	            item.toggleClass('expanded');
	            if(isArray(linePathData.relation)){
	            	relationLineInit(linePathData.relation);
	            }else{
	            	relationLineInit(linePathData.relation[startKey]); 
	            }
	            
	        });
	        item.prepend(expander);
	    }
	}



}


function relationLineInit(ary){
	d3.select("#bipatiteLineSvg").selectAll("path").remove();
	relationLine(ary);
}

function relationLine(ary){
	
	var rootImport = $("#importDiv"),
		rootOutput = $("#outputDiv");
		
	for(let i=0;i<ary.length;i++){

		let fromPath = replacePoint(ary[i].from);
		let	toPath = replacePoint(ary[i].to);
		let	fromDom = rootImport.find("div[data-path="+fromPath+"]");
		let	toDom = rootOutput.find("div[data-path="+toPath+"]");


		
		if(fromDom.hasClass("expanded") && toDom.hasClass("expanded")){
			continue;
		}


		if(fromDom.is(":visible") && toDom.is(":visible")){
			settingOut([
				fromDom.offset().left,
				fromDom.offset().top,
				toDom.offset().left,
				toDom.offset().top	
			],fromPath,toPath,i);	
		}


		if(fromDom.is(":visible") && toDom.is(":hidden")){
			getVisibleToParent(toDom);
		  	if(toParentDom != undefined){
			  	settingOut([
					fromDom.offset().left,
					fromDom.offset().top,
					toParentDom.offset().left,
					toParentDom.offset().top	
				],fromPath,toPath,i);	
			  }
		}

		if(fromDom.is(":hidden") && toDom.is(":visible")){

		  
		  	getVisibleFromParent(fromDom);

		  	if(fromParentDom != undefined){
		  		settingOut([
					fromParentDom.offset().left,
					fromParentDom.offset().top,
					toDom.offset().left,
					toDom.offset().top	
				],fromPath,toPath,i);	
		  	}
		  	
		}

		if(fromDom.is(":hidden") && toDom.is(":hidden")){

		  	getVisibleFromParent(fromDom);
		  	getVisibleToParent(toDom);

		  	if(fromParentDom != undefined && toParentDom != undefined){
		  		settingOut([
					fromParentDom.offset().left,
					fromParentDom.offset().top,
					toParentDom.offset().left,
					toParentDom.offset().top	
				],fromPath,toPath,i);	
		  	}
		  	
		}

	}
}



function jsonTransformation(json){
	var newJsonArray=[];

	for(var key in json){
		newJsonArray.push({
			"key" : key,
			"type" : jsonType(json[key]),
			"path" : "."+key
		});
		if(isObject(json[key]) || isArray(json[key])){
			var child = newJsonArray[newJsonArray.length-1].childNode = [];
			jsonChange(child,json[key],newJsonArray[newJsonArray.length-1].path);
		}
		
	}
	return newJsonArray;

}

function jsonChange(child,json,path){
	
	if(isObject(json)){
		for(var key in json){
			child.push({
				"key" : key,
				"type" : jsonType(json[key]),
				"path" : path+"."+key
			})
			if(isObject(json[key]) || isArray(json[key])){
				var childNode = child[child.length-1].childNode = [];
				jsonChange(childNode,json[key],child[child.length-1].path);
			}
		}
	}else if(isArray(json)){
		for(var i =0;i<json.length;i++){
			if(isObject(json[i])){
				for(var key in json[i]){
					child.push({
						"key" : key,
						"type" : jsonType(json[i][key]),
						"path" : path+"."+i+"."+key
					})
					if(isObject(json[i][key]) || isArray(json[i][key])){
						var childNode = child[child.length-1].childNode = [];
						jsonChange(childNode,json[i][key],child[child.length-1].path);
					}
				}
			}
		}
	}
}

function settingOut(point,fromPath,toPath,index){
	var offsetTop = $("#bipatiteLineSvg").offset().top;
	var offsetLeft = $("#bipatiteLineSvg").offset().left;
	var x1 = point[0]-offsetLeft+51;
	var y1 = point[1]-offsetTop;
	var x2 = point[2]-offsetLeft+5
	var y2 = point[3]-offsetTop;
	var d = getPathData({x:x1,y:y1},{x:x2,y:y2});

	d3.select("#bipatiteLineSvg")
	.append("path")
	.attr("d",d)
	.attr("stroke", "#75c880")
    .attr("stroke-width", 6)
    .attr("fill","none")
    .attr("stroke-opacity", "0.8")
    .attr("class","cursor")
    .attr("from",fromPath)
    .attr("to",toPath)
    .attr("data-index",index)
    .on("click",function(d,i){
    	$("#removeLine").removeClass("hide");
    	$("#bipatiteLineSvg path").attr("stroke","#75c880").removeClass("active");
    	$(this).attr("class","cursor active").attr("stroke","red");
    });

}



function replacePoint(str){
	str = ("start"+str).replace(/\./g,'-');
	return str;
}

function getVisibleFromParent(dom){
	var parent = $(dom).parent();

	if(parent.is(":hidden")){
		getVisibleFromParent(parent);
	}else{
	  fromParentDom = parent;
	}
	
}

function getVisibleToParent(dom){
	var parent = $(dom).parent();

	if(parent.is(":hidden")){
		getVisibleToParent(parent);
	}else{
	  	toParentDom = parent;
	}
	
}

function jsonType(json){
	if(isObject(json)){
		return "object";
	}else if(isArray(json)){
		return "array";
	}else if(isBoolean(json)){
		return "boolean";
	}else if(isString(json)){
		return "string";
	}else if(isNumber(json)){
		return "number";
	}else {
		return "null";
	}
}



function dragDropLine(point){
    
	var offsetTop = $("#bipatiteLineSvg").offset().top;
	var offsetLeft = $("#bipatiteLineSvg").offset().left;
	var x1 = point[0]-offsetLeft+51;
	var y1 = point[1]-offsetTop;
	var x2 = point[2]-offsetLeft+5
	var y2 = point[3]-offsetTop;
	var d = getPathData({x:x1,y:y1},{x:x2,y:y2});

	if($("#bipatiteLineSvg .drag-drop-line").length == 0){
		d3.select("#bipatiteLineSvg")
		.append("path")
		.attr("d",d)
		.attr("stroke", "red")
	    .attr("stroke-width", 3)
	    .attr("fill","none")
	    .attr("stroke-opacity", "0.8")
	    .attr("class","drag-drop-line");
	}else{
		d3.select(".drag-drop-line")
		.attr("d",d);
	}
    $("#removeLine").addClass("hide");

}

function getVisibleInputStr(){
	var str = "";
	
	$("#importDiv div.item").each(function(){
		var path = $(this).attr("data-path").replace(/\-/g,'.');
		str = str + path.substring(5)+";";
	})

	return str;
}

function getVisibleOutputStr(){
	var str = "";
	
	$("#outputDiv div.item").each(function(){
		var path = $(this).attr("data-path").replace(/\-/g,'.');
		str = str + path.substring(5)+";";
	})

	return str;
}



