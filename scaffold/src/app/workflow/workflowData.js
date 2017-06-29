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
import {ajaxCall} from "../common/api";

let allWorkflows = [];

export function getAllWorkflows(){
    return ajaxCall("workflow.list");
}

export function getWorkflow(name,versionid){
    var params = {
        "workflowName" : name,
        "workflowID" : versionid
    }
    return ajaxCall("workflow.detail",params);
}

export function addWorkflow(){
    if(!$('#newpp-form').parsley().validate()){
        return false;
    }

    var reqbody = {
        "name" : $("#pp-name").val(),
        "version" : $("#pp-version").val()
    }

    return ajaxCall("workflow.add",{},reqbody);
}

export function saveWorkflow(name,version,versionid,nodes,lines,setting){
    var params = {
        "workflowName" : name
    }

    var reqbody = {
        "id" : versionid,
        "version" : version.toString(),
        "define":{
            "lineList" : lines,
            "stageList" : nodes,
            "setting" : setting
        }
    }

    return ajaxCall("workflow.save",params,reqbody);
}

export function addWorkflowVersion(name, versionid, nodes,lines,setting){
    if(!$('#newpp-version-form').parsley().validate()){
        return false;
    }else{
        var version = $("#pp-newversion").val();
        return saveWorkflow(name, version, versionid, nodes,lines,setting)
    }
}

export function getEnvs(name,versionid){
    var params = {
        "workflowName" : name,
        "workflowID" : versionid
    }

    return ajaxCall("workflow.getEnv",params);
}

export function setEnvs(name,versionid,envs){
    var params = {
        "workflowName" : name
    }

    var reqbody = {
        "id" : versionid,
        "env" : _.object(envs)
    }

    return ajaxCall("workflow.setEnv",params,reqbody);
}

export function getVars(name,versionid){
    var params = {
        "workflowName" : name,
        "workflowID" : versionid
    }

    return ajaxCall("workflow.getVar",params);
}

export function setVars(name,versionid,vars){
    var params = {
        "workflowName" : name
    }

    var reqbody = {
        "id" : versionid,
        "var" : _.object(vars)
    }

    return ajaxCall("workflow.setVar",params,reqbody);
}

export function changeState(name,versionid,state){
    var params = {
        "workflowName" : name
    }

    var reqbody = {
        "id" : versionid,
        "state" : state
    }

    return ajaxCall("workflow.changeState",params,reqbody);
}

export function getToken(name,versionid){
    var params = {
        "workflowName" : name,
        "workflowID" : versionid
    }

    return ajaxCall("workflow.getToken",params);
}
