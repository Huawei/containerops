/* 
Copyright 2014 Huawei Technologies Co., Ltd. All rights reserved.

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
 
import {jsonEditor} from "../../vendor/jquery.jsoneditor";
import {notify} from "../common/notify";
import {pipelineApi} from "../common/api";
import {loading} from "../common/loading";

var treeEdit_OutputContainer;
var fromEdit_OutputCodeContainer,fromEdit_OutputTreeContainer;
var fromEdit_OutputViewContainer;
var fromEdit_CodeEditor,fromEdit_TreeEditor;

let startIOData;
export function initStartIO(start){
    startIOData = start;

    if(startIOData.outputJson == undefined){
        startIOData.outputJson = {};
    }

    treeEdit_OutputContainer = $('#outputTreeDiv');
    fromEdit_OutputCodeContainer = $("#outputCodeEditor")[0];
    fromEdit_OutputTreeContainer = $("#outputTreeEditor")[0];
    fromEdit_OutputViewContainer = $("#outputTreeViewer")[0];

    initTreeEdit();
    initFromEdit("output");
}

export function initTreeEdit(){
    if(_.isUndefined(startIOData.outputJson) || _.isEmpty(startIOData.outputJson)){
        $("#outputTreeStart").show();
        $("#outputTreeDiv").hide();
        $("#outputStartBtn").on('click',function(){
            startIOData.outputJson = {
                "newKey" : null
            }
            initTreeEdit();
            initFromEdit("output");
        })
    }else{
        try{
            $("#outputTreeStart").hide();
            $("#outputTreeDiv").show();
            jsonEditor(treeEdit_OutputContainer,startIOData.outputJson, {
                change:function(data){
                    startIOData.outputJson = data;
                    initFromEdit("output");
                }
            },"start");
        }catch(e){
            notify("Output Error in parsing json.","error");
        }
    }
}

export function initFromEdit(type){
    if(fromEdit_CodeEditor){
        fromEdit_CodeEditor.destroy();
    }

    if(fromEdit_TreeEditor){
        fromEdit_TreeEditor.destroy();
    }

    var codeOptions = {
        "mode": "code",
        "indentation": 2
    };

    var treeOptions = {
        "mode": "tree",
        "search": true
    };

    if(type == "output"){
        fromEdit_CodeEditor = new JSONEditor(fromEdit_OutputCodeContainer, codeOptions);
        fromEdit_TreeEditor = new JSONEditor(fromEdit_OutputTreeContainer, treeOptions);
        fromEdit_CodeEditor.set(startIOData.outputJson);
        fromEdit_TreeEditor.set(startIOData.outputJson);
        $("#outputFromJson").on('click',function(){
            fromCodeToTree("output");
        })
        $("#outputToJson").on('click',function(){
            fromTreeToCode("output");
        })
    }
    
    fromEdit_TreeEditor.expandAll();
}

function fromCodeToTree(type){
    if(type == "output"){
        try{
            startIOData.outputJson = fromEdit_CodeEditor.get();
            fromEdit_TreeEditor.set(startIOData.outputJson);
            initTreeEdit();
        }catch(e){
            notify("Output Code Changes Error in parsing json.","error");
        } 
    }
    
    fromEdit_TreeEditor.expandAll();
}

function fromTreeToCode(type){
    if(type == "output"){
        try{
            startIOData.outputJson = fromEdit_TreeEditor.get();
            fromEdit_CodeEditor.set(startIOData.outputJson);
            initTreeEdit();
        }catch(e){
            notify("Output Tree Changes Error in parsing json.","error");
        } 
    }
}

function initFromView(){
    if(fromEdit_TreeEditor){
        fromEdit_TreeEditor.destroy();
    }

    var treeOptions = {
        "mode": "view",
        "search": true
    };

    fromEdit_TreeEditor = new JSONEditor(fromEdit_OutputViewContainer, treeOptions);
    fromEdit_TreeEditor.set(startIOData.outputJson);
    
    fromEdit_TreeEditor.expandAll();
}

export function getOutputForEvent(selecetedEvent){
    loading.show();
    var promise = pipelineApi.eventOutput(selecetedEvent);
    promise.done(function(data){
        loading.hide();
        startIOData.outputJson = data.output;
        initFromView();
    });
    promise.fail(function(xhr,status,error){
        loading.hide();
        if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
            notify(xhr.responseJSON.errMsg,"error");
        }else{
            notify("Server is unreachable","error");
        }
    });
}