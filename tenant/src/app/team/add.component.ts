import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import {TeamDataService} from './data.service';
import {OrgDataService} from '../organization/data.service';
import { NotifyService } from '../common/notify.service';

@Component({
  selector: 'add-team',
  template: require('../../template/team/add.html')
})


export class TeamAddComponent implements OnInit { 
	team = {
		"name" : "",
		"desc" : "",
		"org" : 1
	};
	orgs = [];

	constructor(private router: Router, 
				private teamDataService: TeamDataService,
				private orgDataService : OrgDataService,
				private notifyService: NotifyService){

	}

	ngOnInit(): void {
       this.getOrgs();
	}

	addTeam(): void {
		try{
			//fake, to be deleted
			this.teamDataService.addTeam(this.team);
			this.notifyService.notify("Add team '" + this.team.name + "' successfully.","success");

			this.router.navigate(["team"]);
			//fake end
		}catch(e){
			this.notifyService.notify("Fail to add team.","error");
		}
	}

	cancelAdd(): void{
		this.router.navigate(["team"]);
	}	

	getOrgs(): void{
        this.orgs = this.orgDataService.getOrgs();
	}
}