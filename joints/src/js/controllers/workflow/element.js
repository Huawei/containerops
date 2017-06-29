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
    app.controllerProvider.register('WorkflowElementController', ['$scope', '$rootScope', '$state', 'notifyService', 'diagramService', function($scope, $rootScope, $state, notifyService, diagramService) {
        $scope.workflows = [{"name":"workflow1","versions":["v1.0","v1.1"]}]
        // $scope.showCreateWorkflow = function(){
        // 	$state.go("workflow.detail");
        // };
        // $scope.showMoreVersion = function(verionname){
        // 	notifyService.notify("no more workflow","info")
        // };
        var currentStageIndex = diagramService.currentStageIndex;
        $scope.workflowData = diagramService.workflowData;
        $scope.currentStageInfo = $scope.workflowData[currentStageIndex];
        $scope.isShowPage = true;

        $scope.redraw = function(){
            diagramService.resetWorkflowData($scope.workflowData);
            $rootScope.drawWorkflow();
        };

        $scope.stageEvent = {
            delete: function(){
                $scope.workflowData.splice(currentStageIndex,1);
                $scope.redraw();
                $scope.isShowPage = false;
            }
        };

    }]);
})
