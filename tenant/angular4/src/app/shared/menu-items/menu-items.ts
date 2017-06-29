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

import { Injectable } from '@angular/core';

export interface BadgeItem {
  type: string;
  value: string;
}

export interface ChildrenItems {
  state: string;
  name: string;
  type?: string;
}

export interface Menu {
  state: string;
  name: string;
  type: string;
  icon: string;
  badge?: BadgeItem[];
  children?: ChildrenItems[];
}

const MENUITEMS = [
  {
    state: 'home',
    name: 'HOME',
    type: 'link',
    icon: 'explore'
  },
  {
    state: 'project',
    name: 'PROJECT',
    type: 'link',
    icon: 'device_hub'
  },
  {
    state: 'repo',
    name: 'REPOSITORIES',
    type: 'submenu',
    icon: 'apps',
    children: [
      {state: 'hub', name: 'REPOSHUB'},
      {state: 'git', name: 'REPOSGIT'},
      {state: 'binary', name: 'REPOSBINARY'},
      {state: 'docker', name: 'REPOSDOCKER'},
      {state: 'aci', name: 'REPOSACI'}
    ]
  },
  {
    state: 'flow',
    name: 'ORCHESTRATION',
    type: 'submenu',
    icon: 'timeline',
    children: [
      {state: 'overview', name: 'FLOWOVERVIEW'},
      {state: 'status', name: 'FLOWOSTATUS'}
    ]
  },
  {
    state: 'setting',
    name: 'SYSSETTING',
    type: 'link',
    icon: 'settings_applications'
  },
  {
    state: 'https://github.com/Huawei/containerops',
    name: 'DOCS',
    type: 'extraLink',
    icon: 'local_library'
  }
];

@Injectable()
export class MenuItems {
  getAll(): Menu[] {
    return MENUITEMS;
  }

  add(menu: Menu) {
    MENUITEMS.push(menu);
  }
}
