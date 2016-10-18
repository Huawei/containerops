export let data;

export function getStartSetupData(start){
    if(!_.isUndefined(start.setupData) && !_.isEmpty(start.setupData)){
      data = start.setupData;
    }else{
      data = $.extend(true,{},metadata);
      start.setupData = data;
    } 
}

export function setTypeSelect(){
    data.type = $("#type-select").val();
}

export function getTypeSelect(){
    return data.type;
}

export function setEventSelect(){
    data.event = $("#event-select").val();
}

export function getEventSelect(){
    return data.event;
}

var metadata = {
  "type" : "github",
  "url" : "",
  "token" : "",
  "event" : "pull_request"
}



