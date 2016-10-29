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

import * as componentSetupData from "./componentSetupData";
import {notify} from "../common/notify";

// let k8sAdvancedEditor,k8sAdvancedContainer;

// let k8sAD;

export function initComponentSetup(component){
    componentSetupData.getComponentSetupData(component);

    // action part
    $("#action-component-select").val(componentSetupData.data.action.type);
    $("#action-component-select").on('change',function(){
        componentSetupData.setActionType();    
    });

    $("#action-name").val(componentSetupData.data.action.name);
    $("#action-name").on("blur",function(){
        componentSetupData.setActionName();
    });

    $("#action-timeout").val(componentSetupData.data.action.timeout);
    $("#action-timeout").on("blur",function(){
        componentSetupData.setActionTimeout();
    });

    // $("#k8s-ip").val(componentSetupData.data.action.ip);
    // $("#k8s-ip").on("blur",function(){
    //     componentSetupData.setActionIP();
    // });

    // $("#k8s-pod-image").val(componentSetupData.data.pod.spec.containers[0].image);
    // $("#k8s-pod-image").on("blur",function(){
    //     componentSetupData.setK8s(k8sAdvancedEditor);
    // });
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
        $("#setting-way-advanced").attr("checked","checked");
        $("#setting-way-base").removeAttr("checked");
    }else{
        $("#setting-way-advanced").removeAttr("checked");
        $("#setting-way-base").attr("checked","checked");
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
    // $("#k8s-service-port").val(componentSetupData.data.service.spec.ports[0].port);
    // $("#k8s-service-port").on("blur",function(){
    //     // componentSetupData.setK8s(k8sAdvancedEditor);
    //     componentSetupData.setServicePort();
    // });

    $("#k8s-cpu-limits").val(componentSetupData.data.pod.spec.containers[0].resources.limits.cpu);
    $("#k8s-cpu-limits").on("blur",function(){
        // componentSetupData.setK8s(k8sAdvancedEditor);
        componentSetupData.setCPULimit();
    });

    $("#k8s-cpu-requests").val(componentSetupData.data.pod.spec.containers[0].resources.requests.cpu);
    $("#k8s-cpu-requests").on("blur",function(){
        // componentSetupData.setK8s(k8sAdvancedEditor);
        componentSetupData.setCPURequest();
    });

    $("#k8s-memory-limits").val(componentSetupData.data.pod.spec.containers[0].resources.limits.memory);
    $("#k8s-memory-limits").on("blur",function(){
        // componentSetupData.setK8s(k8sAdvancedEditor);
        componentSetupData.setMemoryLimit();
    });

    $("#k8s-memory-requests").val(componentSetupData.data.pod.spec.containers[0].resources.requests.memory);
    $("#k8s-memory-requests").on("blur",function(){
        // componentSetupData.setK8s(k8sAdvancedEditor);
        componentSetupData.setMemoryRequest();
    });

    // ports
    showPorts();

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



    // k8sAD = $.extend(true,{},componentSetupData.data);
    // delete k8sAD.action;
    // delete k8sAD.service.spec.ports[0].port;
    // delete k8sAD.pod.spec.containers[0].resources;
    // delete k8sAD.pod.spec.containers[0].image;

    // initK8sForm();


    // $("#k8s-advanced").on("click",function(){
    //     $("#k8s-advanced").hide();
    //     $("#close-k8s-advanced").show();
    //     $("#advanced").parent().show();
    // })

    // $("#close-k8s-advanced").on("click",function(){
    //     $("#k8s-advanced").show();
    //     $("#close-k8s-advanced").hide();
    //     $("#advanced").parent().hide();
    // })
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

function showPorts(){
    $("#ports-setting").empty();
    _.each(componentSetupData.data.service.spec.ports,function(item,index){
        var row = '<div class="port-row"><div class="port-div">'
                        +'<div>'
                            +'<label for="normal-field" class="col-sm-4 control-label">'
                                +'Target Port'
                            +'</label>'
                            +'<div class="col-sm-7" data-index="' + index + '">'
                                +'<input type="number" name="k8s-service-port" value="' + item.port + '" class="form-control" required min="0" max="65535">'
                            +'</div>'
                        +'</div>'
                    +'</div>'
                    +'<div class="target-port-div">'
                        +'<div>'
                            +'<label for="normal-field" class="col-sm-4 control-label">'
                                +'Node Port'
                            +'</label>'
                            +'<div class="col-sm-7" data-index="' + index + '">' 
                                +'<input type="number" name="k8s-service-node-port" value="' + item.nodePort + '" class="form-control" required min="0" max="65535">'
                            +'</div>'
                        +'</div>'
                    +'</div>'
                    +'<div class="port-remove-div rm-port" data-index="' + index + '">'
                        +'<span class="glyphicon glyphicon-remove"></span>'
                    +'</div></div>';
        $("#ports-setting").append(row);
    });
    
    var addrow = `<button type="button" class="btn btn-success add-port">
                        <i class="glyphicon glyphicon-plus" style="top:1px"></i>&nbsp;&nbsp;Add Port
                    </button>`;
    $("#ports-setting").append(addrow);

    $("input[name=k8s-service-port]").on("blur",function(event){
        componentSetupData.setServicePort(event);
    });

    $("input[name=k8s-service-node-port]").on("blur",function(event){
        componentSetupData.setServiceNodePort(event);
    });

    $(".rm-port").on("click",function(event){
        componentSetupData.removeServicePorts(event);
        showPorts();
    });

    $(".add-port").on("click",function(){
        componentSetupData.addServicePort();
        showPorts();
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

// function initK8sForm(){
//     k8sAdvancedContainer = $("#advanced")[0];
//     initK8sAdvanced();
// }

// function initK8sAdvanced(){
//     if(k8sAdvancedEditor){
//         k8sAdvancedEditor.destroy();
//     }

//     var treeOptions = {
//         "mode": "tree",
//         "search": true,
//         "onChange" : function(){
//             componentSetupData.setK8s(k8sAdvancedEditor);
//         }
//     };

//     k8sAdvancedEditor = new JSONEditor(k8sAdvancedContainer, treeOptions);
//     k8sAdvancedEditor.set(k8sAD);
    
//     k8sAdvancedEditor.expandAll();
// }