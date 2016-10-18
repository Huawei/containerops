import { pipelineData } from "../pipeline/main";
import * as constant from "../common/constant";
import * as util from "../common/util";

export function addStage(data, index) {
    pipelineData.splice(
        pipelineData.length - 2,
        0, {
            id: constant.PIPELINE_STAGE + "-" + uuid.v1(),
            type: constant.PIPELINE_STAGE,
            class: constant.PIPELINE_STAGE,
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
         pipelineData.splice(index, 1);
}
