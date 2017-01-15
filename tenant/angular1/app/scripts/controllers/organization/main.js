auth.controller('OrganizationController', ['$scope', '$location','$state', 'OrganizationService', function($scope, $location, $state, OrganizationService) {
    $scope.get = function() {
        OrganizationService.get().then(function(data) {
                // $scope.orgs = data.data;
                $scope.orgs = [{"name":"org1","desc":"desc1"}];
            },
            function(errorMessage) {
                console.log("error")
                $scope.orgs = [{"name":"org1","desc":"desc1"}];
            });
    };
    $scope.create = function(){
    	$state.go("organization.create");
    }
    $scope.get();


}])
 .controller('OrganizationCreateController', ['$scope', '$location','$state', 'OrganizationService', function($scope, $location, $state, OrganizationService) {
     $scope.back = function(){
     	$state.go("organization");
     }

}]);
