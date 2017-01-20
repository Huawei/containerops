auth.controller('OrganizationController', ['$scope', '$location', '$state', 'OrganizationService', function($scope, $location, $state, OrganizationService) {
        $scope.get = function() {
            OrganizationService.get().then(function(data) {
                    // $scope.orgs = data.data;
                },
                function(errMsg) {
                    console.log("error")
                    $scope.orgs = [{ "name": "org1", "desc": "desc1", "id": "1" }, { "name": "org2", "desc": "desc2", "id": "2" }];
                });
        };
        $scope.create = function() {
            $state.go("organization.create");
        }
        $scope.showTeams = function(orgId) {
            $state.go("team", { "orgId": orgId });
        }

        $scope.get();


    }])
    .controller('OrganizationCreateController', ['$scope', '$location', '$state', 'NotifyService', 'OrganizationService', 'TeamService', function($scope, $location, $state, NotifyService, OrganizationService, TeamService) {
        $scope.org = { name: "", description: "", type: "public" };
        $scope.team = { name: "", description: "", type: "public" };
        $scope.newMember = "";
        $scope.members = [];
        // $scope.step = "basicInfo";
        $scope.steps = [
            { "name": "basicInfo", "index": 1 },
            { "name": "setupTeam", "index": 2 },
            { "name": "inviteMember", "index": 3 }
        ];
        // $scope.step = $scope.steps[0];
        $scope.step = "basicInfo";
        $scope.currentOrg = {};
        $scope.currentTeam = {};
        $scope.back = function() {
            $state.go("organization");
        }
        $scope.save = function(next) {
            OrganizationService.save($scope.org).then(function(data) {
                    console.log("create success")
                        // $scope.currentOrg = $scope.org;
                        // $scope.currentOrg.id = 1;
                        // $scope.step = "setupTeam";
                },
                function(errMsg) {
                    NotifyService.notify("save organization success");
                    if (next) {
                        $scope.currentOrg = $scope.org;
                        $scope.currentOrg.id = 1;
                        // $scope.step = "setupTeam";
                        $scope.step = next;
                    } else {
                        $scope.org = { name: "", description: "", type: "public" };
                    }

                })
        }

        $scope.saveTeam = function(next) {
            TeamService.save($scope.team).then(function(data) {

                },
                function(errMsg) {
                    NotifyService.notify("save team success");
                    if (next) {
                        $scope.currentTeam = $scope.team;
                        $scope.currentTeam.id = 1;
                        $scope.step = next;
                    }else{
                        $scope.team = { name: "", description: "", type: "public" };
                    }

                })
        }
        $scope.saveMember = function() {
            if (!_.isEmpty($scope.newMember)) {
                TeamService.saveMember($scope.currentTeam.id).then(function(data) {
                        $scope.members.push($scope.newMember);
                        $scope.newMember = "";
                    },
                    function(errMsg) {
                        $scope.members.push($scope.newMember);
                        $scope.newMember = "";
                    })
            }
        }
        $scope.gotoStep = function(stepName) {
            var clickedStep = _.find($scope.steps, function(item) {
                return item.name == stepName; });
            if (clickedStep.index <= $scope.step.index) {
                switch (stepName) {
                    case "basicInfo":
                        $scope.org = { name: "", description: "", type: "public" };
                    default:
                        $scope.team = { name: "", description: "", type: "public" };
                        $scope.step = clickedStep;
                }
            }

        }




    }]);
