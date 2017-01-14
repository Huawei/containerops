import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { RepoDataService } from './data.service';
import { NotifyService } from '../common/notify.service';

var _ = require("underscore");

@Component({
  selector: 'repo-add-team',
  template: require('../../template/repository/addTeam.html')
})

export class RepoAddTeamComponent implements OnInit { 
	repo;
	availableteams;
	relation = {
		"teamid" : 0,
		"repoid" : 0,
		"auth":[]
	};
	auth = {
		"read":false,
		"write":false,
		'admin':false
	};

	constructor(private router: Router,
				private route: ActivatedRoute,
				private repoDataService: RepoDataService,
				private notifyService: NotifyService){

	}

	ngOnInit(): void {
		(<any>$('.select2')).select2({minimumResultsForSearch: Infinity})
      						.on("change",(e) => this.relation.teamid = $(e.target).val());
		this.getAvailableTeams();
	}

	getAvailableTeams(): void {
		this.route.params.subscribe(params => {
	      let id = +params['id'];
	      this.relation.repoid = id;
	      this.repo = this.repoDataService.getRepo(id);
	      this.availableteams = this.repoDataService.getAvailableTeams(id);
	      this.relation.teamid = this.availableteams[0].id;
	    });
	}

	addTeam(): void {
		try{
			//fake, to be deleted
			var self = this;
			this.caclAuth();
			this.repoDataService.addTeam(this.relation);
			this.notifyService.notify("Add team '" + _.find(this.availableteams,function(team){
				return team.id == self.relation.teamid;
			}).name + "' successfully.","success");

			this.router.navigate(['/repository', this.relation.repoid]);
			//fake end
		}catch(e){
			this.notifyService.notify("Fail to add team.","error");
		}
	}

	cancelAdd(): void{
		this.router.navigate(['/repository', this.relation.repoid]);
	}

	caclAuth(): void{
    	var self = this;
    	_.each(this.auth, function(value, key){
            if(value){
            	self.relation.auth.push(key);
            }
    	})
    }	
}