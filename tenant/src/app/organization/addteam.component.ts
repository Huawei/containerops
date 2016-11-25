import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import {OrgDataService} from './data.service';
import {TeamDataService} from '../team/data.service';

@Component({
  selector: 'org-add-team',
  templateUrl: '../../template/organization/addTeam.html'
})

export class OrgAddTeamComponent implements OnInit { 
	org;
	team = {
		"name" : "",
		"desc" : ""
	};

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
	}

	addTeam(): void {
		try{
			//fake, to be deleted
			this.team["orgid"] = this.org.id;
			this.teamDataService.addTeam(this.team);
			alert("team "+ this.team.name + " added.");

			this.router.navigate(['/organization', this.org.id]);
			//fake end
		}catch(e){
			alert("failed to add team.")
		}
	}

	cancelAdd(): void{
		this.router.navigate(['/organization', this.org.id]);
	}	
}