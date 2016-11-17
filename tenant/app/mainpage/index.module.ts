import { NgModule }      from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

// import { ContentComponent }  from './content.component';
// import { NavComponent } from './nav.component';
// import { RepoListComponent }  from './repoList.component';
// import { RepositoriesComponent }  from './repositories.component';
// import { RepoCreateComponent }  from './repoCreate.component';
// import { RepoDetailComponent } from './repoDetail.component';
// import { OrgListComponent }      from './orgList.component';
// import { OrgCreateComponent }      from './orgCreate.component';
// import { OrgEditComponent }      from './orgEdit.component';
// import { UserSettingComponent }      from './userSetting.component';
// import { PromptComponent }  from './prompt.component';

// import { PromptModule }  from './prompt.module';

// import { OrgService }      from './org.service';
// import { RepoService }      from './repo.service';

import { FormsModule }   from '@angular/forms';
import { HttpModule }    from '@angular/http';
// import { routing } from './content.routing';


@NgModule({
  imports: [ 
  	BrowserModule,
    FormsModule,
    HttpModule,
  	// routing
  ],
  declarations: [ 
  	// ContentComponent,
   //  NavComponent,
  	// RepoListComponent,
   //  RepositoriesComponent, 
  	// RepoCreateComponent,
   //  RepoDetailComponent,
  	// OrgListComponent,
  	// OrgCreateComponent,
   //  OrgEditComponent,
   //  UserSettingComponent,
   //  PromptComponent
  ],
  providers: [
    // OrgService,
    // RepoService
  ],
  // bootstrap: [ ContentComponent ]
})
export class IndexModule { }