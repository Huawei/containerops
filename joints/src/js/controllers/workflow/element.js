define(["app"], function(app) {
    app.controllerProvider.register('WorkflowElementController', ['$scope', '$state', 'notifyService', function($scope, $state, notifyService) {
        $scope.workflows = [{"name":"workflow1","versions":["v1.0","v1.1"]}]
        $scope.showCreateWorkflow = function(){
        	$state.go("workflow.detail");
        };
        $scope.showMoreVersion = function(verionname){
        	notifyService.notify("no more workflow","info")
        };


    }]);
})
