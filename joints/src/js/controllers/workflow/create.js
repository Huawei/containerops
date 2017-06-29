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

define(["app","services/diagram/main"], function(app) {
    app.controllerProvider.register('WorkflowCreateController', ['$scope', '$state', '$rootScope', 'notifyService', 'diagramService', function($scope, $state, $rootScope, notifyService, diagramService) {
        $scope.navbar = 'workflowInfo';
        $scope.originData = angular.copy(diagramService.workflowData);

        $scope.isShowSetting = {
            showInfo: false,
            startWay: "number",
            settingType: "base",
        };

        $scope.setting = {
            "baseInfo":{
                "id":'',
                "name":"workflow",
                "version":"0.0.1",
                "webhookURL":"https://address/to/exec/workflow",
                "webhookSecret":"37dfeg8efab3",
                "latestRunStatus":1,
                "latestRunTime":"2016-01-02 15:04:05",
                "serverIp":"100.10.10.1",
                "nodeIp":"100.10.10.1"
            },
            "kubeSetting":{
               "apiServerAddr":"https://url/to/kube/api/server:port"
            },
            "runningInstances":{
                "available":true,
                "number":10
            },
            "manualStart":"test",
            "timedTasks":{
                "available":true,
                "tasks":[
                    {
                        "byDesigner":true,
                        "collapse":true,
                        "cronEntry":"* * * * *",
                        "eventName":"test",
                        "eventType":"test",
                        "startJson":{
                            "name":"test"
                        }
                    }
                ]
            },
            "env":{
                "serverIp":"100.10.10.1",
                "nodeIp":"100.10.10.1"
            },
            "globalVar":{
                "serverIp":"100.10.10.2",
                "nodeIp":"100.10.10.2"
            }
        };

        $scope.resetWorkflowData = function(){
            diagramService.resetWorkflowData($scope.originData);
        };

        $scope.changeNav = function(val){
            $scope.navbar = val;
        };

        $scope.backToList = function() {
            $state.go("workflow");
            $scope.resetWorkflowData();
        };

        $scope.saveWorkflow = function(){
            console.log(diagramService.workflowData)
            $state.go("workflow");
        };

        $scope.changeSettingNav = function(nav){
            $scope.isShowSetting.settingType = nav;
        };

        $scope.isShowDialog = function(val){
            $scope.isShowSetting.showInfo = val;
        };

        $scope.timeTaskEvent = {
            timeTask: {
                            "byDesigner":true,
                            "collapse":true,
                            "cronEntry":"* * * * *",
                            "eventName":"test",
                            "eventType":"test",
                            "startJson":{
                                "name":"test"
                            }
                        },
            delete: function(index){
                $scope.setting.timedTasks.tasks.splice(index,1)
            },
            add: function(){
                $scope.setting.timedTasks.tasks.push(angular.copy($scope.timeTaskEvent.timeTask))
            }
        };

        $rootScope.drawWorkflow = function() {
            diagramService.drawWorkflow($scope,'#div-d3-main-svg',diagramService.workflowData);
        };

        $scope.chosedStageIndex = '';
        $scope.chosedActionIndex = '';

        $scope.resetStageInfo = function(index){
            diagramService.currentStageIndex = index;
        };

        $scope.resetActionInfo = function(stageIndex,actionIndex){
            diagramService.currentActionIndex = actionIndex;
            $scope.resetStageInfo(stageIndex);
        };

        $scope.newStage = {
            "name":"",
            "id":"stage-",
            "type":"edit-stage",
            "runMode":"parallel",
            "timeout":'',
            "actions":[
                {
                    "isChosed":false,
                    "name":"",
                    "id":"action-",
                    "type":"action",
                    "timeout":'',
                    "components":[]
                }
            ]
        };

        // $scope.newComponent = {
        //     "name":"",
        //     "id":"component-",
        //     "type":"action",
        //     "version":"",
        //     "inputData":"",
        //     "outputData":""
        // };

        $scope.newAction = {
            "isChosed":false,
            "name":"",
            "id":"action-",
            "type":"action",
            "timeout":'',
            "components":[]
        };
        // add or chosed stage
        $scope.chosedStage = function(d,i){
            event.stopPropagation();
            $scope.clearChosedIndex('stage');
            $scope.clearChosedStageColor();
            $scope.workflowData = diagramService.workflowData;
            if(d.type === 'add-stage'){
                var stage = angular.copy($scope.newStage);
                    stage.id += uuid.v1();
                    stage.actions[0].id += uuid.v1();
                var addstage = angular.copy($scope.workflowData[i]);
                var endstage = angular.copy($scope.workflowData[i+1]);

                $scope.workflowData[i] = stage;
                $scope.workflowData[i+1] = addstage;
                $scope.workflowData[i+2] = endstage;
                // dataset.splice(dataset[i-1],0,stage)
                $rootScope.drawWorkflow();
            };

            if(d.type === 'edit-stage'){
                d3.select(this).attr('href','assets/images/icon-stage.svg');
                // $scope.chosedStageIndex = i;
                $scope.resetStageInfo(i);
                $state.go("workflow.create.stage",{"id": d.id});
            };
        };

        $scope.clearChosedStageColor = function(){
            d3.selectAll('.stage-pic')
                .attr('href',function(d){
                    if(d.type === 'add-stage'){
                        return 'assets/images/icon-add-stage.svg';
                    }else if(d.type === 'end-stage'){
                        return 'assets/images/icon-stage-empty.svg';
                    }else if(d.type === 'edit-stage'){
                        return d.runMode === 'parallel' ? 'assets/images/icon-action-parallel.svg' : 'assets/images/icon-action-serial.svg';
                    }
                })
        };

        $scope.clearChosedIndex = function(type){
            if(type === 'stage'){
                $scope.chosedStageIndex = '';
            }else{
                $scope.chosedActionIndex = '';
            }
        };

        // chosed action
        $scope.chosedAction = function(d,i){
            event.stopPropagation();
            $scope.clearChosedIndex('action');
            $scope.clearAddActionIcon();
            $scope.workflowData = diagramService.workflowData;

            var currentElement = d3.select(this);
            $scope.chosedStageIndex = parseInt(currentElement.attr('data-stageIndex'));
            $scope.chosedActionIndex = parseInt(currentElement.attr('data-actionIndex'));
            var chosedStageIndex = $scope.chosedStageIndex;
            var chosedActionIndex = $scope.chosedActionIndex;
            var isChosed = $scope.workflowData[chosedStageIndex]['actions'][chosedActionIndex]['isChosed'];
            $scope.workflowData[chosedStageIndex]['actions'][chosedActionIndex]['isChosed'] = !isChosed;
            $rootScope.drawWorkflow(); 
            $scope.resetActionInfo(chosedStageIndex,chosedActionIndex);
            $state.go("workflow.create.action",{"id": d.id});
        };

        $scope.clearAddActionIcon = function(){
            angular.forEach($scope.workflowData,function(d,i){
                angular.forEach(d.actions,function(a,ai){
                    a.isChosed = false;
                })
            })
        };

        // add action to bottom
        $scope.addBottomAction = function(){
            $scope.addElement(d3.select(this),'bottom');
        };

        // add action to top
        $scope.addTopAction = function(){
            $scope.addElement(d3.select(this),'top');
        };

        $scope.addElement = function(currentElement,type){
            event.stopPropagation();
            $scope.chosedStageIndex = parseInt(currentElement.attr('data-stageIndex'));
            $scope.chosedActionIndex = parseInt(currentElement.attr('data-actionIndex'));
            var chosedStageIndex = $scope.chosedStageIndex;
            var chosedActionIndex = $scope.chosedActionIndex;
            var actionLength = $scope.workflowData[chosedStageIndex]['actions'].length;
            var action = angular.copy($scope.newAction);
            action.id += uuid.v1();

            if(type === 'bottom'){
                if(actionLength === chosedActionIndex+1){
                    $scope.workflowData[chosedStageIndex]['actions'].push(action);
                }else{
                    $scope.workflowData[chosedStageIndex]['actions'].splice(chosedActionIndex+1,0,action);
                }
            }else{
                $scope.workflowData[chosedStageIndex]['actions'].splice(chosedActionIndex,0,action);
            };
                
            $rootScope.drawWorkflow(); 
        }; 

        // add component
        $scope.addComponent = function(){
            addElement(d3.select(this),'component');
        };

        function init(){
            initUuid();
            $scope.resetWorkflowData();
            $rootScope.drawWorkflow();
        };

        function initUuid(){
            var originData = $scope.originData;
            var length = originData.length;
            if(!originData[0].id){
                originData[0].id = 'stage' + uuid.v1();
                originData[length-1].id += uuid.v1();
                originData[length-2].id += uuid.v1();
                originData[0].actions[0].id += uuid.v1();
            }
        };

        init();
  

    }]);
})
