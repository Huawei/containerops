import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import {OrgDataService} from './data.service';

var _ = require("underscore");

@Component({
  selector: 'org-detail',
  templateUrl: '../../template/organization/detail.html'
})

export class OrgDetailComponent implements OnInit { 
	org: any;

	constructor(private router: Router,
				private route: ActivatedRoute,
				private orgDataService: OrgDataService){

	}

	ngOnInit(): void {
		this.getOrg();
	}

	getOrg(): void {
		var org = this.route.params.subscribe(params => {
	      let id = +params['id'];
	      this.org = this.orgDataService.getOrg(id);
	    });
	}
}