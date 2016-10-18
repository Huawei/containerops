import * as stageSetupData from "./stageSetupData";

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

    // $("#stage-env").val(stageSetupData.data.env);
    // $("#stage-env").on("blur",function(){
    //     stageSetupData.setStageEnv();
    // });

    // $("#stage-callback-url").val(stageSetupData.data.callbackurl);
    // $("#stage-callback-url").on("blur",function(){
    //     stageSetupData.setStageCallbackUrl();
    // });
}