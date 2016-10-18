import {linePathAry} from "../common/constant";

export function removeLinkArray(sd){
	var parentIndex = $("#"+sd.id).attr("data-parent");
	var index = $("#"+sd.id).attr("data-index");


	for(var i =0;i<linePathAry.length;i++){
		
		if(parentIndex == undefined){
			if(linePathAry[i].fromParentIndex == index || linePathAry[i].toParentIndex == index){
				linePathAry.splice(i,1);
				removeLinkArray(sd);
				return;
			}
		}else{
			if(linePathAry[i].fromParentIndex == parentIndex){
				if(linePathAry[i].fromIndex == index){
					linePathAry.splice(i,1);
					removeLinkArray(sd);
					return;
				}
			}

			if(linePathAry[i].toParentIndex == parentIndex){
				if(linePathAry[i].toIndex == index){
					linePathAry.splice(i,1);
					removeLinkArray(sd);
					return;
				}
			}
		}

		
	}

	
}