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

export class Workflow {
  name: string;
  tag: string;
  timeout: number;
  title: string;
  version: number;
  stages: array<Stage>;
}

export class Stage {
  name: string;
  sequencing: string;
  title: string;
  type: string;
  actions: array<Action>;
}

export class Action {
  name: string;
  title: string;
  jobs: array<Job>;
}

export class Job {
  endpoint: string;
  environments: array<any>;
  resources: array<any>;
  timeout: number;
  type: string;
}
