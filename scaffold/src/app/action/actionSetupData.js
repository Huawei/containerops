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

export function setK8s(k8sAdvancedEditor){
    var k8s_ad = k8sAdvancedEditor.get();
    k8s_ad.service.spec.ports[0].port = $("#k8s-service-port").val();
    k8s_ad.pod.spec.containers[0].image = $("#k8s-pod-image").val();
    var resources = {
      "limits":[{
        "cpu": $("#k8s-cpu-limits").val(), 
        "memory": $("#k8s-memory-limits").val()
      }],
      "requests":[{
        "cpu": $("#k8s-cpu-requests").val(), 
        "memory": $("#k8s-memory-requests").val()
      }]
    }
    k8s_ad.pod.spec.containers[0].resources = resources;

    data.service = k8s_ad.service;
    data.pod = k8s_ad.pod;
}

var metadata = {
  "action" : {
    "type" : "Kubernetes",
    "name" : "",
    "timeout" : "",
    "ip" : ""
  },
  "service" : {
    "metadata": {
      "name": "",
      "deletionTimestamp": "",
      "deletionGracePeriodSeconds": 0
    },
    "spec": {
      "ports": [
        {
          "protocol": "tcp",
          "port": "",
          "targetPort": "",
          "nodePort": ""
        }
      ],
      "clusterIP": "",
      "type": "",
      "sessionAffinity": "",
      "loadBalancerIP": ""
    }
  },
  "pod" : {
    "metadata": {
      "name" : "",
      "deletionTimestamp": "",
      "deletionGracePeriodSeconds": 0
    },
    "spec": {
      "containers": [
        {
          "name": "",
          "image" : "",
          "command": [],
          "args": [],
          "workingDir": "",
          "ports": [
            {
              "name": "",
              "hostPort": 0,
              "containerPort": 0,
              "protocol": "",
              "hostIP": ""
            }
          ],
          "env": [
            {
              "name": "",
              "value": ""
            }
          ],
          "resources": {
            "limits":[ {"cpu": 4.0, "memory": "99Mi"} ],
            "requests":[ {"cpu": 4.0, "memory": "99Mi"} ]
          },
          "livenessProbe": {
            "exec": {
              "command": []
            },
            "httpGet": {
              "path": "",
              "port": "",
              "host": "",
              "scheme": ""
            },
            "tcpSocket": {
              "port": ""
            },
            "initialDelaySeconds": 0,
            "timeoutSeconds": 0,
            "periodSeconds": 0,
            "successThreshold": 0,
            "failureThreshold": 0
          },
          "readinessProbe": {
            "exec": {
              "command": [
                "string"
              ]
            },
            "httpGet": {
              "path": "",
              "port": "",
              "host": "",
              "scheme": ""
            },
            "tcpSocket": {
              "port": ""
            },
            "initialDelaySeconds": 0,
            "timeoutSeconds": 0,
            "periodSeconds": 0,
            "successThreshold": 0,
            "failureThreshold": 0
          },
          "lifecycle": {
            "postStart": {
              "exec": {
                "command": []
              },
              "httpGet": {
                "path": "",
                "port": "",
                "host": "",
                "scheme": ""
              },
              "tcpSocket": {
                "port": ""
              }
            },
            "preStop": {
              "exec": {
                "command": [
                  "string"
                ]
              },
              "httpGet": {
                "path": "",
                "port": "",
                "host": "",
                "scheme": ""
              },
              "tcpSocket": {
                "port": ""
              }
            }
          },
          "terminationMessagePath": "",
          "imagePullPolicy": "", 
          "securityContext": {
            "capabilities": {
              "add": [
                {}
              ],
              "drop": [
                {}
              ]
            },
            "privileged": true,
            "seLinuxOptions": {
              "user": "",
              "role": "",
              "type": "",
              "level": ""
            },
            "runAsUser": 0,
            "runAsNonRoot": true
          },
          "stdin": true,
          "stdinOnce": true, 
          "tty": true
        }
      ],
      "restartPolicy": "",
      "terminationGracePeriodSeconds": 0,
      "activeDeadlineSeconds": 0,
      "dnsPolicy": "",
      "nodeSelector":"",
      "serviceAccountName": "",
      "nodeName": "",
      "hostNetwork": true,
      "hostPID": true,
      "hostIPC": true,
      "securityContext": {
        "seLinuxOptions": {
          "user": "",
          "role": "",
          "type": "",
          "level": ""
        },
        "runAsUser": 0,
        "runAsNonRoot": true,
        "supplementalGroups": [
          {}
        ],
        "fsGroup": 0
      },
      "imagePullSecrets": [
        {
          "name": ""
        }
      ]
    }
   }
}



