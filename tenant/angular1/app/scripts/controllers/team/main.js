auth.controller('TeamController', ['$scope', '$location', '$state','$stateParams', 'TeamService', 'OrganizationService', function($scope, $location, $state, $stateParams, TeamService, OrganizationService) {
        $scope.selected = {orgId:""};
        $scope.init = function() {
            OrganizationService.get().then(function(data) {
                
            }, function(errMsg) {
                $scope.orgs = [{ "name": "org1", "desc": "desc1","id":"1"},{ "name": "org2", "desc": "desc2","id":"2" }];
                if($stateParams.orgId){
	               console.log("aaa"+ $stateParams.orgId);
	                $scope.setSelectedOrg($stateParams.orgId);
	        	}else{
	        		$scope.setSelectedOrg($scope.orgs[0].id);
	        	}
                
                $scope.get();
            })

        };
        $scope.get = function() {
            TeamService.get($scope.selected.orgId).then(function(data) {
                    $scope.teams = [{ "name": "team1", "desc": "desc1" }];
                },
                function(errMsg) {
                    $scope.teams = [{ "name": "team1", "desc": "desc1" }];
                    // console.log("from controller")
                    // console.log($scope.selected.orgId)
                });
        };
        $scope.create = function() {
            $state.go("team.create");
        }
        $scope.setSelectedOrg = function(item){
            $scope.selected.orgId = item;
        }
      
        $scope.$on("$stateChangeSuccess", function(){
        	$scope.init();

        })
        // $scope.init();

    }])
    .controller('TeamCreateController', ['$scope', '$location', '$state', 'TeamService', 'OrganizationService', function($scope, $location, $state, TeamService, OrganizationService) {
        OrganizationService.get().then(function(data) {
                $scope.orgs = [{ "name": "org1", "desc": "desc1" }];
                $scope.selectedOrg = $scope.orgs.length > 0 ? $scope.orgs[0] : { "name": "", "desc": "" };
            },
            function(errMsg) {
                $scope.orgs = [{ "name": "org1", "desc": "desc1" }];
                $scope.selectedOrg = $scope.orgs.length > 0 ? $scope.orgs[0] : { "name": "", "desc": "" };
            })
        $scope.team = { name: "", description: "", type: "public" };
        $scope.step = "basicInfo";
        $scope.members = [];
        $scope.newMember = "";
        // $scope.newMemberValidated = true;
        $scope.back = function() {
            $state.go("team");
        }
        $scope.save = function() {
            TeamService.save($scope.team).then(function(data) {
                    console.log("create success")
                    $scope.team = {};
                },
                function(errMsg) {
                    console.log("create error")
                    $scope.team.id = 1;
                    $scope.step = "inviteMember";
                })
        };
        $scope.saveMember = function() {
            if (!_.isEmpty($scope.newMember)) {
                // $scope.newMemberValidated = false;
                // } else {
                TeamService.saveMember($scope.team.id, $scope.newMember).then(function(data) {
                        TeamService.getMembers($scope.team.id).then(function(data) {
                                $scope.members.push($scope.newMember);
                            },
                            function(errMsg) {
                                $scope.members.push($scope.newMember);
                            })
                    },
                    function(errMsg) {
                        $scope.members.push($scope.newMember);
                        $scope.newMember = "";
                    })
            }

        };

        $scope.removeMember = function(member) {
            TeamService.removeMember($scope.team.id, member).then(function(data) {
                    TeamService.getMembers($scope.team.id).then(function(data) {
                            $scope.members = _.without($scope.members, member);
                        },
                        function(errMsg) {
                            $scope.members = _.without($scope.members, member);
                        })
                },
                function(errMsg) {
                    $scope.members = _.without($scope.members, member);
                })
        };
        $scope.gotoCreate = function() {
            $scope.step = "basicInfo";
            $scope.team = { name: "", description: "", type: "public" };
        }

    }])
