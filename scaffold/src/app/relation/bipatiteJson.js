export var bipatiteArray =[];
	

export function bipatiteJson(json1,json2,parent){
	if(typeof json1 != typeof json2){
		return false;
	}
	var ary1 = [];
	var ary2 = [];

	for(var key in json1){
		ary1.push(key);
	}

	for(var key in json2){
		ary2.push(key);
	}
		
	var comAry = arrayIntersection(ary1,ary2);

	for(var i=0;i<comAry.length;i++){
		bipatiteArray.push({
			parent : parent ,
			name : comAry[i]
		});
		
		if(typeof json1[comAry[i]] == typeof json2[comAry[i]]){
			if(isObject(json1[comAry[i]]) || isArray(json1[comAry[i]]) ){
				var thisParent = !parent ? comAry[i] : parent+"-"+comAry[i];

				bipatiteJson(json1[comAry[i]],json2[comAry[i]],thisParent);
			}
			
		}
	}

	return bipatiteArray;

}

function arrayIntersection ( a, b )
{
    var ai=0, bi=0;
    var result = new Array();
    while ( ai < a.length && bi < b.length )
    {
        if      ( a[ai] < b[bi] ) { ai++; }
        else if ( a[ai] > b[bi] ) { bi++; }
        else /* they're equal */
        {
            result.push ( a[ai] );
            ai++;
            bi++;
        }
    }
    return result;
}

function isObject(o) { return Object.prototype.toString.call(o) == '[object Object]'; }
function isArray(o) { return Object.prototype.toString.call(o) == '[object Array]'; }