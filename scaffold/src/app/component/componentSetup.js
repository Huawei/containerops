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

import * as componentSetupData from "./componentSetupData";
import {notify} from "../common/notify";

export function initComponentSetup(component){
    componentSetupData.getComponentSetupData(component);

    // action part
    $("#action-component-select").val(componentSetupData.data.action.type);
    $("#action-component-select").on('change',function(){
        componentSetupData.setActionType();    
    });

    $("#action-timeout").val(componentSetupData.data.action.timeout);
    $("#action-timeout").on("blur",function(){
        componentSetupData.setActionTimeout();
    });

    $("#k8s-pod-image-name").val(componentSetupData.data.action.image.name);
    $("#k8s-pod-image-name").on("blur",function(){
        componentSetupData.setActionImageName();
    });

    $("#k8s-pod-image-tag").val(componentSetupData.data.action.image.tag);
    $("#k8s-pod-image-tag").on("blur",function(){
        componentSetupData.setActionImageTag();
    });

    $("#action-data-from").val(componentSetupData.data.action.datafrom);
    $("#action-data-from").on("blur",function(){
        componentSetupData.setActionDataFrom();
    });

    // setting way
    if(componentSetupData.data.action.useAdvanced){
        $("#setting-way-advanced").prop("checked",true);
        $("#setting-way-base").prop("checked",false);
    }else{
        $("#setting-way-advanced").prop("checked",false);
        $("#setting-way-base").prop("checked",true);
    }
    $("#setting-way-advanced").on("click",function(){
        componentSetupData.setActionUseAdvanced(true);
        showSetting();
    });
    $("#setting-way-base").on("click",function(){
        componentSetupData.setActionUseAdvanced(false);
        showSetting();
    });
    showSetting();

    // base setting
    $("#k8s-cpu-limits").val(componentSetupData.data.pod.spec.containers[0].resources.limits.cpu);
    $("#k8s-cpu-limits").on("blur",function(){
        componentSetupData.setCPULimit();
    });

    $("#k8s-cpu-requests").val(componentSetupData.data.pod.spec.containers[0].resources.requests.cpu);
    $("#k8s-cpu-requests").on("blur",function(){
        componentSetupData.setCPURequest();
    });

    $("#k8s-memory-limits").val(componentSetupData.data.pod.spec.containers[0].resources.limits.memory);
    $("#k8s-memory-limits").on("blur",function(){
        componentSetupData.setMemoryLimit();
    });

    $("#k8s-memory-requests").val(componentSetupData.data.pod.spec.containers[0].resources.requests.memory);
    $("#k8s-memory-requests").on("blur",function(){
        componentSetupData.setMemoryRequest();
    });

    // ports
    $("#service-type-select").val(componentSetupData.data.service.spec.type);
    $("#service-type-select").on('change',function(){
        componentSetupData.setServiceType();    
        showComponentPorts();
    });
    showComponentPorts();
    
    // advanced setting
    $("#serviceCodeEditor").val(JSON.stringify(componentSetupData.data.service_advanced,null,2));
    $("#serviceCodeEditor").on("blur",function(){
        var result = toJsonYaml("service");
        if(result){
            componentSetupData.setServiceAdvanced(result);
        }    
    });

    $("#podCodeEditor").val(JSON.stringify(componentSetupData.data.pod_advanced,null,2));
    $("#podCodeEditor").on("blur",function(){
        var result = toJsonYaml("pod");
        if(result){
            componentSetupData.setPodAdvanced(result);
        }    
    });
}

function showSetting(){
    if(componentSetupData.data.action.useAdvanced){
        $("#basesetting").addClass("hide");
        $("#advancedsetting").removeClass("hide");
    }else{
        $("#basesetting").removeClass("hide");
        $("#advancedsetting").addClass("hide");
    }
}

function showComponentPorts(){
    $("#ports-setting").empty();
    _.each(componentSetupData.data.service.spec.ports,function(item,index){
        var row = `<div class="port-row">`;

        if(componentSetupData.getUseNodePort()){
            row +=  `<div class="port-div">`;
        }else{
            row += `<div class="port-div no-use-node-port">`;
        }    

        row += `<div>
                    <label for="normal-field" class="col-sm-4 control-label">
                        Port
                    </label>
                    <div class="col-sm-7" data-index="` + index + `">
                        <input type="number" name="k8s-service-port" value="` + item.port + `" class="form-control" required min="0" max="65535">
                    </div>
                </div>
                </div>`;
        if(componentSetupData.getUseNodePort()){
            row +=  `<div class="port-div">`;
        }else{
            row += `<div class="port-div no-use-node-port">`;
        }
        
        row +=  `<div>
                    <label for="normal-field" class="col-sm-4 control-label">
                        Target Port
                    </label>
                    <div class="col-sm-7" data-index="` + index + `">
                        <input type="number" name="k8s-service-target-port" value="` + item.targetPort + `" class="form-control" required min="0" max="65535">
                    </div>
                </div>
                </div>`;

        if(componentSetupData.getUseNodePort()){
            row += `<div class="port-div">
                            <div>
                                <label for="normal-field" class="col-sm-4 control-label">
                                    Node Port
                                </label>
                                <div class="col-sm-7" data-index="` + index + `">
                                    <input type="number" name="k8s-service-node-port" value="` + item.nodePort + `" class="form-control" required min="0" max="65535">
                                </div>
                            </div>
                        </div>`;
        }
                        
        row +=  `<div class="port-remove-div rm-port" data-index="` + index + `">
                            <span class="glyphicon glyphicon-remove"></span>
                        </div>
                    </div>`;
        $("#ports-setting").append(row);
    });
    
    var addrow = `<button type="button" class="btn btn-success add-port">
                        <i class="glyphicon glyphicon-plus" style="top:1px"></i>
                        <span style="margin-left:5px">Add Port</span>
                    </button>`;
    $("#ports-setting").append(addrow);

    $("input[name=k8s-service-port]").on("blur",function(event){
        componentSetupData.setServicePort(event);
    });

    $("input[name=k8s-service-target-port]").on("blur",function(event){
        componentSetupData.setServiceTargetPort(event);
    });

    $("input[name=k8s-service-node-port]").on("blur",function(event){
        componentSetupData.setServiceNodePort(event);
    });

    $(".rm-port").on("click",function(event){
        componentSetupData.removeServicePorts(event);
        showComponentPorts();
    });

    $(".add-port").on("click",function(){
        componentSetupData.addServicePort();
        showComponentPorts();
    });
}

function toJsonYaml(type){
    var value,result;
    if(type == "service"){
        value = $("#serviceCodeEditor").val();
    }else if(type == "pod"){
        value = $("#podCodeEditor").val();
    }

    try{
        result = JSON.parse(value);
    }catch(e){
       try{
        result = jsyaml.safeLoad(value);
       }catch(e){
        notify("Your advanced " + type + " setting is not a legal json or yaml.","error");
        result = false;
       }
    }
    if(!result){
        notify("Your advanced " + type + " setting is not a legal json or yaml.","error");
    }
    return result;
}