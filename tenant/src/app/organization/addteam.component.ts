import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import {OrgDataService} from './data.service';
import {TeamDataService} from '../team/data.service';
import { NotifyService } from '../common/notify.service';

var _ = require("underscore");

@Component({
  selector: 'org-add-team',
  template: require('../../template/organization/addTeam.html')
})

export class OrgAddTeamComponent implements OnInit { 
	org;
	team = {
		"name" : "",
		"desc" : "",
		"org" : {"id":1,"name":""},
		"auth":[]
	};
	auth = {
		"read":false,
		"write":false,
		'admin':false
	};

	constructor(private router: Router,
				private route: ActivatedRoute,
				private orgDataService: OrgDataService,
				private teamDataService: TeamDataService,
				private notifyService: NotifyService){

	}

	ngOnInit(): void {
		this.getOrg();
	}

	getOrg(): void {
		this.route.params.subscribe(params => {
	      let id = +params['id'];
	      this.org = this.orgDataService.getOrg(id);
	    });

	    this.team.org.id = this.org.id;
	    this.team.org.name = this.org.name;
	}

	addTeam(): void {
		try{
			//fake, to be deleted
			this.caclAuth();
			this.teamDataService.addTeam(this.team);
			this.notifyService.notify("Add team '" + this.team.name + "' successfully.","success");

			this.router.navigate(['/organization', this.org.id]);
			//fake end
		}catch(e){
			this.notifyService.notify("Fail to add team.","error");
		}
	}

	cancelAdd(): void{
		this.router.navigate(['/organization', this.org.id]);
	}

	caclAuth(): void{
    	var self = this;
    	_.each(this.auth, function(value, key){
            if(value){
            	self.team.auth.push(key);
            }
    	})
    }	
}