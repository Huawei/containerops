import { Injectable } from '@angular/core';

var _ = require("underscore");

Injectable()
export class TeamDataService {
	getTeams(orgid = undefined) {
		var teams;
		if(_.isUndefined(localStorage["teams"])){
			teams = [];
		}else{
			teams = JSON.parse(localStorage["teams"])
		}

		if(orgid){
			teams = _.filter(teams,function(team){
				return team.orgid == orgid;
			});
		}

		return teams;
	}

	getTeam(id: number) {
		var teams = this.getTeams();
	  	return _.find(teams,function(team){
	  				return team.id == id;
	  			});
	}
    removeItem(collection:Array<any>, id:string){
        let index = _.indexOf(collection, id);
        collection.splice(index, 1);
    }
	addTeam(team: any) {
		var teams = this.getTeams();
		team.id = teams.length + 1;
		teams.push(team);
	  	localStorage["teams"] = JSON.stringify(teams);
	}
	updateTeam(users:Array<any>, team:any){
		let teams = JSON.parse(localStorage["teams"]);
	    let selected = _.find(teams, function(item){
	    	return item.id == team.id;
	    })
	    selected.users = users;
	    localStorage["teams"] = JSON.stringify(teams);
	}



}
