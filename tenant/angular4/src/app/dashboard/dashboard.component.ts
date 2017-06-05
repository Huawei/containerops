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

import { Component, ElementRef, ViewChild, OnInit } from '@angular/core';
import { Http, Response } from '@angular/http';

import * as D3 from 'd3';
import * as yaml from 'js-yaml';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  @ViewChild('design') element: ElementRef;
  constructor(private http: Http) {};

  private host: D3.Selection;
  private svg: D3.Selection;
  private svgWidth: number;
  private svgHeight: number;

  private stageGroup: D3.Selection;

  private htmlElement: HTMLElement;
  private pieData = [1, 2, 3, 4, 5];

  private workflowObj: Object;

  getYaml() {
    this.http.get('http://localhost:4200/assets/debug/cncf-demo.yaml').subscribe(data => this.extractData(data));
  }
  private extractData(res: Response) {
    this.workflowObj = yaml.load(res.text());
  }

  ngOnInit() {

    this.getYaml();

    this.htmlElement = this.element.nativeElement;
    this.host = D3.select(this.htmlElement);

    this.buildSVG();
  }

  private buildSVG(): void {
    this.host.html('');
    this.svg = this.host.append('svg');
    this.svgWidth = this.htmlElement.offsetWidth;
    this.svgHeight = (this.htmlElement.parentElement.parentElement.offsetHeight) / 2 ;
    this.svg.attr('width', this.svgWidth).attr('height', this.svgHeight);

    this.stageGroup = this.svg.append('g');

    // init stages
    this.stageGroup.append('image')
      .attr('xlink:href', 'http://localhost:4200/assets/images/workflow/stage_start.svg')
      .attr('x', 60)
      .attr('y', 60)
      .attr('width', 40)
      .attr('height', 40);

  }

  private drawStage(): void {}

  private drawAction(): void {}

  private drawLink(): void {}

  private drawJob(): void {}



}
