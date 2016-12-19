import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import {OrgDataService} from './data.service';
import {TeamDataService} from '../team/data.service';

@Component({
  selector: 'org-detail',
  template: require('../../template/organization/detail.html')
})

export class OrgDetailComponent implements OnInit { 
	org;
	teams;

	constructor(private router: Router,
				private route: ActivatedRoute,
				private orgDataService: OrgDataService,
				private teamDataService: TeamDataService){

	}

	ngOnInit(): void {
		this.getOrg();
	}

	getOrg(): void {
		this.route.params.subscribe(params => {
	      let id = +params['id'];
	      this.org = this.orgDataService.getOrg(id);
	    });

	    this.getTeams();
	}

	getTeams(): void{
		this.teams = this.teamDataService.getTeams(this.org.id);
	}

	showNewTeam(): void{
		this.router.navigate(['/organization', this.org.id, "addteam"]);
	}
}