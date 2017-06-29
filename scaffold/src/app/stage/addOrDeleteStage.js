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

import { workflowData } from "../workflow/main";
import * as constant from "../common/constant";
import * as util from "../common/util";

export function addStage(data, index) {
    workflowData.splice(
        workflowData.length - 2,
        0, {
            id: constant.WORKFLOW_STAGE + "-" + uuid.v1(),
            type: constant.WORKFLOW_STAGE,
            class: constant.WORKFLOW_STAGE,
            drawX: 0,
            drawY: 0,
            width: 0,
            height: 0,
            translateX: 0,
            translateY: 0,
            actions: [],
            setupData: {}
        });


}

export function deleteStage(data, index){
     var relatedActions = util.findAllActionsOfStage(data.id)
         util.removeRelatedLines(relatedActions);
         workflowData.splice(index, 1);
}
