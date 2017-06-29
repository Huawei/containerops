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

export let data;

let selectedTab;

export function getStartIOData(start){
    if(start.outputJson == undefined){
        start.outputJson = [];
        start.outputJson.push($.extend(true,{},metadata));
    }

    data = start.outputJson;
}

// selectedTab
export function setSelectedTab(index){
    selectedTab = index;
}

// type select
export function findTypeSelectDom(){
    return $(".output-type-event-div[data-index="+ selectedTab +"]").find(".output-type-select");
}

export function setTypeSelect(){
    data[selectedTab].type = findTypeSelectDom().val();
}

export function setTypeSelectDom(){
    findTypeSelectDom().val(data[selectedTab].type);
    findTypeSelectDom().select2({
        minimumResultsForSearch: Infinity
    });
}

export function getTypeSelect(){
    return data[selectedTab].type;
}

// event select
export function setEvent(e){
    data[selectedTab].event = e;
}

export function findGitHubEventSelectDom(){
    return $(".output-type-event-div[data-index="+ selectedTab +"]").find(".github-event-select");
}

export function findGitLabEventSelectDom(){
    return $(".output-type-event-div[data-index="+ selectedTab +"]").find(".gitlab-event-select");
}

export function findEventInputDom(){
    return $(".output-type-event-div[data-index="+ selectedTab +"]").find(".output-event-input");
}

export function setEventSelect(type){
    if(type == "github"){
        data[selectedTab].event = findGitHubEventSelectDom().val();
    }else if(type == "gitlab"){
        data[selectedTab].event = findGitLabEventSelectDom().val();
    }  
}

export function setEventInput(){
    data[selectedTab].event = findEventInputDom().val();
}

export function setEventSelectDom(type){
    if(type == "github"){
        findGitHubEventSelectDom().val(data[selectedTab].event);
        findGitHubEventSelectDom().select2({
            minimumResultsForSearch: Infinity
        });
    }else if(type == "gitlab"){
        findGitLabEventSelectDom().val(data[selectedTab].event);
        findGitLabEventSelectDom().select2({
            minimumResultsForSearch: Infinity
        });
    }
}

export function setEventInputDom(){
    findEventInputDom().val(data[selectedTab].event);
}

export function getEventSelect(){
    return data[selectedTab].event;
}

export function findEventSelectDivDom(){
    return $(".output-type-event-div[data-index="+ selectedTab +"]").find(".event-select-div");
}

export function findEventInputDivDom(){
    return $(".output-type-event-div[data-index="+ selectedTab +"]").find(".event-input-div");
}

// viewer
export function findOutputTreeViewerDom(){
    return $(".output-json-div[data-index="+ selectedTab +"]").find(".startOutputTreeViewer");
}

// designer
export function findOutputTreeDesignerDom(){
    return $(".output-json-div[data-index="+ selectedTab +"]").find(".startOutputTreeDesigner");
}

export function findOutputTreeStartDom(){
  return $(".output-json-div[data-index="+ selectedTab +"]").find(".outputTreeStart");
}

export function findOutputTreeDivDom(){
  return $(".output-json-div[data-index="+ selectedTab +"]").find(".outputTreeDiv");
}

export function findOutputCodeEditorDom(){
  return $(".output-json-div[data-index="+ selectedTab +"]").find(".outputCodeEditor");
}

export function findOutputTreeEditorDom(){
  return $(".output-json-div[data-index="+ selectedTab +"]").find(".outputTreeEditor");
}

// json
export function setJson(d){
    data[selectedTab].json = d;
}

export function getJson(){
    return data[selectedTab].json;
}

// new output
export function addOutput(){
    data.push($.extend(true,{},metadata));
}

// tab
export function findSelectedStartOutputTabDom(){
    return $(".start-output-tab[data-index="+ selectedTab +"]");
}

export function findSelectedStartOutputTabContentDom(){
    return $(".output-type-event-div[data-index="+ selectedTab +"]").parent();
}

// check event options
export function isEventOptionAvailable(){
    var numbers = _.filter(data,function(item){
        return item.type == data[selectedTab].type && item.event == data[selectedTab].event;
    }).length;

    if(numbers == 1){
        return true;
    }

    return false;
}

export function deleteOutput(){
    data.splice(selectedTab,1);
}

var metadata = {
  "type" : "github",
  "event" : "PullRequest",
  "json" : {}
}



