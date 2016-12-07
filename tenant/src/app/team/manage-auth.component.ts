import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { TeamDataService } from './data.service';
import { UserService } from '../user/user.service';
import { NotifyService } from '../common/notify.service';

var _ = require("underscore");
// declare _ 
@Component({
  selector: 'team-manage-auth',
  template: require('../../template/team/manage-auth.html')
})

export class TeamManageAuthComponent implements OnInit { 
    team;
	users;
	members;
	others;
    usersToBeAdd = [];
    usersToBeRemove = [];
	constructor(private router: Router,
				private route: ActivatedRoute,
				private teamDataService: TeamDataService,
				private userService: UserService,
				private notifyService: NotifyService){

	}

	ngOnInit(): void {
		this.getTeam();
		this.getUsers();
		this.getMembers();
		this.getOthers();
	}


	getTeam() {
		this.route.params.subscribe(params => {
			let id = +params['id'];
			this.team = this.teamDataService.getTeam(id);
		})
	}
	getUsers(): void {
		this.users = this.userService.getUsers();
	}
	getMembers(): void{
		this.members = this.teamDataService.getMembers(this.team);
	}
	getOthers() : void{
		this.others = this.teamDataService.getOthers(this.users, this.members);
	}
	
    toggleMember(event:any, opt:string, id = undefined){
        let checked = event.currentTarget.checked;
        if(checked){
        	if(id == undefined){
        		if(opt == "add"){
        			this.usersToBeAdd = this.others;
        		}else if(opt == "remove"){
        			this.usersToBeRemove = this.members;
        		}
        		
        	}else {
        		if(opt == "add"){
        			this.usersToBeAdd.push(id);
        		}else if(opt == "remove"){
                    this.usersToBeRemove.push(id);
        		}
        		
        	}
        	
        }else {
            if(id == undefined){
        		if(opt == "add"){
        			this.usersToBeAdd = [];
        		}else if(opt == "remove"){
        			this.usersToBeRemove = [];
        		}
        		
        	}else {
        		if(opt == "add"){
        			this.teamDataService.removeItem(this.usersToBeAdd, id);
        		}else if(opt == "remove"){
        			this.teamDataService.removeItem(this.usersToBeRemove, id);
        		}
        		
        	}

        }
    }

	updateMembers(opt:string) : void{
		var targetArray = opt == "add" ? this.usersToBeAdd : this.usersToBeRemove;
		var msg = opt == "add" ? "Add member to " : "Remove member from ";
			
		try{
			//fake, to be deleted
			
			this.teamDataService.updateTeam(targetArray, this.team, opt);
			this.notifyService.notify(msg + this.team.name + "' successfully.","success");
			targetArray = [];
            if(opt == "add"){
            	this.usersToBeAdd = [];
            }else {
            	this.usersToBeRemove = [];
            }
			this.getMembers();
			this.getOthers();
			//fake end
		}catch(e){
			this.notifyService.notify("Fail to update member.","error");
			targetArray = [];
			if(opt == "add"){
            	this.usersToBeAdd = [];
            }else {
            	this.usersToBeRemove = [];
            }
		}
	}
    clickToOperate(user:any, opt:string){

       var targetArray = opt == "add" ? this.usersToBeAdd : this.usersToBeRemove;
       	   if(!_.contains(targetArray, user)){
       	   	  targetArray.push(user);
       	   }
       	   
       this.updateMembers(opt);
    }
    
	cancelAdd(): void{
		this.router.navigate(['team']);
	}	
}