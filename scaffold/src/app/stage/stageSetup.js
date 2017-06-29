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

import * as stageSetupData from "./stageSetupData";
import {workflowVars} from "../workflow/workflowVar";

export function initStageSetup(stage){
    stageSetupData.getStageSetupData(stage);

    $("#stage-name").val(stageSetupData.data.name);
    $("#stage-name").on("blur",function(){
        stageSetupData.setStageName();
    });

    $("#stage-timeout").val(stageSetupData.data.timeout);
    $("#stage-timeout").on("blur",function(){
        stageSetupData.setStageTimeout();
    });

    // use global vars
    var globalvars = _.map(workflowVars,function(item){
        return "@"+item[0]+"@";
    });
    $(".allowFromVar").autocomplete({
        source:[globalvars],
        limit: 100,
        visibleLimit: 5
    }); 
}