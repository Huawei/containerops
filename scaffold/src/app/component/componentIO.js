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
var fromEdit_OutputViewContainer;
var fromEdit_InputCodeEditor,fromEdit_InputTreeEditor,fromEdit_OutputCodeEditor,fromEdit_OutputTreeEditor;

let componentIOData;
export function initComponentIO(component){
    componentIOData = component;
    if(componentIOData.inputJson == undefined){
        componentIOData.inputJson = {};
    }
    if(componentIOData.outputJson == undefined){
        componentIOData.outputJson = {};
    }
    treeEdit_InputContainer = $('#inputTreeDiv');
    treeEdit_OutputContainer = $('#outputTreeDiv'); 
    fromEdit_InputCodeContainer = $("#inputCodeEditor")[0];
    fromEdit_InputTreeContainer = $("#inputTreeEditor")[0];
    fromEdit_OutputCodeContainer = $("#outputCodeEditor")[0];
    fromEdit_OutputTreeContainer = $("#outputTreeEditor")[0];
    fromEdit_OutputViewContainer = $("#outputTreeViewer")[0];

    $("#jsonMode").hide();

    // input output from edit
    $("#design-tab").on('click',function(){
        initTreeEdit();
        initFromEdit("input");
        initFromEdit("output");
    });

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
    if(_.isUndefined(componentIOData.inputJson) || _.isEmpty(componentIOData.inputJson)){
        $("#inputTreeStart").show();
        $("#inputTreeDiv").hide();
        $("#inputStartBtn").on('click',function(){
            componentIOData.inputJson = {
                "newKey" : null
            }
            initTreeEdit();
        })
    }else{
        try{
            $("#inputTreeStart").hide();
            $("#inputTreeDiv").show();
            jsonEditor(treeEdit_InputContainer,componentIOData.inputJson, {
                change:function(data){
                    componentIOData.inputJson = data;
                }
            },"component");
        }catch(e){
            notify("Input Error in parsing json.","error");
        }
    }
    
    if(_.isUndefined(componentIOData.outputJson) || _.isEmpty(componentIOData.outputJson)){
        $("#outputTreeStart").show();
        $("#outputTreeDiv").hide();
        $("#outputStartBtn").on('click',function(){
            componentIOData.outputJson = {
                "newKey" : null
            }
            initTreeEdit();
        })
    }else{
        try{
            $("#outputTreeStart").hide();
            $("#outputTreeDiv").show();
            jsonEditor(treeEdit_OutputContainer,componentIOData.outputJson, {
                change:function(data){
                    componentIOData.outputJson = data;
                }
            },"component");
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
        fromEdit_InputCodeEditor.set(componentIOData.inputJson);
        fromEdit_InputTreeEditor.set(componentIOData.inputJson);
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
        fromEdit_OutputCodeEditor.set(componentIOData.outputJson);
        fromEdit_OutputTreeEditor.set(componentIOData.outputJson);
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
            componentIOData.inputJson = fromEdit_InputCodeEditor.get();
            fromEdit_InputTreeEditor.set(componentIOData.inputJson);
        }catch(e){
            notify("Input Code Changes Error in parsing json.","error");
        }  
        fromEdit_InputTreeEditor.expandAll();
    }else if(type == "output"){
        try{
            componentIOData.outputJson = fromEdit_OutputCodeEditor.get();
            fromEdit_OutputTreeEditor.set(componentIOData.outputJson);
        }catch(e){
            notify("Output Code Changes Error in parsing json.","error");
        } 
        fromEdit_OutputTreeEditor.expandAll();
    }
}

function fromTreeToCode(type){
    if(type == "input"){
        try{
            componentIOData.inputJson = fromEdit_InputTreeEditor.get();
            fromEdit_InputCodeEditor.set(componentIOData.inputJson);
        }catch(e){
            notify("Input Tree Changes Error in parsing json.","error");
        }  
    }else if(type == "output"){
        try{
            componentIOData.outputJson = fromEdit_OutputTreeEditor.get();
            fromEdit_OutputCodeEditor.set(componentIOData.outputJson);
        }catch(e){
            notify("Output Tree Changes Error in parsing json.","error");
        } 
    }
}