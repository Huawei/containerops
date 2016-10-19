/* Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

import {notify} from "../common/notify";

let allComponents = [];

export function getAllComponents(){
    // call api, return ajax promise

    // to be removed below
    return allComponents;
}

export function getComponent(name,version){
    // call api, return ajax promise

    // to be removed below
    var component = _.find(allComponents,function(item){
        return item.name == name;
    });
    return _.find(component.versions,function(item){
        return item.version = version;
    }).data;
}

export function addComponent(){
    if(!$('#newcomponent-form').parsley().validate()){
        return false;
    }
    var name = $("#c-name").val();
    var version = $("#c-version").val();

    // call api here, return promise

    // to be removed
    var component = _.find(allComponents,function(item){
        return item.name == name;
    })
    if(!_.isUndefined(component)){
        var newversion = {
            "version" : version,
            "data" : [].concat(newComponentData)
        }
        component.versions.push(newversion);
    }else{
        component = {
            "name" : name,
            "versions" : [
                {
                    "version" : version,
                    "data" : [].concat(newComponentData)
                }
            ]
        }
        allComponents.push(component);
    }
    return true;
}

export function addComponentVersion(oldversion){
    if(!$('#newcomponent-version-form').parsley().validate()){
        return false;
    }
    var name = $("#c-name-newversion").val();
    var version = $("#c-newversion").val();

    // call api here, return promise

    // to be removed
    var component = _.find(allComponents,function(item){
        return item.name == name;
    });

    var oldversion = _.find(component.versions,function(item){
        return item.version == oldversion;
    });

    var newversion = {
        "version" : version,
        "data" : [].concat(oldversion.data)
    }
    component.versions.push(newversion);
    return true;
}

export function saveComponent(componentName, componentVersion, componentData){
    if(!$('#component-form').parsley().validate()){
        notify("missed some required base config.","error");
    }else if(_.isEmpty(componentData.inputJson)){
        notify("component input json is empty.","error");
    }else if(_.isEmpty(componentData.outputJson)){
        notify("component output json is empty.","error");
    }else{
        notify("component saved.","success");
    }  
}

export var newComponentData = {
    "setupData" : {},
    "inputJson" : {},
    "outputJson" : {}
}