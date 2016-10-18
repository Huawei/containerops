import {jsonEditor} from "../../vendor/jquery.jsoneditor";
import {notify} from "../common/notify";

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

    // input output from edit
    $("#tree-edit-tab").on('click',function(){
        initTreeEdit();
    })

    $("#output-from-edit-tab").on('click',function(){
        initFromEdit("output");
    });

    initTreeEdit();
}

export function initTreeEdit(){
    try{
        jsonEditor(treeEdit_OutputContainer,startIOData.outputJson, {
            change:function(data){
                startIOData.outputJson = data;
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
    fromEdit_TreeEditor.set(startIOData.outputJson);
    
    fromEdit_TreeEditor.expandAll();
}

export function getOutputForEvent(selecetedEvent){
    // call api

    var fakereturn = {"event":selecetedEvent};
    startIOData.outputJson = fakereturn;
    initFromView();
}