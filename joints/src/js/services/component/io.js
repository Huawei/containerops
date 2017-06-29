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

define(['app'], function(app) {
    app.provide.factory("componentIO", ["notifyService", 'jsonEditor', function(notifyService, jsonEditor) {
        var treeEdit_InputContainer, treeEdit_OutputContainer;
        var fromEdit_InputCodeContainer, fromEdit_InputTreeContainer, fromEdit_OutputCodeContainer, fromEdit_OutputTreeContainer;
        var fromEdit_OutputViewContainer;
        var fromEdit_InputCodeEditor, fromEdit_InputTreeEditor, fromEdit_OutputCodeEditor, fromEdit_OutputTreeEditor;

        var componentIOData;

        function init(component) {
            componentIOData = component;

            treeEdit_InputContainer = $('#inputTreeDiv');
            treeEdit_OutputContainer = $('#outputTreeDiv');
            fromEdit_InputCodeContainer = $("#inputCodeEditor")[0];
            fromEdit_InputTreeContainer = $("#inputTreeEditor")[0];
            fromEdit_OutputCodeContainer = $("#outputCodeEditor")[0];
            fromEdit_OutputTreeContainer = $("#outputTreeEditor")[0];

            initTreeEdit();
            initFromEdit("input");
            initFromEdit("output");
        }

        function initTreeEdit() {
            if (_.isUndefined(componentIOData.input) || _.isEmpty(componentIOData.input)) {
                $("#inputTreeStart").show();
                $("#inputTreeDiv").hide();
                $("#inputStartBtn").on('click', function() {
                    componentIOData.input = {
                        "newKey": null
                    }
                    initTreeEdit();
                })
            } else {
                try {
                    $("#inputTreeStart").hide();
                    $("#inputTreeDiv").show();
                    jsonEditor.init(treeEdit_InputContainer, componentIOData.input, {
                        change: function(data) {
                            componentIOData.input = data;
                        }
                    });
                } catch (e) {
                    notifyService.notify("Input Error in parsing json.", "error");
                }
            }

            if (_.isUndefined(componentIOData.output) || _.isEmpty(componentIOData.output)) {
                $("#outputTreeStart").show();
                $("#outputTreeDiv").hide();
                $("#outputStartBtn").on('click', function() {
                    componentIOData.output = {
                        "newKey": null
                    }
                    initTreeEdit();
                })
            } else {
                try {
                    $("#outputTreeStart").hide();
                    $("#outputTreeDiv").show();
                    jsonEditor.init(treeEdit_OutputContainer, componentIOData.output, {
                        change: function(data) {
                            componentIOData.output = data;
                        }
                    });
                } catch (e) {
                    notifyService.notify("Output Error in parsing json.", "error");
                }
            }
        }

        function initFromEdit(type) {
            var codeOptions = {
                "mode": "code",
                "indentation": 2
            };

            var treeOptions = {
                "mode": "tree",
                "search": true
            };

            if (type == "input") {
                if (fromEdit_InputCodeEditor) {
                    fromEdit_InputCodeEditor.destroy();
                }
                if (fromEdit_InputTreeEditor) {
                    fromEdit_InputTreeEditor.destroy();
                }
                fromEdit_InputCodeEditor = new JSONEditor(fromEdit_InputCodeContainer, codeOptions);
                fromEdit_InputTreeEditor = new JSONEditor(fromEdit_InputTreeContainer, treeOptions);
                fromEdit_InputCodeEditor.set(componentIOData.input);
                fromEdit_InputTreeEditor.set(componentIOData.input);
                $("#inputFromJson").on('click', function() {
                    fromCodeToTree("input");
                })
                $("#inputToJson").on('click', function() {
                    fromTreeToCode("input");
                })

                fromEdit_InputTreeEditor.expandAll();
            } else if (type == "output") {
                if (fromEdit_OutputCodeEditor) {
                    fromEdit_OutputCodeEditor.destroy();
                }
                if (fromEdit_OutputTreeEditor) {
                    fromEdit_OutputTreeEditor.destroy();
                }
                fromEdit_OutputCodeEditor = new JSONEditor(fromEdit_OutputCodeContainer, codeOptions);
                fromEdit_OutputTreeEditor = new JSONEditor(fromEdit_OutputTreeContainer, treeOptions);
                fromEdit_OutputCodeEditor.set(componentIOData.output);
                fromEdit_OutputTreeEditor.set(componentIOData.output);
                $("#outputFromJson").on('click', function() {
                    fromCodeToTree("output");
                })
                $("#outputToJson").on('click', function() {
                    fromTreeToCode("output");
                })

                fromEdit_OutputTreeEditor.expandAll();
            }
        }

        function fromCodeToTree(type) {
            if (type == "input") {
                try {
                    componentIOData.input = fromEdit_InputCodeEditor.get();
                    fromEdit_InputTreeEditor.set(componentIOData.input);
                } catch (e) {
                    notifyService.notify("Input Code Changes Error in parsing json.", "error");
                }
                fromEdit_InputTreeEditor.expandAll();
            } else if (type == "output") {
                try {
                    componentIOData.output = fromEdit_OutputCodeEditor.get();
                    fromEdit_OutputTreeEditor.set(componentIOData.output);
                } catch (e) {
                    notifyService.notify("Output Code Changes Error in parsing json.", "error");
                }
                fromEdit_OutputTreeEditor.expandAll();
            }
        }

        function fromTreeToCode(type) {
            if (type == "input") {
                try {
                    componentIOData.input = fromEdit_InputTreeEditor.get();
                    fromEdit_InputCodeEditor.set(componentIOData.input);
                } catch (e) {
                    notifyService.notify("Input Tree Changes Error in parsing json.", "error");
                }
            } else if (type == "output") {
                try {
                    componentIOData.output = fromEdit_OutputTreeEditor.get();
                    fromEdit_OutputCodeEditor.set(componentIOData.output);
                } catch (e) {
                    notifyService.notify("Output Tree Changes Error in parsing json.", "error");
                }
            }
        }

        function getInputEmpty() {
            return isInputEmpty;
        }

        function getOutputEmpty() {
            return isOutputEmpty;
        }

        return {
            "init": init,
            "initTreeEdit": initTreeEdit,
            "initFromEdit": initFromEdit,
            "getInputEmpty": getInputEmpty,
            "getOutputEmpty": getOutputEmpty
        }
    }])
})
