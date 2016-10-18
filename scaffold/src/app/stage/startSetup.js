import * as startSetupData from "./startSetupData";
import {initStartIO,initTreeEdit,initFromEdit,initFromView,getOutputForEvent} from "./startIO";

export function initStartSetup(start){
    startSetupData.getStartSetupData(start);
    initStartIO(start);

    // type select
    $("#type-select").val(startSetupData.getTypeSelect());
    selectType(startSetupData.getTypeSelect());

    $("#type-select").on("change",function(){
        startSetupData.setTypeSelect();
        selectType(startSetupData.getTypeSelect());
    });

    $("#type-select").select2({
        minimumResultsForSearch: Infinity
    });

    // event select
    $("#event-select").on("change",function(){
        startSetupData.setEventSelect();
        getOutputForEvent(startSetupData.getEventSelect());
    });
}

function selectType(pipelineType){
    if(pipelineType == "github" || pipelineType == "gitlab"){
        $("#event_select").show();
        $("#outputTreeViewer").show();
        $("#outputTreeDesigner").hide();
        
        $("#event-select").val(startSetupData.getEventSelect());
        $("#event-select").select2({
            minimumResultsForSearch: Infinity
        });
        getOutputForEvent(startSetupData.getEventSelect()); 
    }else{
        $("#event_select").hide();
        $("#outputTreeViewer").hide();
        $("#outputTreeDesigner").show();

        initTreeEdit();
        initFromEdit("output");
    }
}