import { Component, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import {RepoDataService} from './data.service';

@Component({
  selector: 'list-repo',
  template: require('../../template/repository/list.html')
})

export class RepoListComponent implements OnInit { 
	repos = [];
	
	constructor(private router: Router, private repoDataService: RepoDataService){

	}

	ngOnInit(): void {
		this.getRepos();
	}

	getRepos(): void {
		this.repos = this.repoDataService.getRepos();
	}

	showNewRepo(): void{
		this.router.navigate(["repository/add"]);
	}	

	showRepoDetail(id): void{
		this.router.navigate(['/repository', id]);
	}
}