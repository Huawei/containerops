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
 
export let data;

export function getComponentSetupData(component){
    if(!_.isUndefined(component.setupData) && !_.isEmpty(component.setupData)){
      data = component.setupData;
    }else{
      data = $.extend(true,{},metadata);
      component.setupData = data;
    } 
}

export function setActionType(){
    data.action.type = $("#action-component-select").val();
}

export function setActionTimeout(){
    data.action.timeout = $("#action-timeout").val();
}

export function setActionImageName(){
  data.action.image.name = $("#k8s-pod-image-name").val();
}

export function setActionImageTag(){
  data.action.image.tag = $("#k8s-pod-image-tag").val();
}

export function setActionDataFrom(){
    data.action.datafrom = $("#action-data-from").val();
}

// setting way
export function setActionUseAdvanced(value){
    data.action.useAdvanced = value;
}

// ports
export function getUseNodePort(){
    if(data.service.spec.type == "NodePort"){
        return true;
    }else{
        return false;
    }
}

export function setServiceType(){
    data.service.spec.type = $("#service-type-select").val();
    if(data.service.spec.type == "NodePort"){
        _.each(data.service.spec.ports,function(item){
          if(_.isUndefined(item.nodePort)){
            item.nodePort = "";
          }
        })
    }else{
        _.each(data.service.spec.ports,function(item){
          if(!_.isUndefined(item.nodePort)){
            delete item.nodePort;
          }
        })
    }
}

export function setServicePort(event){
  var target = $(event.currentTarget);
  var index = target.parent().data("index");
  var value = target.val();
  data.service.spec.ports[index].port = parseInt(value);
}

export function setServiceTargetPort(event){
  var target = $(event.currentTarget);
  var index = target.parent().data("index");
  var value = target.val();
  data.service.spec.ports[index].targetPort = parseInt(value);
}

export function setServiceNodePort(event){
  var target = $(event.currentTarget);
  var index = target.parent().data("index");
  var value = target.val();
  data.service.spec.ports[index].nodePort = parseInt(value);
}

export function removeServicePorts(event){
  var index = $(event.currentTarget).data("index");
  data.service.spec.ports.splice(index,1);
}

export function addServicePort(){
  if(getUseNodePort()){
      data.service.spec.ports.push({
        "port" : "",
        "targetPort" : "",
        "nodePort" : ""
      });
  }else{
      data.service.spec.ports.push({
        "port" : "",
        "targetPort" : ""
      });
  }  
}
// ports end

export function setCPULimit(){
  data.pod.spec.containers[0].resources.limits.cpu = $("#k8s-cpu-limits").val();
}

export function setCPURequest(){
  data.pod.spec.containers[0].resources.requests.cpu = $("#k8s-cpu-requests").val();
}

export function getMemoryLimit(){
  var value = data.pod.spec.containers[0].resources.limits.memory;
  return Number(value.substring(0,value.length-2));
}

export function setMemoryLimit(){
  var value = $("#k8s-memory-limits").val();
  data.pod.spec.containers[0].resources.limits.memory = value.toString() + "Mi";
}

export function getMemoryRequest(){
  var value = data.pod.spec.containers[0].resources.requests.memory;
  return Number(value.substring(0,value.length-2));
}

export function setMemoryRequest(){
  var value = $("#k8s-memory-requests").val();
  data.pod.spec.containers[0].resources.requests.memory = value.toString() + "Mi";
}

export function setServiceAdvanced(value){
  data.service_advanced = value;
}

export function setPodAdvanced(value){
  data.pod_advanced = value;
}

var metadata = {
  "action" : {
    "type" : "Kubernetes",
    "name" : "",
    "timeout" : "",
    "ip" : "",
    "apiserver" : "",
    "image" : {
      "name" : "",
      "tag" : ""
    },
    "useAdvanced" : false
  },
  "service" : {
    "spec": {
      "type":"NodePort",
      "ports": [
        {
          "port": "",
          "targetPort" : "",
          "nodePort" : ""
        }
      ]
    }
  },
  "pod" : {
    "spec": {
      "containers": [
        {
          "resources": {
            "limits":{"cpu": 0.2, "memory": "1024Mi"},
            "requests":{"cpu": 0.1, "memory": "128Mi"}
          }
        }
      ]
    }
   },
   "service_advanced" : {},
   "pod_advanced" : {}
}
