import * as constant from "../common/constant";
import { initPipeline } from "../pipeline/initPipeline";
import { initAction } from "../pipeline/initAction";
import { pipelineData } from "../pipeline/main";
import { resizeWidget } from "../theme/widget";
import { removeLinkArray } from "../relation/removeLinkArray";
import { initStageSetup } from "./stageSetup";

export function clickStage(sd, si) {
    //show stage form
    $.ajax({
        url: "../../templates/stage/stageEdit.html",
        type: "GET",
        cache: false,
        success: function(data) {
            $("#pipeline-info-edit").html($(data));

            initStageSetup(sd);

            $("#uuid").attr("value", sd.id);

            resizeWidget();
        }
    });


}
