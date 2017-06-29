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

export let data;

export function getStageSetupData(stage){
    if(!_.isUndefined(stage.setupData) && !_.isEmpty(stage.setupData)){
      data = stage.setupData;
    }else{
      data = $.extend(true,{},metadata);
      stage.setupData = data;
    } 
}

export function setStageName(){
    data.name = $("#stage-name").val();
}

export function setStageTimeout(){
    data.timeout = $("#stage-timeout").val();
}

var metadata = {
  "name" : "",
  "timeout" : ""
}



