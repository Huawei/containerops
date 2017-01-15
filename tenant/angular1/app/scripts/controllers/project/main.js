auth.controller('ProjectController', ['$scope', '$location', '$state', 'projectService', function($scope, $location, $state, projectService) {
	
	$scope.projectList = [];

	$scope.create = function(){
  	$state.go("project.create");
  }

	$scope.params = {
		startNum: $scope.projectList.length,
		endNum: $scope.projectList.length+10
	}

	$scope.getList = function(params){
		projectService.getList(params)
			.then(function(data){
				// $scope.projectList = data;
				$scope.projectList = [
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
				console.log('获取project list err:',err);
				$scope.projectList = [
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
	};

	$scope.getList($scope.params);

}])
.controller('ProjectCreateController', ['$scope', '$location','$state', 'projectService', function($scope, $location, $state, projectService) {
  $scope.baseInfo = {
		name: '',
		desc: '',
		visible: 'public'
	};


	$scope.saveBaseInfo = function(){
		console.log($scope.baseInfo)
		var baseInfo = $scope.baseInfo;
		if(baseInfo.name){
			projectService.saveBaseInfo(baseInfo)
				.then(function(data){
					console.log('保存成功')
				},function(err){
					console.log('保存失败：',err)
				})
		}else{
			console.log('请填写用户名')
		}
	}

}]);

