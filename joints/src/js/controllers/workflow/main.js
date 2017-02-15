define(["app"], function(app) {
    app.controllerProvider.register('WorkflowController', ['$scope', '$location','notifyService', function($scope, $location,notifyService) {
        $scope.workflow = "workflow controllerProvider";
        // notifyService.notify("zyjtest","success");
    }]);
})
