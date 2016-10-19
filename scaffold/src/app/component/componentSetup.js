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

import * as componentSetupData from "./componentSetupData";

let k8sAdvancedEditor,k8sAdvancedContainer;

let k8sAD;

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

    $("#k8s-pod-image").val(componentSetupData.data.pod.spec.containers[0].image);
    $("#k8s-pod-image").on("blur",function(){
        componentSetupData.setK8s(k8sAdvancedEditor);
    });

    $("#k8s-ip").val(componentSetupData.data.action.ip);
    $("#k8s-ip").on("blur",function(){
        componentSetupData.setActionIP();
    });

    $("#k8s-service-port").val(componentSetupData.data.service.spec.ports[0].port);
    $("#k8s-service-port").on("blur",function(){
        componentSetupData.setK8s(k8sAdvancedEditor);
    });

    $("#k8s-cpu-limits").val(componentSetupData.data.pod.spec.containers[0].resources.limits[0].cpu);
    $("#k8s-cpu-limits").on("blur",function(){
        componentSetupData.setK8s(k8sAdvancedEditor);
    });

    $("#k8s-cpu-requests").val(componentSetupData.data.pod.spec.containers[0].resources.requests[0].cpu);
    $("#k8s-cpu-requests").on("blur",function(){
        componentSetupData.setK8s(k8sAdvancedEditor);
    });

    $("#k8s-memory-limits").val(componentSetupData.data.pod.spec.containers[0].resources.limits[0].memory);
    $("#k8s-memory-limits").on("blur",function(){
        componentSetupData.setK8s(k8sAdvancedEditor);
    });

    $("#k8s-memory-requests").val(componentSetupData.data.pod.spec.containers[0].resources.requests[0].memory);
    $("#k8s-memory-requests").on("blur",function(){
        componentSetupData.setK8s(k8sAdvancedEditor);
    });

    k8sAD = $.extend(true,{},componentSetupData.data);
    delete k8sAD.action;
    delete k8sAD.service.spec.ports[0].port;
    delete k8sAD.pod.spec.containers[0].resources;
    delete k8sAD.pod.spec.containers[0].image;

    initK8sForm();


    $("#k8s-advanced").on("click",function(){
        $("#k8s-advanced").hide();
        $("#close-k8s-advanced").show();
        $("#advanced").parent().show();
    })

    $("#close-k8s-advanced").on("click",function(){
        $("#k8s-advanced").show();
        $("#close-k8s-advanced").hide();
        $("#advanced").parent().hide();
    })
}

function initK8sForm(){
    k8sAdvancedContainer = $("#advanced")[0];
    initK8sAdvanced();
}

function initK8sAdvanced(){
    if(k8sAdvancedEditor){
        k8sAdvancedEditor.destroy();
    }

    var treeOptions = {
        "mode": "tree",
        "search": true,
        "onChange" : function(){
            componentSetupData.setK8s(k8sAdvancedEditor);
        }
    };

    k8sAdvancedEditor = new JSONEditor(k8sAdvancedContainer, treeOptions);
    k8sAdvancedEditor.set(k8sAD);
    
    k8sAdvancedEditor.expandAll();
}