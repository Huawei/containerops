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

// import { Component } from '@angular/core';
import { Component, ElementRef, ViewChild, OnInit } from '@angular/core';
import * as D3 from 'd3';


@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit {
  @ViewChild('design') element: ElementRef;

  private host: D3.Selection;
  private svg: D3.Selection;
  private width: number;
  private height: number;
  private radius: number;
  private htmlElement: HTMLElement;
  private pieData = [1, 2, 3, 4, 5]

  ngOnInit() {
    this.htmlElement = this.element.nativeElement;
    this.host = D3.select(this.htmlElement);
    this.setup();
    this.buildSVG();
  }

  private setup(): void {
    this.width = 250;
    this.height = 250;
    this.radius = Math.min(this.width, this.height) / 2;
  }

  private buildSVG(): void {
    this.host.html('');
    this.svg = this.host.append('svg')
      .attr('viewBox', `0 0 500 500`);

    const myline = this.svg.append('rect');
    myline
      .attr('x', 100)
      .attr('y', 100)
      .attr('width', 200)
      .attr('height', 200);
    myline.style('stroke', 'red');
    myline.style('stroke-width', 5);
    myline.style('fill', 'yellow');
  }


  private buildPie(): void {
    const pie = D3.layout.pie();
    const arcSelection = this.svg.selectAll('.arc')
      .data(pie(this.pieData))
      .enter()
      .append('g')
      .attr('class', 'arc');

    this.populatePie(arcSelection);
  }

  private populatePie(arcSelection: D3.Selection<D3.layout.pie.Arc>): void {
    const innerRadius = this.radius - 50;
    const outerRadius = this.radius - 10;
    const pieColor = D3.scale.category20c();
    const arc = D3.svg.arc<D3.layout.pie.Arc>()
      .outerRadius(outerRadius);
    arcSelection.append('path')
      .attr('d', arc)
      .attr('fill', (datum, index) => {
        return pieColor(`${index}`);
      });

    arcSelection.append('text')
      .attr('transform', (datum: any) => {
        datum.innerRadius = 0;
        datum.outerRadius = outerRadius;
        return 'translate(' + arc.centroid(datum) + ')';
      })
      .text((datum, index) => this.pieData[index])
      .style('text-anchor', 'middle');
  }
}
