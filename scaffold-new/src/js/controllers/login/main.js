login.controller('LoginController', ['$scope', '$rootScope','loginService', 'notifyService',
  function($scope,$rootScope,loginService,notifyService) {

  	$scope.login = function() {
		try{
			//fake, to be deleted
			var self = this;
			if(_.isUndefined(localStorage["users"])){
				notifyService.notify("No such user, please sign up first.","info");
			}else{
				var users = JSON.parse(localStorage["users"]);
				var targetuser = _.find(users,function(item){
					return item.username == self.user.username && item.password == self.user.password;
				});
				if(_.isUndefined(targetuser)){
					notifyService.notify("No such user or password is incorrect.","info");
				}else{
					notifyService.notify("Welcome. " + self.user.username , "success");
					sessionStorage["currentUser"] = self.user.username;
					// $scope.changeNav('index');
				}
			}
			//fake end
		}catch(e){
			notifyService.notify("Failed to Sign in.", "error");
		}
	}

	$scope.changeNav = function(val){

	}
}]);
