import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import {TeamDataService} from './data.service';
import {OrgDataService} from '../organization/data.service';
import { NotifyService } from '../common/notify.service';
// declare var Select2: any;
// interface JQuery {
//     select2(): void;
// }
var _ = require("underscore");
require('select2/dist/js/select2.min.js');
@Component({
  selector: 'add-team',
  template: require('../../template/team/add.html')
})


export class TeamAddComponent implements OnInit { 
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
	orgs = [];

	constructor(private router: Router, 
				private teamDataService: TeamDataService,
				private orgDataService : OrgDataService,
				private notifyService: NotifyService){

	}

	ngOnInit(): void {
      (<any>$('.select2')).select2({minimumResultsForSearch: Infinity});
       this.getOrgs();
	}

	addTeam(): void {
			this.caclAuth();
			this.caclOrg();
			this.teamDataService.addTeam(this.team);
			this.notifyService.notify("Add team '" + this.team.name + "' successfully.","success");
			this.router.navigate(["team"]);
	}

	cancelAdd(): void{
		this.router.navigate(["team"]);
	}	
    caclAuth(): void{
    	var self = this;
    	_.each(this.auth, function(value, key){
            if(value){
            	self.team.auth.push(key);
            }
    	})
    }
    caclOrg(): void{
    	var org = this.orgDataService.getOrg(this.team.org.id);
    	this.team.org.name = org.name;
    }
	getOrgs(): void{
        this.orgs = this.orgDataService.getOrgs();
	}
}