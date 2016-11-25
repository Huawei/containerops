import { Injectable } from '@angular/core';

var _ = require("underscore");

Injectable()
export class OrgDataService {
	getOrgs() {
		var orgs;
		if(_.isUndefined(localStorage["orgs"])){
			orgs = [];
		}else{
			orgs = JSON.parse(localStorage["orgs"])
		}
		return orgs;
	}

	getOrg(id: number) {
		var orgs = this.getOrgs();
	  	return _.find(orgs,function(org){
	  				return org.id == id;
	  			});
	}

	addOrg(org: any) {
		var orgs = this.getOrgs();
		org.id = orgs.length + 1;
		orgs.push(org);
	  	localStorage["orgs"] = JSON.stringify(orgs);
	}
}
