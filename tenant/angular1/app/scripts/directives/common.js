auth.directive('ngEnter', function () {
    return function (scope, element, attrs) {
        $(element).on("keydown keypress", function (event) {
            if(event.which === 13) {
                scope.$apply(function (){
                    scope.$eval(attrs.ngEnter);
                });

                event.preventDefault();
            }
        });
    };
})
.directive('select2', function() {
    return {
        restrict: 'ACEM',
        link: function(scope, element, attrs, ctrl) {
            $(element).select2({ minimumResultsForSearch: Infinity });
        }
    };
});