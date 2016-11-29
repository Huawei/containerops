import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import {OrgDataService} from './data.service';

@Component({
  selector: 'list-org',
  template: require('../../template/organization/list.html')
})

export class OrgListComponent implements OnInit { 
	orgs = [];

	constructor(private router: Router, private orgDataService: OrgDataService){

	}

	ngOnInit(): void {
		this.getOrgs();
	}

	getOrgs(): void {
		this.orgs = this.orgDataService.getOrgs();
	}

	showNewOrg(): void{
		this.router.navigate(["organization/add"]);
	}	

	showOrgDetail(id): void{
		this.router.navigate(['/organization', id]);
	}
}