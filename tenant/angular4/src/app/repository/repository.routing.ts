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

import { HubComponent } from './hub/hub.component';
import { BinaryComponent } from './binary/binary.component';
import { DockerComponent } from './docker/docker.component';
import { GitComponent } from './git/git.component';
import { AciComponent } from './aci/aci.component';
import { CreateComponent } from './create/create.component';

export const RepositoryRoutes: Routes = [{
  path: '',
  children: [{
    path: 'hub',
    component: HubComponent
  }, {
    path: 'create',
    component: CreateComponent
  }, {
    path: 'binary',
    component: BinaryComponent
  }, {
    path: 'git',
    component: GitComponent
  }, {
    path: 'docker',
    component: DockerComponent
  }, {
    path: 'aci',
    component: AciComponent
  }, {
    path: '**',
    component: HubComponent
  }]
}];
