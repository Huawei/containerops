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

// import { Component, ElementRef, ViewChild, OnInit } from '@angular/core';
import { Http, Response } from '@angular/http';

// import * as D3 from 'd3';
import * as yaml from 'js-yaml';

export class Network {

  // getAll(): void {
  //   // return MENUITEMS;
  // }

  getYaml() {
    this.http.get('http://localhost:4200/assets/debug/cncf-demo.yaml').subscribe(data => this.extractData(data));
  }
  private extractData(res: Response) {
    this.workflowObj = yaml.load(res.text());
    console.log(this.workflowObj);
  }
}
