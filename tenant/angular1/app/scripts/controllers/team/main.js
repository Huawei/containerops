auth.controller('TeamController', ['$scope', '$location','$state', 'TeamService','OrganizationService',function($scope, $location,$state,TeamService,OrganizationService) {

     $scope.get = function() {
        OrganizationService.get().then(function(data){
            	 $scope.orgs = data;
            	 $scope.selectedOrg = $scope.orgs[0];
	             TeamService.get(orgid).then(function(data) {
	                $scope.teams = [{"name":"team1","desc":"desc1"}];
	            },
	            function(errMsg) {
	                console.log("error")
	                $scope.teams = [{"name":"team1","desc":"desc1"}];
	            });
        },function(errMsg){
            $scope.orgs = [{"name":"org1","desc":"desc1"}];
            $scope.selectedOrg = $scope.orgs[0];
            TeamService.get($scope.selectedOrg.name).then(function(data) {
	                $scope.teams = [{"name":"team1","desc":"desc1"}];
	            },
	            function(errMsg) {
	                console.log("error")
	                $scope.teams = [{"name":"team1","desc":"desc1"}];
	            });
        })
       
    };
    $scope.create = function(){
    	$state.go("team.create");
    }
    $scope.get();

}])
.controller('TeamCreateController', ['$scope', '$location','$state', 'TeamService','OrganizationService',function($scope, $location,$state,TeamService,OrganizationService) {
     OrganizationService.get().then(function(data){
        $scope.orgs = [{"name":"org1","desc":"desc1"}];
        $scope.selectedOrg = $scope.orgs.length > 0 ? $scope.orgs[0] : {"name":"","desc":""};
     },
     function(errMsg) {
        $scope.orgs = [{"name":"org1","desc":"desc1"}];
        $scope.selectedOrg = $scope.orgs.length > 0 ? $scope.orgs[0] : {"name":"","desc":""};
     })
     $scope.team = {name:"",description:"",type:"public"};
     $scope.step = "basicInfo";
     $scope.members = [];
     $scope.back = function(){
     	$state.go("team");
     }
     $scope.save = function(){
         TeamService.save($scope.team).then(function(data){
            console.log("create success")
         },
         function(errMsg){
         	console.log("create error")
         	$scope.team.id = 1;
         	$scope.step = "inviteMember";
         })
     }
     $scope.addMember = function(){
        
         TeamService.addMember($scope.team.id, $scope.newMember).then(function(data){
             TeamService.getMembers($scope.team.id).then(function(data){
                $scope.members.push($scope.newMember);
             },
             function(errMsg){
                $scope.members.push($scope.newMember);
             })
         },
         function(errMsg){
              $scope.members.push($scope.newMember);
              $scope.newMember="";
         })
     }

}])
