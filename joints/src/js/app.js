define([
	'angular',
	'services/index'
	], function(angular) {
	return angular.module("DevOps", ['ui.router', 'ui.bootstrap','ngWebSocket', 'app.services'
        ]);
})