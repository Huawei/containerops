import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { TeamDataService } from './data.service';

@Component({
  selector: 'list-team',
  template: require('../../template/team/list.html')
})

export class TeamListComponent implements OnInit { 
	teams = [];
	constructor(private router: Router, private teamDataService: TeamDataService){

	}

	ngOnInit(): void {
		this.getTeams();
	}
	getTeams(): void {
		this.teams = this.teamDataService.getTeams();
	}

	showNewTeam(): void{
		this.router.navigate(["team/add"]);
	}	

	showTeamDetail(id): void{
		this.router.navigate(['team', id]);
	}

	manageMembers(id): void{
		this.router.navigate(['team', id, 'member','manage']);
	}
	manageAuth(id) : void{
		this.router.navigate(['team', id, 'auth', 'manage']);
	}
}