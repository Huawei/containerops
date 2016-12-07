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

import {workflowApi} from "../common/api";

let allWorkflows = [];

export function getAllWorkflows(){
    return workflowApi.list();
}

export function getWorkflow(name,versionid){
    return workflowApi.data(name,versionid);
}

export function addWorkflow(){
    if(!$('#newpp-form').parsley().validate()){
        return false;
    }
    var name = $("#pp-name").val();
    var version = $("#pp-version").val();

    return workflowApi.add(name,version);
}

export function saveWorkflow(name,version,versionid,nodes,lines){
    var reqbody = {
        "id" : versionid,
        "version" : version.toString(),
        "define":{
            "lineList" : lines,
            "stageList" : nodes
        }
    }

    return workflowApi.save(name,reqbody);
}

export function addWorkflowVersion(name, versionid, nodes,lines){
    if(!$('#newpp-version-form').parsley().validate()){
        return false;
    }else{
        var version = $("#pp-newversion").val();
        return saveWorkflow(name, version, versionid, nodes,lines)
    }
}

export function getEnvs(name,versionid){
    return workflowApi.getEnv(name,versionid);
}

export function setEnvs(name,versionid,envs){
    var reqbody = {
        "id" : versionid,
        "env" : _.object(envs)
    }

    return workflowApi.setEnv(name,reqbody);
}

export function getVars(name,versionid){
    return workflowApi.getVar(name,versionid);
}

export function setVars(name,versionid,vars){
    var reqbody = {
        "id" : versionid,
        "var" : _.object(vars)
    }

    return workflowApi.setVar(name,reqbody);
}

export function changeState(name,versionid,state){
    var reqbody = {
        "id" : versionid,
        "state" : state
    }

    return workflowApi.changeState(name,reqbody);
}

export function getToken(name,versionid){
    return workflowApi.getToken(name,versionid);
}
