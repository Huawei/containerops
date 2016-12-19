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
}
