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

import { Component, OnInit } from '@angular/core';

import { Http, Response } from '@angular/http';
import 'rxjs/add/operator/toPromise';
// import * as D3 from 'd3';
import * as yaml from 'js-yaml';

@Component({
  selector: 'app-status',
  templateUrl: './status.component.html',
  styleUrls: ['./status.component.scss']
})
export class StatusComponent implements OnInit {
  constructor(private http: Http) {};

  private workflowObj: Object;

  ngOnInit() {
    this.http.get('http://localhost:4200/assets/debug/cncf-demo.yaml')
      .toPromise()
      .then(response => {this.workflowObj = yaml.load(response.text()); return this.workflowObj; })
      .then(workflowData => {console.log('------------------');  console.log(workflowData); })
      ;
      // .subscribe(data => {
      //   this.workflowObj = yaml.load(data.text());
      //   console.log(this.workflowObj);
      // });
  }


}
