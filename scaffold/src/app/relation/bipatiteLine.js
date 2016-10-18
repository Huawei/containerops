
import {getPathData} from "./setPath";


export function bipatiteLine(array){
	
	d3.select("#bipatiteLineSvg").selectAll("path").remove();
	for(var i =0;i<array.length;i++){
		if(array[i].parent){
			var $input = $("div[data-path="+array[i].parent+"]").find(">input[title="+array[i].name+"]");

			if($($input[0]).is(":visible") && $($input[1]).is(":visible")){
				if($input.length == 2){
					var point = [];
					$input.each(function(){
						point.push($(this).offset().left);
						point.push($(this).offset().top);
					});
					settingOut(point);
				}
				
			}

			
		}
		

	}
}



function settingOut(point){
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
    .attr("stroke-width", 1)
    .attr("fill","none")
    .attr("stroke-opacity", "0.2");



}