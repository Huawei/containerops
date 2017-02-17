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
