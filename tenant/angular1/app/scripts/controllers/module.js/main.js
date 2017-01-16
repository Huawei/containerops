auth.controller('ModuleController', ['$scope', '$location', '$state', 'moduleService', function($scope, $location, $state, moduleService) {
	
	$scope.moduleList = [];

	$scope.create = function(){
  	$state.go("project.create");
  }

	$scope.params = {
		startNum: $scope.moduleList.length,
		endNum: $scope.moduleList.length+10
	}

	$scope.getList = function(params){
		moduleService.getList(params)
			.then(function(data){
				// $scope.moduleList = data;
				$scope.moduleList = [
					{
						"id":100,
						"name":"p1",
						"desc":"p1"
					},
					{
						"id":101,
						"name":"p2",
						"desc":"p2"
					}
				]
			}, function(err){
				console.log('获取module list err:',err);
				$scope.moduleList = [
					{
						"id":100,
						"name":"p1",
						"desc":"p1"
					},
					{
						"id":101,
						"name":"p2",
						"desc":"p2"
					}
				]
			})	
	};

	$scope.getList($scope.params);

}])
.controller('ModuleCreateController', ['$scope', '$location','$state', 'moduleService', function($scope, $location, $state, moduleService) {
  $scope.baseInfo = {
		name: '',
		desc: '',
		visible: 'public'
	};


	$scope.saveBaseInfo = function(){
		console.log($scope.baseInfo)
		var params = $scope.baseInfo;
		if(params.name){
			moduleService.saveBaseInfo(params)
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

