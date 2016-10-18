import {jsonEditor} from "../../vendor/jquery.jsoneditor";
import {notify} from "../common/notify";

var treeEdit_InputContainer,treeEdit_OutputContainer;
var fromEdit_InputCodeContainer,fromEdit_InputTreeContainer,fromEdit_OutputCodeContainer,fromEdit_OutputTreeContainer;
var fromEdit_OutputViewContainer;
var fromEdit_CodeEditor,fromEdit_TreeEditor;

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
    fromEdit_OutputViewContainer = $("#outputTreeViewer")[0];

    // input output from edit
    $("#tree-edit-tab").on('click',function(){
        initTreeEdit();
    })

    $("#input-from-edit-tab").on('click',function(){
        initFromEdit("input");
    })

    $("#output-from-edit-tab").on('click',function(){
        initFromEdit("output");
    });

    initTreeEdit();

    // $("#saveActionIO").on('click',function(){
    //     saveActionIOData(action,actionIOData);
    // })
}

export function initTreeEdit(){
    try{
        jsonEditor(treeEdit_InputContainer,actionIOData.inputJson, {
            change:function(data){
                actionIOData.inputJson = data;
            }
        });
    }catch(e){
        notify("Input Error in parsing json.","error");
    }

    try{
        jsonEditor(treeEdit_OutputContainer,actionIOData.outputJson, {
            change:function(data){
                actionIOData.outputJson = data;
            }
        });
    }catch(e){
        notify("Output Error in parsing json.","error");
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

    if(type == "input"){
        fromEdit_CodeEditor = new JSONEditor(fromEdit_InputCodeContainer, codeOptions);
        fromEdit_TreeEditor = new JSONEditor(fromEdit_InputTreeContainer, treeOptions);
        fromEdit_CodeEditor.set(actionIOData.inputJson);
        fromEdit_TreeEditor.set(actionIOData.inputJson);
        $("#inputFromJson").on('click',function(){
            fromCodeToTree("input");
        })  
        $("#inputToJson").on('click',function(){
            fromTreeToCode("input");
        })       
    }else if(type == "output"){
        fromEdit_CodeEditor = new JSONEditor(fromEdit_OutputCodeContainer, codeOptions);
        fromEdit_TreeEditor = new JSONEditor(fromEdit_OutputTreeContainer, treeOptions);
        fromEdit_CodeEditor.set(actionIOData.outputJson);
        fromEdit_TreeEditor.set(actionIOData.outputJson);
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
    if(type == "input"){
        try{
            actionIOData.inputJson = fromEdit_CodeEditor.get();
            fromEdit_TreeEditor.set(actionIOData.inputJson);
        }catch(e){
            notify("Input Code Changes Error in parsing json.","error");
        }  
    }else if(type == "output"){
        try{
            actionIOData.outputJson = fromEdit_CodeEditor.get();
            fromEdit_TreeEditor.set(actionIOData.outputJson);
        }catch(e){
            notify("Output Code Changes Error in parsing json.","error");
        } 
    }
    
    fromEdit_TreeEditor.expandAll();
}

function fromTreeToCode(type){
    if(type == "input"){
        try{
            actionIOData.inputJson = fromEdit_TreeEditor.get();
            fromEdit_CodeEditor.set(actionIOData.inputJson);
        }catch(e){
            notify("Input Tree Changes Error in parsing json.","error");
        }  
    }else if(type == "output"){
        try{
            actionIOData.outputJson = fromEdit_TreeEditor.get();
            fromEdit_CodeEditor.set(actionIOData.outputJson);
        }catch(e){
            notify("Output Tree Changes Error in parsing json.","error");
        } 
    }
}

export function initFromView(){
    if(fromEdit_TreeEditor){
        fromEdit_TreeEditor.destroy();
    }

    var treeOptions = {
        "mode": "view",
        "search": true
    };

    fromEdit_TreeEditor = new JSONEditor(fromEdit_OutputViewContainer, treeOptions);
    fromEdit_TreeEditor.set(actionIOData.outputJson);
    
    fromEdit_TreeEditor.expandAll();
}