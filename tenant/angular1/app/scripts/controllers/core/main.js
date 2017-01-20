auth.controller('CoreController', ['$scope', '$location', function($scope, $location) {

    $scope.sidebars = [
        {
            "name": "Organization",
            "href": "",
            "active": true,
            "icon": "fa fa-users",
            "collapsed":false,
            "children": [{
                "name": "Organization",
                "href": "organization"
            }, {
                "name": "Team",
                "href": "team"
            }, ]
        },
        {
            "name": "Project",
            "href": "",
            "active": false,
            "icon": "fa fa-cubes",
            "collapsed":true,
            "children": [{
                "name": "Project",
                "href": "project"
            }, {
                "name": "Application",
                "href": "application"
            }, {
                "name": "Module",
                "href": "module"
            }]
        }
    ]

    $scope.chooseNav = function(name) {
        _.each($scope.navigators, function(item) {
            if (item.name != name) {
                item.ngclass = "";
            } else {
                item.ngclass = "active";
            }
        })
    }

    function initializeNav() {
        var path = $location.path();
        _.each($scope.navigators, function(item) {
            if (path != "") {
                if (item.href == path) {
                    item.ngclass = "active";
                } else {
                    item.ngclass = "";
                }
            }
        })
    }
    $scope.toggle = function(name){
        // var selectedBar = _.find($scope.sidebars, function(item){return item.name === name});
        // selectedBar.collapsed = !selectedBar.collapsed;
         _.each($scope.sidebars, function(item){
            if(item.name === name){
                item.collapsed = !item.collapsed;
            }else{
                item.collapsed = true;
            }
        })

    }
    $scope.toggleParentStatus = function(name){
        _.each($scope.sidebars, function(item){
            if(item.name === name){
                item.active = true;
            }else{
                item.active = false;
            }
        })
    }
    // initializeNav();

}]);
