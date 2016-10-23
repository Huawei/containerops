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

export let data;

export function getActionSetupData(action){
    if(!_.isUndefined(action.setupData) && !_.isEmpty(action.setupData)){
      data = action.setupData;
    }else{
      data = $.extend(true,{},metadata);
      action.setupData = data;
    } 
}

export function setActionType(){
    data.action.type = $("#action-component-select").val();
}

export function setActionName(){
    data.action.name = $("#action-name").val();
}

export function setActionTimeout(){
    data.action.timeout = $("#action-timeout").val();
}

export function setActionIP(){
    data.action.ip = $("#k8s-ip").val();
}

export function setActionAPIServer(){
    data.action.apiserver = $("#k8s-api-server").val();
}

// export function setActionImage(){
//   data.action.image = $("#k8s-pod-image").val();
// }

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
export function setServicePort(event){
  var target = $(event.currentTarget);
  var index = target.parent().data("index");
  var value = target.val();
  data.service.spec.ports[index].port = parseInt(value);
  data.service.spec.ports[index].targetPort = parseInt(value);
}

export function setServiceNodePort(event){
  var target = $(event.currentTarget);
  var index = target.parent().data("index");
  var value = target.val();
  data.service.spec.ports[index].nodePort = parseInt(value);
}

export function removeServicePorts(event){
  var index = $(event.currentTarget).parent().data("index");
  data.service.spec.ports.splice(index,1);
}

export function addServicePort(){
  data.service.spec.ports.push({
    "port" : "",
    "nodePort" : ""
  })
}
// ports end

export function setCPULimit(){
  data.pod.spec.containers[0].resources.limits.cpu = $("#k8s-cpu-limits").val();
}

export function setCPURequest(){
  data.pod.spec.containers[0].resources.requests.cpu = $("#k8s-cpu-requests").val();
}

export function setMemoryLimit(){
  data.pod.spec.containers[0].resources.limits.memory = $("#k8s-memory-limits").val();
}

export function setMemoryRequest(){
  data.pod.spec.containers[0].resources.requests.memory = $("#k8s-memory-requests").val();
}

export function setServiceAdvanced(value){
  data.service_advanced = value;
}

export function setPodAdvanced(value){
  data.pod_advanced = value;
}

// export function setK8s(k8sAdvancedEditor){
//     var k8s_ad = k8sAdvancedEditor.get();
//     k8s_ad.service.spec.ports[0].port = $("#k8s-service-port").val();
//     k8s_ad.pod.spec.containers[0].image = $("#k8s-pod-image").val();
//     var resources = {
//       "limits":[{
//         "cpu": $("#k8s-cpu-limits").val(), 
//         "memory": $("#k8s-memory-limits").val()
//       }],
//       "requests":[{
//         "cpu": $("#k8s-cpu-requests").val(), 
//         "memory": $("#k8s-memory-requests").val()
//       }]
//     }
//     k8s_ad.pod.spec.containers[0].resources = resources;

//     data.service = k8s_ad.service;
//     data.pod = k8s_ad.pod;
// }

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
            "limits":{"cpu": 4.0, "memory": "99Mi"},
            "requests":{"cpu": 4.0, "memory": "99Mi"}
          }
        }
      ]
    }
   },
   "service_advanced" : {},
   "pod_advanced" : {}
}

// var metadata = {
//   "action" : {
//     "type" : "Kubernetes",
//     "name" : "",
//     "timeout" : "",
//     "ip" : ""
//   },
//   "service" : {
//     "metadata": {
//       "name": "",
//       "deletionTimestamp": "",
//       "deletionGracePeriodSeconds": 0
//     },
//     "spec": {
//       "ports": [
//         {
//           "protocol": "tcp",
//           "port": "",
//           "targetPort": "",
//           "nodePort": ""
//         }
//       ],
//       "clusterIP": "",
//       "type": "",
//       "sessionAffinity": "",
//       "loadBalancerIP": ""
//     }
//   },
//   "pod" : {
//     "metadata": {
//       "name" : "",
//       "deletionTimestamp": "",
//       "deletionGracePeriodSeconds": 0
//     },
//     "spec": {
//       "containers": [
//         {
//           "name": "",
//           "image" : "",
//           "command": [],
//           "args": [],
//           "workingDir": "",
//           "ports": [
//             {
//               "name": "",
//               "hostPort": 0,
//               "containerPort": 0,
//               "protocol": "",
//               "hostIP": ""
//             }
//           ],
//           "env": [
//             {
//               "name": "",
//               "value": ""
//             }
//           ],
//           "resources": {
//             "limits":[ {"cpu": 4.0, "memory": "99Mi"} ],
//             "requests":[ {"cpu": 4.0, "memory": "99Mi"} ]
//           },
//           "livenessProbe": {
//             "exec": {
//               "command": []
//             },
//             "httpGet": {
//               "path": "",
//               "port": "",
//               "host": "",
//               "scheme": ""
//             },
//             "tcpSocket": {
//               "port": ""
//             },
//             "initialDelaySeconds": 0,
//             "timeoutSeconds": 0,
//             "periodSeconds": 0,
//             "successThreshold": 0,
//             "failureThreshold": 0
//           },
//           "readinessProbe": {
//             "exec": {
//               "command": [
//                 "string"
//               ]
//             },
//             "httpGet": {
//               "path": "",
//               "port": "",
//               "host": "",
//               "scheme": ""
//             },
//             "tcpSocket": {
//               "port": ""
//             },
//             "initialDelaySeconds": 0,
//             "timeoutSeconds": 0,
//             "periodSeconds": 0,
//             "successThreshold": 0,
//             "failureThreshold": 0
//           },
//           "lifecycle": {
//             "postStart": {
//               "exec": {
//                 "command": []
//               },
//               "httpGet": {
//                 "path": "",
//                 "port": "",
//                 "host": "",
//                 "scheme": ""
//               },
//               "tcpSocket": {
//                 "port": ""
//               }
//             },
//             "preStop": {
//               "exec": {
//                 "command": [
//                   "string"
//                 ]
//               },
//               "httpGet": {
//                 "path": "",
//                 "port": "",
//                 "host": "",
//                 "scheme": ""
//               },
//               "tcpSocket": {
//                 "port": ""
//               }
//             }
//           },
//           "terminationMessagePath": "",
//           "imagePullPolicy": "", 
//           "securityContext": {
//             "capabilities": {
//               "add": [
//                 {}
//               ],
//               "drop": [
//                 {}
//               ]
//             },
//             "privileged": true,
//             "seLinuxOptions": {
//               "user": "",
//               "role": "",
//               "type": "",
//               "level": ""
//             },
//             "runAsUser": 0,
//             "runAsNonRoot": true
//           },
//           "stdin": true,
//           "stdinOnce": true, 
//           "tty": true
//         }
//       ],
//       "restartPolicy": "",
//       "terminationGracePeriodSeconds": 0,
//       "activeDeadlineSeconds": 0,
//       "dnsPolicy": "",
//       "nodeSelector":"",
//       "serviceAccountName": "",
//       "nodeName": "",
//       "hostNetwork": true,
//       "hostPID": true,
//       "hostIPC": true,
//       "securityContext": {
//         "seLinuxOptions": {
//           "user": "",
//           "role": "",
//           "type": "",
//           "level": ""
//         },
//         "runAsUser": 0,
//         "runAsNonRoot": true,
//         "supplementalGroups": [
//           {}
//         ],
//         "fsGroup": 0
//       },
//       "imagePullSecrets": [
//         {
//           "name": ""
//         }
//       ]
//     }
//    }
// }



