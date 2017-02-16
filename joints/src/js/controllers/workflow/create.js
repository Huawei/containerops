define(["app"], function(app) {
    app.controllerProvider.register('WorkflowCreateController', ['$scope', '$state', 'notifyService', function($scope, $state, notifyService) {
        $scope.backToList = function() {
            $state.go("workflow");
        };
        $scope.saveWorkflow = function(){
            $state.go("workflow");
        };
        $scope.drawWorkflow = function() {
            var svg = d3.select("#div-d3-main-svg")
                .append("svg")
                .attr("width", "100%")
                .attr("height", "100%");
            var g = svg.append("g");
            var svgMainRect = g.append("rect")
                .attr("width", "100%")
                .attr("height", "100%")
                .attr("fill", "white");
            var svgMainRect = g.append("circle")
                .attr("cx", 0) 
                .attr("cy", 0)
                .attr("r", 20)
                .attr("fill", "green")
                .attr("transform", function(d, i) {
                    return "translate(" +150 + "," + 150 + ")";
                })
                .attr("cursor","pointer")
                .on("click", function() {
                   notifyService.notify("click stage","success");
                });
        };

        $scope.drawWorkflow();




    }]);
})
