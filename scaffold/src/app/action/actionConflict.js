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


import {isObject,isArray,isBoolean,isNumber,isString,judgeType} from '../common/util';

export function svgTree(container,data){
	var depthY = 0;
	var jsonArray = [];

	var svg = container.append("svg")
        .attr("width", "100%")
        .attr("height", 600)
        .style("fill", "white");

	for(var i=0;i<data.node.length;i++){
    	transformJson(data.node[i]);
    }

	for(var i=0;i<jsonArray.length;i++){
		construct(svg,jsonArray[i]);
	}


	function transformJson(data){
	
		var depthX = 1;
		depthY++;
		jsonArray.push({
			depthX:depthX,
			depthY:depthY,
			type:"object",
			name : data.name
		});

		
		for(var i=0;i<data.conflicts.length;i++){
			
			var conflicts = data.conflicts[i];

			for(var key in conflicts){
				depthY++;

				jsonArray.push({
					depthX:2,
					depthY:depthY,
					type:judgeType(conflicts[key]),
					name : key
				});

				getChildJson(conflicts[key],3);
			}
		}
	}


	function getChildJson(data,depthX){

		if(isObject(data)){
			for(var key in data){
				depthY++;
				jsonArray.push({
					depthX:depthX,
					depthY:depthY,
					type:judgeType(data[key]),
					name : key
				});
				getChildJson(data[key],depthX+1);
			}
		}

		if(isArray(data) && data.length>0){
			
		}
	}


}




function construct(svg,options){

	var g = svg.append("g")
		.attr("transform","translate("+(options.depthX*20+100)+","+(options.depthY*28)+")");

	var rect = g.append('rect')
		.attr("ry",4)
		.attr("rx",4)
		.attr("y",0)
		.attr("width",135)
		.attr("height",24)
		.attr("fill",function(){
			switch(options.type)
				{
				case "string":
				  	return "#13b5b1";
				  	break;
				case "object":
				  	return "#eb6876";
				  	break;
				case "number":
				  	return "#32b16c";
				  	break;
				case "array":
				  	return "#c490c0";
				  	break;
				case "boolean":
				  	return "#8fc320";
				  	break;
				default:
					return "#cfcfcf"; 
			}
		});

	var clashImage = g.append('image')
		.attr("transform","translate(0,0)")
		.attr("xlink:href","../../assets/svg/conflict.svg")
		.attr("x",2)
		.attr("y",2)
		.attr("width",20)
		.attr("height",20);



	var typeImage = g.append('image')
		.attr("transform","translate(115,0)")
		.attr("xlink:href",function(){
			switch(options.type)
				{
				case "string":
				  	return "../../assets/images/string.png";
				  	break;
				case "object":
				  	return "../../assets/images/object.png";
				  	break;
				case "number":
				  	return "../../assets/images/number.png";
				  	break;
				case "array":
				  	return "../../assets/images/array.png";
				  	break;
				case "boolean":
				  	return "../../assets/images/boolean.png";
				  	break;
				default:
					return ""; 
			}
		})
		.attr("x","0")
		.attr("y","0")
		.attr("width","20")
		.attr("height","24")

	var text = g.append('text')
		.attr("dx",28)
		.attr("dy",17)
		.attr("fill",function(){
			if(options.type == "null"){
				return "#8e8a89";
			}else{
				return "#fff";
			}
		})
		.text(function(){
			if(options.name.length > 12){
				return options.name.substring(0,10) +"...";
			}else{
				return options.name;
			}
		});
}


