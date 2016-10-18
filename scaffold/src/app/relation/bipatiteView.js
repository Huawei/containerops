import {getPathData} from "./setPath";
import {isObject,isArray,isBoolean,isNumber,isString} from "../common/util";
import {addRelation,delRelation,initPipeline} from "./relation";


// var relationArray;

var importTreeJson,outputTreeJson;

export function bipatiteView(importJson,outputJson,linePathData){
	
    var importTree = importTreeJson = jsonTransformation(importJson);
    var outputTree = outputTreeJson = jsonTransformation(outputJson);
    initView(importTree,outputTree,linePathData);
       
}

function getRelationArray(){

	var visibleInputStr = getVisibleInputStr();
	var visibleOutputStr = getVisibleOutputStr();
	var visibleInput = visibleInputStr.split(";");
    var visibleOutput = visibleOutputStr.split(";");

	return initPipeline(importTreeJson,outputTreeJson,visibleInput,visibleOutput);
}


function initView(importTree,outputTree,linePathData){

	
	if(linePathData.relation == undefined){
		linePathData.relation = getRelationArray();
	}

	construct($("#importDiv"),importTree,linePathData.relation);
	construct($("#outputDiv"),outputTree,linePathData.relation);

	relationLineInit(linePathData.relation);

	

	$("span.property").mousedown(function(event){

		var _startX = $(event.target).offset().left,
	        _startY = $(event.target).offset().top,
	        startClass = $(event.target).parent().attr("class"),
	    	fromPath = $(event.target).parent().attr("data-path").replace(/\-/g,'.');
	    	fromPath = fromPath.substring(5);
		
	    document.onmousemove = function(event){
	    	event.pageX
	    	event.pageY
	    	dragDropLine([_startX,_startY,event.pageX,event.pageY]);
	    }

	    document.onmouseup = function(event){
	    	document.onmousemove = null;   
        	document.onmouseup = null; 
        	
	    	var endX = $(event.target).offset().left,
	    		endY = $(event.target).offset().top,
	    		endClass = $(event.target).parent().attr("class"),
	    		toPath = $(event.target).parent().attr("data-path");
	    		
	    	
	    	$("#bipatiteLineSvg .drag-drop-line").remove();
	    	
            
	    	if(toPath != undefined){

	    		if(startClass != endClass){
		    		alert("difference type");
		    		return false;
		    	}

	    		toPath = toPath.replace(/\-/g,'.').substring(5);

	    		linePathData.relation = addRelation(linePathData.relation,true,fromPath,toPath,getVisibleInputStr(),getVisibleOutputStr());
	    	
	    		relationLineInit(linePathData.relation);
	    	}
	    	
	    }
	})




	$("#removeLine").click(function(){
		var path = $("#bipatiteLineSvg path.active");
		
		var fromPath = path.attr("from"); 
		fromPath = fromPath.replace(/\-/g,'.').substring(5);
		
		linePathData.relation = delRelation(linePathData.relation,fromPath);
		relationLineInit(linePathData.relation);
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
	            relationLineInit(linePathData.relation);
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
		
	for(var i=0;i<ary.length;i++){
		var fromPath = replacePoint(ary[i].from),
			toPath = replacePoint(ary[i].to),
			fromDom = rootImport.find("div[data-path="+fromPath+"]"),
			toDom = rootOutput.find("div[data-path="+toPath+"]"),
			parentDom;

		
		if(fromDom.is(":visible") && toDom.is(":visible")){
			settingOut([
				fromDom.offset().left,
				fromDom.offset().top,
				toDom.offset().left,
				toDom.offset().top	
			],fromPath,toPath);	
		}


		if(fromDom.is(":visible") && toDom.is(":hidden")){
			
		  	parentDom = getVisibleParent(toDom);
		  	if(parentDom != undefined){
			  	settingOut([
					fromDom.offset().left,
					fromDom.offset().top,
					parentDom.offset().left,
					parentDom.offset().top	
				],fromPath,toPath);	
			  }
		}

		if(fromDom.is(":hidden") && toDom.is(":visible")){

		  	parentDom = getVisibleParent(fromDom);
		  	if(parentDom != undefined){
		  		settingOut([
					parentDom.offset().left,
					parentDom.offset().top,
					toDom.offset().left,
					toDom.offset().top	
				],fromPath,toPath);	
		  	}
		  	
		}



		if(ary[i].child){
			relationLine(ary[i].child);
		}
	

	}
}



function jsonTransformation(json){
	var newJsonArray=[];

	for(var key in json){
		newJsonArray.push({
			"key" : key,
			"type" : JsonType(json[key]),
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
				"type" : JsonType(json[key]),
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
						"type" : JsonType(json[i][key]),
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

function settingOut(point,fromPath,toPath){
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
	.attr("stroke", "green")
    .attr("stroke-width", 3)
    .attr("fill","none")
    .attr("stroke-opacity", "0.8")
    .attr("class","cursor")
    .attr("from",fromPath)
    .attr("to",toPath)
    .on("click",function(d,i){
    	$("#removeLine").removeClass("hide");
    	$("#bipatiteLineSvg path").attr("stroke","green").removeClass("active");
    	$(this).attr("class","cursor active").attr("stroke","red");
    });

}



function replacePoint(str){
	str = ("start"+str).replace(/\./g,'-');
	return str;
}

function getVisibleParent(dom){

	var parent = dom.parent();

	if(parent.is(":visible")){
		return parent;
	}else{
		getVisibleParent(parent);
	}
}

function JsonType(json){
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

}

function getVisibleInputStr(){
	var str = "";
	
	$("#importDiv div.item").each(function(){
		
		// if($(this).is(":visible")){
			var path = $(this).attr("data-path").replace(/\-/g,'.');
			str = str + path.substring(5)+";";
		// }
	})

	return str;
}

function getVisibleOutputStr(){
	var str = "";
	
	$("#outputDiv div.item").each(function(){
		
		// if($(this).is(":visible")){
			var path = $(this).attr("data-path").replace(/\-/g,'.');
			str = str + path.substring(5)+";";
		// }
	})

	return str;
}



