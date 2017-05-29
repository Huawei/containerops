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

import {notify} from "../common/notify";
import {ajaxCall} from "../common/api";
import {isEnvsLegal} from "../common/check";

let allComponents = [];

export function getAllComponents(){
     return ajaxCall("component.list");
}

export function getComponent(name,versionid){
    var params = {
        "componentName" : name,
        "componentID" : versionid
    }

    return ajaxCall("component.detail",params);
}

export function addComponent(){
    if(!$('#newcomponent-form').parsley().validate()){
        return false;
    }

    var reqbody = {
        "name" : $("#c-name").val(),
        "version" : $("#c-version").val()
    }

    return ajaxCall("component.add",{},reqbody);
}

export function addComponentVersion(name, versionid, componentData){
    if(!$('#newcomponent-version-form').parsley().validate()){
        return false;
    }else{
        var version = $("#c-newversion").val();
        return saveComponent(name, version, versionid, componentData)
    }
}

export function saveComponent(name, version, versionid, componentData){
    var params = {
        "componentName" : name
    }

    var reqbody = {
        "id" : versionid,
        "version" : version.toString(),
        "define": componentData
    }

    return ajaxCall("component.save",params,reqbody);
}

export function validateComponent(componentData){
    if(!$('#component-form').parsley().validate()){
        notify("Missed some required base config.","error");
        return false;
    }else if(!componentData.setupData.action.useAdvanced && !$('#base-setting-form').parsley().validate()){
        notify("Missed some required base setting of kubernetes.","error");
        return false;
    }else if(_.isEmpty(componentData.inputJson)){
        notify("Component input json is empty.","error");
        return false;
    }else if(_.isEmpty(componentData.outputJson)){
        notify("Component output json is empty.","error");
        return false;
    }else if(!_.isEmpty(componentData.env) && !isEnvsLegal(componentData.env)){
        notify("Component env key is not allowed to start with 'CO_'.","error");
        return false;
    }else{
        return true;
    }
}

export var newComponentData = {
    "setupData" : {},
    "inputJson" : {},
    "outputJson" : {},
    "env" :{}
}