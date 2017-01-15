auth.controller('ApplicationController', ['$scope', '$location', '$state', 'applicationService', function($scope, $location, $state, applicationService) {
	$scope.applicationList = [];

	$scope.create = function(){
	  $state.go("application.create");
	}

	$scope.params = {
		startNum: $scope.applicationList.length,
		endNum: $scope.applicationList.length+10
	}

	applicationService.getList($scope.params)
		.then(function(data){
			// $scope.applicationList = data;
			$scope.applicationList = [
				{
					"id":100,
					"name":"p1",
					"desc":"p1",
					"applicationNum":5,
					"moduleNum":8
				},
				{
					"id":101,
					"name":"p2",
					"desc":"p2",
					"applicationNum":7,
					"moduleNum":10
				}
			]
		}, function(err){
			console.log('获取application list err:',err);
			$scope.applicationList = [
				{
					"id":100,
					"name":"p1",
					"desc":"p1",
					"applicationNum":5,
					"moduleNum":8
				},
				{
					"id":101,
					"name":"p2",
					"desc":"p2",
					"applicationNum":7,
					"moduleNum":10
				}
			]
		})	

}])
.controller('ApplicationCreateController', ['$scope', '$location','$state', 'applicationService', function($scope, $location, $state, applicationService) {
    

}]);

