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

export let data,envs;

export function getComponentEnvData(component){
    if(!_.isUndefined(component.env) && !_.isEmpty(component.env)){
      data = component.env;
    }else{
      data = {};
      component.env = data;
    } 
    envs = _.pairs(data);
}

export function getComponentEnvPairs(){
    return envs;
}

export function setEnvKey(event){
    var index = $(event.currentTarget).parent().parent().data("index");
    envs[index][0] = $(event.currentTarget).val();
}

export function setEnvValue(event){
    var index = $(event.currentTarget).parent().parent().data("index");
    envs[index][1] = $(event.currentTarget).val();
}

export function removeEnv(event){
    var index = $(event.currentTarget).parent().parent().data("index");
    envs.splice(index,1);
}

export function addEnv(){
    envs.push(["",""]);
}

export function getWholeEnvs(){
    return _.object(envs);
}