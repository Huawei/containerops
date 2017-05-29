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
 
import {jsonEditor} from "../../vendor/jquery.jsoneditor";
import {notify} from "../common/notify";

var treeEdit_InputContainer,treeEdit_OutputContainer;
var fromEdit_InputCodeContainer,fromEdit_InputTreeContainer,fromEdit_OutputCodeContainer,fromEdit_OutputTreeContainer;
var fromEdit_InputCodeEditor,fromEdit_InputTreeEditor,fromEdit_OutputCodeEditor,fromEdit_OutputTreeEditor;

let actionIOData;
export function initActionIO(action){
    actionIOData = action;
    if(actionIOData.inputJson == undefined){
        actionIOData.inputJson = {};
    }
    if(actionIOData.outputJson == undefined){
        actionIOData.outputJson = {};
    }
    treeEdit_InputContainer = $('#inputTreeDiv');
    treeEdit_OutputContainer = $('#outputTreeDiv');
    fromEdit_InputCodeContainer = $("#inputCodeEditor")[0];
    fromEdit_InputTreeContainer = $("#inputTreeEditor")[0];
    fromEdit_OutputCodeContainer = $("#outputCodeEditor")[0];
    fromEdit_OutputTreeContainer = $("#outputTreeEditor")[0];

    $("#jsonMode").hide();

    // input output from edit
    $("#design-tab").on('click',function(){
        initTreeEdit();
        initFromEdit("input");
        initFromEdit("output");
    })

    // action design way
    $(".action-json").on('click',function(){
        $("#designMode").hide();
        $("#jsonMode").show();
        initFromEdit("input");
        initFromEdit("output");
    });

    $(".action-design").on('click',function(){
        $("#designMode").show();
        $("#jsonMode").hide();
        initTreeEdit();
    });
}

export function initTreeEdit(){
    if(_.isUndefined(actionIOData.inputJson) || _.isEmpty(actionIOData.inputJson)){
        $("#inputTreeStart").show();
        $("#inputTreeDiv").hide();
        $("#inputStartBtn").on('click',function(){
            actionIOData.inputJson = {
                "newKey" : null
            }
            initTreeEdit();
        })
    }else{
        try{
            $("#inputTreeStart").hide();
            $("#inputTreeDiv").show();
            jsonEditor(treeEdit_InputContainer,actionIOData.inputJson, {
                change:function(data){
                    actionIOData.inputJson = data;
                }
            },"action");
        }catch(e){
            notify("Input Error in parsing json.","error");
        }
    }
    
    if(_.isUndefined(actionIOData.outputJson) || _.isEmpty(actionIOData.outputJson)){
        $("#outputTreeStart").show();
        $("#outputTreeDiv").hide();
        $("#outputStartBtn").on('click',function(){
            actionIOData.outputJson = {
                "newKey" : null
            }
            initTreeEdit();
        })
    }else{
        try{
            $("#outputTreeStart").hide();
            $("#outputTreeDiv").show();
            jsonEditor(treeEdit_OutputContainer,actionIOData.outputJson, {
                change:function(data){
                    actionIOData.outputJson = data;
                }
            },"action");
        }catch(e){
            notify("Output Error in parsing json.","error");
        }
    }
}

function initFromEdit(type){
    var codeOptions = {
        "mode": "code",
        "indentation": 2
    };

    var treeOptions = {
        "mode": "tree",
        "search": true
    };

    if(type == "input"){
        if(fromEdit_InputCodeEditor){
            fromEdit_InputCodeEditor.destroy();
        }
        if(fromEdit_InputTreeEditor){
            fromEdit_InputTreeEditor.destroy();
        }
        fromEdit_InputCodeEditor = new JSONEditor(fromEdit_InputCodeContainer, codeOptions);
        fromEdit_InputTreeEditor = new JSONEditor(fromEdit_InputTreeContainer, treeOptions);
        fromEdit_InputCodeEditor.set(actionIOData.inputJson);
        fromEdit_InputTreeEditor.set(actionIOData.inputJson);
        $("#inputFromJson").on('click',function(){
            fromCodeToTree("input");
        })  
        $("#inputToJson").on('click',function(){
            fromTreeToCode("input");
        })       

        fromEdit_InputTreeEditor.expandAll();
    }else if(type == "output"){
        if(fromEdit_OutputCodeEditor){
            fromEdit_OutputCodeEditor.destroy();
        }
        if(fromEdit_OutputTreeEditor){
            fromEdit_OutputTreeEditor.destroy();
        }
        fromEdit_OutputCodeEditor = new JSONEditor(fromEdit_OutputCodeContainer, codeOptions);
        fromEdit_OutputTreeEditor = new JSONEditor(fromEdit_OutputTreeContainer, treeOptions);
        fromEdit_OutputCodeEditor.set(actionIOData.outputJson);
        fromEdit_OutputTreeEditor.set(actionIOData.outputJson);
        $("#outputFromJson").on('click',function(){
            fromCodeToTree("output");
        })
        $("#outputToJson").on('click',function(){
            fromTreeToCode("output");
        })

        fromEdit_OutputTreeEditor.expandAll();
    }
}

function fromCodeToTree(type){
    if(type == "input"){
        try{
            actionIOData.inputJson = fromEdit_InputCodeEditor.get();
            fromEdit_InputTreeEditor.set(actionIOData.inputJson);
        }catch(e){
            notify("Input Code Changes Error in parsing json.","error");
        }  
        fromEdit_InputTreeEditor.expandAll();
    }else if(type == "output"){
        try{
            actionIOData.outputJson = fromEdit_OutputCodeEditor.get();
            fromEdit_OutputTreeEditor.set(actionIOData.outputJson);
        }catch(e){
            notify("Output Code Changes Error in parsing json.","error");
        } 
        fromEdit_OutputTreeEditor.expandAll();
    }
}

function fromTreeToCode(type){
    if(type == "input"){
        try{
            actionIOData.inputJson = fromEdit_InputTreeEditor.get();
            fromEdit_InputCodeEditor.set(actionIOData.inputJson);
        }catch(e){
            notify("Input Tree Changes Error in parsing json.","error");
        }  
    }else if(type == "output"){
        try{
            actionIOData.outputJson = fromEdit_OutputTreeEditor.get();
            fromEdit_OutputCodeEditor.set(actionIOData.outputJson);
        }catch(e){
            notify("Output Tree Changes Error in parsing json.","error");
        } 
    }
}