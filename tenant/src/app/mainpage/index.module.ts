import { NgModule }      from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule }   from '@angular/forms';
import { HttpModule }    from '@angular/http';

import { IndexComponent }  from './index.component';

// Directive
import { CheckAllDirective } from '../common/check-all.directive';

// org
import { OrgListComponent }  from '../organization/list.component';
import { OrgAddComponent }  from '../organization/add.component';
import { OrgDetailComponent }  from '../organization/detail.component';
import { OrgAddTeamComponent }  from '../organization/addteam.component';

// team
import { TeamListComponent } from '../team/list.component';
import { TeamAddComponent }  from '../team/add.component';
import { TeamDetailComponent }  from '../team/detail.component';
import { TeamManageMemberComponent } from '../team/manage-member.component';
import { TeamManageAuthComponent } from '../team/manage-auth.component';
// repo
import { RepoListComponent }  from '../repository/list.component';
import { RepoAddComponent }  from '../repository/add.component';
import { RepoDetailComponent }  from '../repository/detail.component';
import { RepoAddTeamComponent }  from '../repository/addteam.component';

// services
import { OrgDataService } from '../organization/data.service';
import { TeamDataService } from '../team/data.service';
import { RepoDataService } from '../repository/data.service';

import { routing } from './index.router';


@NgModule({
  imports: [ 
  	BrowserModule,
    FormsModule,
    HttpModule,
  	routing
  ],
  declarations: [ 
    IndexComponent,
    CheckAllDirective,
    OrgListComponent,
    OrgAddComponent,
    OrgDetailComponent,
    OrgAddTeamComponent,
    TeamListComponent,
    TeamAddComponent,
    TeamDetailComponent,
    TeamManageMemberComponent,
    TeamManageAuthComponent,
  	RepoListComponent,
    RepoAddComponent,
    RepoDetailComponent,
    RepoAddTeamComponent
  ],
  providers: [
    OrgDataService,
    TeamDataService,
    RepoDataService
  ],
  bootstrap: [ IndexComponent ]
})
export class IndexModule { }