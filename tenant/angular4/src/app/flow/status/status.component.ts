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

import { Component, OnInit, AfterViewInit, OnDestroy, ElementRef, ViewChild } from '@angular/core';

import { Http, Response } from '@angular/http';
import 'rxjs/add/operator/toPromise';

import * as D3 from 'd3';
import * as yaml from 'js-yaml';

import * as workflow from './workflow';

const stageTypeStar = 'start';
const stageTypeEnd = 'end';
const stageTypeNormal = 'normal';
const stageSequencingSequence = 'sequence';
const stageSequencingParallel = 'parallel';


@Component({
  selector: 'app-status',
  templateUrl: './status.component.html',
  styleUrls: ['./status.component.scss']
})
export class StatusComponent implements OnInit, AfterViewInit, OnDestroy {
  @ViewChild('design') element: ElementRef;
  constructor(private http: Http) {};

  private htmlElement: HTMLElement;

  private host: D3.Selection;
  private svg: D3.Selection;
  private svgWidth: number;
  private svgHeight: number;

  private stageGroup: D3.Selection;
  private stageLineGroup: D3.Selection;
  private actionGroups: Map<string, D3.Selection> = new Map();
  private jobGroups: Map<string, D3.Selection> = new Map();


  private workflowObj: workflow.Workflow;
  private refreshInterval: any;
  private refreshIntervalHttp: any;

  nowLogs: Array<string>;


  tabs = [{
    label: 'Tab 1',
    content: 'This is the body of the first tab'
  }, {
    label: 'Tab 2',
    content: 'This is the body of the second tab'
  }, {
    label: 'Tab 3',
    extraContent: true,
    content: 'This is the body of the third tab'
  }, {
    label: 'Tab this',
    content: 'This is the body of the fourth tab'
  }];


  ngOnInit() {
    this.refreshIntervalHttp = this.http;
    this.http.get('http://localhost:4200/flow/v1/cncf/cncf-demo/aaa/latest/1/runtime/yaml')
      .toPromise()
      .then(response => yaml.load(response.text()))
      .then(wfObj => {
        this.htmlElement = this.element.nativeElement;
        this.host = D3.select(this.htmlElement);
        this.host.html('');
        this.svg = this.host.append('svg');
        this.svgWidth = 1600;//this.htmlElement.offsetWidth;
        this.svgHeight = 500;//(this.htmlElement.parentElement.parentElement.offsetHeight) / 2 ;
        this.svg.attr('width', this.svgWidth).attr('height', this.svgHeight);

        this.stageGroup = this.svg.append('g');
        this.stageGroup.attr('id', 'stageGroup');

        this.stageLineGroup = this.svg.append('g');
        this.stageLineGroup.attr('id', 'stageLineGroup');

        this.drawWorkflow(wfObj);
      });
  }

  drawWorkflow(wfObj: any): void {
    this.workflowObj = wfObj;
    wfObj.stages.forEach((stageValue, stageIndex) => {
      // console.log(stageValue);
      let stageImageUrl = '';

      switch (stageValue.type) {
        case stageTypeStar :
          stageImageUrl = 'http://localhost:4200/assets/images/workflow/stage_start.svg';
          break;
        case stageTypeEnd :
          stageImageUrl = 'http://localhost:4200/assets/images/workflow/stage_end.svg';
          break;
        case stageTypeNormal :
          if (stageSequencingSequence === stageValue.sequencing) {
            stageImageUrl = 'http://localhost:4200/assets/images/workflow/stage_sequnce.svg';
          }else if (stageSequencingParallel === stageValue.sequencing) {
            stageImageUrl = 'http://localhost:4200/assets/images/workflow/stage_parallel.svg';
          }
          break;
      };

      // Draw Stage
      this.stageGroup.append('image')
        .attr('xlink:href', stageImageUrl)
        .attr('x', 140 * (stageIndex + 1) + (stageIndex * 100))
        .attr('y', 20 )
        .attr('width', 40)
        .attr('height', 40);

      // Draw Stage line
      if (stageValue.type !== stageTypeEnd) {
        this.stageLineGroup.append('line')
          .attr('x1', 140 * (stageIndex + 1) + (stageIndex * 100) + 40)
          .attr('y1', 40)
          .attr('x2', 140 * (stageIndex + 1) + (stageIndex * 100) + 240)
          .attr('y2', 40)
          .attr('stroke', '#adadad')
          .attr('stroke-width', '5');
      }

      // Draw Action

      if (stageValue.actions) {
        let jobRowCount, jobColCount: number;


        const actionGroupName = stageValue.name +  '-Actions';
        this.actionGroups.set(actionGroupName, this.svg.append('g'));
        this.actionGroups.get(actionGroupName).attr('id', actionGroupName);

        stageValue.actions.forEach((actionValue, actionIndex) => {

          jobRowCount = Math.ceil( actionValue.jobs.length / 4 ) ;
          // jobColCount = actionValue.jobs.length > 4 ? 4 : actionValue.jobs.length;
          // console.log('===rowCount===');
          // console.log( jobRowCount + ':' + jobColCount);

          const actonXBase = 140 * (stageIndex + 1) + (stageIndex * 100);

          this.actionGroups.get(actionGroupName).append('rect')
            .attr('x', actonXBase - 40 )
            .attr('y', (actionIndex + 1) * 150)
            .attr('width', 120)
            .attr('height', 40)
            .attr('rx', 5)
            .attr('ry', 5)
            .attr('stroke', '#000000')
            .attr('stroke-width', '1')
            .attr('fill', '#000000')
            .attr('fill-opacity', '0');

          // let mX, mY, lX, lY, c1x, c1y, c2x, c2y, c3x, c3y, elX, elY, actionLinePath;
          //
          // mX = 202 + 50 + (240 * stageIndex);
          // mY = 40;
          // lX = 202 + 50 + (240 * stageIndex);
          // lY = 75 + 50 + (150 * actionIndex) + 30;
          //
          // c1x = 202 + 50 + (240 * stageIndex);
          // c1y = 83 + 50 + (150 * actionIndex) + 30;
          //
          // c2x = 195 + 50 + (240 * stageIndex);
          // c2y = 90 + 50 + (150 * actionIndex) + 30;
          //
          // c3x = 187 + 50 + (240 * stageIndex);
          // c3y = 90 + 50 + (150 * actionIndex) + 30;
          //
          // elX = 150 + 70 + (240 * stageIndex);
          // elY = 90 + 50 + (150 * actionIndex) + 30;
          //
          // actionLinePath = 'M' + mX + ',' + mY
          //   + ' L' + lX + ',' + lY
          //   + ' C' + c1x + ',' + c1y + ' ' + c2x + ',' + c2y + ' ' + c3x + ',' + c3y
          //   + ' L' + elX + ',' + elY ;
          //
          // this.actionGroups.get(actionGroupName).append('path')
          //   .attr('d', actionLinePath)
          //   // .attr('d', 'M202,15 L202,75 C202,83 195,90 187,90 L150,90' )
          //   .attr('stroke', '#adadad')
          //   .attr('stroke-width', '1')
          //   .attr('fill', '#000000')
          //   .attr('fill-opacity', '0');

          let amX, amY, alX, alY, actionLinkPath;
          amX = 140 * (stageIndex + 1) + (stageIndex * 100) + 20;
          amY = 150 + (150 * actionIndex);
          alX = 140 * (stageIndex + 1) + (stageIndex * 100) + 20;
          alY = 150 + (150 * actionIndex) - (actionIndex === 0 ? 90 : 110);
          actionLinkPath = 'M' + amX + ',' + amY
            + ' L' + alX + ',' + alY;

          this.actionGroups.get(actionGroupName).append('path')
            .attr('d', actionLinkPath)
            // .attr('d', 'M202,15 L202,75' )
            .attr('stroke', '#adadad')
            .attr('stroke-width', '1')
            .attr('fill', '#000000')
            .attr('fill-opacity', '0');

          // Draw Job
          // console.log(actionValue.action.jobs);
          actionValue.jobs.forEach((jobValue, jobIndex) => {

            // console.log('=====jobValue.status=====');
            // console.log(jobValue.status);
            // console.log('=========================');
            let jobColor = '#000000';
            switch (jobValue.status) {
              case 'pending' :
                jobColor = '#000000';
                break;
              case 'running' :
                jobColor = '#ffff00';
                break;
              case 'failure' :
                jobColor = '#ff0000';
                break;
              case 'success' :
                jobColor = '#00ff00';
                break;
            };


            const jobGroupName = stageValue.name + '-' + actionValue.name + '-job';
            this.jobGroups.set(jobGroupName, this.svg.append('g'));
            this.jobGroups.get(jobGroupName).attr('id', jobGroupName);
            this.jobGroups.get(jobGroupName).append('rect')
              .attr('x', actonXBase - 30 + (jobIndex * 27))
              .attr('y', (actionIndex + 1) * 150 + 10 )
              .attr('width', 20)
              .attr('height', 20)
              .attr('rx', 2)
              .attr('ry', 2)
              .attr('stroke', '#000000')
              .attr('stroke-width', '1')
              .attr('fill', jobColor)
              .attr('id', stageIndex + '::' + actionIndex + '::' + jobIndex)
              .on('click', (d, i, t) => {
                const index = t[0].id.split('::');
                // console.log(index);
                const josInfo = this.workflowObj.stages[index[0]].actions[index[1]].jobs[index[2]];

                if (josInfo.logs) {
                  this.nowLogs = josInfo.logs;
                }else {
                  this.nowLogs = null;
                }

                // console.log(josInfo.logs[0]);
                // let nowLog = '';
                // // josInfo['logs'][0].map((logv, logk) => {
                // //   nowLog += logv + logk;
                // // });
                // console.log(nowLog);
              });
              // .attr('fill-opacity', '0');
          });


        });
      }

    });


  }

  ngAfterViewInit() {
    this.refreshInterval = setInterval(() => {
      this.svg.remove();

      this.htmlElement = this.element.nativeElement;
      this.host = D3.select(this.htmlElement);
      this.host.html('');
      this.svg = this.host.append('svg');
      this.svgWidth = 1600;//this.htmlElement.offsetWidth;
      this.svgHeight = 500;//(this.htmlElement.parentElement.parentElement.offsetHeight) / 2 ;
      this.svg.attr('width', this.svgWidth).attr('height', this.svgHeight);

      this.stageGroup = this.svg.append('g');
      this.stageGroup.attr('id', 'stageGroup');

      this.stageLineGroup = this.svg.append('g');
      this.stageLineGroup.attr('id', 'stageLineGroup');

      this.http.get('http://localhost:4200/flow/v1/cncf/cncf-demo/aaa/latest/1/runtime/yaml')
        .toPromise()
        .then(response => yaml.load(response.text()))
        .then(wfObj => {
          this.drawWorkflow(wfObj);
          return;
        });
    }, 10000);
  }

  ngOnDestroy() {
    if (this.refreshInterval) {
      clearInterval(this.refreshInterval);
    }
  }
}
