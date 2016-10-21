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

import {pipelineApi} from "../common/api";

let allPipelines = [];

export function getAllPipelines(){
    return pipelineApi.list();
}

export function getPipeline(name,versionid){
    return pipelineApi.data(name,versionid);
}

export function addPipeline(){
    if(!$('#newpp-form').parsley().validate()){
        return false;
    }
    var name = $("#pp-name").val();
    var version = $("#pp-version").val();

    return pipelineApi.add(name,version);
}

export function savePipeline(name,version,versionid,nodes,lines){
    var reqbody = {
        "id" : versionid,
        "version" : version.toString(),
        "define":{
            "lineList" : lines,
            "stageList" : nodes
        }
    }

    return pipelineApi.save(name,reqbody);
}

export function addPipelineVersion(name, versionid, nodes,lines){
    if(!$('#newpp-version-form').parsley().validate()){
        return false;
    }else{
        var version = $("#pp-newversion").val();
        return savePipeline(name, version, versionid, nodes,lines)
    }
}

export function getEnvs(){
    return [];
}