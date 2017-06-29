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

import {notify} from "../common/notify";
import * as constant from "../common/constant";
import {isAvailableVar,getValue} from "../workflow/workflowVar";

// validate
export function workflowCheck(data,setting){
    var completeness = true;
    for(var index=0;index<data.length;index++){
        var item = data[index];
        if(item.type == constant.WORKFLOW_START){
           completeness = checkWorkflowStart(item);
        }else if(item.type == constant.WORKFLOW_STAGE){
           completeness = checkWorkflowStage(item,index);
        }
        if(!completeness){
            break;
        }
    }

    if(completeness){
        completeness = checkSetting(setting);
    }

    if(completeness){
        notify("Workflow is available.","success");
    }
    return completeness;
}

function checkWorkflowStart(data){
    var completeness = true;
    if(_.isUndefined(data.outputJson) || _.isEmpty(data.outputJson)){
        notify("No any outputs found ---- < Start stage >","info");
        completeness = false;
    }else{
        for(var i=0;i<data.outputJson.length;i++){
            var item = data.outputJson[i];
            completeness = !_.isEmpty(item.json);
            if(!completeness){
                notify("Output json missed ---- < Start stage / Output " + (i+1)+" >","info");
                break;
            }

            completeness = !_.isEmpty(item.event);
            if(!completeness){
                notify("Output event missed ---- < Start stage / Output " + (i+1)+" >","info");
                break;
            }

            // if(isUsingGlobalVar(item.event)){
            //     completeness = isAvailableVar(item.event);
            //     if(!completeness){
            //         notify("Output event is using an unknown global variable ---- < Start stage / Output " + (i+1)+" >","info");
            //         break;
            //     }
            // }    
        }
    }

    return completeness;
}

function checkWorkflowStage(data,index){
    var completeness = true;
    if(_.isUndefined(data.setupData) || _.isEmpty(data.setupData)){
        notify("Base config missed ---- < Stage No. " + index+" >","info");
        completeness = false;
    }else if(_.isEmpty(data.setupData.name)){
        notify("Name missed ---- < Stage No. " + index+" >","info");
        completeness = false;
    }else if(_.isEmpty(data.setupData.timeout)){
        notify("Timeout missed ---- < Stage No. " + index+" >","info");
        completeness = false;
    }

    if(completeness && isUsingGlobalVar(data.setupData.name)){
        completeness = isAvailableVar(data.setupData.name);
        if(!completeness){
            notify("Name is using an unknown global variable ---- < Stage No. " + index+" >","info");
        }
    }

    if(completeness){
        if(isUsingGlobalVar(data.setupData.timeout)){
            completeness = isAvailableVar(data.setupData.timeout);
            if(!completeness){
                notify("Timeout is using an unknown global variable ---- < Stage No. " + index+" >","info");
            }else{
                var realkey = data.setupData.timeout.substring(1,data.setupData.timeout.length-1);
                completeness = checkTimeout(getValue(realkey));
                if(!completeness){
                    notify("Timeout must be equal to or greater than 0---- < Stage No. " + index+" >","info");
                }
            }
        }else{
            completeness = checkTimeout(data.setupData.timeout);
            if(!completeness){
                notify("Timeout must be equal to or greater than 0---- < Stage No. " + index+" >","info");
            }
        } 
    }

    if(!completeness){
        return completeness;
    }
    
    for(var i=0;i<data.actions.length;i++){
        var item = data.actions[i];
        completeness = checkWorkflowAction(item,index,i);
        if(!completeness){
            break;
        }
    }

    return completeness;
}

function checkWorkflowAction(data,stageindex,actionindex){
    var completeness = true;
    if(_.isUndefined(data.outputJson) || _.isEmpty(data.outputJson)){
        notify("Output json missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }else if(_.isUndefined(data.inputJson) || _.isEmpty(data.inputJson)){
        notify("Input json missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }else if(_.isUndefined(data.setupData) || _.isEmpty(data.setupData)){
        notify("Base config missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }else{
        completeness = checkActionCompleteness(data,stageindex,actionindex);
    }

    if(completeness && !_.isEmpty(data.env)){
        completeness = checkActionEnv(data,stageindex,actionindex);
    }

    return completeness;
}

function checkActionCompleteness(data,stageindex,actionindex){
    var completeness = true;
    if(_.isEmpty(data.setupData.action.name)){
        notify("Name missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }
    // else if(_.isEmpty(data.setupData.action.timeout)){
    //     notify("Timeout missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
    //     completeness = false;
    // }
    else if(_.isEmpty(data.setupData.action.image.name)){
        notify("Image name missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }else if(_.isEmpty(data.setupData.action.image.tag)){
        notify("Image tag missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }else if(_.isEmpty(data.setupData.action.datafrom)){
        notify("Data From missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }else if(_.isEmpty(data.setupData.action.ip)){
        notify("Kubernetes IP missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }else if(_.isEmpty(data.setupData.action.apiserver)){
        notify("Kubernetes api server missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }else if(!data.setupData.action.useAdvanced){
        completeness = checkActionBaseSetting(data,stageindex,actionindex);
    }else if(data.setupData.action.useAdvanced){
        completeness = checkActionAdvancedSetting(data,stageindex,actionindex);
    }

    if(completeness && isUsingGlobalVar(data.setupData.action.name)){
        completeness = isAvailableVar(data.setupData.action.name);
        if(!completeness){
            notify("Name is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        }
    }

    if(completeness && !_.isEmpty(data.setupData.action.timeout)){
        if(isUsingGlobalVar(data.setupData.action.timeout)){
            completeness = isAvailableVar(data.setupData.action.timeout);
            if(!completeness){
                notify("Timeout is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
            }else{
                var realkey = data.setupData.action.timeout.substring(1,data.setupData.action.timeout.length-1);
                completeness = checkTimeout(getValue(realkey));
                if(!completeness){
                    notify("Timeout must be equal to or greater than 0---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
                }
            }
        }else{
            completeness = checkTimeout(data.setupData.action.timeout);
            if(!completeness){
                notify("Timeout must be equal to or greater than 0---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
            }
        }    
    }

    if(completeness && isUsingGlobalVar(data.setupData.action.image.name)){
        completeness = isAvailableVar(data.setupData.action.image.name);
        if(!completeness){
            notify("Repository name is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        }
    }

    if(completeness && isUsingGlobalVar(data.setupData.action.image.tag)){
        completeness = isAvailableVar(data.setupData.action.image.tag);
        if(!completeness){
            notify("Image tag is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        }
    }

    if(completeness && isUsingGlobalVar(data.setupData.action.datafrom)){
        completeness = isAvailableVar(data.setupData.action.datafrom);
        if(!completeness){
            notify("External data uri is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        }
    }

    if(completeness && isUsingGlobalVar(data.setupData.action.ip)){
        completeness = isAvailableVar(data.setupData.action.ip);
        if(!completeness){
            notify("Kubernetes IP is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        }
    }

    if(completeness && isUsingGlobalVar(data.setupData.action.apiserver)){
        completeness = isAvailableVar(data.setupData.action.apiserver);
        if(!completeness){
            notify("Kubernetes api server is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        }
    }

    return completeness;
}

function checkActionBaseSetting(data,stageindex,actionindex){
    var completeness = true;
    if(_.isEmpty(data.setupData.pod.spec.containers[0].resources.limits.cpu.toString())){
        notify("CPU limits missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }else if(_.isEmpty(data.setupData.pod.spec.containers[0].resources.limits.memory.toString())){
        notify("Memory limits missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }else if(_.isEmpty(data.setupData.pod.spec.containers[0].resources.requests.cpu.toString())){
        notify("CPU requests missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }else if(_.isEmpty(data.setupData.pod.spec.containers[0].resources.requests.memory.toString())){
        notify("Memory requests missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }else{
        var type = data.setupData.service.spec.type;
        var ports = data.setupData.service.spec.ports;
        // if(ports.length == 0){
        //     notify("No ports setting ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        //     completeness = false;
        // }
        for(var i=0;i<ports.length;i++){
            if(_.compact(_.values(ports[i])).length<3 && type == "NodePort"){
                notify("Ports or target ports or node ports missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
                completeness = false;
                break;
            }else if(_.compact(_.values(ports[i])).length<2 && type == "ClusterIP"){
                notify("Ports or target ports missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
                completeness = false;
                break;
            }

            if(completeness && isUsingGlobalVar(ports[i].port)){
                completeness = isAvailableVar(ports[i].port);
                if(!completeness){
                    notify("Port is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
                    break;
                }
            }

            if(completeness && isUsingGlobalVar(ports[i].targetPort)){
                completeness = isAvailableVar(ports[i].targetPort);
                if(!completeness){
                    notify("Target port is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
                    break;
                }
            }

            if(completeness && type == "NodePort" && isUsingGlobalVar(ports[i].nodePort)){
                completeness = isAvailableVar(ports[i].nodePort);
                if(!completeness){
                    notify("Node port is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
                    break;
                }
            }
        }
    }

    if(completeness && isUsingGlobalVar(data.setupData.pod.spec.containers[0].resources.limits.cpu)){
        completeness = isAvailableVar(data.setupData.pod.spec.containers[0].resources.limits.cpu);
        if(!completeness){
            notify("CPU limits is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        }
    }

    if(completeness && isUsingGlobalVar(data.setupData.pod.spec.containers[0].resources.limits.memory)){
        completeness = isAvailableVar(data.setupData.pod.spec.containers[0].resources.limits.memory);
        if(!completeness){
            notify("Memory limits is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        }
    }

    if(completeness && isUsingGlobalVar(data.setupData.pod.spec.containers[0].resources.requests.cpu)){
        completeness = isAvailableVar(data.setupData.pod.spec.containers[0].resources.requests.cpu);
        if(!completeness){
            notify("CPU requests is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        }
    }

    if(completeness && isUsingGlobalVar(data.setupData.pod.spec.containers[0].resources.requests.memory)){
        completeness = isAvailableVar(data.setupData.pod.spec.containers[0].resources.requests.memory);
        if(!completeness){
            notify("Memory requests is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        }
    }

    return completeness;
}

function checkActionAdvancedSetting(data,stageindex,actionindex){
    var completeness = true;
    if(_.isUndefined(data.setupData.service_advanced) || _.isEmpty(data.setupData.service_advanced)){
        notify("Service advanced setting missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }else if(_.isUndefined(data.setupData.pod_advanced) || _.isEmpty(data.setupData.pod_advanced)){
        notify("Pod advanced setting missed ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
        completeness = false;
    }

    return completeness;
}

export function isUsingGlobalVar(value){
    return !_.isUndefined(value) && value.toString().indexOf("@") == 0 && value.toString().lastIndexOf("@") == value.toString().length-1;
}

function checkActionEnv(data,stageindex,actionindex){
    var completeness = true;
    for(var i=0;i<data.env.length;i++){
        var env = data.env[i];
        if(isUsingGlobalVar(env.key)){
            completeness = isAvailableVar(env.key);
            if(!completeness){
                notify("Env key '" + env.key + "' is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
                break;
            }

            var realkey = getValue(env.key.substring(1,env.key.length-1));
            completeness = isEnvKeyLegal(realkey);
            if(!completeness){
                notify("Env key '" + env.key + "' is illegal. Key is not allowed to start with 'CO_' ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
                break;
            }

        }else{
            completeness = isEnvKeyLegal(env.key);
            if(!completeness){
                notify("Env key '" + env.key + "' is illegal. Key is not allowed to start with 'CO_' ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
                break;
            }
        }

        if(completeness && isUsingGlobalVar(env.value)){
            completeness = isAvailableVar(env.value);
            if(!completeness){
                notify("Env value of key '" + env.key + "' is using an unknown global variable ---- < Stage No. " + stageindex + " / Action No. " + (actionindex+1)+" >","info");
                break;
            }
        }
    }
    return completeness;
}

export function isEnvsLegal(envs){
    var illegalOnes = _.filter(envs,function(env){
        return !/^(.?$|[^C].+|C[^O].+|CO[^_].*)/.test(env.key);
    });

    if(illegalOnes.length>0){
        return false;
    }else{
        return true;
    }
}

export function isEnvKeyLegal(key){
    return /^(.?$|[^C].+|C[^O].+|CO[^_].*)/.test(key);
}

function checkTimeout(value){
    var _value = Number(value);
    return !_.isNaN(_value) && _value >= 0;
}

function checkSetting(setting){
    var completeness = true;

    if(!_.isUndefined(setting.data)){
        if(setting.data.runningInstances.available){
            completeness = !_.isEmpty(setting.data.runningInstances.number.toString()) && checkTimeout(setting.data.runningInstances.number);
            if(!completeness){
                notify("Number of running instances must be equal to or greater than 0 ---- < Workflow setting >","info");
                return completeness;
            }
        }
        
        if(setting.data.timedTasks.available){
            for(var i=0; i<setting.data.timedTasks.tasks.length;i++){
                var cron = setting.data.timedTasks.tasks[i].cronEntry;
                completeness = checkCronTask(cron);
                if(!completeness){
                    notify("Cron entry of timed task No. " + (i+1) + " is illegal ---- < Workflow setting >","info");
                    break;
                }

                if(_.isEmpty(setting.data.timedTasks.tasks[i].eventType)){
                    notify("Event type of timed task No. " + (i+1) + " is required ---- < Workflow setting >","info");
                    completeness = false;
                    break;
                }else if(_.isEmpty(setting.data.timedTasks.tasks[i].eventName)){
                    notify("Event name of timed task No. " + (i+1) + " is required ---- < Workflow setting >","info");
                    completeness = false;
                    break;
                }else if(_.isEmpty(setting.data.timedTasks.tasks[i].startJson)){
                    notify("Start json of timed task No. " + (i+1) + " is required ---- < Workflow setting >","info");
                    completeness = false;
                    break;
                }
            }
        }
    }

    return completeness;
}

function checkCronTask(cron){
    if(!_.isEmpty(cron) && cron.split(" ").length >= 5){
        return true;
    }else{
        return false;
    }
}
