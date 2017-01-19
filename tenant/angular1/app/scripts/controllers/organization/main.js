auth.controller('OrganizationController', ['$scope', '$location','$state', 'OrganizationService', function($scope, $location, $state, OrganizationService) {
    $scope.get = function() {
        OrganizationService.get().then(function(data) {
                // $scope.orgs = data.data;
                $scope.orgs = [{"name":"org1","desc":"desc1"}];
            },
            function(errMsg) {
                console.log("error")
                $scope.orgs = [{"name":"org1","desc":"desc1"}];
            });
    };
    $scope.create = function(){
    	$state.go("organization.create");
    }
    $scope.get();


}])
 .controller('OrganizationCreateController', ['$scope', '$location','$state', 'OrganizationService','TeamService', function($scope, $location, $state, OrganizationService,TeamService) {
     $scope.org = {name:"",description:"",type:"public"};
     $scope.team = {name:"",description:"",type:"public"};
     $scope.newMember = "";
     $scope.members = [];
     $scope.step = "basicInfo";
     $scope.currentOrg = {};
     $scope.currentTeam = {};
     $scope.back = function(){
     	$state.go("organization");
     }
     $scope.save = function(){
         OrganizationService.save($scope.org).then(function(data){
            console.log("create success")
            $scope.currentOrg = $scope.org;
            $scope.currentOrg.id = 1;
            $scope.step = "setupTeam";
         },
         function(errMsg){
         	console.log("create error")
         	$scope.currentOrg = $scope.org;
            $scope.currentOrg.id = 1;
         	$scope.step = "setupTeam";
         })
     }
     $scope.saveTeam = function(){
     	TeamService.save($scope.team).then(function(data){
            $scope.currentTeam = $scope.team;
            $scope.currentTeam.id = 1;
         	$scope.step = "inviteMember";
     	},
     	function(errMsg){
            $scope.currentTeam = $scope.team;
            $scope.currentTeam.id = 1;
         	$scope.step = "inviteMember";
     	})
     }
     $scope.saveMember = function(){
        TeamService.saveMember($scope.currentTeam.id).then(function(data){
            $scope.members.push($scope.newMember);
            $scope.newMember = "";
        },
        function(errMsg){
            $scope.members.push($scope.newMember);
            $scope.newMember = "";
        })
     }

}]);
