import { NgModule }      from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule }   from '@angular/forms';
import { HttpModule }    from '@angular/http';

import { IndexComponent }  from './index.component';

import { OrgListComponent }  from '../organization/list.component';
import { OrgAddComponent }  from '../organization/add.component';
import { OrgDetailComponent }  from '../organization/detail.component';
import { TeamListComponent } from '../team/list.component';
import { RepoListComponent }  from '../repository/list.component';

// services
import { OrgDataService } from '../organization/data.service';

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
    OrgListComponent,
    OrgAddComponent,
    OrgDetailComponent,
    TeamListComponent,
  	RepoListComponent
  ],
  providers: [
    OrgDataService
  ],
  bootstrap: [ IndexComponent ]
})
export class IndexModule { }