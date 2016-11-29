import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import {TeamDataService} from './data.service';

@Component({
  selector: 'team-detail',
  template: require('../../template/team/detail.html')
})

export class TeamDetailComponent implements OnInit { 
	team;

	constructor(private router: Router,
				private route: ActivatedRoute,
				private teamDataService: TeamDataService){

	}

	ngOnInit(): void {
		this.getTeams();
	}


	getTeams(): void{
		this.route.params.subscribe(params => {
	      let id = +params['id'];
	      this.team = this.teamDataService.getTeam(id);
	    });

	}

}