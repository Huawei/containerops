let apiUrlConf = {
	"host" : "https://test-1.containerops.sh",
	"pipeline" : {
		"list" : "/pipeline/v1/demo/demo",
		"data" : "/pipeline/v1/demo/demo/{pipelineName}/json?id={pipelineID}",
		"add" : "/pipeline/v1/demo/demo",
		"newVersion" : "",
		"save" : "",
		"eventOutput" : ""
	},
	"component" : {
		"list" : "",
		"data" : "",
		"add" : "",
		"newVersion" : "",
		"save" : ""
	}
}

// pipeline
export let pipelineApi = {
	"list" : function(){
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.list,
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    }); 
	    return promise;
	},
	"data" : function(name,id){
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.data.replace(/{pipelineName}/g, name).replace(/{pipelineID}/g, id),
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    }); 
	    return promise;
	},
	"add" : function(name,version){
		var data = JSON.stringify({
				"name":name,
				"version":version
			});
		var promise = $.ajax({
	        "url": apiUrlConf.host + apiUrlConf.pipeline.add,
	        "type": "POST",
	        "dataType": "json",
	        "data": data 
	    }); 
	    return promise;
	},
	"newVersion" : function(data){
		var promise = $.ajax({
	        "url": apiUrlConf.pipeline.newVersion,
	        "type": "POST",
	        "dataType": "json",
	        "data": data 
	    }); 
	    return promise;
	},
	"save" : function(data){
		var promise = $.ajax({
	        "url": apiUrlConf.pipeline.save,
	        "type": "PUT",
	        "dataType": "json",
	        "data": data 
	    }); 
	    return promise;
	},
	"eventOutput" : function(){
		var promise = $.ajax({
	        "url": apiUrlConf.pipeline.eventOutput,
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    }); 
	    return promise;
	}
}

// component
export let componentApi = {
	"list" : function(){
		var promise = $.ajax({
	        "url": apiUrlConf.component.list,
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    }); 
	    return promise;
	},
	"data" : function(){
		var promise = $.ajax({
	        "url": apiUrlConf.component.data,
	        "type": "GET",
	        "dataType": "json",
	        "cache": false
	    }); 
	    return promise;
	},
	"add" : function(data){
		var promise = $.ajax({
	        "url": apiUrlConf.component.add,
	        "type": "POST",
	        "dataType": "json",
	        "data": data 
	    }); 
	    return promise;
	},
	"newVersion" : function(data){
		var promise = $.ajax({
	        "url": apiUrlConf.component.newVersion,
	        "type": "POST",
	        "dataType": "json",
	        "data": data 
	    }); 
	    return promise;
	},
	"save" : function(data){
		var promise = $.ajax({
	        "url": apiUrlConf.component.save,
	        "type": "PUT",
	        "dataType": "json",
	        "data": data 
	    }); 
	    return promise;
	}
}