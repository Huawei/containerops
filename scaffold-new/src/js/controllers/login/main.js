login.controller('LoginController', ['$scope', '$location', '$rootScope','loginService', 'notifyService',
  function($scope,$location,$rootScope,loginService,notifyService) {

  	$scope.user = {
		username: '',
		password: ''
	}

  	$scope.login = function() {
		try{
			//fake, to be deleted
			if(_.isUndefined(localStorage["users"])){
				notifyService.notify("No such user, please sign up first.","info");
			}else{
				var users = JSON.parse(localStorage["users"]);
				var targetuser = _.find(users,function(item){
					return item.username == $scope.user.username && item.password == $scope.user.password;
				});
				if(_.isUndefined(targetuser)){
					notifyService.notify("No such user or password is incorrect.","info");
				}else{
					notifyService.notify("Welcome. " + $scope.user.username , "success");
					sessionStorage["currentUser"] = $scope.user.username;
					// changeNav('index');
				}
			}
			//fake end
		}catch(e){
			notifyService.notify("Failed to Sign in.", "error");
		}
	}

	var changeNav = function(val){
		$location.path("index");
	}
}]);
