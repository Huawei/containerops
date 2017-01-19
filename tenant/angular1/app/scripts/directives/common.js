auth.directive('ngEnter', function() {
        return function(scope, element, attrs) {
            $(element).on("keydown keypress", function(event) {
                if (event.which === 13) {
                    scope.$apply(function() {
                        scope.$eval(attrs.ngEnter);
                    });

                    event.preventDefault();
                }
            });
        };
    })
    .directive('select2', function($timeout) {
        return {
            restrict: 'ACEM',
            link: function(scope, element, attrs, ctrl) {
                $(element).select2({ minimumResultsForSearch: Infinity });
                $(element).on("change", function() {
                    if (attrs.ngSelect) {
                        scope[attrs.ngSelect].call(this, $(element).val());
                    }
                    if (attrs.ngChange) {
                        scope[attrs.ngChange].call();
                    }

                })
            }
        };
    })
    .directive('autofocus', ['$timeout', function($timeout) {
        return {
            restrict: 'ACEM',
            link: function(scope, element, attrs, ctrl) {
                $timeout(function() {
                    element[0].focus();
                });
            }
        };
    }]);
