define(["app","services/component/main","services/component/io","services/component/check"], function(app) {
    app.controllerProvider.register('CreateComponentController', ['$scope', '$location', 'componentService', 'componentIO',
        'componentCheck', 'notifyService', 'apiService', 'loading',
        function($scope, $location, componentService, componentIO, componentCheck, notifyService, apiService, loading) {
            // tabs control
            $scope.selectTab = function(index) {
                $scope.tab = index;

                checkTabs();
                if (index == 2) {
                    initScriptEditor();
                    showEventScript();
                }
            }

            function checkTabs() {
                $scope.tabStatus.runtime = componentCheck.tabcheck.runtime($scope.component);
                $scope.tabStatus.newimage = componentCheck.tabcheck.editshell($scope.component) && componentCheck.tabcheck.buildimage($scope.component);
            }

            // runtime tab all functions below
            $scope.changeRuntimeTab = function(index) {
                $scope.runtimeTab = index;
                if ($scope.runtimeTab == 2) {
                    componentIO.init($scope.component);
                }
            }

            $scope.changeNewImageTab = function(index){
                $scope.newImageTab = index;
                showEventScript();
            }

            $scope.baseOrAdvanced = function() {
                if ($scope.component.use_advanced) {
                    $scope.component.service = {};
                    $scope.component.pod = {};

                    setTimeout(function() {
                        $("#serviceCodeEditor").val(JSON.stringify($scope.component.service, null, 2));
                        $("#serviceCodeEditor").on("blur", function() {
                            var result = toJsonYaml("service");
                            if (result) {
                                $scope.component.service = result;
                            }
                        });

                        $("#podCodeEditor").val(JSON.stringify($scope.component.pod, null, 2));
                        $("#podCodeEditor").on("blur", function() {
                            var result = toJsonYaml("pod");
                            if (result) {
                                $scope.component.pod = result;
                            }
                        });
                    }, 500);
                } else {
                    $scope.component.service = $.extend(true, {}, componentService.metadata.base_service);
                    $scope.component.pod = $.extend(true, {}, componentService.metadata.base_pod);
                    if ($scope.component.service.spec.type == "NodePort") {
                        $scope.component.service.spec.ports.push($.extend(true, {}, componentService.metadata.nodeport));
                    } else {
                        $scope.component.service.spec.ports.push($.extend(true, {}, componentService.metadata.clusterip));
                    }
                    $scope.serviceType = $scope.component.service.spec.type;
                }
            }

            $scope.addPort = function() {
                if ($scope.component.service.spec.type == "NodePort") {
                    $scope.component.service.spec.ports.push($.extend(true, {}, componentService.metadata.nodeport));
                } else {
                    $scope.component.service.spec.ports.push($.extend(true, {}, componentService.metadata.clusterip));
                }
            }

            $scope.removePort = function(index) {
                $scope.component.service.spec.ports.splice(index, 1);
            }

            $scope.changeServiceType = function() {
                if ($scope.component.service.spec.type == "NodePort") {
                    $scope.component.service.spec.ports = [];
                    $scope.component.service.spec.ports.push($.extend(true, {}, componentService.metadata.nodeport));
                } else {
                    $scope.component.service.spec.ports = [];
                    $scope.component.service.spec.ports.push($.extend(true, {}, componentService.metadata.clusterip));
                }
                $scope.serviceType = $scope.component.service.spec.type;
            }

            $scope.switchMode = function(value) {
                $scope.jsonMode = value;
                if ($scope.jsonMode) {
                    componentIO.initFromEdit("input");
                    componentIO.initFromEdit("output");
                } else {
                    componentIO.initTreeEdit();
                }
            }

            $scope.addEnv = function() {
                $scope.component.env.push($.extend(true, {}, componentService.metadata.env));
            }

            $scope.removeEnv = function(index) {
                $scope.component.env.splice(index, 1);
            }

            function toJsonYaml(type) {
                var value, result;
                if (type == "service") {
                    value = $("#serviceCodeEditor").val();
                } else if (type == "pod") {
                    value = $("#podCodeEditor").val();
                }

                try {
                    result = JSON.parse(value);
                } catch (e) {
                    try {
                        result = jsyaml.safeLoad(value);
                    } catch (e) {
                        notifyService.notify("Your advanced " + type + " setting is not a legal json or yaml.", "error");
                        result = false;
                    }
                }
                if (!result) {
                    notifyService.notify("Your advanced " + type + " setting is not a legal json or yaml.", "error");
                }
                return result;
            }

            // shell edit page functions
            $scope.selectEvent = function(index) {
                $scope.selectedEvent = index;
                showEventScript();
            }

            function initScriptEditor() {
                $scope.scriptEditor = ace.edit("scriptEditor");
                $scope.scriptEditor.setTheme("ace/theme/dawn");
                $scope.scriptEditor.getSession().setMode("ace/mode/golang");
                $scope.scriptEditor.on("blur", function() {
                    setEventScript();
                })
            }

            function showEventScript() {
                switch ($scope.selectedEvent) {
                    case 1:
                        $scope.scriptEditor.setValue($scope.component.image_setting.events.component_start);
                        break;
                    case 2:
                        $scope.scriptEditor.setValue($scope.component.image_setting.events.component_result);
                        break;
                    case 3:
                        $scope.scriptEditor.setValue($scope.component.image_setting.events.component_stop);
                        break;
                }
                focusEventScript();
            }

            function setEventScript() {
                switch ($scope.selectedEvent) {
                    case 1:
                        $scope.component.image_setting.events.component_start = $scope.scriptEditor.getValue();
                        break;
                    case 2:
                        $scope.component.image_setting.events.component_result = $scope.scriptEditor.getValue();
                        break;
                    case 3:
                        $scope.component.image_setting.events.component_stop = $scope.scriptEditor.getValue();
                        break;
                }
            }

            function focusEventScript() {
                var session = $scope.scriptEditor.getSession();
                var row = session.getLength();
                var col = session.getLine(row - 1).length;
                $scope.scriptEditor.gotoLine(row, col);
                $scope.scriptEditor.focus();
            }

            // build push page functions
            $scope.changeBaseImageType = function() {
                if ($scope.component.image_setting.from.type == "url") {
                    delete $scope.component.image_setting.from.name;
                    delete $scope.component.image_setting.from.tag;
                    delete $scope.component.image_setting.from.dockerfile;
                    $scope.component.image_setting.from["url"] = "";
                } else if ($scope.component.image_setting.from.type == "dockerfile") {
                    delete $scope.component.image_setting.from.name;
                    delete $scope.component.image_setting.from.tag;
                    delete $scope.component.image_setting.from.url;
                    $scope.component.image_setting.from["dockerfile"] = "";
                } else if ($scope.component.image_setting.from.type == "name") {
                    delete $scope.component.image_setting.from.url;
                    delete $scope.component.image_setting.from.dockerfile;
                    $scope.component.image_setting.from["name"] = "";
                    $scope.component.image_setting.from["tag"] = "";
                }
            }

            // save component
            $scope.saveComponent = function() {
                var result = componentCheck.go($scope.component);
                if (result) {
                    var promise = componentService.addComponent($scope.component);
                    promise.done(function(data) {
                        apiService.successToCall(data);
                        $location.path("/component/" + data.component.id);
                    });
                    promise.fail(function(xhr, status, error) {
                        apiService.failToCall(xhr.responseJSON);
                    });
                }
            }

            // back to list
            $scope.backToList = function() {
                $location.path("/component");
            }

            // debug tab functions
            function initDebugEdit() {
                $scope.debugTab = 1;
                $scope.debugLogTab = 1;
                $scope.debugcomponent.input = $.extend(true, {}, $scope.component.input);
                $scope.debugcomponent.env = [];
                _.each($scope.component.env, function(item) {
                    $scope.debugcomponent.env.push($.extend(true, {}, item));
                })
                $scope.hasDebugInput = !_.isEmpty($scope.debugcomponent.input);
                if ($scope.hasDebugInput) {
                    initDebugInput();
                }
            }

            function initDebugInput() {
                if ($scope.inputCodeEditorForDebug) {
                    $scope.inputCodeEditorForDebug.destroy();
                }
                var codeOptions = {
                    "mode": "code",
                    "indentation": 2,
                    "onChange": function() {
                        $scope.debugcomponent.input = $scope.inputCodeEditorForDebug.get();
                    }
                };
                var inputCodeContainer = $("#inputCodeEditorForDebug")[0];
                $scope.inputCodeEditorForDebug = new JSONEditor(inputCodeContainer, codeOptions);
                $scope.inputCodeEditorForDebug.set($scope.debugcomponent.input);
            }

            $scope.changeDebugTab = function(index) {
                $scope.debugTab = index;
            }

            $scope.changeDebugLogTab = function(index) {
                $scope.debugLogTab = index;
            }

            $scope.debug = function() {
                if ($("#component-debug-edit-form").parsley().validate()) {
                    var dataStream = componentService.debugComponent($scope.component.id);
                    dataStream.onMessage(function(message) {
                        apiService.successToCall(data);
                        console.log(message);
                    });
                    dataStream.onError(function(error) {
                        apiService.failToCall(error);
                        console.log(error);
                    })
                    dataStream.send(JSON.stringify($scope.debugcomponent))
                }
            }
            
            // init component create page
            function init() {
                $scope.component = $.extend(true, {}, componentService.metadata.component);

                // for top tabs control
                $scope.tab = 1;
                $scope.tabStatus = {
                    "runtime": false,
                    "newimage": false
                }

                // for runtime config tabs control
                $scope.runtimeTab = 1;

                // determine which editor to use for input output json
                $scope.jsonMode = false;

                // for new image tabs control
                $scope.newImageTab = 1;

                // for event selection
                $scope.selectedEvent = 1;

                // init service pod
                $scope.baseOrAdvanced();

                $scope.$watch("component",function(){
                    if($scope.scriptEditor){
                        showEventScript();
                    }
                });

                $scope.debugcomponent = {
                    "kube_master": "",
                    "input": {},
                    "env": []
                }

                initDebugEdit();
            }

            init();
        }
    ])
})
