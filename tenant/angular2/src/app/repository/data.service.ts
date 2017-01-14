import { Injectable } from '@angular/core';

var _ = require("underscore");

Injectable()
export class RepoDataService {
	getRepos() {
		var repos;
		if(_.isUndefined(localStorage["repos"])){
			repos = [];
		}else{
			repos = JSON.parse(localStorage["repos"])
		}
		return repos;
	}

	getRepo(id: number) {
		var repos = this.getRepos();
	  	return _.find(repos,function(repo){
	  				return repo.id == id;
	  			});
	}

	addRepo(repo: any) {
		var repos = this.getRepos();
		repo.id = repos.length + 1;
		repos.push(repo);
	  	localStorage["repos"] = JSON.stringify(repos);
	}

	getTeams(id:number) {
		var teams;
		if(_.isUndefined(localStorage["repos_teams"])){
			teams = [];
		}else{
			var allrelations = JSON.parse(localStorage["repos_teams"]);
			var relations = _.filter(allrelations,function(relation){
				return relation.repoid == id;
			})

			var allteams;
			if(_.isUndefined(localStorage["teams"])){
				allteams = [];
			}else{
				allteams = JSON.parse(localStorage["teams"])
			}

			teams = _.map(relations,function(relation){
				var team = _.find(allteams,function(team){
					return team.id == relation.teamid;
				});

				team.repoauth = relation.auth;
				return team;
			});
		}

		return teams;
	}

	getAvailableTeams(id:number) {
		var teams,allrelations;

		if(_.isUndefined(localStorage["repos_teams"])){
			allrelations = [];
		}else{
			allrelations = JSON.parse(localStorage["repos_teams"]);
		}
		 
		var relations = _.filter(allrelations,function(relation){
			return relation.repoid == id;
		});

		var allteams;
		if(_.isUndefined(localStorage["teams"])){
			allteams = [];
		}else{
			allteams = JSON.parse(localStorage["teams"])
		}

		teams = _.filter(allteams,function(team){
			return _.pluck(relations,"teamid").indexOf(team.id) < 0;
		});

		return teams;
	}

	addTeam(relation: any) {
		var relations;
		if(_.isUndefined(localStorage["repos_teams"])){
			relations = [];
		}else{
			relations = JSON.parse(localStorage["repos_teams"])
		}
		relation.id = relations.length + 1;
		relations.push(relation);
	  	localStorage["repos_teams"] = JSON.stringify(relations);
	}
}
