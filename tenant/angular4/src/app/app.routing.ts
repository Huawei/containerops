/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import { Routes } from '@angular/router';

import { CoreLayoutComponent } from './layouts/core/core-layout.component';

export const AppRoutes: Routes = [{
  path: '',
  redirectTo: 'home',
  pathMatch: 'full',
}, {
  path: '',
  component: CoreLayoutComponent,
  children: [{
    path: 'home',
    loadChildren: './dashboard/dashboard.module#DashboardModule'
  }, {
    path: 'project',
    loadChildren: './project/project.module#ProjectModule'
  }, {
    path: 'repo',
    loadChildren: './repository/repository.module#RepositoryModule'
  }, {
    path: 'flow',
    loadChildren: './flow/flow.module#FlowModule'
  }]
}, {
  path: '**',
  redirectTo: './dashboard/dashboard.module#DashboardModule'
}];
