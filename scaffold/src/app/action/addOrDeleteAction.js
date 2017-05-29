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

import * as constant from "../common/constant";
import { workflowData } from "../workflow/main";
import * as util from "../common/util";

export function addAction(actions) {
    actions.splice(
        actions.length,
        0, {
            id: constant.WORKFLOW_ACTION + "-" + uuid.v1(),
            type: constant.WORKFLOW_ACTION,
            setupData: {}
        });

}

export function deleteAction(data, index) {
    _.each(workflowData, function(stage){
        if(stage.type == constant.WORKFLOW_STAGE && stage.actions && stage.actions.length > 0){
            _.each(stage.actions, function(action){
                if(action.id == data.id){
                    stage.actions = _.without(stage.actions, action);
                }
            })
        }
    })
    util.removeRelatedLines(data.id);
}
