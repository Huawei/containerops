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
import {ajaxCall} from "../common/api";
import {loading} from "../common/loading";
import * as startIOData from "./startIOData";

var treeEdit_OutputContainer;
var fromEdit_OutputCodeContainer,fromEdit_OutputTreeContainer;
var fromEdit_OutputViewContainer;
var fromEdit_CodeEditor,fromEdit_TreeEditor;

export function initStartIO(start){
    startIOData.getStartIOData(start);
    startIOData.setSelectedTab(0);
    showOutputTabs();

    $(".newStartOutput").on('click',function(){
        startIOData.addOutput();  
        startIOData.setSelectedTab(startIOData.data.length-1);
        showOutputTabs();
    });

    $(".deleteStartOutput").on('click',function(){
        startIOData.deleteOutput();  
        startIOData.setSelectedTab(0);
        showOutputTabs();
    });
}

function showOutputTabs(){
    $("#startOutputTabs").empty();
    $("#startOutputTabsContent").empty();

    _.each(startIOData.data,function(output,index){
        var tabitem = `<li class="nav-item start-output-tab" data-index="`+ index +`">
                            <a class="nav-link" href="#output-` + index + `" data-toggle="tab">Output ` 
                            + (index+1) + `</a></li>`;
        $("#startOutputTabs").append(tabitem);

        var tabcontentitem = `<div class="tab-pane" id="output-`+ index +`">
                                <div class="output-type-event-div" data-index="`+ index +`">
                                    <div class="row col-md-6">
                                        <label class="col-md-4 control-label">Select Type</label>
                                        <div class="col-md-8">
                                            <select class="output-type-select" style="width:100%">
                                                <option value="github">Github</option>
                                                <option value="gitlab">Gitlab</option>
                                                <option value="customize">Customize</option>
                                            </select>
                                        </div>
                                    </div>
                                    <div class="row col-md-6 event-select-div">
                                        <label class="col-md-4 control-label">Event</label>
                                        <div class="col-md-8">
                                            <select class="github-event-select" style="width:100%">
                                                <option value="Create">Create</option>
                                                <option value="Delete">Delete</option>
                                                <option value="Deployment">Deployment</option>
                                                <option value="DeploymentStatus">Deployment Status</option>
                                                <option value="Fork">Fork</option>
                                                <option value="Gollum">Gollum</option>
                                                <option value="IssueComment">Issue Comment</option>
                                                <option value="Issues">Issues</option>
                                                <option value="Member">Member</option>
                                                <option value="PageBuild">Page Build</option>
                                                <option value="Public">Public</option>
                                                <option value="PullRequestReviewComment">Pull Request Review Comment</option>
                                                <option value="PullRequestReview">Pull Request Review</option>
                                                <option value="PullRequest">Pull Request</option>
                                                <option value="Push">Push</option>
                                                <option value="Repository">Repository</option>
                                                <option value="Release">Release</option>
                                                <option value="Status">Status</option>
                                                <option value="TeamAdd">Team Add</option>
                                                <option value="Watch">Watch</option>
                                            </select>
                                            <select class="gitlab-event-select" style="width:100%">
                                                <option value="Push Hook">Push</option>
                                                <option value="Tag Push Hook">Tag Push</option>
                                                <option value="Note Hook">Connents</option>
                                                <option value="Issue Hook">Issues</option>
                                                <option value="Merge Request Hook">Merge Request</option>
                                            </select>
                                        </div>
                                    </div>
                                    <div class="row col-md-6 event-input-div">
                                        <label class="col-md-4 control-label">Event</label>
                                        <div class="col-md-8">
                                            <input type="text" class="form-control output-event-input" required>
                                        </div>
                                    </div>
                                </div>
                                <div class="row col-md-12 output-json-div" data-index="`+ index +`">
                                    <div class="startOutputTreeViewer"></div>
                                    <div class="startOutputTreeDesigner">
                                        <div class="row">
                                            <div class="col-md-6 import-div">
                                                <div class="panel">
                                                    <div class="panel-heading clearfix">
                                                        <i class="glyphicon glyphicon-cloud-download outputicon"></i>
                                                        <span class="panel-title">Output Tree Edit</span>
                                                    </div>
                                                    <div class="panel-body">
                                                        <div class="outputTreeStart tree-add-button">
                                                            <div class="outputStartBtn btn-div">
                                                                <span class="glyphicon glyphicon-plus nohover"></span>
                                                                <div class="desc">
                                                                    <label class="desc-label">Add New Value</label>
                                                                    <div class="desc-btn">
                                                                        <span class="glyphicon glyphicon-plus"></span>
                                                                    </div>
                                                                </div>
                                                            </div>
                                                        </div>
                                                        <div class="outputTreeDiv json-editor"></div>
                                                    </div>
                                                </div>
                                            </div>
                                            <div class="col-md-6 import-div">
                                                <div class="panel">
                                                    <div class="panel-heading clearfix">
                                                        <i class="glyphicon glyphicon-cloud-download outputicon"></i>
                                                        <span class="panel-title">Output From Edit</span>
                                                    </div>
                                                    <div class="panel-body">
                                                        <div class="outputCodeEditor codeEditor"></div>
                                                        <div class="col-md-4 col-md-offset-4 row editor-transfer">
                                                            <span class="outputFromJson col-md-4 code-to-tree"></span>
                                                            <span class="outputToJson col-md-4 tree-to-code"></span>
                                                        </div>
                                                        <div class="outputTreeEditor treeEditor"></div>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>`;
        $("#startOutputTabsContent").append(tabcontentitem);
    });

    // event binding
    $(".start-output-tab").on('click',function(event){
        var index = $(event.currentTarget).data("index");
        startIOData.setSelectedTab(index);
        initOutputDiv();
    });

    $(".output-type-select").on("change",function(){
        startIOData.setTypeSelect();
        selectType(startIOData.getTypeSelect(),true);
    });

    $(".github-event-select").on("change",function(){
        startIOData.setEventSelect("github");
        getOutputForEvent(startIOData.getTypeSelect(),startIOData.getEventSelect());
    });

    $(".gitlab-event-select").on("change",function(){
        startIOData.setEventSelect("gitlab");
        getOutputForEvent(startIOData.getTypeSelect(),startIOData.getEventSelect());
    });

    $(".output-event-input").on("blur",function(){
        startIOData.setEventInput();
        if(!startIOData.isEventOptionAvailable()){
            notify("There's a customize output for event '" + startIOData.getEventSelect() + "', please input another one.","info");
            startIOData.setEvent("");
            startIOData.setEventInputDom();
        }
    });

    $(".outputStartBtn").on('click',function(){
            startIOData.setJson({
                "newKey" : null
            });
            initTreeEdit();
            initFromEdit("output");
    });

    $(".outputFromJson").on('click',function(){
        fromCodeToTree("output");
    });

    $(".outputToJson").on('click',function(){
        fromTreeToCode("output");
    });

    // use global vars
    // var globalvars = _.map(workflowVars,function(item){
    //     return "@"+item[0]+"@";
    // });
    // $(".allowFromVar").autocomplete({
    //     source:[globalvars],
    //     limit: 100,
    //     visibleLimit: 5
    // }); 

    // init trigger
    startIOData.findSelectedStartOutputTabDom().find("a").addClass("active");
    startIOData.findSelectedStartOutputTabContentDom().addClass("active");
    initOutputDiv();
}

function initOutputDiv(){
    startIOData.setTypeSelectDom();
    selectType(startIOData.getTypeSelect());
}

function selectType(workflowType,isTypeChange){
    if(isTypeChange){
        startIOData.setJson({});
        startIOData.setEvent("");
    } 

    if(workflowType == "github" || workflowType == "gitlab"){
        startIOData.findEventSelectDivDom().show();
        startIOData.findEventInputDivDom().hide();
        startIOData.findOutputTreeViewerDom().show();
        startIOData.findOutputTreeDesignerDom().hide();

        if(workflowType == "github"){
            startIOData.findGitHubEventSelectDom().show();
            startIOData.findGitLabEventSelectDom().hide();
            startIOData.findGitLabEventSelectDom().next().hide();
        }else if(workflowType == "gitlab"){
            startIOData.findGitHubEventSelectDom().hide();
            startIOData.findGitHubEventSelectDom().next().hide();
            startIOData.findGitLabEventSelectDom().show();
        }

        if(_.isEmpty(startIOData.getJson()) || isTypeChange){
            if(_.isEmpty(startIOData.getEventSelect())){
                if(workflowType == "github"){
                    startIOData.setEvent("PullRequest");
                }else if(workflowType == "gitlab"){
                    startIOData.setEvent("Push Hook");
                }
                
                startIOData.setEventSelectDom(workflowType);
            }
            getOutputForEvent(startIOData.getTypeSelect(),startIOData.getEventSelect()); 
        }else{
            initFromView();
        }

        startIOData.setEventSelectDom(workflowType);
    }else{
        startIOData.findEventSelectDivDom().hide();
        startIOData.findEventInputDivDom().show();
        startIOData.findOutputTreeViewerDom().hide();
        startIOData.findOutputTreeDesignerDom().show();

        startIOData.setEventInputDom();

        initTreeEdit();
        initFromEdit("output");
    }
}

export function initTreeEdit(){
    if(_.isUndefined(startIOData.getJson()) || _.isEmpty(startIOData.getJson())){
        startIOData.findOutputTreeStartDom().show();
        startIOData.findOutputTreeDivDom().hide();
    }else{
        try{
            startIOData.findOutputTreeStartDom().hide();
            startIOData.findOutputTreeDivDom().show();

            treeEdit_OutputContainer = startIOData.findOutputTreeDivDom();
            jsonEditor(treeEdit_OutputContainer,startIOData.getJson(), {
                change:function(data){
                    startIOData.setJson(data);
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
        fromEdit_OutputCodeContainer = startIOData.findOutputCodeEditorDom()[0];
        fromEdit_CodeEditor = new JSONEditor(fromEdit_OutputCodeContainer, codeOptions);

        fromEdit_OutputTreeContainer = startIOData.findOutputTreeEditorDom()[0];
        fromEdit_TreeEditor = new JSONEditor(fromEdit_OutputTreeContainer, treeOptions);

        fromEdit_CodeEditor.set(startIOData.getJson());
        fromEdit_TreeEditor.set(startIOData.getJson());
    }
    
    fromEdit_TreeEditor.expandAll();
}

function fromCodeToTree(type){
    if(type == "output"){
        try{
            startIOData.setJson(fromEdit_CodeEditor.get());
            fromEdit_TreeEditor.set(startIOData.getJson());
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
            startIOData.setJson(fromEdit_TreeEditor.get());
            fromEdit_CodeEditor.set(startIOData.getJson());
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

    fromEdit_OutputViewContainer = startIOData.findOutputTreeViewerDom()[0];
    fromEdit_TreeEditor = new JSONEditor(fromEdit_OutputViewContainer, treeOptions);
    fromEdit_TreeEditor.set(startIOData.getJson());
    
    fromEdit_TreeEditor.expandAll();
}

export function getOutputForEvent(selectedType,selecetedEvent){
    if(startIOData.isEventOptionAvailable()){
        var params = {
            "site" : selectedType,
            "eventName" : selecetedEvent
        }
        var promise = ajaxCall("workflow.eventOutput",params);
        promise.done(function(data){
            loading.hide();
            startIOData.setJson(data.output);
            initFromView();
        });
        promise.fail(function(xhr,status,error){
            loading.hide();
            if (!_.isUndefined(xhr.responseJSON) && xhr.responseJSON.errMsg) {
                notify(xhr.responseJSON.errMsg,"error");
            }else if(xhr.statusText != "abort"){
                notify("Server is unreachable","error");
            }
        });
    }else{
        if(fromEdit_TreeEditor){
            fromEdit_TreeEditor.destroy();
        }
        startIOData.setJson({});
        notify("There's a " + startIOData.getTypeSelect() + " output for event '" + selecetedEvent + "', please select another one.","info");
    }  
}