import {resizeWidget} from "../theme/widget";
import {initStartSetup} from "./startSetup";

let pipelineType,selectedEvent;
export function clickStart(sd, si) {
    //show git form
    $.ajax({
        url: "../../templates/stage/startEdit.html",
        type: "GET",
        cache: false,
        success: function (data) {
            $("#pipeline-info-edit").html($(data));

            initStartSetup(sd);
            
            // $("#uuid").attr("value", sd.id);

            resizeWidget();
        }
    });
}

