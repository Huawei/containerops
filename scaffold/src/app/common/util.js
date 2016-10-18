import * as constant from "./constant";

export function isObject(o) { return Object.prototype.toString.call(o) == '[object Object]'; }
export function isArray(o) { return Object.prototype.toString.call(o) == '[object Array]'; }
export function isBoolean(o) { return Object.prototype.toString.call(o) == '[object Boolean]'; }
export function isNumber(o) { return Object.prototype.toString.call(o) == '[object Number]'; }
export function isString(o) { return Object.prototype.toString.call(o) == '[object String]'; }

export function findAllRelatedLines(itemId) {
   var relatedLines = _.filter(constant.linePathAry,function(item){return (item.startData != undefined && item.endData != undefined) && (item.startData.id == itemId || item.endData.id == itemId)});
	   return relatedLines;
}
export function findInputLines(itemId){
   var relatedLines = _.filter(constant.linePathAry,function(item){return (item.endData != undefined) && (item.endData.id == itemId)});
}
export function findOutputLines(itemId){
   var relatedLines = _.filter(constant.linePathAry,function(item){return (item.startData != undefined) && (item.startData.id == itemId)});
}
export function removeRelatedLines(args){
   if(isString(args)) {
   	  var relatedLines = findAllRelatedLines(args);
   	  constant.setLinePathAry( _.difference(constant.linePathAry, relatedLines));
   } else {
      _.each(args, function(item) {
      	 removeRelatedLines(item.id);
      })
   }
   
}
export function findAllActionsOfStage(stageId){
   var groupId = "#action" + "-" + stageId;
   var selector = groupId + "> image";
	   return $(selector);
}
export function disappearAnimation(args){
   if(isString(args)) {
   	  d3.selectAll(args)
	     .transition()
	     .duration(200)
	     .style("opacity",0); 
   } else {
      _.each(args, function(selector){
      	 disappearAnimation(selector);
      })
   }
   
}
export function transformAnimation(args,type){
	 _.each(args, function(item){
	 	  d3.selectAll(item.selector)
	        .filter(function(d,i){ return i > item.itemIndex})
	        .transition()
	        .delay(200)
	        .duration(200)
	        .attr("transform", function(d,i){
	        	var translateX=0,translateY=0;
	        	if(type == "action"){
	        		  translateX = item.type == "siblings" ? d.translateX : 0;
	           		  translateY = item.type == "siblings" ? (d.translateY - constant.ActionNodeSpaceSize) : (0 - constant.ActionNodeSpaceSize);
	            	 
	        	}else if(type == "stage"){
                      translateX = item.type == "siblings" ? (d.translateX - constant.PipelineNodeSpaceSize) : (0 - constant.PipelineNodeSpaceSize);
	           		  translateY = item.type == "siblings" ? d.translateY : 0 ;
	           		 
	        	}
	        	return "translate(" + translateX + "," + translateY + ")";
	            
	        });
	 })
	
}
