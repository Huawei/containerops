import { ModuleWithProviders }  from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { IndexComponent }  from './index.component';

import { OrgListComponent }  from '../organization/list.component';
import { TeamListComponent } from '../team/list.component';
import { RepoListComponent }  from '../repository/list.component';

const indexRouting: Routes = [
  {
    path: '',
    component: IndexComponent,
    children: [
      { path: 'organizations', component: OrgListComponent },
      { path: 'teams', component: TeamListComponent },
      { path: 'repositories', component: RepoListComponent }
    ]
  }
];

export const routing: ModuleWithProviders = RouterModule.forChild(indexRouting);