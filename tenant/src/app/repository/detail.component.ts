import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import {RepoDataService} from './data.service';

@Component({
  selector: 'repo-detail',
  template: require('../../template/repository/detail.html')
})

export class RepoDetailComponent implements OnInit { 
	repo;
	teams;

	constructor(private router: Router,
				private route: ActivatedRoute,
				private repoDataService: RepoDataService){

	}

	ngOnInit(): void {
		this.getRepo();
	}

	getRepo(): void {
		this.route.params.subscribe(params => {
	      let id = +params['id'];
	      this.repo = this.repoDataService.getRepo(id);
	    });

	    this.getTeams();
	}

	getTeams(): void{
		this.teams = this.repoDataService.getTeams(this.repo.id);
	}

	showNewTeam(): void{
		this.router.navigate(['/repository', this.repo.id, "addteam"]);
	}
}