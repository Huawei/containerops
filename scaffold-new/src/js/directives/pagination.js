define(['app'], function(app) {
	'use strict';
	app.compileProvider.directive("paginations", function() {
		return {
			template: '<div style="text-align: center;margin-top: 12px;margin-bottom:12px;">' 
			+ '<div onselectstart="return false" class="button_type" ng-click= jump("first") title="first" ng-class="{\'button-disabled\':backButtonDisabled}"><span class="glyphicon glyphicon-backward" aria-hidden="true"></span></div>' 
			+ '<div onselectstart="return false" class="button_type" ng-click= jump("before") title="previous" ng-class="{\'button-disabled\':backButtonDisabled}"><span class="glyphicon glyphicon-triangle-left" aria-hidden="true"></span></div>' 
			+ '<div onselectstart="return false" ng-class="{true: \'choose\', false: \'no_choose\'}[choose[0]]" ng-click="jump(-4)" ng-if="length>=(middle-4)">{{middle-4}}</div>' 
			+ '<div onselectstart="return false" ng-class="{true: \'choose\', false: \'no_choose\'}[choose[1]]" ng-click="jump(-3)" ng-if="length>=(middle-3)">{{middle-3}}</div>' 
			+ '<div onselectstart="return false" ng-class="{true: \'choose\', false: \'no_choose\'}[choose[2]]" ng-click="jump(-2)" ng-if="length>=(middle-2)">{{middle-2}}</div>' 
			+ '<div onselectstart="return false" ng-class="{true: \'choose\', false: \'no_choose\'}[choose[3]]" ng-click="jump(-1)" ng-if="length>=(middle-1)">{{middle-1}}</div>' 
			+ '<div onselectstart="return false" ng-class="{true: \'choose\', false: \'no_choose\'}[choose[4]]" ng-click="jump(0)" ng-if="length>=(middle)">{{middle}}</div>' 
			+ '<div onselectstart="return false" ng-class="{true: \'choose\', false: \'no_choose\'}[choose[5]]" ng-click="jump(1)" ng-if="length>=(middle+1)">{{middle+1}}</div>' 
			+ '<div onselectstart="return false" ng-class="{true: \'choose\', false: \'no_choose\'}[choose[6]]" ng-click="jump(2)" ng-if="length>=(middle+2)">{{middle+2}}</div>' 
			+ '<div onselectstart="return false" ng-class="{true: \'choose\', false: \'no_choose\'}[choose[7]]" ng-click="jump(3)" ng-if="length>=(middle+3)">{{middle+3}}</div>' 
			+ '<div onselectstart="return false" ng-class="{true: \'choose\', false: \'no_choose\'}[choose[8]]" ng-click="jump(4)" ng-if="length>=(middle+4)">{{middle+4}}</div>' 
			+ '<div onselectstart="return false" ng-class="{true: \'choose\', false: \'no_choose\'}[choose[9]]" ng-click="jump(5)" ng-if="length>=(middle+5)">{{middle+5}}</div>' 
			+ '<div onselectstart="return false" class="button_type" ng-click=jump("next") title="next" ng-class="{\'button-disabled\':nextButtonDisabled}"><span class="glyphicon glyphicon-triangle-right" aria-hidden="true"></span></div>' 
			+ '<div onselectstart="return false" class="button_type" ng-click=jump("end") title="last" ng-class="{\'button-disabled\':nextButtonDisabled}"><span class="glyphicon glyphicon-forward" aria-hidden="true"></span></div>'
			+ '<div style="display:inline-block;margin-left:20px;font-size:12px;">{{"Total"}}&nbsp;&nbsp;{{total}}&nbsp;&nbsp;{{"Records"}}</div>' 
			+ '</div>',
			replace: true,
			restrict: 'AECM',
			scope: {
				present: '=',
				length: '=',
				total: '='
			},
			link: function(scope) {
				scope.$watch('present', function(newValue, oldValue, scope) {
					if(newValue == 1){
						scope.backButtonDisabled = true;
					}else{
						scope.backButtonDisabled = false;
					}
					if(newValue == scope.length){
						scope.nextButtonDisabled = true;
					}else{
						scope.nextButtonDisabled = false;
					}
					if (scope.present < 5 || scope.present > scope.length - 5) {
						if (scope.present < 5) {
							scope.middle = 5;
							var po = scope.choose.indexOf(true);
							scope.choose[po] = false;
							scope.choose[scope.present - 1] = true;
						} else {
							if (scope.length > 10) {
								scope.middle = scope.length - 5;
								var po = scope.choose.indexOf(true);
								scope.choose[po] = false;
								scope.choose[scope.present - (scope.length - 5) + 4] = true;
							} else {
								var po = scope.choose.indexOf(true);
								scope.choose[po] = false;
								scope.choose[scope.present - 1] = true;
							}
						}
					} else {
						scope.middle = scope.present;
						var po = scope.choose.indexOf(true);
						scope.choose[po] = false;
						scope.choose[4] = true;
					}
				});
				scope.present = 1;
				scope.button_num = 10;
				scope.middle = 5;
				scope.choose = new Array(scope.button_num);
				for (var i = 1; i < scope.button_num; i++)
					scope.choose[i] = false;
				scope.choose[0] = true;
				scope.jump = function(type) {
					var po = scope.choose.indexOf(true);
					switch (type) {
						case "first":
							if (scope.total == 0) {
								return
							} else {
								if (scope.present != 1) {
									scope.choose[po] = false;
									scope.choose[0] = true;
									scope.middle = 5;
								} else return;
								break;
							}
						case "before":
							if (scope.total == 0) {
								return
							} else {
								if (scope.present != 1) {
									if (scope.middle > 5) {
										if (po > 4 && scope.middle == (scope.length - 5)) {
											scope.choose[po] = false;
											scope.choose[po - 1] = true;
										} else
											scope.middle--;
									} else {
										if (po != 0) {
											scope.choose[po] = false;
											scope.choose[po - 1] = true;
										}
									}
								} else return;
								break;
							}
						case "next":
							if (scope.total == 0) {
								return
							} else {
								if (scope.present != scope.length) {
									if (scope.length >= 10) {
										if (scope.middle < (scope.length - 5)) {
											if (po < 4 && scope.middle == 5) {
												scope.choose[po] = false;
												scope.choose[po + 1] = true;
											} else
												scope.middle++;
										} else {
											if (po + 1 < scope.button_num) {
												scope.choose[po] = false;
												scope.choose[po + 1] = true;
											}
										}
									} else {
										if (po != scope.length - 1) {
											scope.choose[po] = false;
											scope.choose[po + 1] = true;
										} else return;
									}
								} else return;
								break;
							}
						case "end":
							if (scope.total == 0) {
								return
							} else {
								if (scope.present != scope.length) {
									scope.choose[po] = false;
									if (scope.length > scope.button_num) {
										scope.middle = scope.length - 5;
										scope.choose[scope.button_num - 1] = true;
									} else {
										scope.choose[scope.length - 1] = true;
									}
								} else return;
								break;
							}
						default:
							if (scope.length > 10) {
								if (po == 4) {
									if (scope.middle + type >= 5 && scope.middle + type < scope.length - 5) {
										scope.middle = scope.middle + type;
										if (scope.present == scope.middle) return;
									} else {
										if (scope.middle + type < 5) {
											scope.choose[scope.middle - 1 + type] = true;
											scope.middle = 5;
										} else {
											scope.choose[type - scope.length + 9 + scope.middle] = true;
											scope.middle = scope.length - 5;
										}
										if (scope.choose.lastIndexOf(true) != scope.choose.indexOf(true)) {
											scope.choose[po] = false;
										} else return;
									}
								}
								if (po < 4) {
									if (type < 0) {
										scope.choose[type + 4] = true;
									} else {
										scope.choose[4] = true;
										scope.middle = type + 5;
									}
									if (scope.choose.lastIndexOf(true) != scope.choose.indexOf(true)) {
										scope.choose[po] = false;
									} else return;
								}
								if (po > 4) {
									if (type < 0) {
										scope.choose[4] = true;
										scope.middle = scope.middle + type;
									} else {
										scope.choose[type + 4] = true;
									}
									if (scope.choose.lastIndexOf(true) != scope.choose.indexOf(true)) {
										scope.choose[po] = false;
									} else return;
								}
							} else {
								scope.choose[type + 4] = true;
								if (scope.choose.lastIndexOf(true) != scope.choose.indexOf(true)) {
									scope.choose[po] = false;
								} else return;
							}
							break;
					};
					scope.present = scope.choose.indexOf(true) + scope.middle - 4;
				}
			},

		}
	});

});