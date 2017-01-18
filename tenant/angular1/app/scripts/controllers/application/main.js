auth.controller('ApplicationController', ['$scope', '$location', '$state', '$stateParams', 'applicationService', 'projectService', function($scope, $location, $state, $stateParams, applicationService, projectService) {
	$scope.applicationList = [];
	$scope.projectList = [];
	// $scope.params = {
	// 	startNum: $scope.applicationList.length,
	// 	endNum: $scope.applicationList.length+10
	// };
	// filtered project
	$scope.filterItem = {
		id: '',
		name: ''
	}; 

	$scope.create = function(){
	  $state.go("application.create");
	};

	$scope.getList = function(params){
		applicationService.getList(params)
			.then(function(data){
				// $scope.applicationList = data;
				$scope.applicationList = [
					{
						"id":100,
						"name":"a1",
						"desc":"a1",
						"applicationNum":5,
						"moduleNum":8
					},
					{
						"id":101,
						"name":"a2",
						"desc":"a2",
						"applicationNum":7,
						"moduleNum":10
					}
				]
			}, function(err){
				console.log('get list err:',err);
				$scope.applicationList = [
					{
						"id":100,
						"name":"a1",
						"desc":"a1",
						"applicationNum":5,
						"moduleNum":8
					},
					{
						"id":101,
						"name":"a2",
						"desc":"a2",
						"applicationNum":7,
						"moduleNum":10
					}
				]
			});	
	};

	$scope.edit = function(item){
		$state.go('application.edit',{id:item.id,name:item.name})
	};

	$scope.resetFilter = function(item){
		$scope.filterItem = {
			id: item.id,
			name: item.name
		};
		$scope.getList({id: item.id});
	};

	$scope.getProjectList = function(){
		projectService.getList()
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
				];
				$scope.isStateParams();

			}, function(err){
				console.log('get list err:',err);
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
				];
				$scope.isStateParams();
			});
	};

	$scope.isStateParams = function(){
		var item = {
			id: $stateParams.id,
			name: $stateParams.name
		};
		if(!item.id){
			$scope.resetFilter($scope.projectList[0])  
		}else{
			$scope.resetFilter(item)
		};
	};

	$scope.showChildren = function(item){
		$state.go('module',{id:item.id, name:item.name});
	};

	$scope.init = function(){
		$scope.getProjectList();
	};

	$scope.init();
}])

.controller('ApplicationCreateController', ['$scope', '$location','$state', '$stateParams', 'applicationService', function($scope, $location, $state, $stateParams, applicationService) {
  $scope.currentStep = 'baseInfo';
  $scope.action = 'Create';
  $scope.baseInfo = {
		name: '',
		desc: '',
		visible: 'public',
		id: '-1'
	};
	$scope.projectList = [];
	$scope.teamList = [];
	$scope.chooseTeams = [];
	$scope.roleTeams = [];
	$scope.isShowOrgs = false;
	$scope.currentProject = {
		name: '',
		id: ''
	};

	$scope.changeStep = function(val){
		$scope.currentStep = val;
	}

	$scope.saveBaseInfo = function(){
		var params = {
			name: $scope.baseInfo.name,
			desc: $scope.baseInfo.desc,
			id: $scope.baseInfo.id
		};
		if(params.name){
			applicationService.saveBaseInfo(params)
				.then(function(data){
					$scope.baseInfo.id = data.id;
					console.log('success')
				},function(err){
					console.log('failed:',err)
				})
		}else{
			console.log('please input name')
		}
	}

	$scope.saveSetting = function(){
		if($scope.baseInfo.id!=='-1'){
			var params = $scope.chooseTeams;
			if(params.name){
				applicationService.saveSetting(params)
					.then(function(data){
						console.log('success')
					},function(err){
						console.log('failed：',err)
					})
			}else{
				console.log('please input name')
			}
		}else{
			console.log('please create application')
		}
	};

	// get project list
	$scope.getProjectList = function(params){
		applicationService.getProjectList(params)
			.then(function(data){
				// $scope.projectList = data;
				$scope.projectList = [
					{
						"id":100,
						"name":"p1"
					},
					{
						"id":101,
						"name":"p2"
					},
					{
						"id":102,
						"name":"p3"
					}
				]
				$scope.currentProject = $scope.projectList[0];
			},function(err){
				console.log('get list err:',err);
				$scope.projectList = [
					{
						"id":100,
						"name":"project1"
					},
					{
						"id":101,
						"name":"project2"
					},
					{
						"id":102,
						"name":"project3"
					}
				]
				$scope.currentProject = $scope.projectList[0];
			})
	};

	// get team list
	$scope.getTeamList = function(params){
		applicationService.getTeamList(params)
			.then(function(data){
				// $scope.teamList = data;
				$scope.teamList = [
					{
						"id":100,
						"name":"team1",
						"isChosed": false
					},
					{
						"id":101,
						"name":"team2",
						"isChosed": false
					},
					{
						"id":102,
						"name":"team3",
						"isChosed": false
					}
				]
			},function(err){
				console.log('get list err:',err);
				$scope.teamList = [
					{
						"id":100,
						"name":"team1",
						"isChosed": false
					},
					{
						"id":101,
						"name":"team2",
						"isChosed": false
					},
					{
						"id":102,
						"name":"team3",
						"isChosed": false
					}
				]
			})
	};

	// get current project to filter team list
	$scope.getCurrentProject = function(item){
		$scope.currentProject = item;
		$scope.chooseTeams = [];
		$scope.getTeamList($scope.currentProject);
		$scope.isShow('isShowOrgs',false);
	}

	$scope.choseTeam = function(item){
		angular.forEach($scope.teamList,function(obj,i){
			if(item.id === obj.id){
				obj.isChosed = !obj.isChosed
			}
		})
		$scope.getChosedTeams($scope.teamList)
	}

	$scope.getChosedTeams = function(originData){
		$scope.chooseTeams = []
		angular.forEach(originData,function(obj){
			if(obj.isChosed){
				$scope.chooseTeams.push(obj)
			}
		})
	}

	$scope.clearChosedStatus = function(originData){
		angular.forEach(originData,function(obj,i){
			obj.isChosed = false
		})
	}

	$scope.isShow = function(key,val){
		$scope[key] = val;
	}

	$scope.getRole = function(val){
		// console.log($scope.chooseTeams)
		var role = ["Admin","Readonly","ReadWrite"][val];
		var obj = {
			"role": role,
			"teams": $scope.chooseTeams
		}

		// if($scope.roleTeams.length>0){
		// 	angular.forEach($scope.roleTeams,function(itemRole,i){
		// 		if(itemRole.role === role){
		// 			console.log(777)
		// 			$scope.roleTeams.splice(i,1)
		// 			$scope.roleTeams.push(obj)
		// 			return
		// 		}else{
		// 			$scope.roleTeams.push(obj)
		// 			return
		// 		}
		// 	})
		// }else{
		// 	$scope.roleTeams.push(obj)
		// }
		$scope.roleTeams.push(obj)
		$scope.chooseTeams = [];
		$scope.clearChosedStatus($scope.teamList);
	};

	$scope.getEditInfo = function(){
		var id = $stateParams.id;
		if($stateParams.id){
			$scope.action = 'Edit';
  		$scope.baseInfo.id = $stateParams.id;
		}
		if(id){
			applicationService.getEditInfo(id)
				.then(function(data){
					// $scope.baseInfo = data;

					$scope.baseInfo = {
						name: 'tom',
						desc: 'this is test',
						visible: 'public'
					};

				},function(err){
					console.log('get info err',err)
					$scope.baseInfo = {
						name: 'tom',
						desc: 'this is test',
						visible: 'public'
					};
				})
		}
	};

	$scope.getProjectList({user:"small"});
	$scope.getEditInfo();

}]);

