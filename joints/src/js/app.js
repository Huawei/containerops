define([
	'angular',
	'services/index',
	'directives/index'
	], function(angular) {
	return angular.module("DevOps", ['ui.router', 'ui.bootstrap','ngWebSocket', 'app.services','app.directives'
        ]);
})