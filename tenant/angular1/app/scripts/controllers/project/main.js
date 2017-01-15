auth.controller('ProjectController', ['$scope', '$location', 'projectService', function($scope, $location, projectService) {
	$scope.list = [];
	
	$scope.params = {
		startNum: $scope.list.length,
		endNum: $scope.list.length+10
	}
	
	projectService.getProjectList($scope.params,function(data){
		console.log(data)
	})
		// .success(function(data){
		// 	console.log(data)
		// })
		// .error(function(){
		// 	console.log('获取projectList失败')
		// })

}]);
