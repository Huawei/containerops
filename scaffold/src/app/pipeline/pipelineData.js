import {pipelineApi} from "../common/api";

let allPipelines = [];

export function getAllPipelines(){
    return pipelineApi.list();
}

export function getPipeline(name,version){
    // return pipelineApi.data();

    // to be removed below
    // var pipeline = _.find(allPipelines,function(item){
    //     return item.name == name;
    // });
    // return _.find(pipeline.versions,function(item){
    //     return item.version = version;
    // }).data;
    return newPipelineData;
}

export function addPipeline(){
    if(!$('#newpp-form').parsley().validate()){
        return false;
    }
    var name = $("#pp-name").val();
    var version = $("#pp-version").val();

    return pipelineApi.add(name,version);
}

import {PIPELINE_START , PIPELINE_END, PIPELINE_ADD_STAGE, PIPELINE_ADD_ACTION,PIPELINE_STAGE,PIPELINE_ACTION} from "../common/constant";

export function addPipelineVersion(oldversion){
    if(!$('#newpp-version-form').parsley().validate()){
        return false;
    }
    var name = $("#pp-name-newversion").val();
    var version = $("#pp-newversion").val();

    // call api here, return promise

    // to be removed
    var pipeline = _.find(allPipelines,function(item){
        return item.name == name;
    });

    var oldversion = _.find(pipeline.versions,function(item){
        return item.version == oldversion;
    });

    var newversion = {
        "version" : version,
        "data" : [].concat(oldversion.data)
    }
    pipeline.versions.push(newversion);
    return true;
}

export function getEnvs(){
    return [];
}

export var newPipelineData = [
    {
        id: "pipeline-start" + "-" + uuid.v1(),
        type: PIPELINE_START,
        setupData: {}
    },
    
    {
        id: "pipeline-add-stage" + "-" + uuid.v1(),
        type: PIPELINE_ADD_STAGE
    },
    {
        id: "pipeline-end" + "-" + uuid.v1(),
        type: PIPELINE_END,
        setupData: {}
    }
]