import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';

var _ = require("underscore");

@Component({
  selector: 'list-org',
  templateUrl: '../../template/organization/list.html'
})

export class OrgListComponent implements OnInit { 
	orgs = [];

	constructor(private router: Router){

	}

	ngOnInit(): void {
		this.getOrgs();
	}

	getOrgs(): void {
		if(_.isUndefined(localStorage["orgs"])){
			this.orgs = [];
		}else{
			this.orgs = JSON.parse(localStorage["orgs"])
		}
	}

	showNewOrg(): void{
		this.router.navigate(["organization/add"]);
	}	
}