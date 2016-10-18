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

// export function setStageEnv(){
//     data.env = $("#stage-env").val();
// }

// export function setStageCallbackUrl(){
//     data.callbackurl = $("#stage-callback-url").val();
// }


var metadata = {
  "name" : "",
  "timeout" : ""
}



