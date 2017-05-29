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

define(["app"], function(app) {
    app.controllerProvider.register('WorkflowController', ['$scope', '$state', 'notifyService', function($scope, $state, notifyService) {
        $scope.workflows = [{
        		"name":"workflow1",
        		"versions":[
        			{
        				"id" : 1,
        				"version" : "v1.0"
        			},
        			{
        				"id" : 2,
        				"version" : "v1.0"
        			}
        		]
        }]
        $scope.showCreateWorkflow = function(){
        	$state.go("workflow.create");
        };
        $scope.showMoreVersion = function(verionname){
        	notifyService.notify("no more workflow","info")
        };


    }]);
})
