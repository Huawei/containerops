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
	updateTeam(users:Array<any>, team:any, type:string){
		let teams = JSON.parse(localStorage["teams"]);
	    let selected = _.find(teams, function(item){
	    	return item.id == team.id;
	    })
	    selected.users = selected.users || [];
	    if(type == "add"){
	    	selected.users = selected.users.concat(users);
	    }else{
	    	selected.users = _.difference(selected.users,users);
	    }
	    
	    localStorage["teams"] = JSON.stringify(teams);
	}
    getMembers(team:any) : Array<any> {
        let currentTeam  = this.getTeam(team.id);
        return currentTeam.users;
    }
    getOthers(all:Array<any>, contained: Array<any>): Array<any> {
        var others:Array<any> = [];
        _.each(all, function(item){
        	if(!_.contains(contained,item.username)){
              others.push(item); 
        	}
        })
       var others_1:Array<string> = _.map(others, function(item){return item.username;});
       return others_1;

    }

}
