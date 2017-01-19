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

function componentCheck(notifyService){

	var componentData;
	
	// runtime config
	function name(){
		if($('#component-name-form').length>0){
			return $('#component-name-form').parsley().validate();
		}else{
			return true;
		}	
	}

	function version(){
		if($('#component-version-form').length>0){
			return $('#component-version-form').parsley().validate();
		}else{
			return true;
		}
	}

	function existingimage(){
		return $('#component-form').parsley().validate();
	}

	function basesetting(){
		return $('#base-setting-form').parsley().validate();
	}

	function advancedsetting(){
		return $('#advanced-setting-form').parsley().validate();
	}

	function env(){
		return $('#component-envs').parsley().validate();
	}
	
	var image_setting = {
		"build" : function(){
			return $('#build-image-form').parsley().validate();
		},
		"from" : function(){
			return $('#base-image-form').parsley().validate();
		},
		"events" : {
			"component_start" : function(){
				return _.isEmpty(componentData.image_setting.events.component_start);
			},
			"component_result" : function(){
				return _.isEmpty(componentData.image_setting.events.component_result);
			},
			"component_stop" : function(){
				return _.isEmpty(componentData.image_setting.events.component_stop);
			}
		},
		"push" : function(){
			return $('#push-image-form').parsley().validate();
		}
	}

	function init(component){
		componentData = component;
	}

	function go(isnewimage){
		var check_result = true;

		if(!name()){
			check_result = false;
			notifyService.notify("Component name is required.","error");
			return check_result;
		}

		if(!version()){
			check_result = false;
			notifyService.notify("Component version is required.","error");
			return check_result;
		}

		if(!isnewimage){
			if(!existingimage()){
				check_result = false;
				notifyService.notify("Repository name of existing image is required.","error");
				return check_result;
			}
		}else{
			if(image_setting.events.component_start()){
				check_result = false;
				notifyService.notify("Script of component start event is required.","error");
				return check_result;
			}
			if(image_setting.events.component_result()){
				check_result = false;
				notifyService.notify("Script of component result event is required.","error");
				return check_result;
			}
			if(image_setting.events.component_stop()){
				check_result = false;
				notifyService.notify("Script of component stop event is required.","error");
				return check_result;
			}
			if(!image_setting.from()){
				check_result = false;
				notifyService.notify("Base image is not complete.","error");
				return check_result;
			}
			if(!image_setting.build()){
				check_result = false;
				notifyService.notify("Build image is not complete.","error");
				return check_result;
			}
			if(!image_setting.push()){
				check_result = false;
				notifyService.notify("Push image is not complete.","error");
				return check_result;
			}
		}

		if(!componentData.use_advanced){
			if(!basesetting()){
				check_result = false;
				notifyService.notify("Kubernetes base setting is not complete.","error");
				return check_result;
			}
		}else{
			if(!advancedsetting()){
				check_result = false;
				notifyService.notify("Kubernetes advanced setting is not complete.","error");
				return check_result;
			}
		}

		if(componentData.env.length>0){
			if(!env()){
				check_result = false;
				notifyService.notify("Component env is not complete.","error");
				return check_result;
			}
		}

		return check_result;
	}

	var tabcheck = {
		"runtime" : function(isnewimage){
			var result = true;

			if(!isnewimage && !existingimage()){
				result = false;
			}

			if(!componentData.use_advanced){
				if(!basesetting()){
					result = false;
				}
			}else{
				if(!advancedsetting()){
					result = false;
				}
			}

			if(componentData.env.length>0){
				if(!env()){
					result = false;
				}
			} 

			return result;
		},
		"editshell" : function(){
			var result = true;
			if(image_setting.events.component_start()){
				result = false;
			}else if(image_setting.events.component_result()){
				result = false;
			}else if(image_setting.events.component_stop()){
				result = false;
			}
			return result;
		},
		"buildimage" : function(){
			var result = image_setting.from() && image_setting.build() && image_setting.push();
			return result;
		}
	}
	
	return {
		"version" : version,
		"init" : init,
		"go" : go,
		"tabcheck" : tabcheck
	}
}
   
devops.factory('componentCheck', ['notifyService', componentCheck]);
