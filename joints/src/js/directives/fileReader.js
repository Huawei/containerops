define(['directives/module'], function(dirmodule) {
    dirmodule.directive("ngFileModel", [function () {
        return {
            scope: {
                ngFileModel: "="
            },
            link: function (scope, element, attributes) {
                element.bind("change", function (changeEvent) {
                    var reader = new FileReader();
                    reader.onload = function (loadEvent) {
                        scope.$apply(function () {
                            scope.ngFileModel = JSON.parse(window.atob(loadEvent.target.result.replace('data:application/json;base64,', '')));
                        });
                    }
                    reader.readAsDataURL(changeEvent.target.files[0]);
                });
            }
        }
    }]);
});