import {jsonEditor} from "../../vendor/jquery.jsoneditor";
import {notify} from "../common/notify";

var treeEdit_InputContainer,treeEdit_OutputContainer;
var fromEdit_InputCodeContainer,fromEdit_InputTreeContainer,fromEdit_OutputCodeContainer,fromEdit_OutputTreeContainer;
var fromEdit_OutputViewContainer;
var fromEdit_CodeEditor,fromEdit_TreeEditor;

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
}

function initTreeEdit(){
    try{
        jsonEditor(treeEdit_InputContainer,componentIOData.inputJson, {
            change:function(data){
                componentIOData.inputJson = data;
            }
        });
    }catch(e){
        notify("Input Error in parsing json.","error");
    }

    try{
        jsonEditor(treeEdit_OutputContainer,componentIOData.outputJson, {
            change:function(data){
                componentIOData.outputJson = data;
            }
        });
    }catch(e){
        notify("Output Error in parsing json.","error");
    }
}

function initFromEdit(type){
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
        fromEdit_CodeEditor.set(componentIOData.inputJson);
        fromEdit_TreeEditor.set(componentIOData.inputJson);
        $("#inputFromJson").on('click',function(){
            fromCodeToTree("input");
        })  
        $("#inputToJson").on('click',function(){
            fromTreeToCode("input");
        })       
    }else if(type == "output"){
        fromEdit_CodeEditor = new JSONEditor(fromEdit_OutputCodeContainer, codeOptions);
        fromEdit_TreeEditor = new JSONEditor(fromEdit_OutputTreeContainer, treeOptions);
        fromEdit_CodeEditor.set(componentIOData.outputJson);
        fromEdit_TreeEditor.set(componentIOData.outputJson);
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
            componentIOData.inputJson = fromEdit_CodeEditor.get();
            fromEdit_TreeEditor.set(componentIOData.inputJson);
        }catch(e){
            notify("Input Code Changes Error in parsing json.","error");
        }  
    }else if(type == "output"){
        try{
            componentIOData.outputJson = fromEdit_CodeEditor.get();
            fromEdit_TreeEditor.set(componentIOData.outputJson);
        }catch(e){
            notify("Output Code Changes Error in parsing json.","error");
        } 
    }
    
    fromEdit_TreeEditor.expandAll();
}

function fromTreeToCode(type){
    if(type == "input"){
        try{
            componentIOData.inputJson = fromEdit_TreeEditor.get();
            fromEdit_CodeEditor.set(componentIOData.inputJson);
        }catch(e){
            notify("Input Tree Changes Error in parsing json.","error");
        }  
    }else if(type == "output"){
        try{
            componentIOData.outputJson = fromEdit_TreeEditor.get();
            fromEdit_CodeEditor.set(componentIOData.outputJson);
        }catch(e){
            notify("Output Tree Changes Error in parsing json.","error");
        } 
    }
}