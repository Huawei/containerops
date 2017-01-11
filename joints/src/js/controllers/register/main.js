login.controller('RegisterController', ['$scope', '$location', '$rootScope','registerService', 'notifyService',
  function($scope,$location,$rootScope,registerService,notifyService) {

  	$scope.user = {
		username: '',
		password: ''
	}

  	$scope.signup = function() {
		try{
			//fake, to be deleted
			var users = localStorage["users"];
			if(_.isUndefined(users)){
				users = [];
				users.push($scope.user);
				localStorage["users"] = JSON.stringify(users);
			}else{
				users = JSON.parse(users);
				users.push($scope.user);
				localStorage["users"] = JSON.stringify(users);
			}
			notifyService.notify("Sign up successfully. please login with your username.","success");
			//fake end
			changeNav();
		}catch(e){
			notifyService.notify("Fail to sign up.","error");
		}
	}

	var changeNav = function(val){
		$location.path('login');
	}
}]);
