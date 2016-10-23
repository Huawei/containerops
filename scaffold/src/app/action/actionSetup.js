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

import * as actionSetupData from "./actionSetupData";
import {notify} from "../common/notify";

// let k8sAdvancedEditor,k8sAdvancedContainer;

// let k8sAD;

export function initActionSetup(action){
    actionSetupData.getActionSetupData(action);

    // action part
    $("#action-component-select").val(actionSetupData.data.action.type);
    $("#action-component-select").on('change',function(){
        actionSetupData.setActionType();    
    });

    $("#action-name").val(actionSetupData.data.action.name);
    $("#action-name").on("blur",function(){
        actionSetupData.setActionName();
    });

    $("#action-timeout").val(actionSetupData.data.action.timeout);
    $("#action-timeout").on("blur",function(){
        actionSetupData.setActionTimeout();
    });

    // $("#k8s-pod-image").val(actionSetupData.data.pod.spec.containers[0].image);
    // $("#k8s-pod-image").on("blur",function(){
    //     actionSetupData.setK8s(k8sAdvancedEditor);
    // });
    // $("#k8s-pod-image").val(actionSetupData.data.action.image);
    // $("#k8s-pod-image").on("blur",function(){
    //     actionSetupData.setActionImage();
    // });
    $("#k8s-pod-image-name").val(actionSetupData.data.action.image.name);
    $("#k8s-pod-image-name").on("blur",function(){
        actionSetupData.setActionImageName();
    });

    $("#k8s-pod-image-tag").val(actionSetupData.data.action.image.tag);
    $("#k8s-pod-image-tag").on("blur",function(){
        actionSetupData.setActionImageTag();
    });

    $("#action-data-from").val(actionSetupData.data.action.datafrom);
    $("#action-data-from").on("blur",function(){
        actionSetupData.setActionDataFrom();
    });

    $("#k8s-ip").val(actionSetupData.data.action.ip);
    $("#k8s-ip").on("blur",function(){
        actionSetupData.setActionIP();
    });

    $("#k8s-api-server").val(actionSetupData.data.action.apiserver);
    $("#k8s-api-server").on("blur",function(){
        actionSetupData.setActionAPIServer();
    });

    // setting way
    if(actionSetupData.data.action.useAdvanced){
        $("#setting-way-advanced").attr("checked","checked");
        $("#setting-way-base").removeAttr("checked");
    }else{
        $("#setting-way-advanced").removeAttr("checked");
        $("#setting-way-base").attr("checked","checked");
    }
    $("#setting-way-advanced").on("click",function(){
        actionSetupData.setActionUseAdvanced(true);
        showSetting();
    });
    $("#setting-way-base").on("click",function(){
        actionSetupData.setActionUseAdvanced(false);
        showSetting();
    });
    showSetting();

    // base setting
    // $("#k8s-service-port").val(actionSetupData.data.service.spec.ports[0].port);
    // $("#k8s-service-port").on("blur",function(){
    //     // actionSetupData.setK8s(k8sAdvancedEditor);
    //     actionSetupData.setServicePort();
    // });

    $("#k8s-cpu-limits").val(actionSetupData.data.pod.spec.containers[0].resources.limits.cpu);
    $("#k8s-cpu-limits").on("blur",function(){
        // actionSetupData.setK8s(k8sAdvancedEditor);
        actionSetupData.setCPULimit();
    });

    $("#k8s-cpu-requests").val(actionSetupData.data.pod.spec.containers[0].resources.requests.cpu);
    $("#k8s-cpu-requests").on("blur",function(){
        // actionSetupData.setK8s(k8sAdvancedEditor);
        actionSetupData.setCPURequest();
    });

    $("#k8s-memory-limits").val(actionSetupData.data.pod.spec.containers[0].resources.limits.memory);
    $("#k8s-memory-limits").on("blur",function(){
        // actionSetupData.setK8s(k8sAdvancedEditor);
        actionSetupData.setMemoryLimit();
    });

    $("#k8s-memory-requests").val(actionSetupData.data.pod.spec.containers[0].resources.requests.memory);
    $("#k8s-memory-requests").on("blur",function(){
        // actionSetupData.setK8s(k8sAdvancedEditor);
        actionSetupData.setMemoryRequest();
    });

    // ports
    showPorts();

    // advanced setting
    $("#serviceCodeEditor").val(JSON.stringify(actionSetupData.data.service_advanced,null,2));
    $("#serviceCodeEditor").on("blur",function(){
        var result = toJsonYaml("service");
        if(result){
            actionSetupData.setServiceAdvanced(result);
        }    
    });

    $("#podCodeEditor").val(JSON.stringify(actionSetupData.data.pod_advanced,null,2));
    $("#podCodeEditor").on("blur",function(){
        var result = toJsonYaml("pod");
        if(result){
            actionSetupData.setPodAdvanced(result);
        }    
    });

    // k8sAD = $.extend(true,{},actionSetupData.data);
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
    if(actionSetupData.data.action.useAdvanced){
        $("#basesetting").addClass("hide");
        $("#advancedsetting").removeClass("hide");
    }else{
        $("#basesetting").removeClass("hide");
        $("#advancedsetting").addClass("hide");
    }
}

function showPorts(){
    $("#ports-setting").empty();
    _.each(actionSetupData.data.service.spec.ports,function(item,index){
        var row = '<div class="row"><div class="col-md-5">'
                        +'<div class="form-group">'
                            +'<label for="normal-field" class="col-sm-4 control-label">'
                                +'Service Port'
                            +'</label>'
                            +'<div class="col-sm-6" data-index="' + index + '">'
                                +'<input type="number" name="k8s-service-port" value="' + item.port + '" class="form-control" required min="0" max="65535">'
                            +'</div>'
                        +'</div>'
                    +'</div>'
                    +'<div class="col-md-5">'
                        +'<div class="form-group">'
                            +'<label for="normal-field" class="col-sm-4 control-label">'
                                +'Node Port'
                            +'</label>'
                            +'<div class="col-sm-6" data-index="' + index + '">' 
                                +'<input type="number" name="k8s-service-node-port" value="' + item.nodePort + '" class="form-control" required min="0" max="65535">'
                            +'</div>'
                        +'</div>'
                    +'</div>'
                    +'<div class="col-md-2" data-index="' + index + '">'
                        +'<span class="glyphicon glyphicon-minus rm-port"></span>'
                    +'</div></div>';
        $("#ports-setting").append(row);
    });
    
    var addrow = `<div class="row"><div class="col-md-2 col-md-offset-10">
                        <span class="glyphicon glyphicon-plus add-port"></span>
                    </div></div>`;
    $("#ports-setting").append(addrow);

    $("input[name=k8s-service-port]").on("blur",function(event){
        actionSetupData.setServicePort(event);
    });

    $("input[name=k8s-service-node-port]").on("blur",function(event){
        actionSetupData.setServiceNodePort(event);
    });

    $(".rm-port").on("click",function(event){
        actionSetupData.removeServicePorts(event);
        showPorts();
    });

    $(".add-port").on("click",function(){
        actionSetupData.addServicePort();
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
//             actionSetupData.setK8s(k8sAdvancedEditor);
//         }
//     };

//     k8sAdvancedEditor = new JSONEditor(k8sAdvancedContainer, treeOptions);
//     k8sAdvancedEditor.set(k8sAD);
    
//     k8sAdvancedEditor.expandAll();
// }