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
import { initWorkflow } from "../workflow/initWorkflow";
import { initAction } from "../workflow/initAction";
import { workflowData } from "../workflow/main";
import { resizeWidget } from "../theme/widget";
import { initStageSetup } from "./stageSetup";

export function clickStage(sd, si) {
    //show stage form
    $.ajax({
        url: "../../templates/stage/stageEdit.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#workflow-info-edit").html($(data));

            initStageSetup(sd);

            $("#uuid").attr("value", sd.id);

            resizeWidget();
        }
    });


}
